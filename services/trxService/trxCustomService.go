package trxLocalService

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/raj847/togrpc/config"
	"github.com/raj847/togrpc/constans"
	"github.com/raj847/togrpc/helpers"
	"github.com/raj847/togrpc/models"
	"github.com/raj847/togrpc/proto/trxLocal"
	"github.com/raj847/togrpc/services"
	"github.com/raj847/togrpc/services/helperService"
	"github.com/raj847/togrpc/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/anypb"
	"gopkg.in/mgo.v2/bson"
)

type trxCustomService struct {
	Service services.UsecaseService
}

func NewTrxCustomService(service services.UsecaseService) trxCustomService {
	return trxCustomService{
		Service: service,
	}
}

func (svc trxCustomService) AddTrxWithoutCardCustom(ctx context.Context, input *trxLocal.RequestTrxCheckInWithoutCard) (*trxLocal.MyResponse, error) {
	var result *trxLocal.Response
	var resultProduct *models.PolicyOuProductWithRules
	resultProduct = nil
	productName := constans.EMPTY_VALUE
	productCode := constans.EMPTY_VALUE

	if err := utils.ValBlankOrNull(input, "refId"); err != nil {
		result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
		return &trxLocal.MyResponse{
			Response: result,
		}, err
	}

	ID, _ := primitive.ObjectIDFromHex(input.RefId)

	resultTrxOutstanding, err := svc.Service.TrxMongoRepo.FindTrxOutstandingByID(ID)
	if err != nil {
		result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
		return &trxLocal.MyResponse{
			Response: result,
		}, err
	}

	merchantKey, err := utils.DecryptMerchantKey(config.MERCHANT_KEY)
	if err != nil {
		result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
		return &trxLocal.MyResponse{
			Response: result,
		}, err
	}

	if input.ProductCode != constans.EMPTY_VALUE {
		outputProduct, err := svc.Service.ProductRepo.FindProductByProductCode(input.ProductCode, merchantKey.OuId)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
			return &trxLocal.MyResponse{
				Response: result,
			}, err
		}

		resultProduct = &outputProduct
		productName = resultProduct.ProductName
		productCode = resultProduct.ProductCode
	}

	document, errDocument := make(chan map[string]interface{}), make(chan error)
	go func() {
		defer close(document)
		close(errDocument)

		var trxInvoiceItemList []models.TrxInvoiceItem
		var productList []models.PolicyOuProductWithRules
		var productData []byte

		checkinDatetimeParse, _ := time.Parse("2006-01-02 15:04:05", input.CheckInDatetime)
		yearMonth := checkinDatetimeParse.Format("060102")
		HourCheckIn := checkinDatetimeParse.Format("15")
		prefixDocNo := utils.RandStringBytesMaskImprSrcChr(4)
		prefix := fmt.Sprintf("%s%s", yearMonth, HourCheckIn)
		autoNumber, err := svc.Service.GenAutoNumRepo.AutonumberValueWithDatatype(constans.DATATYPE_TRX_LOCAL, prefix, 4)
		if err != nil {
			log.Println("ERROR AutonumberValueWithDatatype :  ", err)
			errDocument <- err
			return
		}

		docNo := fmt.Sprintf("%s%s%d", prefixDocNo, autoNumber, merchantKey.OuId)
		randomKey := utils.RandStringBytesMaskImprSrc(16)
		encQrTxt, _ := utils.Encrypt(docNo, randomKey)
		qrCode := fmt.Sprintf("%s%s", encQrTxt.Result, randomKey)

		timeCheckInUnix := checkinDatetimeParse.Unix()
		resultProductList, err := svc.Service.ProductRepo.GetPolicyOuProductList()
		if err != nil {
			log.Println("ERROR GetPolicyOuProductList : ", err)
			errDocument <- err
			return
		}

		requestData, _ := json.Marshal(input)

		if resultProduct != nil {
			productData, _ = json.Marshal(resultProduct)

			trxItems := models.TrxInvoiceItem{
				DocNo:            docNo,
				ProductId:        resultProduct.ProductId,
				ProductCode:      resultProduct.ProductCode,
				ProductName:      resultProduct.ProductName,
				IsPctServiceFee:  resultProduct.IsPctServiceFee,
				ServiceFee:       resultProduct.ServiceFee,
				ServiceFeeMember: resultProduct.ServiceFeeMember,
				Price:            resultProduct.ProductOuWithRules.Price,
				BaseTime:         resultProduct.ProductOuWithRules.BaseTime,
				ProgressiveTime:  resultProduct.ProductOuWithRules.ProgressiveTime,
				ProgressivePrice: resultProduct.ProductOuWithRules.ProgressivePrice,
				IsPct:            resultProduct.ProductOuWithRules.IsPct,
				ProgressivePct:   resultProduct.ProductOuWithRules.ProgressivePct,
				MaxPrice:         resultProduct.ProductOuWithRules.MaxPrice,
				Is24H:            resultProduct.ProductOuWithRules.Is24H,
				OvernightTime:    resultProduct.ProductOuWithRules.OvernightTime,
				OvernightPrice:   resultProduct.ProductOuWithRules.OvernightPrice,
				GracePeriod:      resultProduct.ProductOuWithRules.GracePeriod,
				FlgRepeat:        resultProduct.ProductOuWithRules.FlgRepeat,
			}

			trxInvoiceItemList = append(trxInvoiceItemList, trxItems)

		} else {

			for _, row := range resultProductList {

				// Exclude sync product for scheduling messages
				excProductCodeList := strings.Split(config.EXCSyncProductCode, ",")
				exists, _ := helpers.InArray(row.ProductCode, excProductCodeList)
				if !exists {
					trxItems := models.TrxInvoiceItem{
						DocNo:            docNo,
						ProductId:        row.ProductId,
						ProductCode:      row.ProductCode,
						ProductName:      row.ProductName,
						IsPctServiceFee:  row.IsPctServiceFee,
						ServiceFee:       row.ServiceFee,
						ServiceFeeMember: row.ServiceFeeMember,
						Price:            row.ProductOuWithRules.Price,
						BaseTime:         row.ProductOuWithRules.BaseTime,
						ProgressiveTime:  row.ProductOuWithRules.ProgressiveTime,
						ProgressivePrice: row.ProductOuWithRules.ProgressivePrice,
						IsPct:            row.ProductOuWithRules.IsPct,
						ProgressivePct:   row.ProductOuWithRules.ProgressivePct,
						MaxPrice:         row.ProductOuWithRules.MaxPrice,
						Is24H:            row.ProductOuWithRules.Is24H,
						OvernightTime:    row.ProductOuWithRules.OvernightTime,
						OvernightPrice:   row.ProductOuWithRules.OvernightPrice,
						GracePeriod:      row.ProductOuWithRules.GracePeriod,
						FlgRepeat:        row.ProductOuWithRules.FlgRepeat,
					}

					productList = append(productList, row)
					trxInvoiceItemList = append(trxInvoiceItemList, trxItems)
				}

			}

			productData, _ = json.Marshal(productList)
		}

		trx := models.Trx{
			DocNo:            docNo,
			CheckInDatetime:  input.CheckInDatetime,
			CheckInTime:      timeCheckInUnix,
			GateIn:           input.IpTerminal,
			ExtLocalDatetime: utils.Timestamp(),
			ProductData:      string(productData),
			ProductCode:      productCode,
			RequestData:      string(requestData),
			OuId:             merchantKey.OuId,
			OuName:           merchantKey.OuName,
			OuCode:           merchantKey.OuCode,
			OuSubBranchId:    merchantKey.OuSubBranchId,
			OuSubBranchName:  merchantKey.OuSubBranchName,
			OuSubBranchCode:  merchantKey.OuSubBranchCode,
			MainOuId:         merchantKey.MainOuId,
			MainOuName:       merchantKey.MainOuName,
			MainOuCode:       merchantKey.MainOuCode,
			QrText:           qrCode,
			TrxInvoiceItem:   trxInvoiceItemList,
		}

		_, err = svc.Service.TrxMongoRepo.AddTrxCheckIn(trx)
		if err != nil {
			log.Println("ERROR AddTrxCheckIn :  ", err)
			errDocument <- err
			return
		}

		filter := bson.M{
			"docNo": docNo,
		}

		updateSet := bson.M{
			"$set": bson.M{
				"docNo":          docNo,
				"refDocNo":       resultTrxOutstanding.DocNo,
				"cardNumber":     resultTrxOutstanding.CardNumber,
				"cardNumberUuid": resultTrxOutstanding.CardNumberUUID,
			},
		}

		err = svc.Service.TrxMongoRepo.AddTrxAddInfoInterfaces(filter, updateSet)
		if err != nil {
			log.Println("ERROR AddTrxAddInfoInterfaces :  ", err)
			errDocument <- err
			return
		}

		data := make(map[string]interface{})
		data["docNo"] = docNo
		data["qrCode"] = qrCode

		document <- data

		//log.Println("Default Publish Redis:", config.DEFAULTChannelRedisParking)
		trxHeader, err := json.Marshal(trx)
		//svc.Service.RedisClient.Publish(config.DEFAULTChannelRedisParking, trxHeader)
		helperService.ConsumeTrxForScheduling(svc.Service, string(trxHeader))

	}()

	err = <-errDocument
	if err != nil {
		result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
		return &trxLocal.MyResponse{
			Response: result,
		}, err
	}

	documentTrx := <-document
	responseTrx := &trxLocal.ResponseTrxTicket{
		CheckInDateTime: input.CheckInDatetime,
		QrCode:          documentTrx["qrCode"].(string),
		ProductName:     productName,
		VehicleNumberIn: constans.EMPTY_VALUE,
		DocNo:           documentTrx["docNo"].(string),
	}

	anyResponseTrx, _ := anypb.New(responseTrx)

	if err := anyResponseTrx.UnmarshalTo(responseTrx); err != nil {
		log.Fatalf("Failed to decode: %v", err)
	}
	log.Println("abc", &responseTrx)

	result = helpers.ResponseJSON(true, constans.SUCCESS_CODE, constans.EMPTY_VALUE, anyResponseTrx)
	return &trxLocal.MyResponse{
		Response: result,
	}, nil
}

func (svc trxCustomService) InquiryTrxCardCustom(ctx context.Context, input *trxLocal.RequestInquiryWithCardCustom) (*trxLocal.MyResponse, error) {
	var result *trxLocal.Response

	documentInq := make(map[string]interface{})
	documentInq["memberCode"] = constans.EMPTY_VALUE
	documentInq["memberName"] = constans.EMPTY_VALUE
	documentInq["memberType"] = constans.EMPTY_VALUE
	documentInq["excludeSf"] = false
	documentInq["lastCheckInDatetime"] = constans.EMPTY_VALUE

	inquiryDate := input.InquiryDateTime[:len(input.InquiryDateTime)-9]

	resultMember, exists, err := svc.Service.MemberRepo.IsMemberByAdvanceIndex(input.UuidCard, input.VehicleNumber, inquiryDate, config.MEMBER_BY, false)
	if err != nil {
		result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
		return &trxLocal.MyResponse{
			Response: result,
		}, err
	}

	isFreePass := false
	if exists {
		if resultMember.TypePartner == constans.TYPE_PARTNER_FREE_PASS {
			isFreePass = true
		}
	}

	if !isFreePass {
		resultTrx, err := svc.Service.TrxMongoRepo.FindTrxOutstandingByUUID(input.UuidCard, input.ProductCode)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
			return &trxLocal.MyResponse{
				Response: result,
			}, err
		}
		checkinDate := resultTrx.CheckInDatetime[:len(resultTrx.CheckInDatetime)-9]
		grandTotal := resultTrx.TrxInvoiceItem[0].TotalAmount
		log.Println("Grand Total Before Calc", grandTotal)

		var resultTrxMemberList []models.TrxMember
		if exists {

			resultActiveMemberList, err := svc.Service.MemberRepo.GetMemberActiveListByPeriod(input.UuidCard, input.VehicleNumber, checkinDate, inquiryDate, config.MEMBER_BY, false)
			if err != nil {
				result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
				return &trxLocal.MyResponse{
					Response: result,
				}, err
			}

			for _, rows := range resultActiveMemberList {
				if rows.ProductCode == input.ProductCode {

					trxMemberList := models.TrxMember{
						DocNo:              resultTrx.DocNo,
						PartnerCode:        rows.PartnerCode,
						FirstName:          rows.FirstName,
						LastName:           rows.LastName,
						RoleType:           rows.RoleType,
						PhoneNumber:        rows.PhoneNumber,
						Email:              rows.Email,
						Active:             rows.Active,
						ActiveAt:           rows.ActiveAt,
						NonActiveAt:        rows.NonActiveAt,
						OuId:               rows.OuId,
						TypePartner:        rows.TypePartner,
						CardNumber:         rows.CardNumber,
						VehicleNumber:      rows.VehicleNumber,
						RegisteredDatetime: rows.RegisteredDatetime,
						DateFrom:           rows.DateFrom,
						DateTo:             rows.DateTo,
						ProductId:          rows.ProductId,
						ProductCode:        rows.ProductCode,
					}

					resultTrxMemberList = append(resultTrxMemberList, trxMemberList)
				}
			}

		} else {
			trxMemberList, err := svc.Service.TrxMongoRepo.GetTrxMemberByPeriodListByProductCode(resultTrx.DocNo, input.ProductCode, inquiryDate)
			if err != nil {
				result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
				return &trxLocal.MyResponse{
					Response: result,
				}, err
			}

			resultTrxMemberList = trxMemberList
		}

		if len(resultTrxMemberList) > 0 {
			// Jika dalam transaksi terdapat member, maka saat inquiry
			// data member yg terbaru disimpan ke redis, pada saat konfirmasi
			// Dilakukan pengecekan, apakah nomor transaksi tsb menggunakan member atau bukan?
			memberData := resultTrxMemberList[len(resultTrxMemberList)-1]
			//trxInvoiceDetailList, err := svc.Service.TrxMongoRepo.GetTrxInvoiceDetailsItemsList(resultTrx.DocNo, request.ProductCode)
			//if err != nil {
			//	result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, "Sesi Tidak Ditemukan", nil)
			//	return ctx.JSON(http.StatusBadRequest, result)
			//}

			//if len(trxInvoiceDetailList) > 0 {
			//	invoiceAmount := helpers.CalculateTrxInvoiceDetail(resultTrxMemberList, trxInvoiceDetailList)
			//	log.Println("Invoice Amount:", invoiceAmount)
			//	log.Println("Total Progressive Amount:", resultTrx.TrxInvoiceItem[0].TotalProgressiveAmount)
			//
			//	grandTotal -= invoiceAmount
			//	grandTotal -= resultTrx.TrxInvoiceItem[0].TotalProgressiveAmount
			//} else {
			grandTotal = 0
			//}

			log.Println("Grand Total After Calc Invoice Amount - Total Progressive Amount:", grandTotal)
			if grandTotal == 0 {
				trxMemberStr, _ := json.Marshal(memberData)
				svc.Service.RedisClientLocal.Set(fmt.Sprintf("%s-%s", resultTrx.DocNo, constans.MEMBER), trxMemberStr, 5*time.Minute)

				documentInq["memberCode"] = memberData.PartnerCode
				documentInq["memberName"] = strings.TrimSpace(fmt.Sprintf("%s %s", memberData.FirstName, memberData.LastName))
				documentInq["memberType"] = memberData.TypePartner
			}

		} else {
			expired24h := fmt.Sprintf("%s-%s-%s-%s", resultTrx.OuCode, input.UuidCard, input.ProductCode, input.VehicleNumber)
			resultGetRedis := svc.Service.RedisClientLocal.Get(expired24h)
			grandTotal = 0
			documentInq["excludeSf"] = true
			documentInq["lastCheckInDatetime"] = resultGetRedis.Val()
			if resultGetRedis.Val() == constans.EMPTY_VALUE {
				resultTrxInvoiceItem, err := helpers.FilterTrxInvoiceByIndex(resultTrx.DocNo, input.ProductCode, resultTrx.TrxInvoiceItem)
				if err != nil {
					result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
					return &trxLocal.MyResponse{
						Response: result,
					}, err
				}
				grandTotal = resultTrxInvoiceItem.TotalAmount

				resultHistoryTrx, exists, err := svc.Service.TrxMongoRepo.IsHistoryTrxListByCheckInAndCheckoutDate(input.UuidCard, input.ProductCode, checkinDate, input.VehicleNumber)
				if err != nil {
					result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
					return &trxLocal.MyResponse{
						Response: result,
					}, err
				}

				if exists {
					grandTotal -= resultHistoryTrx.GrandTotal
				}

				documentInq["excludeSf"] = false
			}
		}

		documentInq["id"] = resultTrx.ID.Hex()
		documentInq["docNo"] = resultTrx.DocNo
		documentInq["checkinDatetime"] = resultTrx.CheckInDatetime
		documentInq["cardNumberUuId"] = input.UuidCard
		documentInq["cardNumber"] = constans.EMPTY_VALUE
		documentInq["productCode"] = resultTrx.TrxInvoiceItem[0].ProductCode
		documentInq["productName"] = resultTrx.TrxInvoiceItem[0].ProductName
		documentInq["vehicleNumberIn"] = constans.EMPTY_VALUE
		documentInq["totalAmount"] = grandTotal
		documentInq["ouCode"] = resultTrx.OuCode

	} else {
		merchantKey, err := utils.DecryptMerchantKey(config.MERCHANT_KEY)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
			return &trxLocal.MyResponse{
				Response: result,
			}, err
		}

		resultTrx, exists, err := svc.Service.TrxMongoRepo.IsTrxOutstandingByUUID(input.UuidCard)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
			return &trxLocal.MyResponse{
				Response: result,
			}, err
		}

		resultProduct, err := svc.Service.ProductRepo.FindProductByProductCode(input.ProductCode, merchantKey.OuId)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
			return &trxLocal.MyResponse{
				Response: result,
			}, err
		}

		checkinDatetimeParse, _ := time.Parse("2006-01-02 15:04:05", input.InquiryDateTime)
		yearMonth := checkinDatetimeParse.Format("060102")
		HourCheckIn := checkinDatetimeParse.Format("15")
		prefixDocNo := utils.RandStringBytesMaskImprSrcChr(4)
		prefix := fmt.Sprintf("%s%s", yearMonth, HourCheckIn)
		autoNumber, err := svc.Service.GenAutoNumRepo.AutonumberValueWithDatatype(constans.DATATYPE_TRX_LOCAL, prefix, 4)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
			return &trxLocal.MyResponse{
				Response: result,
			}, err
		}

		docNo := fmt.Sprintf("%s%s%d", prefixDocNo, autoNumber, merchantKey.OuId)
		randomKey := utils.RandStringBytesMaskImprSrc(16)
		encQrTxt, _ := utils.Encrypt(docNo, randomKey)
		qrCode := fmt.Sprintf("%s%s", encQrTxt.Result, randomKey)

		trx := resultTrx
		if !exists {
			var trxInvoiceItemList []models.TrxInvoiceItem
			trxInvoiceItem := models.TrxInvoiceItem{
				DocNo:            docNo,
				ProductId:        resultProduct.ProductId,
				ProductCode:      resultProduct.ProductCode,
				ProductName:      resultProduct.ProductName,
				IsPctServiceFee:  resultProduct.IsPctServiceFee,
				ServiceFee:       resultProduct.ServiceFee,
				ServiceFeeMember: resultProduct.ServiceFeeMember,
				Price:            resultProduct.ProductOuWithRules.Price,
				BaseTime:         resultProduct.ProductOuWithRules.BaseTime,
				ProgressiveTime:  resultProduct.ProductOuWithRules.ProgressiveTime,
				ProgressivePrice: resultProduct.ProductOuWithRules.ProgressivePrice,
				IsPct:            resultProduct.ProductOuWithRules.IsPct,
				ProgressivePct:   resultProduct.ProductOuWithRules.ProgressivePct,
				MaxPrice:         resultProduct.ProductOuWithRules.MaxPrice,
				Is24H:            resultProduct.ProductOuWithRules.Is24H,
				OvernightTime:    resultProduct.ProductOuWithRules.OvernightTime,
				OvernightPrice:   resultProduct.ProductOuWithRules.OvernightPrice,
				GracePeriod:      resultProduct.ProductOuWithRules.GracePeriod,
				FlgRepeat:        resultProduct.ProductOuWithRules.FlgRepeat,
				TotalAmount:      0,
			}
			trxInvoiceItemList = append(trxInvoiceItemList, trxInvoiceItem)

			trx = &models.Trx{
				DocNo:            docNo,
				DocDate:          constans.EMPTY_VALUE,
				CheckInDatetime:  input.InquiryDateTime,
				CheckOutDatetime: constans.EMPTY_VALUE,
				DeviceIdIn:       constans.EMPTY_VALUE,
				DeviceId:         constans.EMPTY_VALUE,
				GateIn:           constans.EMPTY_VALUE,
				GateOut:          constans.EMPTY_VALUE,
				CardNumberUUIDIn: input.UuidCard,
				CardNumberIn:     constans.EMPTY_VALUE,
				CardNumberUUID:   input.UuidCard,
				CardNumber:       constans.EMPTY_VALUE,
				TypeCard:         constans.EMPTY_VALUE,
				BeginningBalance: 0,
				ExtLocalDatetime: input.InquiryDateTime,
				GrandTotal:       0,
				ProductCode:      resultProduct.ProductCode,
				ProductName:      resultProduct.ProductName,
				ProductData:      utils.ToString(resultProduct),
				RequestData:      utils.ToString(input),
				RequestOutData:   utils.ToString(input),
				OuId:             merchantKey.OuId,
				OuName:           merchantKey.OuName,
				OuCode:           merchantKey.OuCode,
				OuSubBranchId:    merchantKey.OuSubBranchId,
				OuSubBranchName:  merchantKey.OuSubBranchName,
				OuSubBranchCode:  merchantKey.OuSubBranchCode,
				MainOuId:         merchantKey.MainOuId,
				MainOuCode:       merchantKey.MainOuCode,
				MainOuName:       merchantKey.MainOuName,
				MemberCode:       resultMember.PartnerCode,
				MemberName:       strings.TrimSpace(fmt.Sprintf("%s %s", resultMember.FirstName, resultMember.LastName)),
				MemberType:       resultMember.TypePartner,
				CheckInTime:      0,
				CheckOutTime:     0,
				DurationTime:     0,
				VehicleNumberIn:  constans.EMPTY_VALUE,
				VehicleNumberOut: constans.EMPTY_VALUE,
				LogTrans:         constans.EMPTY_VALUE,
				MerchantKey:      config.MERCHANT_KEY,
				QrText:           qrCode,
				TrxInvoiceItem:   trxInvoiceItemList,
				FlagSyncData:     false,
				MemberData:       nil,
				TrxAddInfo:       nil,
			}

		} else {
			trx.MemberCode = resultMember.PartnerCode
			trx.MemberName = strings.TrimSpace(fmt.Sprintf("%s %s", resultMember.FirstName, resultMember.LastName))
			trx.MemberType = resultMember.TypePartner
		}

		redisStatus := svc.Service.RedisClientLocal.Set(input.UuidCard, utils.ToString(trx), 1*time.Minute)
		if redisStatus.Err() != nil {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
			return &trxLocal.MyResponse{
				Response: result,
			}, redisStatus.Err()
		}

		documentInq["id"] = constans.TYPE_PARTNER_FREE_PASS
		documentInq["docNo"] = trx.DocNo
		documentInq["checkinDatetime"] = trx.CheckInDatetime
		documentInq["cardNumberUuId"] = input.UuidCard
		documentInq["cardNumber"] = constans.EMPTY_VALUE
		documentInq["productCode"] = resultProduct.ProductCode
		documentInq["productName"] = resultProduct.ProductName
		documentInq["vehicleNumberIn"] = constans.EMPTY_VALUE
		documentInq["memberCode"] = trx.MemberCode
		documentInq["memberName"] = trx.MemberName
		documentInq["memberType"] = trx.MemberType
		documentInq["excludeSf"] = true
		documentInq["totalAmount"] = float64(0)
		documentInq["ouCode"] = trx.OuCode
	}

	responseTrx := &trxLocal.ResultInquiryTrxWithCard{
		Id:              documentInq["id"].(string),
		DocNo:           documentInq["docNo"].(string),
		CardNumberUUID:  documentInq["cardNumberUuId"].(string),
		CardNumber:      documentInq["cardNumber"].(string),
		Nominal:         documentInq["totalAmount"].(float64),
		ProductCode:     documentInq["productCode"].(string),
		ProductName:     documentInq["productName"].(string),
		VehicleNumberIn: input.VehicleNumber,
		QrCode:          constans.EMPTY_VALUE,
		ExcludeSf:       documentInq["totalAmount"].(float64) == 0,
		MemberCode:      documentInq["memberCode"].(string),
		MemberName:      documentInq["memberName"].(string),
		MemberType:      documentInq["memberType"].(string),
		OuCode:          documentInq["ouCode"].(string),
	}

	checkinDatetimeParse, _ := time.Parse("2006-01-02 15:04", documentInq["checkinDatetime"].(string)[:len(documentInq["checkinDatetime"].(string))-3])
	checkoutDatetimeParse, _ := time.Parse("2006-01-02 15:04", input.InquiryDateTime[:len(input.InquiryDateTime)-3])
	duration := utils.ConvDiffTime(checkinDatetimeParse, checkoutDatetimeParse)
	anyDuration, _ := anypb.New(duration)
	responseTrx.Duration = anyDuration

	if documentInq["lastCheckInDatetime"].(string) != constans.EMPTY_VALUE {
		layout := "2006-01-02 15:04"
		lastCheckInDatetimeFormat, _ := time.Parse(layout, documentInq["lastCheckInDatetime"].(string)[:len(documentInq["lastCheckInDatetime"].(string))-3])
		inquiryDatetimeFormat, _ := time.Parse(layout, input.InquiryDateTime[:len(input.InquiryDateTime)-3])
		duration24h := utils.ConvDiffTime(lastCheckInDatetimeFormat, inquiryDatetimeFormat)
		anyDuration, _ := anypb.New(duration24h)
		responseTrx.Duration24H = anyDuration

	}

	anyResponseTrx, _ := anypb.New(responseTrx)

	if err := anyResponseTrx.UnmarshalTo(responseTrx); err != nil {
		log.Fatalf("Failed to decode: %v", err)
	}
	log.Println("abc", &responseTrx)

	result = helpers.ResponseJSON(true, constans.SUCCESS_CODE, constans.EMPTY_VALUE, anyResponseTrx)
	return &trxLocal.MyResponse{
		Response: result,
	}, nil
}

// INI HARUSNYA GAUSA DI GRPCKAN ???
func (svc trxCustomService) ConfirmTrxCustom(ctx context.Context, input *trxLocal.RequestConfirmTrx) (*trxLocal.MyResponse, error) {
	var result *trxLocal.Response
	var r *http.Request
	var resultTrx models.Trx
	datetimeNow := utils.CurrDatetimeNow()
	username := r.Context().Value("username").(string)

	if err := utils.ValBlankOrNull(input, "vehicleNumber", "shiftCode", "excludeSf"); err != nil {
		result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
		return &trxLocal.MyResponse{
			Response: result,
		}, err
	}

	additionalData := make(map[string]interface{})
	additionalData["username"] = username

	if input.Id == constans.TYPE_PARTNER_FREE_PASS {
		redisStatus := svc.Service.RedisClientLocal.Get(input.UuidCard)
		if redisStatus.Err() != nil {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, "", nil)
			return &trxLocal.MyResponse{
				Response: result,
			}, redisStatus.Err()
		}

		if err := json.Unmarshal([]byte(redisStatus.Val()), &resultTrx); err != nil {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
			return &trxLocal.MyResponse{
				Response: result,
			}, err
		}

		request := new(models.RequestConfirmTrx)
		if err := helpers.BindValidateStruct(request); err != nil {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
			return &trxLocal.MyResponse{
				Response: result,
			}, err
		}

		go helperService.CallSyncConfirmTrxForMemberFreePassCustom(*request, resultTrx, svc.Service, additionalData)
	} else if strings.Contains(input.Id, "SERVER") {
		redisStatus := svc.Service.RedisClientLocal.Get(input.Id)
		if redisStatus.Err() != nil {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, "", nil)
			return &trxLocal.MyResponse{
				Response: result,
			}, redisStatus.Err()
		}

		request := new(models.RequestConfirmTrx)
		if err := helpers.BindValidateStruct(request); err != nil {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
			return &trxLocal.MyResponse{
				Response: result,
			}, err
		}

		if err := json.Unmarshal([]byte(redisStatus.Val()), &resultTrx); err != nil {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
			return &trxLocal.MyResponse{
				Response: result,
			}, err
		}

		go helperService.CallSyncConfirmTrxToCloudCustom(nil, *request, resultTrx, svc.Service, datetimeNow, additionalData, constans.NO)
	} else if strings.Contains(input.Id, "P3") {
		request := new(models.RequestConfirmTrx)
		if err := helpers.BindValidateStruct(request); err != nil {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
			return &trxLocal.MyResponse{
				Response: result,
			}, err
		}
		redisStatus := svc.Service.RedisClientLocal.Get(request.ID)
		if redisStatus.Err() != nil {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, "", nil)
			return &trxLocal.MyResponse{
				Response: result,
			}, redisStatus.Err()
		}

		if err := json.Unmarshal([]byte(redisStatus.Val()), &resultTrx); err != nil {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
			return &trxLocal.MyResponse{
				Response: result,
			}, err
		}

		go helperService.CallSyncConfirmTrxForMemberFreePassCustom(*request, resultTrx, svc.Service, additionalData)
	} else {
		ID, _ := primitive.ObjectIDFromHex(input.Id)
		resultTrxOutstanding, err := svc.Service.TrxMongoRepo.FindTrxOutstandingByID(ID)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
			return &trxLocal.MyResponse{
				Response: result,
			}, err
		}

		request := new(models.RequestConfirmTrx)
		if err := helpers.BindValidateStruct(request); err != nil {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
			return &trxLocal.MyResponse{
				Response: result,
			}, err
		}

		resultTrx = resultTrxOutstanding

		go helperService.CallSyncConfirmTrxToCloudCustom(&ID, *request, resultTrx, svc.Service, datetimeNow, additionalData, constans.NO)
	}

	responseConfirm := &trxLocal.ResponseConfirm{
		DocNo:            resultTrx.DocNo,
		ProductData:      resultTrx.ProductData,
		CardNumber:       input.CardNumber,
		CheckInDatetime:  resultTrx.CheckInDatetime,
		CheckOutDatetime: input.CheckOutDatetime,
		VehicleNumberIn:  input.VehicleNumber,
		VehicleNumberOut: input.VehicleNumber,
		UuidCard:         input.UuidCard,
		ShowQRISArea:     constans.EMPTY_VALUE,
		CurrentBalance:   input.CurrentBalance,
		GrandTotal:       input.GrandTotal,
		OuCode:           resultTrx.OuCode,
		OuName:           resultTrx.OuName,
		Address:          config.ADDRESS,
	}

	anyResponseTrx, _ := anypb.New(responseConfirm)

	result = helpers.ResponseJSON(true, constans.SUCCESS_CODE, constans.EMPTY_VALUE, anyResponseTrx)
	return &trxLocal.MyResponse{
		Response: result,
	}, nil
}
