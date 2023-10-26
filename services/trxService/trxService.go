package trxService

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/raj847/togrpc/config"
	"github.com/raj847/togrpc/constans"
	"github.com/raj847/togrpc/helpers"
	"github.com/raj847/togrpc/models"
	prod "github.com/raj847/togrpc/proto/product"
	trx "github.com/raj847/togrpc/proto/trxLocal"
	"github.com/raj847/togrpc/proto/user"
	"github.com/raj847/togrpc/services"
	"github.com/raj847/togrpc/services/helperService"
	"github.com/raj847/togrpc/utils"
	"strings"
	"time"

	goproto "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"

	"google.golang.org/protobuf/types/known/wrapperspb"
	"gopkg.in/mgo.v2/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type trxService struct {
	Service services.UsecaseService
}

func NewTrxService(service services.UsecaseService) trxService {
	return trxService{
		Service: service,
	}
}

func (svc trxService) AddTrxWithCard(ctx context.Context, input *trx.RequestTrxCheckin) (*trx.MyResponse, error) {
	var result *trx.Response
	var resultProduct *prod.PolicyOuProductWithRules
	resultProduct = nil
	productName := constans.EMPTY_VALUE
	productCode := constans.EMPTY_VALUE

	if err := helpers.BindValidateStruct(input); err != nil {
		result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, nil
	}

	merchantKey, err := utils.DecryptMerchantKey(config.MERCHANT_KEY)
	if err != nil {
		result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	tempMerchantKey := &trx.MerchantKey{
		ID:              merchantKey.ID,
		OuId:            merchantKey.OuId,
		OuName:          merchantKey.OuName,
		OuCode:          merchantKey.OuCode,
		OuSubBranchId:   merchantKey.OuSubBranchId,
		OuSubBranchName: merchantKey.OuSubBranchName,
		OuSubBranchCode: merchantKey.OuSubBranchCode,
		MainOuId:        merchantKey.MainOuId,
		MainOuCode:      merchantKey.MainOuCode,
		MainOuName:      merchantKey.MainOuName,
	}

	if input.ProductCode != constans.EMPTY_VALUE {
		outputProduct, err := svc.Service.ProductRepo.FindProductByProductCode(input.ProductCode, tempMerchantKey.OuId)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, err.Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}
		tempProd := &prod.ProductOuWithRules{
			OuId:             outputProduct.ProductOuWithRules.OuId,
			ProductId:        outputProduct.ProductOuWithRules.ProductId,
			Price:            outputProduct.ProductOuWithRules.Price,
			BaseTime:         outputProduct.ProductOuWithRules.BaseTime,
			ProgressiveTime:  outputProduct.ProductOuWithRules.ProgressiveTime,
			ProgressivePrice: outputProduct.ProductOuWithRules.ProgressivePrice,
			IsPct:            outputProduct.ProductOuWithRules.IsPct,
			ProgressivePct:   outputProduct.ProductOuWithRules.ProgressivePct,
			MaxPrice:         outputProduct.ProductOuWithRules.MaxPrice,
			Is24H:            outputProduct.ProductOuWithRules.Is24H,
			OvernightTime:    outputProduct.ProductOuWithRules.OvernightTime,
			OvernightPrice:   outputProduct.ProductOuWithRules.OvernightPrice,
			GracePeriod:      outputProduct.ProductOuWithRules.GracePeriod,
			FlgRepeat:        outputProduct.ProductOuWithRules.FlgRepeat,
		}
		resultProduct = &prod.PolicyOuProductWithRules{
			OuId:                  outputProduct.OuId,
			OuCode:                outputProduct.OuCode,
			OuName:                outputProduct.OuName,
			ProductId:             outputProduct.ProductId,
			ProductCode:           outputProduct.ProductCode,
			ProductName:           outputProduct.ProductName,
			ServiceFee:            outputProduct.ServiceFee,
			IsPctServiceFee:       outputProduct.IsPctServiceFee,
			IsPctServiceFeeMember: outputProduct.IsPctServiceFeeMember,
			ServiceFeeMember:      outputProduct.ServiceFeeMember,
			ProductRules:          tempProd,
		}
		productName = resultProduct.ProductName
		productCode = resultProduct.ProductCode
	}

	// Digunakan utk melakukan pengecekan apakah kartu tsb sedang digunakan atau tdk?
	//responseCode, err := svc.Service.TrxMongoRepo.ValTrxExistByUUIDCard(request.UUIDCard)
	//if err != nil {
	//	result = helpers.ResponseJSON(false, responseCode, err.Error(), nil)
	//	return ctx.JSON(http.StatusBadRequest, result)
	//}

	isSessionAlreadyExists := false
	errCode := constans.EMPTY_VALUE
	tappingDatetime := utils.Timestamp()
	dateNow := utils.DateNow()
	resultTrxOutstanding, exists := svc.Service.TrxMongoRepo.IsTrxOutstandingExistByUUIDCard(input.UuidCard)
	var tempTrxInvoiceItems []*trx.TrxInvoiceItem
	for _, v := range resultTrxOutstanding.TrxInvoiceItem {
		tempTrxInvoiceItem := &trx.TrxInvoiceItem{
			DocNo:                  v.DocNo,
			ProductId:              v.ProductId,
			ProductCode:            v.ProductCode,
			ProductName:            v.ProductName,
			IsPctServiceFee:        v.IsPctServiceFee,
			ServiceFee:             v.ServiceFee,
			ServiceFeeMember:       v.ServiceFeeMember,
			Price:                  v.Price,
			BaseTime:               v.BaseTime,
			ProgressiveTime:        v.ProgressiveTime,
			ProgressivePrice:       v.ProgressivePrice,
			IsPct:                  v.IsPct,
			ProgressivePct:         v.ProgressivePct,
			MaxPrice:               v.MaxPrice,
			Is24H:                  v.Is24H,
			OvernightTime:          v.OvernightTime,
			OvernightPrice:         v.OvernightPrice,
			GracePeriod:            v.GracePeriod,
			FlgRepeat:              v.FlgRepeat,
			TotalAmount:            v.TotalAmount,
			TotalProgressiveAmount: v.TotalProgressiveAmount,
		}
		tempTrxInvoiceItems = append(tempTrxInvoiceItems, tempTrxInvoiceItem)
	}
	tempMemberData := &trx.TrxMember{
		DocNo:       resultTrxOutstanding.MemberData.DocNo,
		PartnerCode: resultTrxOutstanding.MemberData.PartnerCode,
		FirstName:   resultTrxOutstanding.MemberData.FirstName,
		LastName:    resultTrxOutstanding.MemberData.LastName,
		RoleType:    resultTrxOutstanding.MemberData.RoleType,
		PhoneNumber: resultTrxOutstanding.MemberData.PhoneNumber,
		Email:       resultTrxOutstanding.MemberData.Email,
		Active:      resultTrxOutstanding.MemberData.Active,
		ActiveAt:    resultTrxOutstanding.MemberData.ActiveAt,
		NonActiveAt: func() *wrapperspb.StringValue {
			if resultTrxOutstanding.MemberData.NonActiveAt != nil {
				return &wrapperspb.StringValue{Value: *resultTrxOutstanding.MemberData.NonActiveAt}
			}
			return nil
		}(),
		OuId:               resultTrxOutstanding.MemberData.OuId,
		TypePartner:        resultTrxOutstanding.MemberData.TypePartner,
		CardNumber:         resultTrxOutstanding.MemberData.CardNumber,
		VehicleNumber:      resultTrxOutstanding.MemberData.VehicleNumber,
		RegisteredDatetime: resultTrxOutstanding.MemberData.RegisteredDatetime,
		DateFrom:           resultTrxOutstanding.MemberData.DateFrom,
		DateTo:             resultTrxOutstanding.MemberData.DateTo,
		ProductId:          resultTrxOutstanding.MemberData.ProductId,
		ProductCode:        resultTrxOutstanding.MemberData.ProductCode,
	}

	convertedTrxAddInfo := make(map[string]*anypb.Any)
	for k, v := range resultTrxOutstanding.TrxAddInfo {
		if msg, ok := v.(goproto.Message); ok {
			anyVal, err := anypb.New(msg)
			if err != nil {
				continue
			}
			convertedTrxAddInfo[k] = anyVal
		}
	}

	tempRequestAddTrxInvoiceDetailItem := &trx.TrxInvoiceDetailItem{
		DocNo:         resultTrxOutstanding.RequestAddTrxInvoiceDetailItem.DocNo,
		ProductCode:   resultTrxOutstanding.RequestAddTrxInvoiceDetailItem.ProductCode,
		InvoiceAmount: resultTrxOutstanding.RequestAddTrxInvoiceDetailItem.InvoiceAmount,
		CreatedAt:     resultTrxOutstanding.RequestAddTrxInvoiceDetailItem.CreatedAt,
		CreatedDate:   resultTrxOutstanding.RequestAddTrxInvoiceDetailItem.CreatedDate,
	}

	tempTrxOutstanding := &trx.Trx{
		DocNo:                          resultTrxOutstanding.DocNo,
		DocDate:                        resultTrxOutstanding.DocDate,
		PaymentRefDocno:                resultTrxOutstanding.PaymentRefDocNo,
		CheckinDateTime:                resultTrxOutstanding.CheckInDatetime,
		CheckoutDateTime:               resultTrxOutstanding.CheckOutDatetime,
		DeviceIdIn:                     resultTrxOutstanding.DeviceIdIn,
		DeviceId:                       resultTrxOutstanding.DeviceId,
		GetIn:                          resultTrxOutstanding.GateIn,
		GetOut:                         resultTrxOutstanding.GateOut,
		CardNumberUUIDIn:               resultTrxOutstanding.CardNumberUUIDIn,
		CardNumberIn:                   resultTrxOutstanding.CardNumberIn,
		CardNumberUUID:                 resultTrxOutstanding.CardNumberUUID,
		CardNumber:                     resultTrxOutstanding.CardNumber,
		TypeCard:                       resultTrxOutstanding.TypeCard,
		BeginningBalance:               resultTrxOutstanding.BeginningBalance,
		ExtLocalDateTime:               resultTrxOutstanding.ExtLocalDatetime,
		ChargeAmount:                   resultTrxOutstanding.ChargeAmount,
		GrandTotal:                     resultTrxOutstanding.GrandTotal,
		ProductCode:                    resultTrxOutstanding.ProductCode,
		ProductName:                    resultTrxOutstanding.ProductName,
		ProductData:                    resultTrxOutstanding.ProductData,
		RequestData:                    resultTrxOutstanding.RequestData,
		RequestOutData:                 resultTrxOutstanding.RequestOutData,
		OuId:                           resultTrxOutstanding.OuId,
		OuName:                         resultTrxOutstanding.OuName,
		OuCode:                         resultTrxOutstanding.OuCode,
		OuSubBranchId:                  resultTrxOutstanding.OuSubBranchId,
		USubBranchName:                 resultTrxOutstanding.OuSubBranchName,
		OuSubBranchCode:                resultTrxOutstanding.OuSubBranchCode,
		MainOuId:                       resultTrxOutstanding.MainOuId,
		MainOuCode:                     resultTrxOutstanding.MainOuCode,
		MainOuName:                     resultTrxOutstanding.MainOuName,
		MemberCode:                     resultTrxOutstanding.MemberCode,
		MemberName:                     resultTrxOutstanding.MemberName,
		MemberType:                     resultTrxOutstanding.MemberType,
		MemberStatus:                   resultTrxOutstanding.MemberStatus,
		MemberExpiredDate:              resultTrxOutstanding.MemberExpiredDate,
		CheckInTime:                    resultTrxOutstanding.CheckInTime,
		CheckOutTime:                   resultTrxOutstanding.CheckOutTime,
		DurationTime:                   resultTrxOutstanding.DurationTime,
		VehicleNumberIn:                resultTrxOutstanding.VehicleNumberIn,
		VehicleNumberOut:               resultTrxOutstanding.VehicleNumberOut,
		LogTrans:                       resultTrxOutstanding.LogTrans,
		MerchantKey:                    resultTrxOutstanding.MerchantKey,
		QrText:                         resultTrxOutstanding.QrText,
		QrA2P:                          resultTrxOutstanding.QrA2P,
		QrTextPaymentOnline:            resultTrxOutstanding.QrTextPaymentOnline,
		TrxInvoiceItem:                 tempTrxInvoiceItems,
		FlagSyncData:                   resultTrxOutstanding.FlagSyncData,
		MemberData:                     tempMemberData,
		TrxAddInfo:                     convertedTrxAddInfo,
		FlagTrxFromCloud:               resultTrxOutstanding.FlagTrxFromCloud,
		IsRsyncDataTrx:                 resultTrxOutstanding.IsRsyncDataTrx,
		ExcludeSf:                      resultTrxOutstanding.ExcludeSf,
		FlagCharge:                     resultTrxOutstanding.FlagCharge,
		ChargeType:                     resultTrxOutstanding.ChargeType,
		RequestAddTrxInvoiceDetailItem: tempRequestAddTrxInvoiceDetailItem,
		LastUpdatedAt:                  resultTrxOutstanding.LastUpdatedAt,
	}
	if exists {
		isSessionAlreadyExists = exists
		errCode, err = helperService.CallCheckTrxAlreadyExists(tempTrxOutstanding.DocNo, svc.Service)
		if err != nil {
			isSessionAlreadyExists = true
		}
	}

	if isSessionAlreadyExists && errCode != constans.SUCCESS_CODE {

		checkInDate := input.CheckInDatetime[:len(input.CheckInDatetime)-9]
		resultMember, exists, err := svc.Service.MemberRepo.IsMemberByAdvanceIndex(input.UuidCard, constans.EMPTY_VALUE, checkInDate, config.MEMBER_BY, false)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, err.Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}

		if exists {
			if resultMember.TypePartner == constans.TYPE_PARTNER_FREE_PASS {
				resultTrx, err := svc.Service.TrxMongoRepo.FindIDTrxOutstandingByCard(input.UuidCard)
				if err != nil {
					result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, err.Error(), nil)
					return &trx.MyResponse{
						Response: result,
					}, err
				}
				tempId := &trx.TrxWithId{
					Id: resultTrx.ID,
				}

				// ID, _ := primitive.ObjectIDFromHex(tempId.Id)
				responseTrx := &trx.ResponseTrxTicket{
					Id:              tempId.Id,
					CheckInDateTime: input.CheckInDatetime,
					QrCode:          tempTrxOutstanding.QrText,
					ProductName:     productName,
					VehicleNumberIn: constans.EMPTY_VALUE,
					DocNo:           tempTrxOutstanding.DocNo,
					OuCode:          tempMerchantKey.OuCode,
					OuName:          tempMerchantKey.OuName,
					Address:         config.ADDRESS,
				}
				// Convert responseTrx to anypb.Any
				anyResponseTrx, _ := anypb.New(responseTrx)

				result = helpers.ResponseJSON(true, constans.SUCCESS_CODE, constans.EMPTY_VALUE, anyResponseTrx)
				return &trx.MyResponse{
					Response: result,
				}, nil
			}
		}

		// Get Auto Clear status
		resultRedis := svc.Service.RedisClientLocal.Get("AUTO_CLEAR")
		if resultRedis.Err() != nil || resultRedis.Val() == constans.EMPTY_VALUE || resultRedis.Val() == "0" {
			log.Println("Skip auto clear transaction")
			result = helpers.ResponseJSON(false, errCode, "Sesi Kartu Masih Digunakan", nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}

		checkInDatetime, _ := time.Parse("2006-01-02 15:04", tempTrxOutstanding.CheckinDateTime[:len(tempTrxOutstanding.CheckinDateTime)-3])
		tappingDatetimeParse, _ := time.Parse("2006-01-02 15:04", tappingDatetime[:len(tappingDatetime)-3])

		// Check checkin datetime must be lower than date in tapping datetime
		// If yes do clear session for data transaction
		//if !checkInDatetime.Equal(tappingDatetimeParse) && checkInDatetime.Before(tappingDatetimeParse) {
		//

		log.Println("[DO] CheckIn Datetime:", checkInDatetime, "Tapping Datetime:", tappingDatetimeParse)
		//_, exists, err := svc.Service.TrxMongoRepo.IsTrxOutstandingForClearSession(resultTrxOutstanding.DocNo)
		//if err != nil {
		//	result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, err.Error(), nil)
		//	return ctx.JSON(http.StatusBadRequest, result)
		//}

		//if !exists {
		//ppp
		trxOutstandingForClearSession := models.TrxOutstandingForClearSession{
			RefDocNo:               tempTrxOutstanding.DocNo,
			TappingDate:            dateNow,
			TappingDatetime:        tappingDatetime,
			CardNumberUuid:         tempTrxOutstanding.CardNumberUUIDIn,
			FlagClearSession:       false,
			ClearDatetime:          constans.EMPTY_VALUE,
			TrxOutstandingSnapshot: resultTrxOutstanding,
		}

		log.Println("AddTrxOutstandingForClearSession", utils.ToString(trxOutstandingForClearSession))
		_, err = svc.Service.TrxMongoRepo.AddTrxOutstandingForClearSession(trxOutstandingForClearSession)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, err.Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}

		err = svc.Service.TrxMongoRepo.RemoveTrxByDocNo(tempTrxOutstanding.DocNo)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, err.Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}

		//}

		err = svc.Service.TrxMongoRepo.RemoveTrxByDocNo(tempTrxOutstanding.DocNo)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, err.Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}
		//}

		//result = helpers.ResponseJSON(false, errCode, "Sesi Kartu Masih Digunakan", nil)
		//return ctx.JSON(http.StatusBadRequest, result)
		//}

	}

	//CREATE QR PARKING
	document, errDocument := make(chan map[string]interface{}), make(chan error)
	go func() {
		defer close(document)
		close(errDocument)
		var trxInvoiceItemLists []models.TrxInvoiceItem
		var trxInvoiceItemList []*trx.TrxInvoiceItem
		var productList []*prod.PolicyOuProductWithRules
		// var trxMember *trxLocal.TrxMember
		var trxMembers *models.TrxMember
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

		docNo := fmt.Sprintf("%s%s%d", prefixDocNo, autoNumber, tempMerchantKey.OuId)
		randomKey := utils.RandStringBytesMaskImprSrc(16)
		encQrTxt, _ := utils.Encrypt(docNo, randomKey)
		qrCode := fmt.Sprintf("%s%s", encQrTxt.Result, randomKey)

		timeCheckInUnix := checkinDatetimeParse.Unix()

		requestData, _ := json.Marshal(input)
		if resultProduct != nil {
			productData, _ = json.Marshal(resultProduct)

			trxItems := &trx.TrxInvoiceItem{
				DocNo:            docNo,
				ProductId:        resultProduct.ProductId,
				ProductCode:      resultProduct.ProductCode,
				ProductName:      resultProduct.ProductName,
				IsPctServiceFee:  resultProduct.IsPctServiceFee,
				ServiceFee:       resultProduct.ServiceFee,
				ServiceFeeMember: resultProduct.ServiceFeeMember,
				Price:            resultProduct.ProductRules.Price,
				BaseTime:         resultProduct.ProductRules.BaseTime,
				ProgressiveTime:  resultProduct.ProductRules.ProgressiveTime,
				ProgressivePrice: resultProduct.ProductRules.ProgressivePrice,
				IsPct:            resultProduct.ProductRules.IsPct,
				ProgressivePct:   resultProduct.ProductRules.ProgressivePct,
				MaxPrice:         resultProduct.ProductRules.MaxPrice,
				Is24H:            resultProduct.ProductRules.Is24H,
				OvernightTime:    resultProduct.ProductRules.OvernightTime,
				OvernightPrice:   resultProduct.ProductRules.OvernightPrice,
				GracePeriod:      resultProduct.ProductRules.GracePeriod,
				FlgRepeat:        resultProduct.ProductRules.FlgRepeat,
			}

			trxInvoiceItemList = append(trxInvoiceItemList, trxItems)

		} else {

			resultProductList, err := svc.Service.ProductRepo.GetPolicyOuProductList()
			if err != nil {
				return
			}
			resProdLists := []*prod.PolicyOuProductWithRules{}
			// var resProdList *prod.PolicyOuProductWithRules
			for _, v := range resultProductList {
				resProdListt := &prod.ProductOuWithRules{
					OuId:             v.ProductOuWithRules.OuId,
					ProductId:        v.ProductOuWithRules.ProductId,
					Price:            v.ProductOuWithRules.Price,
					BaseTime:         v.ProductOuWithRules.BaseTime,
					ProgressiveTime:  v.ProductOuWithRules.ProgressiveTime,
					ProgressivePrice: v.ProductOuWithRules.ProgressivePrice,
					IsPct:            v.ProductOuWithRules.IsPct,
					ProgressivePct:   v.ProductOuWithRules.ProgressivePct,
					MaxPrice:         v.ProductOuWithRules.MaxPrice,
					Is24H:            v.ProductOuWithRules.Is24H,
					OvernightTime:    v.ProductOuWithRules.OvernightTime,
					OvernightPrice:   v.ProductOuWithRules.OvernightPrice,
					GracePeriod:      v.ProductOuWithRules.GracePeriod,
					FlgRepeat:        v.ProductOuWithRules.FlgRepeat,
				}
				resProdList := &prod.PolicyOuProductWithRules{
					OuId:                  v.OuId,
					OuCode:                v.OuCode,
					OuName:                v.OuName,
					ProductId:             v.ProductId,
					ProductCode:           v.ProductCode,
					ProductName:           v.ProductName,
					ServiceFee:            v.ServiceFee,
					IsPctServiceFee:       v.IsPctServiceFee,
					IsPctServiceFeeMember: v.IsPctServiceFeeMember,
					ServiceFeeMember:      v.ServiceFeeMember,
					ProductRules:          resProdListt,
				}
				resProdLists = append(resProdLists, resProdList)
			}

			for _, row := range resProdLists {

				// Exclude sync product for scheduling messages
				excProductCodeList := strings.Split(config.EXCSyncProductCode, ",")
				exists, _ = helpers.InArray(row.ProductCode, excProductCodeList)

				if !exists {
					trxItems := trx.TrxInvoiceItem{
						DocNo:            docNo,
						ProductId:        row.ProductId,
						ProductCode:      row.ProductCode,
						ProductName:      row.ProductName,
						IsPctServiceFee:  row.IsPctServiceFee,
						ServiceFee:       row.ServiceFee,
						ServiceFeeMember: row.ServiceFeeMember,
						Price:            row.ProductRules.Price,
						BaseTime:         row.ProductRules.BaseTime,
						ProgressiveTime:  row.ProductRules.ProgressiveTime,
						ProgressivePrice: row.ProductRules.ProgressivePrice,
						IsPct:            row.ProductRules.IsPct,
						ProgressivePct:   row.ProductRules.ProgressivePct,
						MaxPrice:         row.ProductRules.MaxPrice,
						Is24H:            row.ProductRules.Is24H,
						OvernightTime:    row.ProductRules.OvernightTime,
						OvernightPrice:   row.ProductRules.OvernightPrice,
						GracePeriod:      row.ProductRules.GracePeriod,
						FlgRepeat:        row.ProductRules.FlgRepeat,
					}

					productList = append(productList, row)
					trxInvoiceItemList = append(trxInvoiceItemList, &trxItems)
				}
			}

			productData, _ = json.Marshal(productList)
		}
		//ppp
		trxMembers = nil

		//if config.MEMBER_BY == constans.VALIDATE_MEMBER_CARD || config.MEMBER_BY == constans.VALIDATE_MEMBER_MIX {
		//resultMember, exists := svc.Service.MemberRepo.IsMemberActiveExistsByUUIDCard(request.UUIDCard, utils.DateNow())
		//if exists {
		//	trxMember = &models.TrxMember{
		//		DocNo:              docNo,
		//		PartnerCode:        resultMember.PartnerCode,
		//		FirstName:          resultMember.FirstName,
		//		LastName:           resultMember.LastName,
		//		RoleType:           resultMember.RoleType,
		//		PhoneNumber:        resultMember.PhoneNumber,
		//		Email:              resultMember.Email,
		//		Active:             resultMember.Active,
		//		ActiveAt:           resultMember.ActiveAt,
		//		NonActiveAt:        resultMember.NonActiveAt,
		//		OuId:               resultMember.OuId,
		//		TypePartner:        resultMember.TypePartner,
		//		CardNumber:         resultMember.CardNumber,
		//		VehicleNumber:      resultMember.VehicleNumber,
		//		RegisteredDatetime: resultMember.RegisteredDatetime,
		//		DateFrom:           resultMember.DateFrom,
		//		DateTo:             resultMember.DateTo,
		//		ProductId:          resultMember.ProductId,
		//		ProductCode:        resultMember.ProductCode,
		//	}
		//}
		//}

		//ppp
		trx := models.Trx{
			DocNo:            docNo,
			DocDate:          constans.EMPTY_VALUE,
			CheckInDatetime:  input.CheckInDatetime,
			CheckInTime:      timeCheckInUnix,
			DeviceIdIn:       input.DeviceId,
			GateIn:           input.IpTerminal,
			CardNumberUUIDIn: input.UuidCard,
			CardNumberIn:     input.CardNumber,
			CardNumberUUID:   input.UuidCard,
			CardNumber:       input.CardNumber,
			TypeCard:         input.TypeCard,
			BeginningBalance: input.BeginningBalance,
			ExtLocalDatetime: utils.Timestamp(),
			ProductData:      string(productData),
			RequestData:      string(requestData),
			ProductCode:      productCode,
			OuId:             tempMerchantKey.OuId,
			OuName:           tempMerchantKey.OuName,
			OuCode:           tempMerchantKey.OuCode,
			OuSubBranchId:    tempMerchantKey.OuSubBranchId,
			OuSubBranchName:  tempMerchantKey.OuSubBranchName,
			OuSubBranchCode:  tempMerchantKey.OuSubBranchCode,
			MainOuId:         tempMerchantKey.MainOuId,
			MainOuName:       tempMerchantKey.MainOuName,
			MainOuCode:       tempMerchantKey.MainOuCode,
			QrText:           qrCode,
			TrxInvoiceItem:   trxInvoiceItemLists,
			MemberData:       trxMembers,
		}

		// trxLocal := trxLocal.Trx{
		// 	DocNo:            docNo,
		// 	DocDate:          constans.EMPTY_VALUE,
		// 	CheckinDateTime:  input.CheckInDatetime,
		// 	CheckInTime:      timeCheckInUnix,
		// 	DeviceIdIn:       input.DeviceId,
		// 	GetIn:           input.IpTerminal,
		// 	CardNumberUUIDIn: input.UuidCard,
		// 	CardNumberIn:     input.CardNumber,
		// 	CardNumberUUID:   input.UuidCard,
		// 	CardNumber:       input.CardNumber,
		// 	TypeCard:         input.TypeCard,
		// 	BeginningBalance: input.BeginningBalance,
		// 	ExtLocalDateTime: utils.Timestamp(),
		// 	ProductData:      string(productData),
		// 	RequestData:      string(requestData),
		// 	ProductCode:      productCode,
		// 	OuId:             tempMerchantKey.OuId,
		// 	OuName:           tempMerchantKey.OuName,
		// 	OuCode:           tempMerchantKey.OuCode,
		// 	OuSubBranchId:    tempMerchantKey.OuSubBranchId,
		// 	USubBranchName:  tempMerchantKey.OuSubBranchName,
		// 	OuSubBranchCode:  tempMerchantKey.OuSubBranchCode,
		// 	MainOuId:         tempMerchantKey.MainOuId,
		// 	MainOuName:       tempMerchantKey.MainOuName,
		// 	MainOuCode:       tempMerchantKey.MainOuCode,
		// 	QrText:           qrCode,
		// 	TrxInvoiceItem:   trxInvoiceItemList,
		// 	MemberData:       trxMember,
		// }

		//ppp
		idObject, err := svc.Service.TrxMongoRepo.AddTrxCheckIn(trx)
		if err != nil {
			log.Println("ERROR AddTrxCheckIn :  ", err)
			errDocument <- err
			return
		}

		data := make(map[string]interface{})
		data["ID"] = idObject
		data["docNo"] = docNo
		data["qrCode"] = qrCode

		document <- data

		// Send data transaction to scheduling message
		// if there connection internet data will be send to cloud server
		// if there not connection internet data not send to cloud server
		// but data will processed to local server
		trxHeader, _ := json.Marshal(trx)
		helperService.ConsumeTrxForScheduling(svc.Service, string(trxHeader))

	}()

	err = <-errDocument
	if err != nil {
		result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	documentTrx := <-document
	responseTrx := trx.ResponseTrxTicket{
		Id:              documentTrx["ID"].(*primitive.ObjectID).Hex(),
		CheckInDateTime: input.CheckInDatetime,
		QrCode:          documentTrx["qrCode"].(string),
		ProductName:     productName,
		VehicleNumberIn: constans.EMPTY_VALUE,
		DocNo:           documentTrx["docNo"].(string),
		OuCode:          tempMerchantKey.OuCode,
		OuName:          tempMerchantKey.OuName,
		Address:         config.ADDRESS,
	}

	anyResponseTrx, _ := anypb.New(&responseTrx)

	result = helpers.ResponseJSON(true, constans.SUCCESS_CODE, constans.EMPTY_VALUE, anyResponseTrx)
	return &trx.MyResponse{
		Response: result,
	}, nil
}

func (svc trxService) AddTrxWithoutCard(ctx context.Context, input *trx.RequestTrxCheckInWithoutCard) (*trx.MyResponse, error) {
	var result *trx.Response
	var resultProduct *prod.PolicyOuProductWithRules
	resultProduct = nil
	productName := constans.EMPTY_VALUE
	productCode := constans.EMPTY_VALUE

	if err := helpers.BindValidateStruct(input); err != nil {
		result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	merchantKey, err := utils.DecryptMerchantKey(config.MERCHANT_KEY)
	if err != nil {
		result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	tempMerchantKey := &trx.MerchantKey{
		ID:              merchantKey.ID,
		OuId:            merchantKey.OuId,
		OuName:          merchantKey.OuName,
		OuCode:          merchantKey.OuCode,
		OuSubBranchId:   merchantKey.OuSubBranchId,
		OuSubBranchName: merchantKey.OuSubBranchName,
		OuSubBranchCode: merchantKey.OuSubBranchCode,
		MainOuId:        merchantKey.MainOuId,
		MainOuCode:      merchantKey.MainOuCode,
		MainOuName:      merchantKey.MainOuName,
	}

	if input.ProductCode != constans.EMPTY_VALUE {
		outputProduct, err := svc.Service.ProductRepo.FindProductByProductCode(input.ProductCode, tempMerchantKey.OuId)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, err.Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}
		tempProd := &prod.ProductOuWithRules{
			OuId:             outputProduct.ProductOuWithRules.OuId,
			ProductId:        outputProduct.ProductOuWithRules.ProductId,
			Price:            outputProduct.ProductOuWithRules.Price,
			BaseTime:         outputProduct.ProductOuWithRules.BaseTime,
			ProgressiveTime:  outputProduct.ProductOuWithRules.ProgressiveTime,
			ProgressivePrice: outputProduct.ProductOuWithRules.ProgressivePrice,
			IsPct:            outputProduct.ProductOuWithRules.IsPct,
			ProgressivePct:   outputProduct.ProductOuWithRules.ProgressivePct,
			MaxPrice:         outputProduct.ProductOuWithRules.MaxPrice,
			Is24H:            outputProduct.ProductOuWithRules.Is24H,
			OvernightTime:    outputProduct.ProductOuWithRules.OvernightTime,
			OvernightPrice:   outputProduct.ProductOuWithRules.OvernightPrice,
			GracePeriod:      outputProduct.ProductOuWithRules.GracePeriod,
			FlgRepeat:        outputProduct.ProductOuWithRules.FlgRepeat,
		}
		resultProduct = &prod.PolicyOuProductWithRules{
			OuId:                  outputProduct.OuId,
			OuCode:                outputProduct.OuCode,
			OuName:                outputProduct.OuName,
			ProductId:             outputProduct.ProductId,
			ProductCode:           outputProduct.ProductCode,
			ProductName:           outputProduct.ProductName,
			ServiceFee:            outputProduct.ServiceFee,
			IsPctServiceFee:       outputProduct.IsPctServiceFee,
			IsPctServiceFeeMember: outputProduct.IsPctServiceFeeMember,
			ServiceFeeMember:      outputProduct.ServiceFeeMember,
			ProductRules:          tempProd,
		}
		productName = resultProduct.ProductName
		productCode = resultProduct.ProductCode
	}

	document, errDocument := make(chan map[string]interface{}), make(chan error)
	go func() {
		defer close(document)
		close(errDocument)

		var trxInvoiceItemList []models.TrxInvoiceItem
		var productList []*prod.PolicyOuProductWithRules
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

		docNo := fmt.Sprintf("%s%s%d", prefixDocNo, autoNumber, tempMerchantKey.OuId)
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
		resProdLists := []*prod.PolicyOuProductWithRules{}
		// var resProdList *prod.PolicyOuProductWithRules
		for _, v := range resultProductList {
			resProdListt := &prod.ProductOuWithRules{
				OuId:             v.ProductOuWithRules.OuId,
				ProductId:        v.ProductOuWithRules.ProductId,
				Price:            v.ProductOuWithRules.Price,
				BaseTime:         v.ProductOuWithRules.BaseTime,
				ProgressiveTime:  v.ProductOuWithRules.ProgressiveTime,
				ProgressivePrice: v.ProductOuWithRules.ProgressivePrice,
				IsPct:            v.ProductOuWithRules.IsPct,
				ProgressivePct:   v.ProductOuWithRules.ProgressivePct,
				MaxPrice:         v.ProductOuWithRules.MaxPrice,
				Is24H:            v.ProductOuWithRules.Is24H,
				OvernightTime:    v.ProductOuWithRules.OvernightTime,
				OvernightPrice:   v.ProductOuWithRules.OvernightPrice,
				GracePeriod:      v.ProductOuWithRules.GracePeriod,
				FlgRepeat:        v.ProductOuWithRules.FlgRepeat,
			}
			resProdList := &prod.PolicyOuProductWithRules{
				OuId:                  v.OuId,
				OuCode:                v.OuCode,
				OuName:                v.OuName,
				ProductId:             v.ProductId,
				ProductCode:           v.ProductCode,
				ProductName:           v.ProductName,
				ServiceFee:            v.ServiceFee,
				IsPctServiceFee:       v.IsPctServiceFee,
				IsPctServiceFeeMember: v.IsPctServiceFeeMember,
				ServiceFeeMember:      v.ServiceFeeMember,
				ProductRules:          resProdListt,
			}
			resProdLists = append(resProdLists, resProdList)
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
				Price:            resultProduct.ProductRules.Price,
				BaseTime:         resultProduct.ProductRules.BaseTime,
				ProgressiveTime:  resultProduct.ProductRules.ProgressiveTime,
				ProgressivePrice: resultProduct.ProductRules.ProgressivePrice,
				IsPct:            resultProduct.ProductRules.IsPct,
				MaxPrice:         resultProduct.ProductRules.MaxPrice,
				ProgressivePct:   resultProduct.ProductRules.ProgressivePct,
				Is24H:            resultProduct.ProductRules.Is24H,
				OvernightTime:    resultProduct.ProductRules.OvernightTime,
				OvernightPrice:   resultProduct.ProductRules.OvernightPrice,
				GracePeriod:      resultProduct.ProductRules.GracePeriod,
				FlgRepeat:        resultProduct.ProductRules.FlgRepeat,
			}

			trxInvoiceItemList = append(trxInvoiceItemList, trxItems)

		} else {

			for _, row := range resProdLists {

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
						Price:            row.ProductRules.Price,
						BaseTime:         row.ProductRules.BaseTime,
						ProgressiveTime:  row.ProductRules.ProgressiveTime,
						ProgressivePrice: row.ProductRules.ProgressivePrice,
						IsPct:            row.ProductRules.IsPct,
						ProgressivePct:   row.ProductRules.ProgressivePct,
						MaxPrice:         row.ProductRules.MaxPrice,
						Is24H:            row.ProductRules.Is24H,
						OvernightTime:    row.ProductRules.OvernightTime,
						OvernightPrice:   row.ProductRules.OvernightPrice,
						GracePeriod:      row.ProductRules.GracePeriod,
						FlgRepeat:        row.ProductRules.FlgRepeat,
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
			RequestData:      string(requestData),
			OuId:             tempMerchantKey.OuId,
			OuName:           tempMerchantKey.OuName,
			OuCode:           tempMerchantKey.OuCode,
			ProductCode:      productCode,
			OuSubBranchId:    tempMerchantKey.OuSubBranchId,
			OuSubBranchName:  tempMerchantKey.OuSubBranchName,
			OuSubBranchCode:  tempMerchantKey.OuSubBranchCode,
			MainOuId:         tempMerchantKey.MainOuId,
			MainOuName:       tempMerchantKey.MainOuName,
			MainOuCode:       tempMerchantKey.MainOuCode,
			QrText:           qrCode,
			TrxInvoiceItem:   trxInvoiceItemList,
		}

		_, err = svc.Service.TrxMongoRepo.AddTrxCheckIn(trx)
		if err != nil {
			log.Println("ERROR AddTrxCheckIn :  ", err)
			errDocument <- err
			return
		}

		data := make(map[string]interface{})
		data["docNo"] = docNo
		data["qrCode"] = qrCode

		document <- data

		// Send data transaction to scheduling message
		// if there connection internet data will be send to cloud server
		// if there not connection internet data not send to cloud server
		// but data will processed to local server
		trxHeader, _ := json.Marshal(trx)
		helperService.ConsumeTrxForScheduling(svc.Service, string(trxHeader))

	}()

	err = <-errDocument
	if err != nil {
		result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	documentTrx := <-document
	responseTrx := trx.ResponseTrxTicket{
		CheckInDateTime: input.CheckInDatetime,
		QrCode:          documentTrx["qrCode"].(string),
		ProductName:     productName,
		VehicleNumberIn: constans.EMPTY_VALUE,
		DocNo:           documentTrx["docNo"].(string),
		OuCode:          tempMerchantKey.OuCode,
		OuName:          tempMerchantKey.OuName,
		Address:         config.ADDRESS,
	}

	anyResponseTrx, _ := anypb.New(&responseTrx)

	result = helpers.ResponseJSON(true, constans.SUCCESS_CODE, constans.EMPTY_VALUE, anyResponseTrx)
	return &trx.MyResponse{
		Response: result,
	}, nil
}

func (svc trxService) InquiryTrxWithoutCard(ctx context.Context, input *trx.RequestInquiryWithoutCard) (*trx.MyResponse, error) {
	var result *trx.Response
	var responseTrx *trx.ResultInquiryTrx
	var ID string

	if err := helpers.BindValidateStruct(input); err != nil {
		result = helpers.ResponseJSON(constans.FALSE_VALUE, constans.VALIDATE_ERROR_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	if strings.ToUpper(input.QrCode) == "P3" {
		// Get P3/Manual status
		resultRedis := svc.Service.RedisClientLocal.Get("P3")
		if resultRedis.Err() != nil || resultRedis.Val() == constans.EMPTY_VALUE || resultRedis.Val() == "0" {
			result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, "Maaf untuk saat ini sistem P3 sudah tidak aktif,Keterangan lebih lanjut silahkan untuk menghubungi administrator", nil)
			return &trx.MyResponse{
				Response: result,
			}, resultRedis.Err()
		}

		merchantKey, err := utils.DecryptMerchantKey(config.MERCHANT_KEY)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}

		tempMerchantKey := &trx.MerchantKey{
			ID:              merchantKey.ID,
			OuId:            merchantKey.OuId,
			OuName:          merchantKey.OuName,
			OuCode:          merchantKey.OuCode,
			OuSubBranchId:   merchantKey.OuSubBranchId,
			OuSubBranchName: merchantKey.OuSubBranchName,
			OuSubBranchCode: merchantKey.OuSubBranchCode,
			MainOuId:        merchantKey.MainOuId,
			MainOuCode:      merchantKey.MainOuCode,
			MainOuName:      merchantKey.MainOuName,
		}

		resultProduct, err := svc.Service.ProductRepo.FindProductByProductCode(input.ProductCode, merchantKey.OuId)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, err.Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}

		tempProd := &prod.ProductOuWithRules{
			OuId:             resultProduct.ProductOuWithRules.OuId,
			ProductId:        resultProduct.ProductOuWithRules.ProductId,
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
		outputProduct := &prod.PolicyOuProductWithRules{
			OuId:                  resultProduct.OuId,
			OuCode:                resultProduct.OuCode,
			OuName:                resultProduct.OuName,
			ProductId:             resultProduct.ProductId,
			ProductCode:           resultProduct.ProductCode,
			ProductName:           resultProduct.ProductName,
			ServiceFee:            resultProduct.ServiceFee,
			IsPctServiceFee:       resultProduct.IsPctServiceFee,
			IsPctServiceFeeMember: resultProduct.IsPctServiceFeeMember,
			ServiceFeeMember:      resultProduct.ServiceFeeMember,
			ProductRules:          tempProd,
		}

		// productName = outputProduct.ProductName
		// productCode = outputProduct.ProductCode

		checkinDatetimeParse, _ := time.Parse("2006-01-02 15:04:05", input.InquiryDateTime)
		checkoutDatetimeParse, _ := time.Parse("2006-01-02 15:04:05", input.InquiryDateTime)

		yearMonth := checkinDatetimeParse.Format("060102")
		HourCheckIn := checkinDatetimeParse.Format("15")
		prefixDocNo := utils.RandStringBytesMaskImprSrcChr(4)
		prefix := fmt.Sprintf("%s%s", yearMonth, HourCheckIn)
		autoNumber, err := svc.Service.GenAutoNumRepo.AutonumberValueWithDatatype(constans.DATATYPE_TRX_LOCAL, prefix, 4)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, err.Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}

		docNo := fmt.Sprintf("%s%s%d", prefixDocNo, autoNumber, tempMerchantKey.OuId)
		randomKey := utils.RandStringBytesMaskImprSrc(16)
		encQrTxt, _ := utils.Encrypt(docNo, randomKey)
		qrCode := fmt.Sprintf("%s%s", encQrTxt.Result, randomKey)

		var trxInvoiceItemList []*trx.TrxInvoiceItem
		trxInvoiceItem := &trx.TrxInvoiceItem{
			DocNo:                  docNo,
			ProductId:              outputProduct.ProductId,
			ProductCode:            outputProduct.ProductCode,
			ProductName:            outputProduct.ProductName,
			IsPctServiceFee:        outputProduct.IsPctServiceFee,
			ServiceFee:             outputProduct.ServiceFee,
			ServiceFeeMember:       outputProduct.ServiceFeeMember,
			Price:                  outputProduct.ProductRules.Price,
			BaseTime:               outputProduct.ProductRules.BaseTime,
			ProgressiveTime:        outputProduct.ProductRules.ProgressiveTime,
			ProgressivePrice:       outputProduct.ProductRules.ProgressivePrice,
			IsPct:                  outputProduct.ProductRules.IsPct,
			ProgressivePct:         outputProduct.ProductRules.ProgressivePct,
			MaxPrice:               outputProduct.ProductRules.MaxPrice,
			Is24H:                  outputProduct.ProductRules.Is24H,
			OvernightTime:          outputProduct.ProductRules.OvernightTime,
			OvernightPrice:         outputProduct.ProductRules.OvernightPrice,
			GracePeriod:            outputProduct.ProductRules.GracePeriod,
			FlgRepeat:              outputProduct.ProductRules.FlgRepeat,
			TotalAmount:            outputProduct.ProductRules.Price,
			TotalProgressiveAmount: outputProduct.ProductRules.Price,
		}
		trxInvoiceItemList = append(trxInvoiceItemList, trxInvoiceItem)

		trxs := &trx.Trx{
			DocNo:            docNo,
			DocDate:          constans.EMPTY_VALUE,
			CheckinDateTime:  input.InquiryDateTime,
			CheckoutDateTime: constans.EMPTY_VALUE,
			DeviceIdIn:       constans.EMPTY_VALUE,
			DeviceId:         constans.EMPTY_VALUE,
			GetIn:            constans.EMPTY_VALUE,
			GetOut:           constans.EMPTY_VALUE,
			CardNumberUUIDIn: constans.EMPTY_VALUE,
			CardNumberIn:     constans.EMPTY_VALUE,
			CardNumberUUID:   constans.EMPTY_VALUE,
			CardNumber:       constans.EMPTY_VALUE,
			TypeCard:         constans.EMPTY_VALUE,
			BeginningBalance: 0,
			ExtLocalDateTime: input.InquiryDateTime,
			GrandTotal:       outputProduct.ProductRules.Price,
			ProductCode:      outputProduct.ProductCode,
			ProductName:      outputProduct.ProductName,
			ProductData:      utils.ToString(outputProduct),
			RequestData:      utils.ToString(input),
			RequestOutData:   utils.ToString(input),
			OuId:             tempMerchantKey.OuId,
			OuName:           tempMerchantKey.OuName,
			OuCode:           tempMerchantKey.OuCode,
			OuSubBranchId:    tempMerchantKey.OuSubBranchId,
			USubBranchName:   tempMerchantKey.OuSubBranchName,
			OuSubBranchCode:  tempMerchantKey.OuSubBranchCode,
			MainOuId:         tempMerchantKey.MainOuId,
			MainOuCode:       tempMerchantKey.MainOuCode,
			MainOuName:       tempMerchantKey.MainOuName,
			MemberCode:       constans.MANUAL,
			MemberName:       constans.MANUAL,
			MemberType:       constans.TYPE_PARTNER_FREE_PASS,
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

		redisStatus := svc.Service.RedisClientLocal.Set(fmt.Sprintf("%s#%s", "P3", docNo), utils.ToString(trxs), 1*time.Minute)
		if redisStatus.Err() != nil {
			result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, redisStatus.Err().Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}

		duration := utils.ConvDiffTime(checkinDatetimeParse, checkoutDatetimeParse)
		anyDuration, _ := anypb.New(duration)

		responseTrx.Id = fmt.Sprintf("%s#%s", "P3", docNo)
		responseTrx.DocNo = docNo
		responseTrx.Nominal = trxs.GrandTotal
		responseTrx.ProductCode = outputProduct.ProductCode
		responseTrx.ProductName = outputProduct.ProductName
		responseTrx.VehicleNumberIn = constans.EMPTY_VALUE
		responseTrx.QrCode = constans.EMPTY_VALUE
		responseTrx.ExcludeSf = true
		responseTrx.Duration = anyDuration
		responseTrx.OuCode = merchantKey.OuCode

	} else {

		if len(input.QrCode) < 68 {
			input.QrCode = constans.NONE_QRCODE
		}

		saltKey := input.QrCode[len(input.QrCode)-16:]
		keyword := input.QrCode[:len(input.QrCode)-16]

		docNo, err := utils.Decrypt(keyword, saltKey)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}

		if input.ProductCode == constans.EMPTY_VALUE {
			resultTrxOutstanding, exists, err := svc.Service.TrxMongoRepo.IsTrxOutstandingByDocNoForCustom(docNo)
			if err != nil {
				result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
				return &trx.MyResponse{
					Response: result,
				}, err
			}

			var tempTrxInvoiceItems []*trx.TrxInvoiceItem
			for _, v := range resultTrxOutstanding.TrxInvoiceItem {
				tempTrxInvoiceItem := &trx.TrxInvoiceItem{
					DocNo:                  v.DocNo,
					ProductId:              v.ProductId,
					ProductCode:            v.ProductCode,
					ProductName:            v.ProductName,
					IsPctServiceFee:        v.IsPctServiceFee,
					ServiceFee:             v.ServiceFee,
					ServiceFeeMember:       v.ServiceFeeMember,
					Price:                  v.Price,
					BaseTime:               v.BaseTime,
					ProgressiveTime:        v.ProgressiveTime,
					ProgressivePrice:       v.ProgressivePrice,
					IsPct:                  v.IsPct,
					ProgressivePct:         v.ProgressivePct,
					MaxPrice:               v.MaxPrice,
					Is24H:                  v.Is24H,
					OvernightTime:          v.OvernightTime,
					OvernightPrice:         v.OvernightPrice,
					GracePeriod:            v.GracePeriod,
					FlgRepeat:              v.FlgRepeat,
					TotalAmount:            v.TotalAmount,
					TotalProgressiveAmount: v.TotalProgressiveAmount,
				}
				tempTrxInvoiceItems = append(tempTrxInvoiceItems, tempTrxInvoiceItem)
			}
			tempMemberData := &trx.TrxMember{
				DocNo:       resultTrxOutstanding.MemberData.DocNo,
				PartnerCode: resultTrxOutstanding.MemberData.PartnerCode,
				FirstName:   resultTrxOutstanding.MemberData.FirstName,
				LastName:    resultTrxOutstanding.MemberData.LastName,
				RoleType:    resultTrxOutstanding.MemberData.RoleType,
				PhoneNumber: resultTrxOutstanding.MemberData.PhoneNumber,
				Email:       resultTrxOutstanding.MemberData.Email,
				Active:      resultTrxOutstanding.MemberData.Active,
				ActiveAt:    resultTrxOutstanding.MemberData.ActiveAt,
				NonActiveAt: func() *wrapperspb.StringValue {
					if resultTrxOutstanding.MemberData.NonActiveAt != nil {
						return &wrapperspb.StringValue{Value: *resultTrxOutstanding.MemberData.NonActiveAt}
					}
					return nil
				}(),
				OuId:               resultTrxOutstanding.MemberData.OuId,
				TypePartner:        resultTrxOutstanding.MemberData.TypePartner,
				CardNumber:         resultTrxOutstanding.MemberData.CardNumber,
				VehicleNumber:      resultTrxOutstanding.MemberData.VehicleNumber,
				RegisteredDatetime: resultTrxOutstanding.MemberData.RegisteredDatetime,
				DateFrom:           resultTrxOutstanding.MemberData.DateFrom,
				DateTo:             resultTrxOutstanding.MemberData.DateTo,
				ProductId:          resultTrxOutstanding.MemberData.ProductId,
				ProductCode:        resultTrxOutstanding.MemberData.ProductCode,
			}
			convertedTrxAddInfo := make(map[string]*anypb.Any)
			for k, v := range resultTrxOutstanding.TrxAddInfo {
				if msg, ok := v.(goproto.Message); ok {
					anyVal, err := anypb.New(msg)
					if err != nil {
						continue
					}
					convertedTrxAddInfo[k] = anyVal
				}
			}

			tempRequestAddTrxInvoiceDetailItem := &trx.TrxInvoiceDetailItem{
				DocNo:         resultTrxOutstanding.RequestAddTrxInvoiceDetailItem.DocNo,
				ProductCode:   resultTrxOutstanding.RequestAddTrxInvoiceDetailItem.ProductCode,
				InvoiceAmount: resultTrxOutstanding.RequestAddTrxInvoiceDetailItem.InvoiceAmount,
				CreatedAt:     resultTrxOutstanding.RequestAddTrxInvoiceDetailItem.CreatedAt,
				CreatedDate:   resultTrxOutstanding.RequestAddTrxInvoiceDetailItem.CreatedDate,
			}

			tempTrxOutstanding := &trx.Trx{
				DocNo:                          resultTrxOutstanding.DocNo,
				DocDate:                        resultTrxOutstanding.DocDate,
				PaymentRefDocno:                resultTrxOutstanding.PaymentRefDocNo,
				CheckinDateTime:                resultTrxOutstanding.CheckInDatetime,
				CheckoutDateTime:               resultTrxOutstanding.CheckOutDatetime,
				DeviceIdIn:                     resultTrxOutstanding.DeviceIdIn,
				DeviceId:                       resultTrxOutstanding.DeviceId,
				GetIn:                          resultTrxOutstanding.GateIn,
				GetOut:                         resultTrxOutstanding.GateOut,
				CardNumberUUIDIn:               resultTrxOutstanding.CardNumberUUIDIn,
				CardNumberIn:                   resultTrxOutstanding.CardNumberIn,
				CardNumberUUID:                 resultTrxOutstanding.CardNumberUUID,
				CardNumber:                     resultTrxOutstanding.CardNumber,
				TypeCard:                       resultTrxOutstanding.TypeCard,
				BeginningBalance:               resultTrxOutstanding.BeginningBalance,
				ExtLocalDateTime:               resultTrxOutstanding.ExtLocalDatetime,
				ChargeAmount:                   resultTrxOutstanding.ChargeAmount,
				GrandTotal:                     resultTrxOutstanding.GrandTotal,
				ProductCode:                    resultTrxOutstanding.ProductCode,
				ProductName:                    resultTrxOutstanding.ProductName,
				ProductData:                    resultTrxOutstanding.ProductData,
				RequestData:                    resultTrxOutstanding.RequestData,
				RequestOutData:                 resultTrxOutstanding.RequestOutData,
				OuId:                           resultTrxOutstanding.OuId,
				OuName:                         resultTrxOutstanding.OuName,
				OuCode:                         resultTrxOutstanding.OuCode,
				OuSubBranchId:                  resultTrxOutstanding.OuSubBranchId,
				USubBranchName:                 resultTrxOutstanding.OuSubBranchName,
				OuSubBranchCode:                resultTrxOutstanding.OuSubBranchCode,
				MainOuId:                       resultTrxOutstanding.MainOuId,
				MainOuCode:                     resultTrxOutstanding.MainOuCode,
				MainOuName:                     resultTrxOutstanding.MainOuName,
				MemberCode:                     resultTrxOutstanding.MemberCode,
				MemberName:                     resultTrxOutstanding.MemberName,
				MemberType:                     resultTrxOutstanding.MemberType,
				MemberStatus:                   resultTrxOutstanding.MemberStatus,
				MemberExpiredDate:              resultTrxOutstanding.MemberExpiredDate,
				CheckInTime:                    resultTrxOutstanding.CheckInTime,
				CheckOutTime:                   resultTrxOutstanding.CheckOutTime,
				DurationTime:                   resultTrxOutstanding.DurationTime,
				VehicleNumberIn:                resultTrxOutstanding.VehicleNumberIn,
				VehicleNumberOut:               resultTrxOutstanding.VehicleNumberOut,
				LogTrans:                       resultTrxOutstanding.LogTrans,
				MerchantKey:                    resultTrxOutstanding.MerchantKey,
				QrText:                         resultTrxOutstanding.QrText,
				QrA2P:                          resultTrxOutstanding.QrA2P,
				QrTextPaymentOnline:            resultTrxOutstanding.QrTextPaymentOnline,
				TrxInvoiceItem:                 tempTrxInvoiceItems,
				FlagSyncData:                   resultTrxOutstanding.FlagSyncData,
				MemberData:                     tempMemberData,
				TrxAddInfo:                     convertedTrxAddInfo,
				FlagTrxFromCloud:               resultTrxOutstanding.FlagTrxFromCloud,
				IsRsyncDataTrx:                 resultTrxOutstanding.IsRsyncDataTrx,
				ExcludeSf:                      resultTrxOutstanding.ExcludeSf,
				FlagCharge:                     resultTrxOutstanding.FlagCharge,
				ChargeType:                     resultTrxOutstanding.ChargeType,
				RequestAddTrxInvoiceDetailItem: tempRequestAddTrxInvoiceDetailItem,
				LastUpdatedAt:                  resultTrxOutstanding.LastUpdatedAt,
			}

			if !exists {
				result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, "Trx outstanding not found", nil)
				return &trx.MyResponse{
					Response: result,
				}, err
			}

			input.ProductCode = tempTrxOutstanding.ProductCode
		}

		productData, existsProduct := svc.Service.TrxMongoRepo.IsTrxProductCustomExistsByKeyword(docNo)

		tempProdDat := &trx.TrxProductCustom{
			Keyword:     productData.Keyword,
			ProductName: productData.ProductName,
			ProductCode: productData.ProductCode,
		}
		if existsProduct {
			input.ProductCode = tempProdDat.ProductCode
		}

		resultTrx, exists, err := svc.Service.TrxMongoRepo.IsTrxOutstandingByDocNo(docNo, input.ProductCode)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}

		tempTrxInvoiceItems := []*trx.TrxInvoiceItem{}
		for _, v := range resultTrx.TrxInvoiceItem {
			tempTrxInvoiceItem := &trx.TrxInvoiceItem{
				DocNo:                  v.DocNo,
				ProductId:              v.ProductId,
				ProductCode:            v.ProductCode,
				ProductName:            v.ProductName,
				IsPctServiceFee:        v.IsPctServiceFee,
				ServiceFee:             v.ServiceFee,
				ServiceFeeMember:       v.ServiceFeeMember,
				Price:                  v.Price,
				BaseTime:               v.BaseTime,
				ProgressiveTime:        v.ProgressiveTime,
				ProgressivePrice:       v.ProgressivePrice,
				IsPct:                  v.IsPct,
				ProgressivePct:         v.ProgressivePct,
				MaxPrice:               v.MaxPrice,
				Is24H:                  v.Is24H,
				OvernightTime:          v.OvernightTime,
				OvernightPrice:         v.OvernightPrice,
				GracePeriod:            v.GracePeriod,
				FlgRepeat:              v.FlgRepeat,
				TotalAmount:            v.TotalAmount,
				TotalProgressiveAmount: v.TotalProgressiveAmount,
			}
			tempTrxInvoiceItems = append(tempTrxInvoiceItems, tempTrxInvoiceItem)
		}

		tempResultTrx := &trx.ResultFindTrxOutstanding{
			Id:              resultTrx.ID.Hex(),
			DocNo:           resultTrx.DocNo,
			GrandTotal:      resultTrx.GrandTotal,
			CheckInDatetime: resultTrx.CheckInDatetime,
			OverNightPrice:  resultTrx.OverNightPrice,
			Is24H:           resultTrx.Is24H,
			CardNumber:      resultTrx.CardNumber,
			CardNumberUUID:  resultTrx.CardNumberUUID,
			OuCode:          resultTrx.OuCode,
			VehicleNumberIn: resultTrx.VehicleNumberIn,
			TrxInvoiceItem:  tempTrxInvoiceItems,
		}

		if !exists {

			// If internet connected check trxLocal to cloud server
			if utils.IsConnected() {
				responseTrxCloud, err := helperService.CheckTrxCloudServer(docNo)
				if err != nil {
					result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, "Sesi Tidak Ditemukan", nil)
					return &trx.MyResponse{
						Response: result,
					}, err
				}
				var tempTrxInvoiceItems []*trx.TrxInvoiceItem
				for _, v := range responseTrxCloud.Result.Trx.TrxInvoiceItem {
					tempTrxInvoiceItem := &trx.TrxInvoiceItem{
						DocNo:                  v.DocNo,
						ProductId:              v.ProductId,
						ProductCode:            v.ProductCode,
						ProductName:            v.ProductName,
						IsPctServiceFee:        v.IsPctServiceFee,
						ServiceFee:             v.ServiceFee,
						ServiceFeeMember:       v.ServiceFeeMember,
						Price:                  v.Price,
						BaseTime:               v.BaseTime,
						ProgressiveTime:        v.ProgressiveTime,
						ProgressivePrice:       v.ProgressivePrice,
						IsPct:                  v.IsPct,
						ProgressivePct:         v.ProgressivePct,
						MaxPrice:               v.MaxPrice,
						Is24H:                  v.Is24H,
						OvernightTime:          v.OvernightTime,
						OvernightPrice:         v.OvernightPrice,
						GracePeriod:            v.GracePeriod,
						FlgRepeat:              v.FlgRepeat,
						TotalAmount:            v.TotalAmount,
						TotalProgressiveAmount: v.TotalProgressiveAmount,
					}
					tempTrxInvoiceItems = append(tempTrxInvoiceItems, tempTrxInvoiceItem)
				}
				tempMemberData := &trx.TrxMember{
					DocNo:       responseTrxCloud.Result.Trx.MemberData.DocNo,
					PartnerCode: responseTrxCloud.Result.Trx.MemberData.PartnerCode,
					FirstName:   responseTrxCloud.Result.Trx.MemberData.FirstName,
					LastName:    responseTrxCloud.Result.Trx.MemberData.LastName,
					RoleType:    responseTrxCloud.Result.Trx.MemberData.RoleType,
					PhoneNumber: responseTrxCloud.Result.Trx.MemberData.PhoneNumber,
					Email:       responseTrxCloud.Result.Trx.MemberData.Email,
					Active:      responseTrxCloud.Result.Trx.MemberData.Active,
					ActiveAt:    responseTrxCloud.Result.Trx.MemberData.ActiveAt,
					NonActiveAt: func() *wrapperspb.StringValue {
						if responseTrxCloud.Result.Trx.MemberData.NonActiveAt != nil {
							return &wrapperspb.StringValue{Value: *responseTrxCloud.Result.Trx.MemberData.NonActiveAt}
						}
						return nil
					}(),
					OuId:               responseTrxCloud.Result.Trx.MemberData.OuId,
					TypePartner:        responseTrxCloud.Result.Trx.MemberData.TypePartner,
					CardNumber:         responseTrxCloud.Result.Trx.MemberData.CardNumber,
					VehicleNumber:      responseTrxCloud.Result.Trx.MemberData.VehicleNumber,
					RegisteredDatetime: responseTrxCloud.Result.Trx.MemberData.RegisteredDatetime,
					DateFrom:           responseTrxCloud.Result.Trx.MemberData.DateFrom,
					DateTo:             responseTrxCloud.Result.Trx.MemberData.DateTo,
					ProductId:          responseTrxCloud.Result.Trx.MemberData.ProductId,
					ProductCode:        responseTrxCloud.Result.Trx.MemberData.ProductCode,
				}
				convertedTrxAddInfo := make(map[string]*anypb.Any)
				for k, v := range responseTrxCloud.Result.Trx.TrxAddInfo {
					if msg, ok := v.(goproto.Message); ok {
						anyVal, err := anypb.New(msg)
						if err != nil {
							continue
						}
						convertedTrxAddInfo[k] = anyVal
					}
				}

				tempRequestAddTrxInvoiceDetailItem := &trx.TrxInvoiceDetailItem{
					DocNo:         responseTrxCloud.Result.Trx.RequestAddTrxInvoiceDetailItem.DocNo,
					ProductCode:   responseTrxCloud.Result.Trx.RequestAddTrxInvoiceDetailItem.ProductCode,
					InvoiceAmount: responseTrxCloud.Result.Trx.RequestAddTrxInvoiceDetailItem.InvoiceAmount,
					CreatedAt:     responseTrxCloud.Result.Trx.RequestAddTrxInvoiceDetailItem.CreatedAt,
					CreatedDate:   responseTrxCloud.Result.Trx.RequestAddTrxInvoiceDetailItem.CreatedDate,
				}
				result1 := &trx.Trx{
					DocNo:                          responseTrxCloud.Result.Trx.DocNo,
					DocDate:                        responseTrxCloud.Result.Trx.DocDate,
					PaymentRefDocno:                responseTrxCloud.Result.Trx.PaymentRefDocNo,
					CheckinDateTime:                responseTrxCloud.Result.Trx.CheckInDatetime,
					CheckoutDateTime:               responseTrxCloud.Result.Trx.CheckOutDatetime,
					DeviceIdIn:                     responseTrxCloud.Result.Trx.DeviceIdIn,
					DeviceId:                       responseTrxCloud.Result.Trx.DeviceId,
					GetIn:                          responseTrxCloud.Result.Trx.GateIn,
					GetOut:                         responseTrxCloud.Result.Trx.GateOut,
					CardNumberUUIDIn:               responseTrxCloud.Result.Trx.CardNumberUUIDIn,
					CardNumberIn:                   responseTrxCloud.Result.Trx.CardNumberIn,
					CardNumberUUID:                 responseTrxCloud.Result.Trx.CardNumberUUID,
					CardNumber:                     responseTrxCloud.Result.Trx.CardNumber,
					TypeCard:                       responseTrxCloud.Result.Trx.TypeCard,
					BeginningBalance:               responseTrxCloud.Result.Trx.BeginningBalance,
					ExtLocalDateTime:               responseTrxCloud.Result.Trx.ExtLocalDatetime,
					ChargeAmount:                   responseTrxCloud.Result.Trx.ChargeAmount,
					GrandTotal:                     responseTrxCloud.Result.Trx.GrandTotal,
					ProductCode:                    responseTrxCloud.Result.Trx.ProductCode,
					ProductName:                    responseTrxCloud.Result.Trx.ProductName,
					ProductData:                    responseTrxCloud.Result.Trx.ProductData,
					RequestData:                    responseTrxCloud.Result.Trx.RequestData,
					RequestOutData:                 responseTrxCloud.Result.Trx.RequestOutData,
					OuId:                           responseTrxCloud.Result.Trx.OuId,
					OuName:                         responseTrxCloud.Result.Trx.OuName,
					OuCode:                         responseTrxCloud.Result.Trx.OuCode,
					OuSubBranchId:                  responseTrxCloud.Result.Trx.OuSubBranchId,
					USubBranchName:                 responseTrxCloud.Result.Trx.OuSubBranchName,
					OuSubBranchCode:                responseTrxCloud.Result.Trx.OuSubBranchCode,
					MainOuId:                       responseTrxCloud.Result.Trx.MainOuId,
					MainOuCode:                     responseTrxCloud.Result.Trx.MainOuCode,
					MainOuName:                     responseTrxCloud.Result.Trx.MainOuName,
					MemberCode:                     responseTrxCloud.Result.Trx.MemberCode,
					MemberName:                     responseTrxCloud.Result.Trx.MemberName,
					MemberType:                     responseTrxCloud.Result.Trx.MemberType,
					MemberStatus:                   responseTrxCloud.Result.Trx.MemberStatus,
					MemberExpiredDate:              responseTrxCloud.Result.Trx.MemberExpiredDate,
					CheckInTime:                    responseTrxCloud.Result.Trx.CheckInTime,
					CheckOutTime:                   responseTrxCloud.Result.Trx.CheckOutTime,
					DurationTime:                   responseTrxCloud.Result.Trx.DurationTime,
					VehicleNumberIn:                responseTrxCloud.Result.Trx.VehicleNumberIn,
					VehicleNumberOut:               responseTrxCloud.Result.Trx.VehicleNumberOut,
					LogTrans:                       responseTrxCloud.Result.Trx.LogTrans,
					MerchantKey:                    responseTrxCloud.Result.Trx.MerchantKey,
					QrText:                         responseTrxCloud.Result.Trx.QrText,
					QrA2P:                          responseTrxCloud.Result.Trx.QrA2P,
					QrTextPaymentOnline:            responseTrxCloud.Result.Trx.QrTextPaymentOnline,
					FlagSyncData:                   responseTrxCloud.Result.Trx.FlagSyncData,
					MemberData:                     tempMemberData,
					TrxAddInfo:                     convertedTrxAddInfo,
					FlagTrxFromCloud:               responseTrxCloud.Result.Trx.FlagTrxFromCloud,
					IsRsyncDataTrx:                 responseTrxCloud.Result.Trx.IsRsyncDataTrx,
					ExcludeSf:                      responseTrxCloud.Result.Trx.ExcludeSf,
					FlagCharge:                     responseTrxCloud.Result.Trx.FlagCharge,
					ChargeType:                     responseTrxCloud.Result.Trx.ChargeType,
					RequestAddTrxInvoiceDetailItem: tempRequestAddTrxInvoiceDetailItem,
					LastUpdatedAt:                  responseTrxCloud.Result.Trx.LastUpdatedAt,
					TrxInvoiceItem:                 tempTrxInvoiceItems,
				}
				var tempTrxInvoiceItemss []*trx.TrxInvoiceItem
				tempTII := responseTrxCloud.Result.TrxInvoiceItem.TrxInvoiceItem
				for _, x := range tempTII {
					tempTrxInvoiceItem := &trx.TrxInvoiceItem{
						DocNo:                  x.DocNo,
						ProductId:              x.ProductId,
						ProductCode:            x.ProductCode,
						ProductName:            x.ProductName,
						IsPctServiceFee:        x.IsPctServiceFee,
						ServiceFee:             x.ServiceFee,
						ServiceFeeMember:       x.ServiceFeeMember,
						Price:                  x.Price,
						BaseTime:               x.BaseTime,
						ProgressiveTime:        x.ProgressiveTime,
						ProgressivePrice:       x.ProgressivePrice,
						IsPct:                  x.IsPct,
						ProgressivePct:         x.ProgressivePct,
						MaxPrice:               x.MaxPrice,
						Is24H:                  x.Is24H,
						OvernightTime:          x.OvernightTime,
						OvernightPrice:         x.OvernightPrice,
						GracePeriod:            x.GracePeriod,
						FlgRepeat:              x.FlgRepeat,
						TotalAmount:            x.TotalAmount,
						TotalProgressiveAmount: x.TotalProgressiveAmount,
					}
					tempTrxInvoiceItemss = append(tempTrxInvoiceItemss, tempTrxInvoiceItem)
				}
				TrxInvTemp := &trx.ResultFindTrxOutstanding{
					Id:              responseTrxCloud.Result.TrxInvoiceItem.ID.Hex(),
					DocNo:           responseTrxCloud.Result.TrxInvoiceItem.DocNo,
					GrandTotal:      responseTrxCloud.Result.TrxInvoiceItem.GrandTotal,
					CheckInDatetime: responseTrxCloud.Result.TrxInvoiceItem.CheckInDatetime,
					OverNightPrice:  responseTrxCloud.Result.TrxInvoiceItem.OverNightPrice,
					Is24H:           responseTrxCloud.Result.TrxInvoiceItem.Is24H,
					CardNumber:      responseTrxCloud.Result.TrxInvoiceItem.CardNumber,
					CardNumberUUID:  responseTrxCloud.Result.TrxInvoiceItem.CardNumberUUID,
					OuCode:          responseTrxCloud.Result.TrxInvoiceItem.OuCode,
					VehicleNumberIn: responseTrxCloud.Result.TrxInvoiceItem.VehicleNumberIn,
					TrxInvoiceItem:  tempTrxInvoiceItemss,
				}
				resultx := &trx.Result{
					Trx:            result1,
					TrxInvoiceItem: TrxInvTemp,
				}
				tempResponseTrxCloud := &trx.ResponseTrxCloud{
					StatusCode: responseTrxCloud.StatusCode,
					Result:     resultx,
					Success:    responseTrxCloud.Success,
				}

				if !tempResponseTrxCloud.Success {
					result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, "Sesi Tidak Ditemukan", nil)
					return &trx.MyResponse{
						Response: result,
					}, err
				}

				dataTrx := utils.ToString(tempResponseTrxCloud.Result.Trx)
				ID = fmt.Sprintf("%s#%s", "SERVER", tempResponseTrxCloud.Result.TrxInvoiceItem.Id)
				redisStatus := svc.Service.RedisClientLocal.Set(ID, dataTrx, 5*time.Minute)
				if redisStatus.Err() != nil {
					result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, redisStatus.Err().Error(), nil)
					return &trx.MyResponse{
						Response: result,
					}, err
				}

				tempResultTrx = tempResponseTrxCloud.Result.TrxInvoiceItem
				input.ProductCode = tempResultTrx.TrxInvoiceItem[0].ProductCode
			}
		} else {
			ID = tempResultTrx.Id
		}

		checkinDatetimeParse, _ := time.Parse("2006-01-02 15:04", tempResultTrx.CheckInDatetime[:len(tempResultTrx.CheckInDatetime)-3])
		checkoutDatetimeParse, _ := time.Parse("2006-01-02 15:04", input.InquiryDateTime[:len(input.InquiryDateTime)-3])
		duration := utils.ConvDiffTime(checkinDatetimeParse, checkoutDatetimeParse)

		// Request QRIS to A2P
		if config.QRISPayment == constans.YES && tempResultTrx.TrxInvoiceItem[0].TotalAmount > 0 {
			go helperService.CallQRPayment(tempResultTrx, input, duration, svc.Service)
		}
		anyDuration, _ := anypb.New(duration)

		responseTrx.Id = ID
		responseTrx.DocNo = docNo
		responseTrx.Nominal = tempResultTrx.TrxInvoiceItem[0].TotalAmount
		responseTrx.ProductCode = tempResultTrx.TrxInvoiceItem[0].ProductCode
		responseTrx.ProductName = tempResultTrx.TrxInvoiceItem[0].ProductName
		responseTrx.VehicleNumberIn = tempResultTrx.VehicleNumberIn
		responseTrx.QrCode = constans.EMPTY_VALUE
		responseTrx.ExcludeSf = tempResultTrx.TrxInvoiceItem[0].TotalAmount == 0
		responseTrx.Duration = anyDuration
		responseTrx.OuCode = tempResultTrx.OuCode
	}
	anyResponseTrx, _ := anypb.New(responseTrx)

	result = helpers.ResponseJSON(true, constans.SUCCESS_CODE, constans.EMPTY_VALUE, anyResponseTrx)
	return &trx.MyResponse{
		Response: result,
	}, nil
}

func (svc trxService) InquiryTrxWithCard(ctx context.Context, input *trx.RequestInquiryWithCard) (*trx.MyResponse, error) {
	var result *trx.Response
	resultInquiryTrxWithCard := make(map[string]interface{})
	resultInquiryTrxWithCard["memberCode"] = constans.EMPTY_VALUE
	resultInquiryTrxWithCard["memberName"] = constans.EMPTY_VALUE
	resultInquiryTrxWithCard["memberType"] = constans.EMPTY_VALUE

	if err := helpers.BindValidateStruct(input); err != nil {
		result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	inquiryDate := input.InquiryDateTime[:len(input.InquiryDateTime)-9]
	resultMember, exists, err := svc.Service.MemberRepo.IsMemberByAdvanceIndex(input.UuidCard, constans.EMPTY_VALUE, inquiryDate, config.MEMBER_BY, false)
	if err != nil {
		result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	tempMemberData := &trx.Member{
		Id:          resultMember.ID,
		PartnerCode: resultMember.PartnerCode,
		FirstName:   resultMember.FirstName,
		LastName:    resultMember.LastName,
		RoleType:    resultMember.RoleType,
		PhoneNumber: resultMember.PhoneNumber,
		Email:       resultMember.Email,
		Active:      resultMember.Active,
		ActiveAt:    resultMember.ActiveAt,
		NonActiveAt: func() *wrapperspb.StringValue {
			if resultMember.NonActiveAt != nil {
				return &wrapperspb.StringValue{Value: *resultMember.NonActiveAt}
			}
			return nil
		}(),
		OuId:                resultMember.OuId,
		TypePartner:         resultMember.TypePartner,
		CardNumber:          resultMember.CardNumber,
		VehicleNumber:       resultMember.VehicleNumber,
		RegisteredDatetime:  resultMember.RegisteredDatetime,
		DateFrom:            resultMember.DateFrom,
		DateTo:              resultMember.DateTo,
		ProductId:           resultMember.ProductId,
		ProductCode:         resultMember.ProductCode,
		ProductMembershipId: resultMember.ProductMembershipId,
		Price:               resultMember.Price,
		ServiceFee:          resultMember.ServiceFee,
		IsPctSfee:           resultMember.IsPctSfee,
		DueDate:             resultMember.DueDate,
		DiscType:            resultMember.DiscType,
		DiscAmount:          resultMember.DiscAmount,
		DiscPct:             resultMember.DiscPct,
		GracePeriodDate:     resultMember.GracePeriodDate,
		Username:            resultMember.Username,
		IsExtendMember:      resultMember.IsExtendMember,
		CreatedAt:           resultMember.CreatedAt,
		CreatedBy:           resultMember.CreatedBy,
		UpdatedAt:           resultMember.UpdatedAt,
		UpdatedBy:           resultMember.UpdatedBy,
	}

	isFreePass := false
	isSpecialMember := false
	if exists {
		if tempMemberData.TypePartner == constans.TYPE_PARTNER_FREE_PASS {
			isFreePass = true
		} else if tempMemberData.TypePartner == constans.TYPE_PARTNER_SPECIAL_MEMBER {
			isSpecialMember = true
		}
	}

	if !isFreePass {
		productData, existsProduct := svc.Service.TrxMongoRepo.IsTrxProductCustomExistsByKeyword(input.UuidCard)
		tempProdDat := &trx.TrxProductCustom{
			Keyword:     productData.Keyword,
			ProductName: productData.ProductName,
			ProductCode: productData.ProductCode,
		}
		if existsProduct {
			input.ProductCode = tempProdDat.ProductCode
		}

		resultTrx, err := svc.Service.TrxMongoRepo.FindTrxOutstandingByUUID(input.UuidCard, input.ProductCode)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, "Sesi Tidak Ditemukan", nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}

		if isSpecialMember {
			bodySplit := strings.Split(input.ProductCode, "SPM")
			if len(bodySplit) > 1 {
				if input.ProductCode == tempMemberData.ProductCode {
					resultTrx, err = svc.Service.TrxMongoRepo.FindTrxOutstandingByUUID(input.UuidCard, input.ProductCode)
					if err != nil {
						result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, "Sesi Tidak Ditemukan", nil)
						return &trx.MyResponse{
							Response: result,
						}, err
					}
				}
			} else {
				input.ProductCode = fmt.Sprintf("%s-%s", "SPM", input.ProductCode)
				if input.ProductCode == resultMember.ProductCode {
					resultTrx, err = svc.Service.TrxMongoRepo.FindTrxOutstandingByUUID(input.UuidCard, input.ProductCode)
					if err != nil {
						result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, "Sesi Tidak Ditemukan", nil)
						return &trx.MyResponse{
							Response: result,
						}, err
					}
				}
			}

		}
		var tempTrxInvoiceItemss []*trx.TrxInvoiceItem
		tempTII := resultTrx.TrxInvoiceItem
		for _, x := range tempTII {
			tempTrxInvoiceItem := &trx.TrxInvoiceItem{
				DocNo:                  x.DocNo,
				ProductId:              x.ProductId,
				ProductCode:            x.ProductCode,
				ProductName:            x.ProductName,
				IsPctServiceFee:        x.IsPctServiceFee,
				ServiceFee:             x.ServiceFee,
				ServiceFeeMember:       x.ServiceFeeMember,
				Price:                  x.Price,
				BaseTime:               x.BaseTime,
				ProgressiveTime:        x.ProgressiveTime,
				ProgressivePrice:       x.ProgressivePrice,
				IsPct:                  x.IsPct,
				ProgressivePct:         x.ProgressivePct,
				MaxPrice:               x.MaxPrice,
				Is24H:                  x.Is24H,
				OvernightTime:          x.OvernightTime,
				OvernightPrice:         x.OvernightPrice,
				GracePeriod:            x.GracePeriod,
				FlgRepeat:              x.FlgRepeat,
				TotalAmount:            x.TotalAmount,
				TotalProgressiveAmount: x.TotalProgressiveAmount,
			}
			tempTrxInvoiceItemss = append(tempTrxInvoiceItemss, tempTrxInvoiceItem)
		}
		TrxInvTemp := &trx.ResultFindTrxOutstanding{
			Id:              resultTrx.ID.Hex(),
			DocNo:           resultTrx.DocNo,
			GrandTotal:      resultTrx.GrandTotal,
			CheckInDatetime: resultTrx.CheckInDatetime,
			OverNightPrice:  resultTrx.OverNightPrice,
			Is24H:           resultTrx.Is24H,
			CardNumber:      resultTrx.CardNumber,
			CardNumberUUID:  resultTrx.CardNumberUUID,
			OuCode:          resultTrx.OuCode,
			VehicleNumberIn: resultTrx.VehicleNumberIn,
			TrxInvoiceItem:  tempTrxInvoiceItemss,
		}

		totalAmount := float64(0)
		grandTotal := TrxInvTemp.TrxInvoiceItem[0].TotalAmount
		checkinDate := TrxInvTemp.CheckInDatetime[:len(TrxInvTemp.CheckInDatetime)-9]

		svc.Service.RedisClientLocal.Del(fmt.Sprintf("%s-%s", TrxInvTemp.DocNo, constans.MEMBER))
		var resultTrxMemberList []*trx.TrxMember
		if exists && !isSpecialMember {
			resultActiveMemberList, err := svc.Service.MemberRepo.GetMemberActiveListByPeriod(input.UuidCard, constans.EMPTY_VALUE, checkinDate, inquiryDate, config.MEMBER_BY, false)
			if err != nil {
				result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, "Sesi Tidak Ditemukan", nil)
				return &trx.MyResponse{
					Response: result,
				}, err
			}

			for _, rows := range resultActiveMemberList {
				if rows.ProductCode == input.ProductCode {
					trxMemberList := &trx.TrxMember{
						DocNo:       TrxInvTemp.DocNo,
						PartnerCode: rows.PartnerCode,
						FirstName:   rows.FirstName,
						LastName:    rows.LastName,
						RoleType:    rows.RoleType,
						PhoneNumber: rows.PhoneNumber,
						Email:       rows.Email,
						Active:      rows.Active,
						ActiveAt:    rows.ActiveAt,
						NonActiveAt: func() *wrapperspb.StringValue {
							if rows.NonActiveAt != nil {
								return &wrapperspb.StringValue{Value: *rows.NonActiveAt}
							}
							return nil
						}(),
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
		}

		if len(resultTrxMemberList) > 0 {
			memberData := resultTrxMemberList[len(resultTrxMemberList)-1]

			grandTotal = 0

			log.Println("Grand Total After Calc Invoice Amount - Total Progressive Amount:", grandTotal)
			if grandTotal == 0 {
				trxMemberStr, _ := json.Marshal(memberData)
				svc.Service.RedisClientLocal.Set(fmt.Sprintf("%s-%s", TrxInvTemp.DocNo, constans.MEMBER), trxMemberStr, 5*time.Minute)

				resultInquiryTrxWithCard["memberCode"] = memberData.PartnerCode
				resultInquiryTrxWithCard["memberName"] = strings.TrimSpace(fmt.Sprintf("%s %s", memberData.FirstName, memberData.LastName))
				resultInquiryTrxWithCard["memberType"] = memberData.TypePartner
			}
		} else if isSpecialMember && input.ProductCode == resultMember.ProductCode {
			trxMemberList := &trx.TrxMember{
				DocNo:       TrxInvTemp.DocNo,
				PartnerCode: resultMember.PartnerCode,
				FirstName:   resultMember.FirstName,
				LastName:    resultMember.LastName,
				RoleType:    resultMember.RoleType,
				PhoneNumber: resultMember.PhoneNumber,
				Email:       resultMember.Email,
				Active:      resultMember.Active,
				ActiveAt:    resultMember.ActiveAt,
				NonActiveAt: func() *wrapperspb.StringValue {
					if resultMember.NonActiveAt != nil {
						return &wrapperspb.StringValue{Value: *resultMember.NonActiveAt}
					}
					return nil
				}(),
				OuId:               resultMember.OuId,
				TypePartner:        resultMember.TypePartner,
				CardNumber:         resultMember.CardNumber,
				VehicleNumber:      resultMember.VehicleNumber,
				RegisteredDatetime: resultMember.RegisteredDatetime,
				DateFrom:           resultMember.DateFrom,
				DateTo:             resultMember.DateTo,
				ProductId:          resultMember.ProductId,
				ProductCode:        resultMember.ProductCode,
			}

			trxMemberStr, _ := json.Marshal(trxMemberList)
			svc.Service.RedisClientLocal.Set(fmt.Sprintf("%s-%s", TrxInvTemp.DocNo, constans.MEMBER), trxMemberStr, 5*time.Minute)

			resultInquiryTrxWithCard["memberCode"] = trxMemberList.PartnerCode
			resultInquiryTrxWithCard["memberName"] = strings.TrimSpace(fmt.Sprintf("%s %s", trxMemberList.FirstName, trxMemberList.LastName))
			resultInquiryTrxWithCard["memberType"] = trxMemberList.TypePartner
		}

		totalAmount = grandTotal
		resultInquiryTrxWithCard["id"] = TrxInvTemp.Id
		resultInquiryTrxWithCard["docNo"] = TrxInvTemp.DocNo
		resultInquiryTrxWithCard["checkinDatetime"] = TrxInvTemp.CheckInDatetime
		resultInquiryTrxWithCard["cardNumberUuId"] = TrxInvTemp.CardNumberUUID
		resultInquiryTrxWithCard["cardNumber"] = TrxInvTemp.CardNumber
		resultInquiryTrxWithCard["productCode"] = TrxInvTemp.TrxInvoiceItem[0].ProductCode
		resultInquiryTrxWithCard["productName"] = TrxInvTemp.TrxInvoiceItem[0].ProductName
		resultInquiryTrxWithCard["vehicleNumberIn"] = TrxInvTemp.VehicleNumberIn
		resultInquiryTrxWithCard["totalAmount"] = totalAmount
		resultInquiryTrxWithCard["ouCode"] = TrxInvTemp.OuCode

	} else {
		merchantKey, err := utils.DecryptMerchantKey(config.MERCHANT_KEY)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}

		tempMerchantKey := &trx.MerchantKey{
			ID:              merchantKey.ID,
			OuId:            merchantKey.OuId,
			OuName:          merchantKey.OuName,
			OuCode:          merchantKey.OuCode,
			OuSubBranchId:   merchantKey.OuSubBranchId,
			OuSubBranchName: merchantKey.OuSubBranchName,
			OuSubBranchCode: merchantKey.OuSubBranchCode,
			MainOuId:        merchantKey.MainOuId,
			MainOuCode:      merchantKey.MainOuCode,
			MainOuName:      merchantKey.MainOuName,
		}

		resultTrx, exists, err := svc.Service.TrxMongoRepo.IsTrxOutstandingByUUID(input.UuidCard)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, err.Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}

		var tempTrxInvoiceItems []*trx.TrxInvoiceItem
		for _, v := range resultTrx.TrxInvoiceItem {
			tempTrxInvoiceItem := &trx.TrxInvoiceItem{
				DocNo:                  v.DocNo,
				ProductId:              v.ProductId,
				ProductCode:            v.ProductCode,
				ProductName:            v.ProductName,
				IsPctServiceFee:        v.IsPctServiceFee,
				ServiceFee:             v.ServiceFee,
				ServiceFeeMember:       v.ServiceFeeMember,
				Price:                  v.Price,
				BaseTime:               v.BaseTime,
				ProgressiveTime:        v.ProgressiveTime,
				ProgressivePrice:       v.ProgressivePrice,
				IsPct:                  v.IsPct,
				ProgressivePct:         v.ProgressivePct,
				MaxPrice:               v.MaxPrice,
				Is24H:                  v.Is24H,
				OvernightTime:          v.OvernightTime,
				OvernightPrice:         v.OvernightPrice,
				GracePeriod:            v.GracePeriod,
				FlgRepeat:              v.FlgRepeat,
				TotalAmount:            v.TotalAmount,
				TotalProgressiveAmount: v.TotalProgressiveAmount,
			}
			tempTrxInvoiceItems = append(tempTrxInvoiceItems, tempTrxInvoiceItem)
		}
		tempMemberData := &trx.TrxMember{
			DocNo:       resultTrx.MemberData.DocNo,
			PartnerCode: resultTrx.MemberData.PartnerCode,
			FirstName:   resultTrx.MemberData.FirstName,
			LastName:    resultTrx.MemberData.LastName,
			RoleType:    resultTrx.MemberData.RoleType,
			PhoneNumber: resultTrx.MemberData.PhoneNumber,
			Email:       resultTrx.MemberData.Email,
			Active:      resultTrx.MemberData.Active,
			ActiveAt:    resultTrx.MemberData.ActiveAt,
			NonActiveAt: func() *wrapperspb.StringValue {
				if resultTrx.MemberData.NonActiveAt != nil {
					return &wrapperspb.StringValue{Value: *resultTrx.MemberData.NonActiveAt}
				}
				return nil
			}(),
			OuId:               resultTrx.MemberData.OuId,
			TypePartner:        resultTrx.MemberData.TypePartner,
			CardNumber:         resultTrx.MemberData.CardNumber,
			VehicleNumber:      resultTrx.MemberData.VehicleNumber,
			RegisteredDatetime: resultTrx.MemberData.RegisteredDatetime,
			DateFrom:           resultTrx.MemberData.DateFrom,
			DateTo:             resultTrx.MemberData.DateTo,
			ProductId:          resultTrx.MemberData.ProductId,
			ProductCode:        resultTrx.MemberData.ProductCode,
		}

		convertedTrxAddInfo := make(map[string]*anypb.Any)
		for k, v := range resultTrx.TrxAddInfo {
			if msg, ok := v.(goproto.Message); ok {
				anyVal, err := anypb.New(msg)
				if err != nil {
					continue
				}
				convertedTrxAddInfo[k] = anyVal
			}
		}

		tempRequestAddTrxInvoiceDetailItem := &trx.TrxInvoiceDetailItem{
			DocNo:         resultTrx.RequestAddTrxInvoiceDetailItem.DocNo,
			ProductCode:   resultTrx.RequestAddTrxInvoiceDetailItem.ProductCode,
			InvoiceAmount: resultTrx.RequestAddTrxInvoiceDetailItem.InvoiceAmount,
			CreatedAt:     resultTrx.RequestAddTrxInvoiceDetailItem.CreatedAt,
			CreatedDate:   resultTrx.RequestAddTrxInvoiceDetailItem.CreatedDate,
		}

		tempTrxOutstanding := &trx.Trx{
			DocNo:                          resultTrx.DocNo,
			DocDate:                        resultTrx.DocDate,
			PaymentRefDocno:                resultTrx.PaymentRefDocNo,
			CheckinDateTime:                resultTrx.CheckInDatetime,
			CheckoutDateTime:               resultTrx.CheckOutDatetime,
			DeviceIdIn:                     resultTrx.DeviceIdIn,
			DeviceId:                       resultTrx.DeviceId,
			GetIn:                          resultTrx.GateIn,
			GetOut:                         resultTrx.GateOut,
			CardNumberUUIDIn:               resultTrx.CardNumberUUIDIn,
			CardNumberIn:                   resultTrx.CardNumberIn,
			CardNumberUUID:                 resultTrx.CardNumberUUID,
			CardNumber:                     resultTrx.CardNumber,
			TypeCard:                       resultTrx.TypeCard,
			BeginningBalance:               resultTrx.BeginningBalance,
			ExtLocalDateTime:               resultTrx.ExtLocalDatetime,
			ChargeAmount:                   resultTrx.ChargeAmount,
			GrandTotal:                     resultTrx.GrandTotal,
			ProductCode:                    resultTrx.ProductCode,
			ProductName:                    resultTrx.ProductName,
			ProductData:                    resultTrx.ProductData,
			RequestData:                    resultTrx.RequestData,
			RequestOutData:                 resultTrx.RequestOutData,
			OuId:                           resultTrx.OuId,
			OuName:                         resultTrx.OuName,
			OuCode:                         resultTrx.OuCode,
			OuSubBranchId:                  resultTrx.OuSubBranchId,
			USubBranchName:                 resultTrx.OuSubBranchName,
			OuSubBranchCode:                resultTrx.OuSubBranchCode,
			MainOuId:                       resultTrx.MainOuId,
			MainOuCode:                     resultTrx.MainOuCode,
			MainOuName:                     resultTrx.MainOuName,
			MemberCode:                     resultTrx.MemberCode,
			MemberName:                     resultTrx.MemberName,
			MemberType:                     resultTrx.MemberType,
			MemberStatus:                   resultTrx.MemberStatus,
			MemberExpiredDate:              resultTrx.MemberExpiredDate,
			CheckInTime:                    resultTrx.CheckInTime,
			CheckOutTime:                   resultTrx.CheckOutTime,
			DurationTime:                   resultTrx.DurationTime,
			VehicleNumberIn:                resultTrx.VehicleNumberIn,
			VehicleNumberOut:               resultTrx.VehicleNumberOut,
			LogTrans:                       resultTrx.LogTrans,
			MerchantKey:                    resultTrx.MerchantKey,
			QrText:                         resultTrx.QrText,
			QrA2P:                          resultTrx.QrA2P,
			QrTextPaymentOnline:            resultTrx.QrTextPaymentOnline,
			TrxInvoiceItem:                 tempTrxInvoiceItems,
			FlagSyncData:                   resultTrx.FlagSyncData,
			MemberData:                     tempMemberData,
			TrxAddInfo:                     convertedTrxAddInfo,
			FlagTrxFromCloud:               resultTrx.FlagTrxFromCloud,
			IsRsyncDataTrx:                 resultTrx.IsRsyncDataTrx,
			ExcludeSf:                      resultTrx.ExcludeSf,
			FlagCharge:                     resultTrx.FlagCharge,
			ChargeType:                     resultTrx.ChargeType,
			RequestAddTrxInvoiceDetailItem: tempRequestAddTrxInvoiceDetailItem,
			LastUpdatedAt:                  resultTrx.LastUpdatedAt,
		}

		resultProduct, err := svc.Service.ProductRepo.FindProductByProductCode(input.ProductCode, tempMerchantKey.OuId)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, err.Error(), nil)
			return &trx.MyResponse{
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
			result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, err.Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}

		docNo := fmt.Sprintf("%s%s%d", prefixDocNo, autoNumber, tempMerchantKey.OuId)
		randomKey := utils.RandStringBytesMaskImprSrc(16)
		encQrTxt, _ := utils.Encrypt(docNo, randomKey)
		qrCode := fmt.Sprintf("%s%s", encQrTxt.Result, randomKey)

		trxo := tempTrxOutstanding
		if !exists {
			var trxInvoiceItemList []*trx.TrxInvoiceItem
			trxInvoiceItem := &trx.TrxInvoiceItem{
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

			trxo = &trx.Trx{
				DocNo:            docNo,
				DocDate:          constans.EMPTY_VALUE,
				CheckinDateTime:  input.InquiryDateTime,
				CheckoutDateTime: constans.EMPTY_VALUE,
				DeviceIdIn:       constans.EMPTY_VALUE,
				DeviceId:         constans.EMPTY_VALUE,
				GetIn:            constans.EMPTY_VALUE,
				GetOut:           constans.EMPTY_VALUE,
				CardNumberUUIDIn: input.UuidCard,
				CardNumberIn:     constans.EMPTY_VALUE,
				CardNumberUUID:   input.UuidCard,
				CardNumber:       constans.EMPTY_VALUE,
				TypeCard:         constans.EMPTY_VALUE,
				BeginningBalance: 0,
				ExtLocalDateTime: input.InquiryDateTime,
				GrandTotal:       0,
				ProductCode:      resultProduct.ProductCode,
				ProductName:      resultProduct.ProductName,
				ProductData:      utils.ToString(resultProduct),
				RequestData:      utils.ToString(input),
				RequestOutData:   utils.ToString(input),
				OuId:             tempMerchantKey.OuId,
				OuName:           tempMerchantKey.OuName,
				OuCode:           tempMerchantKey.OuCode,
				OuSubBranchId:    tempMerchantKey.OuSubBranchId,
				USubBranchName:   tempMerchantKey.OuSubBranchName,
				OuSubBranchCode:  tempMerchantKey.OuSubBranchCode,
				MainOuId:         tempMerchantKey.MainOuId,
				MainOuCode:       tempMerchantKey.MainOuCode,
				MainOuName:       tempMerchantKey.MainOuName,
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
			trxo.MemberCode = resultMember.PartnerCode
			trxo.MemberName = strings.TrimSpace(fmt.Sprintf("%s %s", resultMember.FirstName, resultMember.LastName))
			trxo.MemberType = resultMember.TypePartner
		}

		resultInquiryTrxWithCard["id"] = constans.TYPE_PARTNER_FREE_PASS
		resultInquiryTrxWithCard["docNo"] = tempTrxOutstanding.DocNo
		resultInquiryTrxWithCard["checkinDatetime"] = tempTrxOutstanding.CheckinDateTime
		resultInquiryTrxWithCard["cardNumberUuId"] = input.UuidCard
		resultInquiryTrxWithCard["cardNumber"] = constans.EMPTY_VALUE
		resultInquiryTrxWithCard["productCode"] = resultProduct.ProductCode
		resultInquiryTrxWithCard["productName"] = resultProduct.ProductName
		resultInquiryTrxWithCard["vehicleNumberIn"] = constans.EMPTY_VALUE
		resultInquiryTrxWithCard["totalAmount"] = float64(0)
		resultInquiryTrxWithCard["memberCode"] = trxo.MemberCode
		resultInquiryTrxWithCard["memberName"] = trxo.MemberName
		resultInquiryTrxWithCard["memberType"] = trxo.MemberType
		resultInquiryTrxWithCard["ouCode"] = trxo.OuCode

		redisStatus := svc.Service.RedisClientLocal.Set(input.UuidCard, utils.ToString(trxo), 30*time.Second)
		if redisStatus.Err() != nil {
			result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, redisStatus.Err().Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}

	}

	checkinDatetimeParse, _ := time.Parse("2006-01-02 15:04:05", resultInquiryTrxWithCard["checkinDatetime"].(string))
	checkoutDatetimeParse, _ := time.Parse("2006-01-02 15:04:05", input.InquiryDateTime)
	duration := utils.ConvDiffTime(checkinDatetimeParse, checkoutDatetimeParse)

	anyDuration, _ := anypb.New(duration)

	responseTrx := &trx.ResultInquiryTrxWithCard{
		Id:              resultInquiryTrxWithCard["id"].(string),
		DocNo:           resultInquiryTrxWithCard["docNo"].(string),
		CardNumberUUID:  resultInquiryTrxWithCard["cardNumberUuId"].(string),
		CardNumber:      resultInquiryTrxWithCard["cardNumber"].(string),
		Nominal:         resultInquiryTrxWithCard["totalAmount"].(float64),
		ProductCode:     resultInquiryTrxWithCard["productCode"].(string),
		ProductName:     resultInquiryTrxWithCard["productName"].(string),
		VehicleNumberIn: resultInquiryTrxWithCard["vehicleNumberIn"].(string),
		MemberCode:      resultInquiryTrxWithCard["memberCode"].(string),
		MemberName:      resultInquiryTrxWithCard["memberName"].(string),
		MemberType:      resultInquiryTrxWithCard["memberType"].(string),
		QrCode:          constans.EMPTY_VALUE,
		ExcludeSf:       false,
		Duration:        anyDuration,
		OuCode:          resultInquiryTrxWithCard["ouCode"].(string),
	}

	anyResponseTrx, _ := anypb.New(responseTrx)

	result = helpers.ResponseJSON(true, constans.SUCCESS_CODE, constans.EMPTY_VALUE, anyResponseTrx)
	return &trx.MyResponse{
		Response: result,
	}, nil
}

func (svc trxService) ConfirmTrx(ctx context.Context, input *trx.RequestConfirmTrx) (*trx.MyResponse, error) {
	var result *trx.Response
	var resultTrx *trx.Trx
	var resultTrxs models.Trx

	request := new(models.RequestConfirmTrx)
	if err := helpers.BindValidateStruct(request); err != nil {
		result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}
	if err := helpers.BindValidateStruct(input); err != nil {
		result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	if input.Id == constans.TYPE_PARTNER_FREE_PASS {
		redisStatus := svc.Service.RedisClientLocal.Get(input.UuidCard)
		if redisStatus.Err() != nil {
			result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, redisStatus.Err().Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, redisStatus.Err()
		}

		if err := json.Unmarshal([]byte(redisStatus.Val()), &resultTrx); err != nil {
			result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, err.Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}

		if err := json.Unmarshal([]byte(redisStatus.Val()), &resultTrxs); err != nil {
			result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, err.Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}

		go helperService.CallSyncConfirmTrxForMemberFreePass(*request, resultTrxs, svc.Service)
	} else if strings.Contains(input.Id, "SERVER") {
		redisStatus := svc.Service.RedisClientLocal.Get(input.Id)
		if redisStatus.Err() != nil {
			result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, redisStatus.Err().Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, redisStatus.Err()
		}

		if err := json.Unmarshal([]byte(redisStatus.Val()), &resultTrx); err != nil {
			result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, err.Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}

		if err := json.Unmarshal([]byte(redisStatus.Val()), &resultTrxs); err != nil {
			result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, err.Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}

		go helperService.CallSyncConfirmTrxToCloud(nil, *request, resultTrxs, svc.Service)
	} else {
		ID, _ := primitive.ObjectIDFromHex(request.ID)

		resultDataTrx, err := svc.Service.TrxMongoRepo.FindTrxOutstandingByID(ID)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, "Sesi Tidak Ditemukan", nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}

		resultTrxs = resultDataTrx

		go helperService.CallSyncConfirmTrxToCloud(&ID, *request, resultTrxs, svc.Service)
	}

	ipAddr, _ := helpers.GetPrivateIPLocal()

	responseConfirm := &trx.ResponseConfirm{
		DocNo:            resultTrx.DocNo,
		ProductData:      resultTrx.ProductData,
		ProductName:      input.ProductName,
		CardType:         resultTrx.TypeCard,
		CardNumber:       input.CardNumber,
		CheckInDatetime:  resultTrx.CheckinDateTime,
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
		IpAddr:           ipAddr,
	}

	anyResponseTrx, _ := anypb.New(responseConfirm)

	result = helpers.ResponseJSON(true, constans.SUCCESS_CODE, constans.EMPTY_VALUE, anyResponseTrx)
	return &trx.MyResponse{
		Response: result,
	}, nil
}

func (svc trxService) ConfirmTrxByPass(ctx context.Context, input *trx.ConfirmTrxByPassMessage) (*trx.MyResponse, error) {
	var result *trx.Response
	var trxs *trx.Trx

	if err := helpers.BindValidateStruct(input); err != nil {
		result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	merchantKey, err := utils.DecryptMerchantKey(config.MERCHANT_KEY)
	if err != nil {
		result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	tempMerchantKey := &trx.MerchantKey{
		ID:              merchantKey.ID,
		OuId:            merchantKey.OuId,
		OuName:          merchantKey.OuName,
		OuCode:          merchantKey.OuCode,
		OuSubBranchId:   merchantKey.OuSubBranchId,
		OuSubBranchName: merchantKey.OuSubBranchName,
		OuSubBranchCode: merchantKey.OuSubBranchCode,
		MainOuId:        merchantKey.MainOuId,
		MainOuCode:      merchantKey.MainOuCode,
		MainOuName:      merchantKey.MainOuName,
	}

	resultProduct, err := svc.Service.ProductRepo.FindProductByProductCode(input.ProductCode, tempMerchantKey.OuId)
	if err != nil {
		result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	resultTrx, exists, err := svc.Service.TrxMongoRepo.IsTrxOutstandingByCardNumber(input.CardNumber)
	if err != nil {
		result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	var tempTrxInvoiceItems []*trx.TrxInvoiceItem
	for _, v := range resultTrx.TrxInvoiceItem {
		tempTrxInvoiceItem := &trx.TrxInvoiceItem{
			DocNo:                  v.DocNo,
			ProductId:              v.ProductId,
			ProductCode:            v.ProductCode,
			ProductName:            v.ProductName,
			IsPctServiceFee:        v.IsPctServiceFee,
			ServiceFee:             v.ServiceFee,
			ServiceFeeMember:       v.ServiceFeeMember,
			Price:                  v.Price,
			BaseTime:               v.BaseTime,
			ProgressiveTime:        v.ProgressiveTime,
			ProgressivePrice:       v.ProgressivePrice,
			IsPct:                  v.IsPct,
			ProgressivePct:         v.ProgressivePct,
			MaxPrice:               v.MaxPrice,
			Is24H:                  v.Is24H,
			OvernightTime:          v.OvernightTime,
			OvernightPrice:         v.OvernightPrice,
			GracePeriod:            v.GracePeriod,
			FlgRepeat:              v.FlgRepeat,
			TotalAmount:            v.TotalAmount,
			TotalProgressiveAmount: v.TotalProgressiveAmount,
		}
		tempTrxInvoiceItems = append(tempTrxInvoiceItems, tempTrxInvoiceItem)
	}
	tempMemberData := &trx.TrxMember{
		DocNo:       resultTrx.MemberData.DocNo,
		PartnerCode: resultTrx.MemberData.PartnerCode,
		FirstName:   resultTrx.MemberData.FirstName,
		LastName:    resultTrx.MemberData.LastName,
		RoleType:    resultTrx.MemberData.RoleType,
		PhoneNumber: resultTrx.MemberData.PhoneNumber,
		Email:       resultTrx.MemberData.Email,
		Active:      resultTrx.MemberData.Active,
		ActiveAt:    resultTrx.MemberData.ActiveAt,
		NonActiveAt: func() *wrapperspb.StringValue {
			if resultTrx.MemberData.NonActiveAt != nil {
				return &wrapperspb.StringValue{Value: *resultTrx.MemberData.NonActiveAt}
			}
			return nil
		}(),
		OuId:               resultTrx.MemberData.OuId,
		TypePartner:        resultTrx.MemberData.TypePartner,
		CardNumber:         resultTrx.MemberData.CardNumber,
		VehicleNumber:      resultTrx.MemberData.VehicleNumber,
		RegisteredDatetime: resultTrx.MemberData.RegisteredDatetime,
		DateFrom:           resultTrx.MemberData.DateFrom,
		DateTo:             resultTrx.MemberData.DateTo,
		ProductId:          resultTrx.MemberData.ProductId,
		ProductCode:        resultTrx.MemberData.ProductCode,
	}

	convertedTrxAddInfo := make(map[string]*anypb.Any)
	for k, v := range resultTrx.TrxAddInfo {
		if msg, ok := v.(goproto.Message); ok {
			anyVal, err := anypb.New(msg)
			if err != nil {
				continue
			}
			convertedTrxAddInfo[k] = anyVal
		}
	}

	tempRequestAddTrxInvoiceDetailItem := &trx.TrxInvoiceDetailItem{
		DocNo:         resultTrx.RequestAddTrxInvoiceDetailItem.DocNo,
		ProductCode:   resultTrx.RequestAddTrxInvoiceDetailItem.ProductCode,
		InvoiceAmount: resultTrx.RequestAddTrxInvoiceDetailItem.InvoiceAmount,
		CreatedAt:     resultTrx.RequestAddTrxInvoiceDetailItem.CreatedAt,
		CreatedDate:   resultTrx.RequestAddTrxInvoiceDetailItem.CreatedDate,
	}

	tempTrxOutstanding := &trx.Trx{
		DocNo:                          resultTrx.DocNo,
		DocDate:                        resultTrx.DocDate,
		PaymentRefDocno:                resultTrx.PaymentRefDocNo,
		CheckinDateTime:                resultTrx.CheckInDatetime,
		CheckoutDateTime:               resultTrx.CheckOutDatetime,
		DeviceIdIn:                     resultTrx.DeviceIdIn,
		DeviceId:                       resultTrx.DeviceId,
		GetIn:                          resultTrx.GateIn,
		GetOut:                         resultTrx.GateOut,
		CardNumberUUIDIn:               resultTrx.CardNumberUUIDIn,
		CardNumberIn:                   resultTrx.CardNumberIn,
		CardNumberUUID:                 resultTrx.CardNumberUUID,
		CardNumber:                     resultTrx.CardNumber,
		TypeCard:                       resultTrx.TypeCard,
		BeginningBalance:               resultTrx.BeginningBalance,
		ExtLocalDateTime:               resultTrx.ExtLocalDatetime,
		ChargeAmount:                   resultTrx.ChargeAmount,
		GrandTotal:                     resultTrx.GrandTotal,
		ProductCode:                    resultTrx.ProductCode,
		ProductName:                    resultTrx.ProductName,
		ProductData:                    resultTrx.ProductData,
		RequestData:                    resultTrx.RequestData,
		RequestOutData:                 resultTrx.RequestOutData,
		OuId:                           resultTrx.OuId,
		OuName:                         resultTrx.OuName,
		OuCode:                         resultTrx.OuCode,
		OuSubBranchId:                  resultTrx.OuSubBranchId,
		USubBranchName:                 resultTrx.OuSubBranchName,
		OuSubBranchCode:                resultTrx.OuSubBranchCode,
		MainOuId:                       resultTrx.MainOuId,
		MainOuCode:                     resultTrx.MainOuCode,
		MainOuName:                     resultTrx.MainOuName,
		MemberCode:                     resultTrx.MemberCode,
		MemberName:                     resultTrx.MemberName,
		MemberType:                     resultTrx.MemberType,
		MemberStatus:                   resultTrx.MemberStatus,
		MemberExpiredDate:              resultTrx.MemberExpiredDate,
		CheckInTime:                    resultTrx.CheckInTime,
		CheckOutTime:                   resultTrx.CheckOutTime,
		DurationTime:                   resultTrx.DurationTime,
		VehicleNumberIn:                resultTrx.VehicleNumberIn,
		VehicleNumberOut:               resultTrx.VehicleNumberOut,
		LogTrans:                       resultTrx.LogTrans,
		MerchantKey:                    resultTrx.MerchantKey,
		QrText:                         resultTrx.QrText,
		QrA2P:                          resultTrx.QrA2P,
		QrTextPaymentOnline:            resultTrx.QrTextPaymentOnline,
		TrxInvoiceItem:                 tempTrxInvoiceItems,
		FlagSyncData:                   resultTrx.FlagSyncData,
		MemberData:                     tempMemberData,
		TrxAddInfo:                     convertedTrxAddInfo,
		FlagTrxFromCloud:               resultTrx.FlagTrxFromCloud,
		IsRsyncDataTrx:                 resultTrx.IsRsyncDataTrx,
		ExcludeSf:                      resultTrx.ExcludeSf,
		FlagCharge:                     resultTrx.FlagCharge,
		ChargeType:                     resultTrx.ChargeType,
		RequestAddTrxInvoiceDetailItem: tempRequestAddTrxInvoiceDetailItem,
		LastUpdatedAt:                  resultTrx.LastUpdatedAt,
	}

	if exists {
		trxs = tempTrxOutstanding
		trxs.GrandTotal = input.GrandTotal
		trxs.ProductCode = resultProduct.ProductCode
		trxs.ProductName = resultProduct.ProductName
		trxs.DeviceId = input.DeviceId
		trxs.CardNumber = input.CardNumber
	} else {
		checkinDatetimeParse, _ := time.Parse("2006-01-02 15:04:05", utils.Timestamp())
		yearMonth := checkinDatetimeParse.Format("060102")
		HourCheckIn := checkinDatetimeParse.Format("15")
		prefixDocNo := utils.RandStringBytesMaskImprSrcChr(4)
		prefix := fmt.Sprintf("%s%s%s", prefixDocNo, yearMonth, HourCheckIn)
		autoNumber, err := svc.Service.GenAutoNumRepo.AutonumberValueWithDatatype(constans.DATATYPE_TRX_LOCAL, prefix, 4)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, err.Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}

		docNo := fmt.Sprintf("%s%d", autoNumber, merchantKey.OuId)
		randomKey := utils.RandStringBytesMaskImprSrc(16)
		encQrTxt, _ := utils.Encrypt(docNo, randomKey)
		qrCode := fmt.Sprintf("%s%s", encQrTxt.Result, randomKey)

		trxs = &trx.Trx{
			DocNo:            docNo,
			DocDate:          utils.DateNow(),
			CheckinDateTime:  utils.Timestamp(),
			CheckoutDateTime: utils.Timestamp(),
			DeviceIdIn:       input.DeviceId,
			DeviceId:         input.DeviceId,
			GetIn:            input.IpTerminal,
			GetOut:           input.IpTerminal,
			CardNumberUUIDIn: constans.EMPTY_VALUE,
			CardNumberIn:     input.CardNumber,
			CardNumberUUID:   constans.VALIDATE_ERROR_CODE,
			CardNumber:       input.CardNumber,
			TypeCard:         input.CardType,
			BeginningBalance: 0,
			ExtLocalDateTime: utils.Timestamp(),
			GrandTotal:       input.GrandTotal,
			ProductCode:      resultProduct.ProductCode,
			ProductName:      resultProduct.ProductName,
			ProductData:      utils.ToString(resultProduct),
			RequestData:      utils.ToString(input),
			RequestOutData:   utils.ToString(input),
			OuId:             merchantKey.OuId,
			OuName:           merchantKey.OuName,
			OuCode:           merchantKey.OuCode,
			OuSubBranchId:    merchantKey.OuSubBranchId,
			USubBranchName:   merchantKey.OuSubBranchName,
			OuSubBranchCode:  merchantKey.OuSubBranchCode,
			MainOuId:         merchantKey.MainOuId,
			MainOuCode:       merchantKey.MainOuCode,
			MainOuName:       merchantKey.MainOuName,
			MemberCode:       constans.EMPTY_VALUE,
			MemberName:       constans.EMPTY_VALUE,
			MemberType:       constans.EMPTY_VALUE,
			CheckInTime:      0,
			CheckOutTime:     0,
			DurationTime:     0,
			VehicleNumberIn:  constans.EMPTY_VALUE,
			VehicleNumberOut: constans.EMPTY_VALUE,
			LogTrans:         constans.EMPTY_VALUE,
			MerchantKey:      config.MERCHANT_KEY,
			QrText:           qrCode,
			TrxInvoiceItem:   nil,
			FlagSyncData:     false,
			MemberData:       nil,
			TrxAddInfo:       nil,
		}

	}

	var trxInvoiceItemList []*trx.TrxInvoiceItem
	trxInvoiceItem := &trx.TrxInvoiceItem{
		DocNo:            trxs.DocNo,
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
		TotalAmount:      input.GrandTotal,
	}

	trxInvoiceItemList = append(trxInvoiceItemList, trxInvoiceItem)
	trxs.TrxInvoiceItem = trxInvoiceItemList

	checkoutDatetimeParse, _ := time.Parse("2006-01-02 15:04:05", trxs.CheckoutDateTime)
	timeCheckOutUnix := checkoutDatetimeParse.Unix()
	durationTime := (timeCheckOutUnix - tempTrxOutstanding.CheckInTime) / 60
	trxs.CheckOutTime = timeCheckOutUnix
	trxs.DurationTime = durationTime

	dataStr, _ := json.Marshal(trxs)
	if utils.IsConnected() {
		redisStatus := svc.Service.RedisClient.Publish(constans.CHANNEL_REDIS_PG_PARKING_CHECKOUT, dataStr)
		if redisStatus.Err() != nil || redisStatus.Val() == 0 {
			tempTrxOutstanding.FlagSyncData = false
		}
	} else {
		tempTrxOutstanding.FlagSyncData = false
	}
	var tempTrxInvoiceItemss []models.TrxInvoiceItem
	for _, v := range resultTrx.TrxInvoiceItem {
		tempTrxInvoiceItemm := models.TrxInvoiceItem{
			DocNo:                  v.DocNo,
			ProductId:              v.ProductId,
			ProductCode:            v.ProductCode,
			ProductName:            v.ProductName,
			IsPctServiceFee:        v.IsPctServiceFee,
			ServiceFee:             v.ServiceFee,
			ServiceFeeMember:       v.ServiceFeeMember,
			Price:                  v.Price,
			BaseTime:               v.BaseTime,
			ProgressiveTime:        v.ProgressiveTime,
			ProgressivePrice:       v.ProgressivePrice,
			IsPct:                  v.IsPct,
			ProgressivePct:         v.ProgressivePct,
			MaxPrice:               v.MaxPrice,
			Is24H:                  v.Is24H,
			OvernightTime:          v.OvernightTime,
			OvernightPrice:         v.OvernightPrice,
			GracePeriod:            v.GracePeriod,
			FlgRepeat:              v.FlgRepeat,
			TotalAmount:            v.TotalAmount,
			TotalProgressiveAmount: v.TotalProgressiveAmount,
		}
		tempTrxInvoiceItemss = append(tempTrxInvoiceItemss, tempTrxInvoiceItemm)
	}
	tempMemberDatas := &models.TrxMember{
		DocNo:              trxs.MemberData.DocNo,
		PartnerCode:        trxs.MemberData.PartnerCode,
		FirstName:          trxs.MemberData.FirstName,
		LastName:           trxs.MemberData.LastName,
		RoleType:           trxs.MemberData.RoleType,
		PhoneNumber:        trxs.MemberData.PhoneNumber,
		Email:              trxs.MemberData.Email,
		Active:             trxs.MemberData.Active,
		ActiveAt:           trxs.MemberData.ActiveAt,
		NonActiveAt:        &trxs.MemberData.NonActiveAt.Value,
		OuId:               resultTrx.MemberData.OuId,
		TypePartner:        resultTrx.MemberData.TypePartner,
		CardNumber:         resultTrx.MemberData.CardNumber,
		VehicleNumber:      resultTrx.MemberData.VehicleNumber,
		RegisteredDatetime: resultTrx.MemberData.RegisteredDatetime,
		DateFrom:           resultTrx.MemberData.DateFrom,
		DateTo:             resultTrx.MemberData.DateTo,
		ProductId:          resultTrx.MemberData.ProductId,
		ProductCode:        resultTrx.MemberData.ProductCode,
	}

	convertedTrxAddInfos := make(map[string]interface{})
	for k, v := range resultTrx.TrxAddInfo {
		if msg, ok := v.(goproto.Message); ok {
			anyVal, _ := msg.(interface{})

			convertedTrxAddInfos[k] = anyVal
		}
	}

	tempRequestAddTrxInvoiceDetailItems := &models.TrxInvoiceDetailItem{
		DocNo:         resultTrx.RequestAddTrxInvoiceDetailItem.DocNo,
		ProductCode:   resultTrx.RequestAddTrxInvoiceDetailItem.ProductCode,
		InvoiceAmount: resultTrx.RequestAddTrxInvoiceDetailItem.InvoiceAmount,
		CreatedAt:     resultTrx.RequestAddTrxInvoiceDetailItem.CreatedAt,
		CreatedDate:   resultTrx.RequestAddTrxInvoiceDetailItem.CreatedDate,
	}
	trxm := models.Trx{
		DocNo:                          trxs.DocNo,
		DocDate:                        trxs.DocDate,
		PaymentRefDocNo:                trxs.PaymentRefDocno,
		CheckInDatetime:                trxs.CheckinDateTime,
		CheckOutDatetime:               trxs.CheckoutDateTime,
		DeviceIdIn:                     trxs.DeviceIdIn,
		DeviceId:                       trxs.DeviceId,
		GateIn:                         trxs.GetIn,
		GateOut:                        trxs.GetOut,
		CardNumberUUIDIn:               trxs.CardNumberUUIDIn,
		CardNumberIn:                   trxs.CardNumberIn,
		CardNumberUUID:                 trxs.CardNumberUUID,
		CardNumber:                     trxs.CardNumber,
		TypeCard:                       trxs.TypeCard,
		BeginningBalance:               trxs.BeginningBalance,
		ExtLocalDatetime:               trxs.ExtLocalDateTime,
		ChargeAmount:                   trxs.ChargeAmount,
		GrandTotal:                     trxs.GrandTotal,
		ProductCode:                    trxs.ProductCode,
		ProductName:                    trxs.ProductName,
		ProductData:                    trxs.ProductData,
		RequestData:                    trxs.RequestData,
		RequestOutData:                 trxs.RequestOutData,
		OuId:                           trxs.OuId,
		OuName:                         trxs.OuName,
		OuCode:                         trxs.OuCode,
		OuSubBranchId:                  trxs.OuSubBranchId,
		OuSubBranchName:                trxs.USubBranchName,
		OuSubBranchCode:                trxs.OuSubBranchCode,
		MainOuId:                       trxs.MainOuId,
		MainOuCode:                     trxs.MainOuCode,
		MainOuName:                     trxs.MainOuName,
		MemberCode:                     trxs.MemberCode,
		MemberName:                     trxs.MemberName,
		MemberType:                     trxs.MemberType,
		MemberStatus:                   trxs.MemberStatus,
		MemberExpiredDate:              trxs.MemberExpiredDate,
		CheckInTime:                    trxs.CheckInTime,
		CheckOutTime:                   trxs.CheckOutTime,
		DurationTime:                   trxs.DurationTime,
		VehicleNumberIn:                trxs.VehicleNumberIn,
		VehicleNumberOut:               trxs.VehicleNumberOut,
		LogTrans:                       trxs.LogTrans,
		MerchantKey:                    trxs.MerchantKey,
		QrText:                         trxs.QrText,
		QrA2P:                          trxs.QrA2P,
		QrTextPaymentOnline:            trxs.QrTextPaymentOnline,
		TrxInvoiceItem:                 tempTrxInvoiceItemss,
		FlagSyncData:                   trxs.FlagSyncData,
		MemberData:                     tempMemberDatas,
		TrxAddInfo:                     convertedTrxAddInfos,
		FlagTrxFromCloud:               trxs.FlagTrxFromCloud,
		IsRsyncDataTrx:                 trxs.IsRsyncDataTrx,
		ExcludeSf:                      trxs.ExcludeSf,
		FlagCharge:                     trxs.FlagCharge,
		ChargeType:                     trxs.ChargeType,
		RequestAddTrxInvoiceDetailItem: tempRequestAddTrxInvoiceDetailItems,
		LastUpdatedAt:                  trxs.LastUpdatedAt,
	}

	_, err = svc.Service.TrxMongoRepo.AddTrx(trxm)
	if err != nil {
		log.Println("ERROR AddTrx : ", err.Error())
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	err = svc.Service.TrxMongoRepo.RemoveTrxByDocNo(trxm.DocNo)
	if err != nil {
		log.Println("ERROR RemoveTrxByDocNo : ", err.Error())
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	result = helpers.ResponseJSON(true, constans.SUCCESS_CODE, constans.EMPTY_VALUE, nil)
	return &trx.MyResponse{
		Response: result,
	}, nil
}

func (svc trxService) ConfirmSyncTrxToCloud(ctx context.Context, input *trx.Empty) (*trx.MyResponse, error) {
	var result *trx.Response
	var trxLists []*trx.Trx

	trxList, err := svc.Service.TrxMongoRepo.GetTrxListForSyncDataFailed()
	if err != nil {
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	for _, x := range trxList {
		var tempTrxInvoiceItems []*trx.TrxInvoiceItem
		for _, v := range x.TrxInvoiceItem {
			tempTrxInvoiceItem := &trx.TrxInvoiceItem{
				DocNo:                  v.DocNo,
				ProductId:              v.ProductId,
				ProductCode:            v.ProductCode,
				ProductName:            v.ProductName,
				IsPctServiceFee:        v.IsPctServiceFee,
				ServiceFee:             v.ServiceFee,
				ServiceFeeMember:       v.ServiceFeeMember,
				Price:                  v.Price,
				BaseTime:               v.BaseTime,
				ProgressiveTime:        v.ProgressiveTime,
				ProgressivePrice:       v.ProgressivePrice,
				IsPct:                  v.IsPct,
				ProgressivePct:         v.ProgressivePct,
				MaxPrice:               v.MaxPrice,
				Is24H:                  v.Is24H,
				OvernightTime:          v.OvernightTime,
				OvernightPrice:         v.OvernightPrice,
				GracePeriod:            v.GracePeriod,
				FlgRepeat:              v.FlgRepeat,
				TotalAmount:            v.TotalAmount,
				TotalProgressiveAmount: v.TotalProgressiveAmount,
			}
			tempTrxInvoiceItems = append(tempTrxInvoiceItems, tempTrxInvoiceItem)
		}
		tempMemberData := &trx.TrxMember{
			DocNo:       x.MemberData.DocNo,
			PartnerCode: x.MemberData.PartnerCode,
			FirstName:   x.MemberData.FirstName,
			LastName:    x.MemberData.LastName,
			RoleType:    x.MemberData.RoleType,
			PhoneNumber: x.MemberData.PhoneNumber,
			Email:       x.MemberData.Email,
			Active:      x.MemberData.Active,
			ActiveAt:    x.MemberData.ActiveAt,
			NonActiveAt: func() *wrapperspb.StringValue {
				if x.MemberData.NonActiveAt != nil {
					return &wrapperspb.StringValue{Value: *x.MemberData.NonActiveAt}
				}
				return nil
			}(),
			OuId:               x.MemberData.OuId,
			TypePartner:        x.MemberData.TypePartner,
			CardNumber:         x.MemberData.CardNumber,
			VehicleNumber:      x.MemberData.VehicleNumber,
			RegisteredDatetime: x.MemberData.RegisteredDatetime,
			DateFrom:           x.MemberData.DateFrom,
			DateTo:             x.MemberData.DateTo,
			ProductId:          x.MemberData.ProductId,
			ProductCode:        x.MemberData.ProductCode,
		}

		convertedTrxAddInfo := make(map[string]*anypb.Any)
		for k, v := range x.TrxAddInfo {
			if msg, ok := v.(goproto.Message); ok {
				anyVal, err := anypb.New(msg)
				if err != nil {
					continue
				}
				convertedTrxAddInfo[k] = anyVal
			}
		}

		tempRequestAddTrxInvoiceDetailItem := &trx.TrxInvoiceDetailItem{
			DocNo:         x.RequestAddTrxInvoiceDetailItem.DocNo,
			ProductCode:   x.RequestAddTrxInvoiceDetailItem.ProductCode,
			InvoiceAmount: x.RequestAddTrxInvoiceDetailItem.InvoiceAmount,
			CreatedAt:     x.RequestAddTrxInvoiceDetailItem.CreatedAt,
			CreatedDate:   x.RequestAddTrxInvoiceDetailItem.CreatedDate,
		}

		tempTrxOutstanding := &trx.Trx{
			DocNo:                          x.DocNo,
			DocDate:                        x.DocDate,
			PaymentRefDocno:                x.PaymentRefDocNo,
			CheckinDateTime:                x.CheckInDatetime,
			CheckoutDateTime:               x.CheckOutDatetime,
			DeviceIdIn:                     x.DeviceIdIn,
			DeviceId:                       x.DeviceId,
			GetIn:                          x.GateIn,
			GetOut:                         x.GateOut,
			CardNumberUUIDIn:               x.CardNumberUUIDIn,
			CardNumberIn:                   x.CardNumberIn,
			CardNumberUUID:                 x.CardNumberUUID,
			CardNumber:                     x.CardNumber,
			TypeCard:                       x.TypeCard,
			BeginningBalance:               x.BeginningBalance,
			ExtLocalDateTime:               x.ExtLocalDatetime,
			ChargeAmount:                   x.ChargeAmount,
			GrandTotal:                     x.GrandTotal,
			ProductCode:                    x.ProductCode,
			ProductName:                    x.ProductName,
			ProductData:                    x.ProductData,
			RequestData:                    x.RequestData,
			RequestOutData:                 x.RequestOutData,
			OuId:                           x.OuId,
			OuName:                         x.OuName,
			OuCode:                         x.OuCode,
			OuSubBranchId:                  x.OuSubBranchId,
			USubBranchName:                 x.OuSubBranchName,
			OuSubBranchCode:                x.OuSubBranchCode,
			MainOuId:                       x.MainOuId,
			MainOuCode:                     x.MainOuCode,
			MainOuName:                     x.MainOuName,
			MemberCode:                     x.MemberCode,
			MemberName:                     x.MemberName,
			MemberType:                     x.MemberType,
			MemberStatus:                   x.MemberStatus,
			MemberExpiredDate:              x.MemberExpiredDate,
			CheckInTime:                    x.CheckInTime,
			CheckOutTime:                   x.CheckOutTime,
			DurationTime:                   x.DurationTime,
			VehicleNumberIn:                x.VehicleNumberIn,
			VehicleNumberOut:               x.VehicleNumberOut,
			LogTrans:                       x.LogTrans,
			MerchantKey:                    x.MerchantKey,
			QrText:                         x.QrText,
			QrA2P:                          x.QrA2P,
			QrTextPaymentOnline:            x.QrTextPaymentOnline,
			TrxInvoiceItem:                 tempTrxInvoiceItems,
			FlagSyncData:                   x.FlagSyncData,
			MemberData:                     tempMemberData,
			TrxAddInfo:                     convertedTrxAddInfo,
			FlagTrxFromCloud:               x.FlagTrxFromCloud,
			IsRsyncDataTrx:                 x.IsRsyncDataTrx,
			ExcludeSf:                      x.ExcludeSf,
			FlagCharge:                     x.FlagCharge,
			ChargeType:                     x.ChargeType,
			RequestAddTrxInvoiceDetailItem: tempRequestAddTrxInvoiceDetailItem,
			LastUpdatedAt:                  x.LastUpdatedAt,
		}
		trxLists = append(trxLists, tempTrxOutstanding)
	}

	log.Println("ResultTrxList:", utils.ToString(trxLists))
	for _, rows := range trxLists {

		flagSyncStatus := true
		dataStr := utils.ToString(rows)
		log.Println("ConfirmSyncTrxToCloud (IsConnected):", utils.IsConnected())
		if utils.IsConnected() {
			redisStatus := svc.Service.RedisClient.Publish(constans.CHANNEL_REDIS_PG_PARKING_CHECKOUT, dataStr)
			if redisStatus.Err() != nil || redisStatus.Val() == 0 {
				flagSyncStatus = false
			}
		} else {
			flagSyncStatus = false
		}

		if flagSyncStatus {

			filter := bson.M{
				"docNo": rows.DocNo,
			}

			updateSet := bson.M{
				"$set": bson.M{
					"flagSyncData": flagSyncStatus,
				},
			}

			err = svc.Service.TrxMongoRepo.UpdateTrxByInterface(filter, updateSet)
			if err != nil {
				result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
				return &trx.MyResponse{
					Response: result,
				}, err
			}

		}

	}

	result = helpers.ResponseJSON(true, constans.SUCCESS_CODE, constans.EMPTY_VALUE, nil)
	return &trx.MyResponse{
		Response: result,
	}, nil
}

func (svc trxService) InquiryPayment(ctx context.Context, input *trx.RequestInquiryPayment) (*trx.MyResponse, error) {
	var result *trx.Response

	if err := helpers.BindValidateStruct(input); err != nil {
		result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	resultTrx, exists, err := svc.Service.TrxMongoRepo.IsTrxOutstandingByDocNoForCustom(input.DocNo)
	if err != nil {
		result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}
	var tempTrxInvoiceItems []*trx.TrxInvoiceItem
	for _, v := range resultTrx.TrxInvoiceItem {
		tempTrxInvoiceItem := &trx.TrxInvoiceItem{
			DocNo:                  v.DocNo,
			ProductId:              v.ProductId,
			ProductCode:            v.ProductCode,
			ProductName:            v.ProductName,
			IsPctServiceFee:        v.IsPctServiceFee,
			ServiceFee:             v.ServiceFee,
			ServiceFeeMember:       v.ServiceFeeMember,
			Price:                  v.Price,
			BaseTime:               v.BaseTime,
			ProgressiveTime:        v.ProgressiveTime,
			ProgressivePrice:       v.ProgressivePrice,
			IsPct:                  v.IsPct,
			ProgressivePct:         v.ProgressivePct,
			MaxPrice:               v.MaxPrice,
			Is24H:                  v.Is24H,
			OvernightTime:          v.OvernightTime,
			OvernightPrice:         v.OvernightPrice,
			GracePeriod:            v.GracePeriod,
			FlgRepeat:              v.FlgRepeat,
			TotalAmount:            v.TotalAmount,
			TotalProgressiveAmount: v.TotalProgressiveAmount,
		}
		tempTrxInvoiceItems = append(tempTrxInvoiceItems, tempTrxInvoiceItem)
	}
	tempMemberData := &trx.TrxMember{
		DocNo:       resultTrx.MemberData.DocNo,
		PartnerCode: resultTrx.MemberData.PartnerCode,
		FirstName:   resultTrx.MemberData.FirstName,
		LastName:    resultTrx.MemberData.LastName,
		RoleType:    resultTrx.MemberData.RoleType,
		PhoneNumber: resultTrx.MemberData.PhoneNumber,
		Email:       resultTrx.MemberData.Email,
		Active:      resultTrx.MemberData.Active,
		ActiveAt:    resultTrx.MemberData.ActiveAt,
		NonActiveAt: func() *wrapperspb.StringValue {
			if resultTrx.MemberData.NonActiveAt != nil {
				return &wrapperspb.StringValue{Value: *resultTrx.MemberData.NonActiveAt}
			}
			return nil
		}(),
		OuId:               resultTrx.MemberData.OuId,
		TypePartner:        resultTrx.MemberData.TypePartner,
		CardNumber:         resultTrx.MemberData.CardNumber,
		VehicleNumber:      resultTrx.MemberData.VehicleNumber,
		RegisteredDatetime: resultTrx.MemberData.RegisteredDatetime,
		DateFrom:           resultTrx.MemberData.DateFrom,
		DateTo:             resultTrx.MemberData.DateTo,
		ProductId:          resultTrx.MemberData.ProductId,
		ProductCode:        resultTrx.MemberData.ProductCode,
	}

	convertedTrxAddInfo := make(map[string]*anypb.Any)
	for k, v := range resultTrx.TrxAddInfo {
		if msg, ok := v.(goproto.Message); ok {
			anyVal, err := anypb.New(msg)
			if err != nil {
				continue
			}
			convertedTrxAddInfo[k] = anyVal
		}
	}

	tempRequestAddTrxInvoiceDetailItem := &trx.TrxInvoiceDetailItem{
		DocNo:         resultTrx.RequestAddTrxInvoiceDetailItem.DocNo,
		ProductCode:   resultTrx.RequestAddTrxInvoiceDetailItem.ProductCode,
		InvoiceAmount: resultTrx.RequestAddTrxInvoiceDetailItem.InvoiceAmount,
		CreatedAt:     resultTrx.RequestAddTrxInvoiceDetailItem.CreatedAt,
		CreatedDate:   resultTrx.RequestAddTrxInvoiceDetailItem.CreatedDate,
	}

	tempTrxOutstanding := &trx.Trx{
		DocNo:                          resultTrx.DocNo,
		DocDate:                        resultTrx.DocDate,
		PaymentRefDocno:                resultTrx.PaymentRefDocNo,
		CheckinDateTime:                resultTrx.CheckInDatetime,
		CheckoutDateTime:               resultTrx.CheckOutDatetime,
		DeviceIdIn:                     resultTrx.DeviceIdIn,
		DeviceId:                       resultTrx.DeviceId,
		GetIn:                          resultTrx.GateIn,
		GetOut:                         resultTrx.GateOut,
		CardNumberUUIDIn:               resultTrx.CardNumberUUIDIn,
		CardNumberIn:                   resultTrx.CardNumberIn,
		CardNumberUUID:                 resultTrx.CardNumberUUID,
		CardNumber:                     resultTrx.CardNumber,
		TypeCard:                       resultTrx.TypeCard,
		BeginningBalance:               resultTrx.BeginningBalance,
		ExtLocalDateTime:               resultTrx.ExtLocalDatetime,
		ChargeAmount:                   resultTrx.ChargeAmount,
		GrandTotal:                     resultTrx.GrandTotal,
		ProductCode:                    resultTrx.ProductCode,
		ProductName:                    resultTrx.ProductName,
		ProductData:                    resultTrx.ProductData,
		RequestData:                    resultTrx.RequestData,
		RequestOutData:                 resultTrx.RequestOutData,
		OuId:                           resultTrx.OuId,
		OuName:                         resultTrx.OuName,
		OuCode:                         resultTrx.OuCode,
		OuSubBranchId:                  resultTrx.OuSubBranchId,
		USubBranchName:                 resultTrx.OuSubBranchName,
		OuSubBranchCode:                resultTrx.OuSubBranchCode,
		MainOuId:                       resultTrx.MainOuId,
		MainOuCode:                     resultTrx.MainOuCode,
		MainOuName:                     resultTrx.MainOuName,
		MemberCode:                     resultTrx.MemberCode,
		MemberName:                     resultTrx.MemberName,
		MemberType:                     resultTrx.MemberType,
		MemberStatus:                   resultTrx.MemberStatus,
		MemberExpiredDate:              resultTrx.MemberExpiredDate,
		CheckInTime:                    resultTrx.CheckInTime,
		CheckOutTime:                   resultTrx.CheckOutTime,
		DurationTime:                   resultTrx.DurationTime,
		VehicleNumberIn:                resultTrx.VehicleNumberIn,
		VehicleNumberOut:               resultTrx.VehicleNumberOut,
		LogTrans:                       resultTrx.LogTrans,
		MerchantKey:                    resultTrx.MerchantKey,
		QrText:                         resultTrx.QrText,
		QrA2P:                          resultTrx.QrA2P,
		QrTextPaymentOnline:            resultTrx.QrTextPaymentOnline,
		TrxInvoiceItem:                 tempTrxInvoiceItems,
		FlagSyncData:                   resultTrx.FlagSyncData,
		MemberData:                     tempMemberData,
		TrxAddInfo:                     convertedTrxAddInfo,
		FlagTrxFromCloud:               resultTrx.FlagTrxFromCloud,
		IsRsyncDataTrx:                 resultTrx.IsRsyncDataTrx,
		ExcludeSf:                      resultTrx.ExcludeSf,
		FlagCharge:                     resultTrx.FlagCharge,
		ChargeType:                     resultTrx.ChargeType,
		RequestAddTrxInvoiceDetailItem: tempRequestAddTrxInvoiceDetailItem,
		LastUpdatedAt:                  resultTrx.LastUpdatedAt,
	}

	if !exists {
		result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, "Sesi Tidak Ditemukan", nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	requestInquiryQRIS := make(map[string]interface{})
	requestInquiryQRIS["docNo"] = tempTrxOutstanding.DocNo
	requestInquiryQRIS["paymentMethod"] = input.PaymentMethod
	requestInquiryQRIS["productCode"] = input.ProductCode
	requestInquiryQRIS["productName"] = input.ProductName
	requestInquiryQRIS["grandTotal"] = input.GrandTotal
	requestInquiryQRIS["mKey"] = config.MERCHANT_KEY_APPS2PAY

	body, err := json.Marshal(requestInquiryQRIS)
	if err != nil {
		result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	responseInquiryPaymentStr, err := helpers.CallHttpCloudServerParking("POST", body, "/mpos/local/inquiry-payment")
	if err != nil {
		result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	var responseInquiryPayment *trx.ResponseInquiryPayment
	err = json.Unmarshal([]byte(*responseInquiryPaymentStr), &responseInquiryPayment)
	if err != nil {
		result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	if !responseInquiryPayment.Success {
		result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, responseInquiryPayment.Message, nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	response := make(map[string]interface{})
	response["qrCode"] = responseInquiryPayment.Result.QrCode
	response["type"] = responseInquiryPayment.Result.Type
	response["mKey"] = config.MERCHANT_KEY_APPS2PAY
	response["paymentRefDocNo"] = responseInquiryPayment.Result.PaymentRefDocNo

	s, err := structpb.NewStruct(response)

	anyResponseTrx, _ := anypb.New(s)

	result = helpers.ResponseJSON(true, constans.SUCCESS_CODE, constans.EMPTY_VALUE, anyResponseTrx)
	return &trx.MyResponse{
		Response: result,
	}, nil
}

func (svc trxService) InquiryWithCardP3(ctx context.Context, input *trx.RequestInquiryWithCardP3) (*trx.MyResponse, error) {
	var result *trx.Response
	var resultProducts *prod.PolicyOuProductWithRules
	var resultProduct models.PolicyOuProductWithRules
	resultInquiryTrxWithCard := make(map[string]interface{})
	resultInquiryTrxWithCard["memberCode"] = constans.EMPTY_VALUE
	resultInquiryTrxWithCard["memberName"] = constans.EMPTY_VALUE
	resultInquiryTrxWithCard["memberType"] = constans.EMPTY_VALUE

	if err := helpers.BindValidateStruct(input); err != nil {
		result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	resultRedis := svc.Service.RedisClientLocal.Get("P3")
	if resultRedis.Err() != nil || resultRedis.Val() == constans.EMPTY_VALUE || resultRedis.Val() == "0" {
		result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, "Maaf untuk saat ini sistem P3 sudah tidak aktif,Keterangan lebih lanjut silahkan untuk menghubungi administrator", nil)
		return &trx.MyResponse{
			Response: result,
		}, resultRedis.Err()
	}

	merchantKey, err := utils.DecryptMerchantKey(config.MERCHANT_KEY)
	if err != nil {
		result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}
	tempMerchantKey := &trx.MerchantKey{
		ID:              merchantKey.ID,
		OuId:            merchantKey.OuId,
		OuName:          merchantKey.OuName,
		OuCode:          merchantKey.OuCode,
		OuSubBranchId:   merchantKey.OuSubBranchId,
		OuSubBranchName: merchantKey.OuSubBranchName,
		OuSubBranchCode: merchantKey.OuSubBranchCode,
		MainOuId:        merchantKey.MainOuId,
		MainOuCode:      merchantKey.MainOuCode,
		MainOuName:      merchantKey.MainOuName,
	}

	resultProduct, err = svc.Service.ProductRepo.FindProductByProductCode(input.ProductCode, tempMerchantKey.OuId)
	if err != nil {
		result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}
	tempProd := &prod.ProductOuWithRules{
		OuId:             resultProduct.ProductOuWithRules.OuId,
		ProductId:        resultProduct.ProductOuWithRules.ProductId,
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
	resultProducts = &prod.PolicyOuProductWithRules{
		OuId:                  resultProduct.OuId,
		OuCode:                resultProduct.OuCode,
		OuName:                resultProduct.OuName,
		ProductId:             resultProduct.ProductId,
		ProductCode:           resultProduct.ProductCode,
		ProductName:           resultProduct.ProductName,
		ServiceFee:            resultProduct.ServiceFee,
		IsPctServiceFee:       resultProduct.IsPctServiceFee,
		IsPctServiceFeeMember: resultProduct.IsPctServiceFeeMember,
		ServiceFeeMember:      resultProduct.ServiceFeeMember,
		ProductRules:          tempProd,
	}

	checkinDatetimeParse, _ := time.Parse("2006-01-02 15:04:05", input.CheckInDateTime)
	checkoutDatetimeParse, _ := time.Parse("2006-01-02 15:04:05", input.CheckInDateTime)

	yearMonth := checkinDatetimeParse.Format("060102")
	HourCheckIn := checkinDatetimeParse.Format("15")
	prefixDocNo := utils.RandStringBytesMaskImprSrcChr(4)
	prefix := fmt.Sprintf("%s%s", yearMonth, HourCheckIn)
	autoNumber, err := svc.Service.GenAutoNumRepo.AutonumberValueWithDatatype(constans.DATATYPE_TRX_LOCAL, prefix, 4)
	if err != nil {
		result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	docNo := fmt.Sprintf("%s%s%d", prefixDocNo, autoNumber, tempMerchantKey.OuId)
	randomKey := utils.RandStringBytesMaskImprSrc(16)
	encQrTxt, _ := utils.Encrypt(docNo, randomKey)
	qrCode := fmt.Sprintf("%s%s", encQrTxt.Result, randomKey)

	inquiryDate := input.CheckInDateTime[:len(input.CheckInDateTime)-9]
	resultMember, exists, err := svc.Service.MemberRepo.IsMemberByAdvanceIndex(input.UuidCard, constans.EMPTY_VALUE, inquiryDate, config.MEMBER_BY, false)
	tempMemberData := &trx.Member{
		Id:          resultMember.ID,
		PartnerCode: resultMember.PartnerCode,
		FirstName:   resultMember.FirstName,
		LastName:    resultMember.LastName,
		RoleType:    resultMember.RoleType,
		PhoneNumber: resultMember.PhoneNumber,
		Email:       resultMember.Email,
		Active:      resultMember.Active,
		ActiveAt:    resultMember.ActiveAt,
		NonActiveAt: func() *wrapperspb.StringValue {
			if resultMember.NonActiveAt != nil {
				return &wrapperspb.StringValue{Value: *resultMember.NonActiveAt}
			}
			return nil
		}(),
		OuId:                resultMember.OuId,
		TypePartner:         resultMember.TypePartner,
		CardNumber:          resultMember.CardNumber,
		VehicleNumber:       resultMember.VehicleNumber,
		RegisteredDatetime:  resultMember.RegisteredDatetime,
		DateFrom:            resultMember.DateFrom,
		DateTo:              resultMember.DateTo,
		ProductId:           resultMember.ProductId,
		ProductCode:         resultMember.ProductCode,
		ProductMembershipId: resultMember.ProductMembershipId,
		Price:               resultMember.Price,
		ServiceFee:          resultMember.ServiceFee,
		IsPctSfee:           resultMember.IsPctSfee,
		DueDate:             resultMember.DueDate,
		DiscType:            resultMember.DiscType,
		DiscAmount:          resultMember.DiscAmount,
		DiscPct:             resultMember.DiscPct,
		GracePeriodDate:     resultMember.GracePeriodDate,
		Username:            resultMember.Username,
		IsExtendMember:      resultMember.IsExtendMember,
		CreatedAt:           resultMember.CreatedAt,
		CreatedBy:           resultMember.CreatedBy,
		UpdatedAt:           resultMember.UpdatedAt,
		UpdatedBy:           resultMember.UpdatedBy,
	}
	if err != nil {
		result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	log.Println(tempMemberData)
	isFreePass := false
	isSpecialMember := false
	if exists {
		if tempMemberData.TypePartner == constans.TYPE_PARTNER_FREE_PASS {
			isFreePass = true
		} else if tempMemberData.TypePartner == constans.TYPE_PARTNER_SPECIAL_MEMBER {
			isSpecialMember = true
		}
	}

	var totalAmount float64
	if !isFreePass {
		log.Println("DUDU123")
		productData, existsProduct := svc.Service.TrxMongoRepo.IsTrxProductCustomExistsByKeyword(input.UuidCard)
		tempProdDat := &trx.TrxProductCustom{
			Keyword:     productData.Keyword,
			ProductName: productData.ProductName,
			ProductCode: productData.ProductCode,
		}
		if existsProduct {
			input.ProductCode = tempProdDat.ProductCode
		}

		resultProduct, err = svc.Service.ProductRepo.FindProductByProductCode(input.ProductCode, tempMerchantKey.OuId)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, err.Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}

		tempProd = &prod.ProductOuWithRules{
			OuId:             resultProduct.ProductOuWithRules.OuId,
			ProductId:        resultProduct.ProductOuWithRules.ProductId,
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
		resultProducts = &prod.PolicyOuProductWithRules{
			OuId:                  resultProduct.OuId,
			OuCode:                resultProduct.OuCode,
			OuName:                resultProduct.OuName,
			ProductId:             resultProduct.ProductId,
			ProductCode:           resultProduct.ProductCode,
			ProductName:           resultProduct.ProductName,
			ServiceFee:            resultProduct.ServiceFee,
			IsPctServiceFee:       resultProduct.IsPctServiceFee,
			IsPctServiceFeeMember: resultProduct.IsPctServiceFeeMember,
			ServiceFeeMember:      resultProduct.ServiceFeeMember,
			ProductRules:          tempProd,
		}

		if isSpecialMember {
			bodySplit := strings.Split(input.ProductCode, "SPM")
			if len(bodySplit) > 1 {
				if input.ProductCode == tempMemberData.ProductCode {
					resultProduct, err = svc.Service.ProductRepo.FindProductByProductCode(input.ProductCode, tempMerchantKey.OuId)
					if err != nil {
						result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, err.Error(), nil)
						return &trx.MyResponse{
							Response: result,
						}, err
					}
					tempProd = &prod.ProductOuWithRules{
						OuId:             resultProduct.ProductOuWithRules.OuId,
						ProductId:        resultProduct.ProductOuWithRules.ProductId,
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
					resultProducts = &prod.PolicyOuProductWithRules{
						OuId:                  resultProduct.OuId,
						OuCode:                resultProduct.OuCode,
						OuName:                resultProduct.OuName,
						ProductId:             resultProduct.ProductId,
						ProductCode:           resultProduct.ProductCode,
						ProductName:           resultProduct.ProductName,
						ServiceFee:            resultProduct.ServiceFee,
						IsPctServiceFee:       resultProduct.IsPctServiceFee,
						IsPctServiceFeeMember: resultProduct.IsPctServiceFeeMember,
						ServiceFeeMember:      resultProduct.ServiceFeeMember,
						ProductRules:          tempProd,
					}

				}
			} else {
				input.ProductCode = fmt.Sprintf("%s-%s", "SPM", input.ProductCode)
				if input.ProductCode == tempMemberData.ProductCode {
					resultProduct, err = svc.Service.ProductRepo.FindProductByProductCode(input.ProductCode, tempMerchantKey.OuId)
					if err != nil {
						result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, err.Error(), nil)
						return &trx.MyResponse{
							Response: result,
						}, err
					}
					tempProd = &prod.ProductOuWithRules{
						OuId:             resultProduct.ProductOuWithRules.OuId,
						ProductId:        resultProduct.ProductOuWithRules.ProductId,
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
					resultProducts = &prod.PolicyOuProductWithRules{
						OuId:                  resultProduct.OuId,
						OuCode:                resultProduct.OuCode,
						OuName:                resultProduct.OuName,
						ProductId:             resultProduct.ProductId,
						ProductCode:           resultProduct.ProductCode,
						ProductName:           resultProduct.ProductName,
						ServiceFee:            resultProduct.ServiceFee,
						IsPctServiceFee:       resultProduct.IsPctServiceFee,
						IsPctServiceFeeMember: resultProduct.IsPctServiceFeeMember,
						ServiceFeeMember:      resultProduct.ServiceFeeMember,
						ProductRules:          tempProd,
					}
				}
			}

		}

		totalAmount = float64(0)
		grandTotal := resultProducts.ProductRules.Price
		checkinDate := input.CheckInDateTime[:len(input.CheckInDateTime)-9]

		svc.Service.RedisClientLocal.Del(fmt.Sprintf("%s-%s", docNo, constans.MEMBER))
		var resultTrxMemberList []*trx.TrxMember
		if exists && !isSpecialMember {
			resultActiveMemberList, err := svc.Service.MemberRepo.GetMemberActiveListByPeriod(input.UuidCard, constans.EMPTY_VALUE, checkinDate, inquiryDate, config.MEMBER_BY, false)
			if err != nil {
				result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, "Sesi Tidak Ditemukan", nil)
				return &trx.MyResponse{
					Response: result,
				}, err
			}

			for _, rows := range resultActiveMemberList {
				if rows.ProductCode == input.ProductCode {
					trxMemberList := &trx.TrxMember{
						DocNo:       docNo,
						PartnerCode: rows.PartnerCode,
						FirstName:   rows.FirstName,
						LastName:    rows.LastName,
						RoleType:    rows.RoleType,
						PhoneNumber: rows.PhoneNumber,
						Email:       rows.Email,
						Active:      rows.Active,
						ActiveAt:    rows.ActiveAt,
						NonActiveAt: func() *wrapperspb.StringValue {
							if rows.NonActiveAt != nil {
								return &wrapperspb.StringValue{Value: *rows.NonActiveAt}
							}
							return nil
						}(),
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
		}

		if len(resultTrxMemberList) > 0 {
			memberData := resultTrxMemberList[len(resultTrxMemberList)-1]

			grandTotal = 0

			log.Println("Grand Total After Calc Invoice Amount - Total Progressive Amount:", grandTotal)
			if grandTotal == 0 {
				trxMemberStr, _ := json.Marshal(memberData)
				svc.Service.RedisClientLocal.Set(fmt.Sprintf("%s-%s", docNo, constans.MEMBER), trxMemberStr, 5*time.Minute)

				resultInquiryTrxWithCard["memberCode"] = memberData.PartnerCode
				resultInquiryTrxWithCard["memberName"] = strings.TrimSpace(fmt.Sprintf("%s %s", memberData.FirstName, memberData.LastName))
				resultInquiryTrxWithCard["memberType"] = memberData.TypePartner
			}
		} else if isSpecialMember && input.ProductCode == resultMember.ProductCode {
			trxMemberList := models.TrxMember{
				DocNo:              docNo,
				PartnerCode:        resultMember.PartnerCode,
				FirstName:          resultMember.FirstName,
				LastName:           resultMember.LastName,
				RoleType:           resultMember.RoleType,
				PhoneNumber:        resultMember.PhoneNumber,
				Email:              resultMember.Email,
				Active:             resultMember.Active,
				ActiveAt:           resultMember.ActiveAt,
				NonActiveAt:        resultMember.NonActiveAt,
				OuId:               resultMember.OuId,
				TypePartner:        resultMember.TypePartner,
				CardNumber:         resultMember.CardNumber,
				VehicleNumber:      resultMember.VehicleNumber,
				RegisteredDatetime: resultMember.RegisteredDatetime,
				DateFrom:           resultMember.DateFrom,
				DateTo:             resultMember.DateTo,
				ProductId:          resultMember.ProductId,
				ProductCode:        resultMember.ProductCode,
			}

			trxMemberStr, _ := json.Marshal(trxMemberList)
			svc.Service.RedisClientLocal.Set(fmt.Sprintf("%s-%s", docNo, constans.MEMBER), trxMemberStr, 5*time.Minute)

			resultInquiryTrxWithCard["memberCode"] = resultMember.PartnerCode
			resultInquiryTrxWithCard["memberName"] = strings.TrimSpace(fmt.Sprintf("%s %s", resultMember.FirstName, resultMember.LastName))
			resultInquiryTrxWithCard["memberType"] = resultMember.TypePartner
		}

		totalAmount = grandTotal
		resultInquiryTrxWithCard["totalAmount"] = totalAmount
	} else {
		log.Println("MEMBER : ", utils.ToString(resultMember))
		totalAmount = constans.EMPTY_VALUE_INT
		resultInquiryTrxWithCard["totalAmount"] = totalAmount
		resultInquiryTrxWithCard["memberCode"] = resultMember.PartnerCode
		resultInquiryTrxWithCard["memberName"] = strings.TrimSpace(fmt.Sprintf("%s %s", resultMember.FirstName, resultMember.LastName))
		resultInquiryTrxWithCard["memberType"] = resultMember.TypePartner
		log.Println("BERHASIL ISI MEMBER")
	}

	log.Println("PRODUCT : ", utils.ToString(resultProducts))
	var trxInvoiceItemList []models.TrxInvoiceItem
	trxInvoiceItem := models.TrxInvoiceItem{
		DocNo:                  docNo,
		ProductId:              resultProducts.ProductId,
		ProductCode:            resultProducts.ProductCode,
		ProductName:            resultProducts.ProductName,
		IsPctServiceFee:        resultProducts.IsPctServiceFee,
		ServiceFee:             resultProducts.ServiceFee,
		ServiceFeeMember:       resultProducts.ServiceFeeMember,
		Price:                  resultProducts.ProductRules.Price,
		BaseTime:               resultProducts.ProductRules.BaseTime,
		ProgressiveTime:        resultProducts.ProductRules.ProgressiveTime,
		ProgressivePrice:       resultProducts.ProductRules.ProgressivePrice,
		IsPct:                  resultProducts.ProductRules.IsPct,
		ProgressivePct:         resultProducts.ProductRules.ProgressivePct,
		MaxPrice:               resultProducts.ProductRules.MaxPrice,
		Is24H:                  resultProducts.ProductRules.Is24H,
		OvernightTime:          resultProducts.ProductRules.OvernightTime,
		OvernightPrice:         resultProducts.ProductRules.OvernightPrice,
		GracePeriod:            resultProducts.ProductRules.GracePeriod,
		FlgRepeat:              resultProducts.ProductRules.FlgRepeat,
		TotalAmount:            resultProducts.ProductRules.Price,
		TotalProgressiveAmount: resultProducts.ProductRules.Price,
	}
	trxInvoiceItemList = append(trxInvoiceItemList, trxInvoiceItem)

	log.Println("BERHASIL ISI INVOICE")

	trxm := &models.Trx{
		DocNo:            docNo,
		DocDate:          constans.EMPTY_VALUE,
		CheckInDatetime:  input.CheckInDateTime,
		CheckOutDatetime: input.CheckInDateTime,
		DeviceIdIn:       input.DeviceId,
		DeviceId:         input.DeviceId,
		GateIn:           input.TerminalId,
		GateOut:          input.TerminalId,
		CardNumberUUIDIn: input.UuidCard,
		CardNumberIn:     input.CardNumber,
		CardNumberUUID:   input.UuidCard,
		CardNumber:       input.CardNumber,
		TypeCard:         input.TypeCard,
		BeginningBalance: input.BeginningBalance,
		ExtLocalDatetime: input.CheckInDateTime,
		GrandTotal:       totalAmount,
		ProductCode:      resultProducts.ProductCode,
		ProductName:      resultProducts.ProductName,
		ProductData:      utils.ToString(resultProducts),
		RequestData:      utils.ToString(input),
		RequestOutData:   utils.ToString(input),
		OuId:             tempMerchantKey.OuId,
		OuName:           tempMerchantKey.OuName,
		OuCode:           tempMerchantKey.OuCode,
		OuSubBranchId:    tempMerchantKey.OuSubBranchId,
		OuSubBranchName:  tempMerchantKey.OuSubBranchName,
		OuSubBranchCode:  tempMerchantKey.OuSubBranchCode,
		MainOuId:         tempMerchantKey.MainOuId,
		MainOuCode:       tempMerchantKey.MainOuCode,
		MainOuName:       tempMerchantKey.MainOuName,
		MemberCode:       resultInquiryTrxWithCard["memberCode"].(string),
		MemberName:       resultInquiryTrxWithCard["memberName"].(string),
		MemberType:       resultInquiryTrxWithCard["memberType"].(string),
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

	idTrx, err := svc.Service.TrxMongoRepo.AddTrxCheckIn(*trxm)
	if err != nil {
		log.Println("ERROR AddTrxCheckIn :  ", err)
		result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	duration := utils.ConvDiffTime(checkinDatetimeParse, checkoutDatetimeParse)
	anyDuration, _ := anypb.New(duration)
	responseTrx := &trx.ResultInquiryTrxWithCard{
		Id:              idTrx.Hex(),
		DocNo:           docNo,
		CardNumberUUID:  input.UuidCard,
		CardNumber:      input.CardNumber,
		Nominal:         resultInquiryTrxWithCard["totalAmount"].(float64),
		ProductCode:     resultProduct.ProductCode,
		ProductName:     resultProduct.ProductName,
		VehicleNumberIn: constans.EMPTY_VALUE,
		MemberCode:      resultInquiryTrxWithCard["memberCode"].(string),
		MemberName:      resultInquiryTrxWithCard["memberName"].(string),
		MemberType:      resultInquiryTrxWithCard["memberType"].(string),
		QrCode:          constans.EMPTY_VALUE,
		ExcludeSf:       false,
		Duration:        anyDuration,
		OuCode:          merchantKey.OuCode,
	}
	anyResponseTrx, _ := anypb.New(responseTrx)

	result = helpers.ResponseJSON(true, constans.SUCCESS_CODE, constans.EMPTY_VALUE, anyResponseTrx)
	return &trx.MyResponse{
		Response: result,
	}, nil
}

func (svc trxService) InquiryPaymentP3(ctx context.Context, input *trx.RequestInquiryPaymentP3) (*trx.MyResponse, error) {
	var result *trx.Response
	//var responseTrx models.ResultInquiryTrx
	requestInquiry := make(map[string]interface{})
	channelCallback := make(map[string]interface{})

	if err := helpers.BindValidateStruct(input); err != nil {
		result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	merchantKey, err := utils.DecryptMerchantKey(config.MERCHANT_KEY)
	if err != nil {
		result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	tempMerchantKey := &trx.MerchantKey{
		ID:              merchantKey.ID,
		OuId:            merchantKey.OuId,
		OuName:          merchantKey.OuName,
		OuCode:          merchantKey.OuCode,
		OuSubBranchId:   merchantKey.OuSubBranchId,
		OuSubBranchName: merchantKey.OuSubBranchName,
		OuSubBranchCode: merchantKey.OuSubBranchCode,
		MainOuId:        merchantKey.MainOuId,
		MainOuCode:      merchantKey.MainOuCode,
		MainOuName:      merchantKey.MainOuName,
	}

	resultProduct, err := svc.Service.ProductRepo.FindProductByProductCode(input.ProductCode, tempMerchantKey.OuId)
	if err != nil {
		result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}
	tempProd := &prod.ProductOuWithRules{
		OuId:             resultProduct.ProductOuWithRules.OuId,
		ProductId:        resultProduct.ProductOuWithRules.ProductId,
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
	resultProducts := &prod.PolicyOuProductWithRules{
		OuId:                  resultProduct.OuId,
		OuCode:                resultProduct.OuCode,
		OuName:                resultProduct.OuName,
		ProductId:             resultProduct.ProductId,
		ProductCode:           resultProduct.ProductCode,
		ProductName:           resultProduct.ProductName,
		ServiceFee:            resultProduct.ServiceFee,
		IsPctServiceFee:       resultProduct.IsPctServiceFee,
		IsPctServiceFeeMember: resultProduct.IsPctServiceFeeMember,
		ServiceFeeMember:      resultProduct.ServiceFeeMember,
		ProductRules:          tempProd,
	}

	checkinDatetimeParse, _ := time.Parse("2006-01-02 15:04:05", input.InquiryDatetime)
	checkoutDatetimeParse, _ := time.Parse("2006-01-02 15:04:05", input.InquiryDatetime)

	yearMonth := checkinDatetimeParse.Format("060102")
	HourCheckIn := checkinDatetimeParse.Format("15")
	prefixDocNo := utils.RandStringBytesMaskImprSrcChr(4)
	prefix := fmt.Sprintf("%s%s", yearMonth, HourCheckIn)
	autoNumber, err := svc.Service.GenAutoNumRepo.AutonumberValueWithDatatype(constans.DATATYPE_TRX_LOCAL, prefix, 4)
	if err != nil {
		result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	docNo := fmt.Sprintf("%s%s%d", prefixDocNo, autoNumber, tempMerchantKey.OuId)
	randomKey := utils.RandStringBytesMaskImprSrc(16)
	encQrTxt, _ := utils.Encrypt(docNo, randomKey)
	qrCode := fmt.Sprintf("%s%s", encQrTxt.Result, randomKey)

	duration := utils.ConvDiffTime(checkinDatetimeParse, checkoutDatetimeParse)
	keyInquiryPayment := fmt.Sprintf("%s", "PAYMENT-INQUIRY-QRIS-PASS")
	checkExpInquiryPayment := svc.Service.RedisClientLocal.Get(keyInquiryPayment)
	log.Println("Check:", checkExpInquiryPayment)

	if checkExpInquiryPayment.Val() != constans.EMPTY_VALUE {
		redisResult := strings.Split(checkExpInquiryPayment.Val(), "#")

		resultTrxPayment, exists := svc.Service.TrxMongoRepo.IsTrxOutstandingByDocNoNew(redisResult[0])
		//resultTrxOnlineOutstanding = resultTrxPayment
		var tempTrxInvoiceItems []*trx.TrxInvoiceItem
		for _, v := range resultTrxPayment.TrxInvoiceItem {
			tempTrxInvoiceItem := &trx.TrxInvoiceItem{
				DocNo:                  v.DocNo,
				ProductId:              v.ProductId,
				ProductCode:            v.ProductCode,
				ProductName:            v.ProductName,
				IsPctServiceFee:        v.IsPctServiceFee,
				ServiceFee:             v.ServiceFee,
				ServiceFeeMember:       v.ServiceFeeMember,
				Price:                  v.Price,
				BaseTime:               v.BaseTime,
				ProgressiveTime:        v.ProgressiveTime,
				ProgressivePrice:       v.ProgressivePrice,
				IsPct:                  v.IsPct,
				ProgressivePct:         v.ProgressivePct,
				MaxPrice:               v.MaxPrice,
				Is24H:                  v.Is24H,
				OvernightTime:          v.OvernightTime,
				OvernightPrice:         v.OvernightPrice,
				GracePeriod:            v.GracePeriod,
				FlgRepeat:              v.FlgRepeat,
				TotalAmount:            v.TotalAmount,
				TotalProgressiveAmount: v.TotalProgressiveAmount,
			}
			tempTrxInvoiceItems = append(tempTrxInvoiceItems, tempTrxInvoiceItem)
		}
		tempMemberData := &trx.TrxMember{
			DocNo:       resultTrxPayment.MemberData.DocNo,
			PartnerCode: resultTrxPayment.MemberData.PartnerCode,
			FirstName:   resultTrxPayment.MemberData.FirstName,
			LastName:    resultTrxPayment.MemberData.LastName,
			RoleType:    resultTrxPayment.MemberData.RoleType,
			PhoneNumber: resultTrxPayment.MemberData.PhoneNumber,
			Email:       resultTrxPayment.MemberData.Email,
			Active:      resultTrxPayment.MemberData.Active,
			ActiveAt:    resultTrxPayment.MemberData.ActiveAt,
			NonActiveAt: func() *wrapperspb.StringValue {
				if resultTrxPayment.MemberData.NonActiveAt != nil {
					return &wrapperspb.StringValue{Value: *resultTrxPayment.MemberData.NonActiveAt}
				}
				return nil
			}(),
			OuId:               resultTrxPayment.MemberData.OuId,
			TypePartner:        resultTrxPayment.MemberData.TypePartner,
			CardNumber:         resultTrxPayment.MemberData.CardNumber,
			VehicleNumber:      resultTrxPayment.MemberData.VehicleNumber,
			RegisteredDatetime: resultTrxPayment.MemberData.RegisteredDatetime,
			DateFrom:           resultTrxPayment.MemberData.DateFrom,
			DateTo:             resultTrxPayment.MemberData.DateTo,
			ProductId:          resultTrxPayment.MemberData.ProductId,
			ProductCode:        resultTrxPayment.MemberData.ProductCode,
		}

		convertedTrxAddInfo := make(map[string]*anypb.Any)
		for k, v := range resultTrxPayment.TrxAddInfo {
			if msg, ok := v.(goproto.Message); ok {
				anyVal, err := anypb.New(msg)
				if err != nil {
					continue
				}
				convertedTrxAddInfo[k] = anyVal
			}
		}

		tempRequestAddTrxInvoiceDetailItem := &trx.TrxInvoiceDetailItem{
			DocNo:         resultTrxPayment.RequestAddTrxInvoiceDetailItem.DocNo,
			ProductCode:   resultTrxPayment.RequestAddTrxInvoiceDetailItem.ProductCode,
			InvoiceAmount: resultTrxPayment.RequestAddTrxInvoiceDetailItem.InvoiceAmount,
			CreatedAt:     resultTrxPayment.RequestAddTrxInvoiceDetailItem.CreatedAt,
			CreatedDate:   resultTrxPayment.RequestAddTrxInvoiceDetailItem.CreatedDate,
		}

		tempTrxOutstanding := &trx.Trx{
			DocNo:                          resultTrxPayment.DocNo,
			DocDate:                        resultTrxPayment.DocDate,
			PaymentRefDocno:                resultTrxPayment.PaymentRefDocNo,
			CheckinDateTime:                resultTrxPayment.CheckInDatetime,
			CheckoutDateTime:               resultTrxPayment.CheckOutDatetime,
			DeviceIdIn:                     resultTrxPayment.DeviceIdIn,
			DeviceId:                       resultTrxPayment.DeviceId,
			GetIn:                          resultTrxPayment.GateIn,
			GetOut:                         resultTrxPayment.GateOut,
			CardNumberUUIDIn:               resultTrxPayment.CardNumberUUIDIn,
			CardNumberIn:                   resultTrxPayment.CardNumberIn,
			CardNumberUUID:                 resultTrxPayment.CardNumberUUID,
			CardNumber:                     resultTrxPayment.CardNumber,
			TypeCard:                       resultTrxPayment.TypeCard,
			BeginningBalance:               resultTrxPayment.BeginningBalance,
			ExtLocalDateTime:               resultTrxPayment.ExtLocalDatetime,
			ChargeAmount:                   resultTrxPayment.ChargeAmount,
			GrandTotal:                     resultTrxPayment.GrandTotal,
			ProductCode:                    resultTrxPayment.ProductCode,
			ProductName:                    resultTrxPayment.ProductName,
			ProductData:                    resultTrxPayment.ProductData,
			RequestData:                    resultTrxPayment.RequestData,
			RequestOutData:                 resultTrxPayment.RequestOutData,
			OuId:                           resultTrxPayment.OuId,
			OuName:                         resultTrxPayment.OuName,
			OuCode:                         resultTrxPayment.OuCode,
			OuSubBranchId:                  resultTrxPayment.OuSubBranchId,
			USubBranchName:                 resultTrxPayment.OuSubBranchName,
			OuSubBranchCode:                resultTrxPayment.OuSubBranchCode,
			MainOuId:                       resultTrxPayment.MainOuId,
			MainOuCode:                     resultTrxPayment.MainOuCode,
			MainOuName:                     resultTrxPayment.MainOuName,
			MemberCode:                     resultTrxPayment.MemberCode,
			MemberName:                     resultTrxPayment.MemberName,
			MemberType:                     resultTrxPayment.MemberType,
			MemberStatus:                   resultTrxPayment.MemberStatus,
			MemberExpiredDate:              resultTrxPayment.MemberExpiredDate,
			CheckInTime:                    resultTrxPayment.CheckInTime,
			CheckOutTime:                   resultTrxPayment.CheckOutTime,
			DurationTime:                   resultTrxPayment.DurationTime,
			VehicleNumberIn:                resultTrxPayment.VehicleNumberIn,
			VehicleNumberOut:               resultTrxPayment.VehicleNumberOut,
			LogTrans:                       resultTrxPayment.LogTrans,
			MerchantKey:                    resultTrxPayment.MerchantKey,
			QrText:                         resultTrxPayment.QrText,
			QrA2P:                          resultTrxPayment.QrA2P,
			QrTextPaymentOnline:            resultTrxPayment.QrTextPaymentOnline,
			TrxInvoiceItem:                 tempTrxInvoiceItems,
			FlagSyncData:                   resultTrxPayment.FlagSyncData,
			MemberData:                     tempMemberData,
			TrxAddInfo:                     convertedTrxAddInfo,
			FlagTrxFromCloud:               resultTrxPayment.FlagTrxFromCloud,
			IsRsyncDataTrx:                 resultTrxPayment.IsRsyncDataTrx,
			ExcludeSf:                      resultTrxPayment.ExcludeSf,
			FlagCharge:                     resultTrxPayment.FlagCharge,
			ChargeType:                     resultTrxPayment.ChargeType,
			RequestAddTrxInvoiceDetailItem: tempRequestAddTrxInvoiceDetailItem,
			LastUpdatedAt:                  resultTrxPayment.LastUpdatedAt,
		}
		if exists {
			log.Println("TRX OUTSTANDING : ", utils.ToString(tempTrxOutstanding))
			resultOutstanding, err := svc.Service.TrxMongoRepo.FindTrxOutstandingByDocNoCustom(tempTrxOutstanding.DocNo)
			if err != nil {
				result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
				return &trx.MyResponse{
					Response: result,
				}, err
			}

			tempTrxInvoiceItems := []*trx.TrxInvoiceItem{}
			for _, v := range resultOutstanding.TrxInvoiceItem {
				tempTrxInvoiceItem := &trx.TrxInvoiceItem{
					DocNo:                  v.DocNo,
					ProductId:              v.ProductId,
					ProductCode:            v.ProductCode,
					ProductName:            v.ProductName,
					IsPctServiceFee:        v.IsPctServiceFee,
					ServiceFee:             v.ServiceFee,
					ServiceFeeMember:       v.ServiceFeeMember,
					Price:                  v.Price,
					BaseTime:               v.BaseTime,
					ProgressiveTime:        v.ProgressiveTime,
					ProgressivePrice:       v.ProgressivePrice,
					IsPct:                  v.IsPct,
					ProgressivePct:         v.ProgressivePct,
					MaxPrice:               v.MaxPrice,
					Is24H:                  v.Is24H,
					OvernightTime:          v.OvernightTime,
					OvernightPrice:         v.OvernightPrice,
					GracePeriod:            v.GracePeriod,
					FlgRepeat:              v.FlgRepeat,
					TotalAmount:            v.TotalAmount,
					TotalProgressiveAmount: v.TotalProgressiveAmount,
				}
				tempTrxInvoiceItems = append(tempTrxInvoiceItems, tempTrxInvoiceItem)
			}

			tempResultTrx := &trx.ResultFindTrxOutstanding{
				Id:              resultOutstanding.ID.Hex(),
				DocNo:           resultOutstanding.DocNo,
				GrandTotal:      resultOutstanding.GrandTotal,
				CheckInDatetime: resultOutstanding.CheckInDatetime,
				OverNightPrice:  resultOutstanding.OverNightPrice,
				Is24H:           resultOutstanding.Is24H,
				CardNumber:      resultOutstanding.CardNumber,
				CardNumberUUID:  resultOutstanding.CardNumberUUID,
				OuCode:          resultOutstanding.OuCode,
				VehicleNumberIn: resultOutstanding.VehicleNumberIn,
				TrxInvoiceItem:  tempTrxInvoiceItems,
			}

			requestA2PConfirm := models.RequestConfirmA2P{
				Merchantkey:     config.MERCHANT_KEY_APPS2PAY,
				MerchantNoRef:   redisResult[1],
				PaymentCategory: constans.PAYMENT_METHOD_QRIS,
			}

			resultConfirmA2P, err := helperService.WorkerExtCloudPayment(requestA2PConfirm, fmt.Sprintf("%s%s", config.UrlPaymentA2P, "/public/payment/confirmation"), "bWtwbW9iaWxlOm1rcG1vYmlsZTEyMw==", constans.EMPTY_VALUE)
			if err != nil {
				result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
				return &trx.MyResponse{
					Response: result,
				}, err
			}

			var responseA2PConfirm models.ResponseA2PConfirmation
			json.Unmarshal(resultConfirmA2P, &responseA2PConfirm)

			if responseA2PConfirm.Result.StatusCode == constans.DOCUMENT_SUCCESS_CODE && responseA2PConfirm.Result.StatusPayment == constans.SUCCESS {

				requestInquiry["_id"] = tempResultTrx.Id
				requestInquiry["docNo"] = tempTrxOutstanding.DocNo
				requestInquiry["nominal"] = tempTrxOutstanding.GrandTotal
				requestInquiry["productCode"] = tempTrxOutstanding.ProductCode
				requestInquiry["productName"] = tempTrxOutstanding.ProductName
				requestInquiry["duration"] = duration
				requestInquiry["ouCode"] = tempMerchantKey.OuCode
				requestInquiry["status"] = constans.CONFIRM

				s, _ := structpb.NewStruct(requestInquiry)

				anyResponseTrx, _ := anypb.New(s)

				result = helpers.ResponseJSON(true, constans.SUCCESS_CODE, constans.EMPTY_VALUE, anyResponseTrx)
				return &trx.MyResponse{
					Response: result,
				}, err
			} else {

				if config.QRISPayment == constans.YES && tempTrxOutstanding.GrandTotal > 0 {
					if tempTrxOutstanding.GrandTotal == resultProducts.ProductRules.Price {

						channelCallback["docNo"] = tempTrxOutstanding.DocNo
						channelCallback["qrCode"] = redisResult[2]
						channelCallback["type"] = "INQUIRY"
						channelCallback["grandTotal"] = tempTrxOutstanding.GrandTotal
						channelCallback["productCode"] = tempTrxOutstanding.ProductCode
						channelCallback["productName"] = tempTrxOutstanding.ProductName
						channelCallbackStr, _ := json.Marshal(channelCallback)

						redisClientLocal := svc.Service.RedisClientLocal.Publish(fmt.Sprintf("%s%s", input.TerminalId, "-DISPLAY"), string(channelCallbackStr))
						log.Println(redisClientLocal)
						requestInquiry["_id"] = tempResultTrx.Id
						requestInquiry["docNo"] = tempTrxOutstanding.DocNo
						requestInquiry["nominal"] = tempTrxOutstanding.GrandTotal
						requestInquiry["productCode"] = tempTrxOutstanding.ProductCode
						requestInquiry["productName"] = tempTrxOutstanding.ProductName
						requestInquiry["duration"] = duration
						requestInquiry["ouCode"] = tempMerchantKey.OuCode
						requestInquiry["status"] = constans.INQUIRY
						s, _ := structpb.NewStruct(requestInquiry)

						anyResponseTrx, _ := anypb.New(s)

						result = helpers.ResponseJSON(true, constans.SUCCESS_CODE, constans.EMPTY_VALUE, anyResponseTrx)
						return &trx.MyResponse{
							Response: result,
						}, err

					}
				}

			}
		} else {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, "Transaction Already Success!", nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}
	}

	var trxInvoiceItemList []*trx.TrxInvoiceItem
	trxInvoiceItem := &trx.TrxInvoiceItem{
		DocNo:                  docNo,
		ProductId:              resultProducts.ProductId,
		ProductCode:            resultProducts.ProductCode,
		ProductName:            resultProducts.ProductName,
		IsPctServiceFee:        resultProducts.IsPctServiceFee,
		ServiceFee:             resultProducts.ServiceFee,
		ServiceFeeMember:       resultProducts.ServiceFeeMember,
		Price:                  resultProducts.ProductRules.Price,
		BaseTime:               resultProducts.ProductRules.BaseTime,
		ProgressiveTime:        resultProducts.ProductRules.ProgressiveTime,
		ProgressivePrice:       resultProducts.ProductRules.ProgressivePrice,
		IsPct:                  resultProducts.ProductRules.IsPct,
		ProgressivePct:         resultProducts.ProductRules.ProgressivePct,
		MaxPrice:               resultProducts.ProductRules.MaxPrice,
		Is24H:                  resultProducts.ProductRules.Is24H,
		OvernightTime:          resultProducts.ProductRules.OvernightTime,
		OvernightPrice:         resultProducts.ProductRules.OvernightPrice,
		GracePeriod:            resultProducts.ProductRules.GracePeriod,
		FlgRepeat:              resultProducts.ProductRules.FlgRepeat,
		TotalAmount:            resultProducts.ProductRules.Price,
		TotalProgressiveAmount: resultProducts.ProductRules.Price,
	}
	trxInvoiceItemList = append(trxInvoiceItemList, trxInvoiceItem)

	trxm := &trx.Trx{
		DocNo:            docNo,
		DocDate:          constans.EMPTY_VALUE,
		CheckinDateTime:  input.InquiryDatetime,
		CheckoutDateTime: constans.EMPTY_VALUE,
		DeviceIdIn:       constans.EMPTY_VALUE,
		DeviceId:         constans.EMPTY_VALUE,
		GetIn:            input.TerminalId,
		GetOut:           input.TerminalId,
		CardNumberUUIDIn: constans.EMPTY_VALUE,
		CardNumberIn:     constans.EMPTY_VALUE,
		CardNumberUUID:   constans.EMPTY_VALUE,
		CardNumber:       constans.EMPTY_VALUE,
		TypeCard:         constans.EMPTY_VALUE,
		BeginningBalance: 0,
		ExtLocalDateTime: input.InquiryDatetime,
		GrandTotal:       resultProducts.ProductRules.Price,
		ProductCode:      resultProducts.ProductCode,
		ProductName:      resultProducts.ProductName,
		ProductData:      utils.ToString(resultProduct),
		RequestData:      utils.ToString(input),
		RequestOutData:   utils.ToString(input),
		OuId:             tempMerchantKey.OuId,
		OuName:           tempMerchantKey.OuName,
		OuCode:           tempMerchantKey.OuCode,
		OuSubBranchId:    tempMerchantKey.OuSubBranchId,
		USubBranchName:   tempMerchantKey.OuSubBranchName,
		OuSubBranchCode:  tempMerchantKey.OuSubBranchCode,
		MainOuId:         tempMerchantKey.MainOuId,
		MainOuCode:       tempMerchantKey.MainOuCode,
		MainOuName:       tempMerchantKey.MainOuName,
		MemberCode:       constans.EMPTY_VALUE,
		MemberName:       constans.EMPTY_VALUE,
		MemberType:       constans.TYPE_PARTNER_FREE_PASS,
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

	qrCodePayment := constans.EMPTY_VALUE
	paymentRefDocNo := constans.EMPTY_VALUE
	if config.QRISPayment == constans.YES && trxm.GrandTotal > 0 {
		if utils.IsConnected() {
			qrCodePayment, paymentRefDocNo, err = helperService.CallQRPaymentP3(trxm, duration, svc.Service)
			if err != nil {
				result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
				return &trx.MyResponse{
					Response: result,
				}, err
			}
		} else {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, "Theres No Internet Connection!", nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}

	}
	var trxInvoiceItemLists []models.TrxInvoiceItem
	trxInvoiceItems := models.TrxInvoiceItem{
		DocNo:                  docNo,
		ProductId:              resultProducts.ProductId,
		ProductCode:            resultProducts.ProductCode,
		ProductName:            resultProducts.ProductName,
		IsPctServiceFee:        resultProducts.IsPctServiceFee,
		ServiceFee:             resultProducts.ServiceFee,
		ServiceFeeMember:       resultProducts.ServiceFeeMember,
		Price:                  resultProducts.ProductRules.Price,
		BaseTime:               resultProducts.ProductRules.BaseTime,
		ProgressiveTime:        resultProducts.ProductRules.ProgressiveTime,
		ProgressivePrice:       resultProducts.ProductRules.ProgressivePrice,
		IsPct:                  resultProducts.ProductRules.IsPct,
		ProgressivePct:         resultProducts.ProductRules.ProgressivePct,
		MaxPrice:               resultProducts.ProductRules.MaxPrice,
		Is24H:                  resultProducts.ProductRules.Is24H,
		OvernightTime:          resultProducts.ProductRules.OvernightTime,
		OvernightPrice:         resultProducts.ProductRules.OvernightPrice,
		GracePeriod:            resultProducts.ProductRules.GracePeriod,
		FlgRepeat:              resultProducts.ProductRules.FlgRepeat,
		TotalAmount:            resultProducts.ProductRules.Price,
		TotalProgressiveAmount: resultProducts.ProductRules.Price,
	}
	trxInvoiceItemLists = append(trxInvoiceItemLists, trxInvoiceItems)

	trxx := models.Trx{
		DocNo:            docNo,
		DocDate:          constans.EMPTY_VALUE,
		CheckInDatetime:  input.InquiryDatetime,
		CheckOutDatetime: constans.EMPTY_VALUE,
		DeviceIdIn:       constans.EMPTY_VALUE,
		DeviceId:         constans.EMPTY_VALUE,
		GateIn:           input.TerminalId,
		GateOut:          input.TerminalId,
		CardNumberUUIDIn: constans.EMPTY_VALUE,
		CardNumberIn:     constans.EMPTY_VALUE,
		CardNumberUUID:   constans.EMPTY_VALUE,
		CardNumber:       constans.EMPTY_VALUE,
		TypeCard:         constans.EMPTY_VALUE,
		BeginningBalance: 0,
		ExtLocalDatetime: input.InquiryDatetime,
		GrandTotal:       resultProducts.ProductRules.Price,
		ProductCode:      resultProducts.ProductCode,
		ProductName:      resultProducts.ProductName,
		ProductData:      utils.ToString(resultProduct),
		RequestData:      utils.ToString(input),
		RequestOutData:   utils.ToString(input),
		OuId:             tempMerchantKey.OuId,
		OuName:           tempMerchantKey.OuName,
		OuCode:           tempMerchantKey.OuCode,
		OuSubBranchId:    tempMerchantKey.OuSubBranchId,
		OuSubBranchName:  tempMerchantKey.OuSubBranchName,
		OuSubBranchCode:  tempMerchantKey.OuSubBranchCode,
		MainOuId:         tempMerchantKey.MainOuId,
		MainOuCode:       tempMerchantKey.MainOuCode,
		MainOuName:       tempMerchantKey.MainOuName,
		MemberCode:       constans.EMPTY_VALUE,
		MemberName:       constans.EMPTY_VALUE,
		MemberType:       constans.TYPE_PARTNER_FREE_PASS,
		CheckInTime:      0,
		CheckOutTime:     0,
		DurationTime:     0,
		VehicleNumberIn:  constans.EMPTY_VALUE,
		VehicleNumberOut: constans.EMPTY_VALUE,
		LogTrans:         constans.EMPTY_VALUE,
		MerchantKey:      config.MERCHANT_KEY,
		QrText:           qrCode,
		TrxInvoiceItem:   trxInvoiceItemLists,
		FlagSyncData:     false,
		MemberData:       nil,
		TrxAddInfo:       nil,
	}

	trxm.QrA2P = qrCodePayment
	log.Println("TRX: ", utils.ToString(trxm))
	idTrx, err := svc.Service.TrxMongoRepo.AddTrxCheckIn(trxx)
	if err != nil {
		log.Println("ERROR AddTrxCheckIn :  ", err)
		result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	if qrCodePayment != constans.EMPTY_VALUE {
		trxHeader, _ := json.Marshal(trxm)
		helperService.ConsumeTrxForScheduling(svc.Service, string(trxHeader))
		channelCallback["docNo"] = trxm.DocNo
		channelCallback["qrCode"] = qrCodePayment
		channelCallback["type"] = "INQUIRY"
		channelCallback["grandTotal"] = trxm.GrandTotal
		channelCallback["productCode"] = trxm.ProductCode
		channelCallback["productName"] = trxm.ProductName
		channelCallbackStr, _ := json.Marshal(channelCallback)

		_ = svc.Service.RedisClientLocal.Publish(fmt.Sprintf("%s%s", input.TerminalId, "-DISPLAY"), string(channelCallbackStr))
		statusRedis := svc.Service.RedisClientLocal.Set(keyInquiryPayment, fmt.Sprintf("%s#%s#%s", trxm.DocNo, paymentRefDocNo, qrCodePayment), constans.EMPTY_VALUE_INT)
		if statusRedis.Err() != nil {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, "Error Set Redis Local!", nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}

	}

	requestInquiry["_id"] = idTrx.Hex()
	requestInquiry["docNo"] = docNo
	requestInquiry["nominal"] = trxm.GrandTotal
	requestInquiry["productCode"] = input.ProductCode
	requestInquiry["productName"] = resultProducts.ProductName
	requestInquiry["duration"] = duration
	requestInquiry["ouCode"] = tempMerchantKey.OuCode
	requestInquiry["status"] = constans.INQUIRY

	s, _ := structpb.NewStruct(requestInquiry)

	anyResponseTrx, _ := anypb.New(s)

	result = helpers.ResponseJSON(true, constans.SUCCESS_CODE, constans.EMPTY_VALUE, anyResponseTrx)
	return &trx.MyResponse{
		Response: result,
	}, nil
}

func (svc trxService) ConfirmTrxP3(ctx context.Context, input *trx.RequestConfirmTrx) (*trx.MyResponse, error) { //ppp
	var result *trx.Response
	var resultTrx models.Trx

	request := new(models.RequestConfirmTrx)
	if err := helpers.BindValidateStruct(request); err != nil {
		result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	if request.ID == constans.TYPE_PARTNER_FREE_PASS {
		redisStatus := svc.Service.RedisClientLocal.Get(request.UUIDCard)
		if redisStatus.Err() != nil {
			result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, redisStatus.Err().Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, redisStatus.Err()
		}

		if err := json.Unmarshal([]byte(redisStatus.Val()), &resultTrx); err != nil {
			result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, err.Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}

		go helperService.CallSyncConfirmTrxForMemberFreePass(*request, resultTrx, svc.Service)
	} else if strings.Contains(request.ID, "SERVER") {
		redisStatus := svc.Service.RedisClientLocal.Get(request.ID)
		if redisStatus.Err() != nil {
			result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, redisStatus.Err().Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, redisStatus.Err()
		}

		if err := json.Unmarshal([]byte(redisStatus.Val()), &resultTrx); err != nil {
			result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, err.Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}

		go helperService.CallSyncConfirmTrxToCloud(nil, *request, resultTrx, svc.Service)
	} else {
		ID, _ := primitive.ObjectIDFromHex(request.ID)

		resultDataTrx, err := svc.Service.TrxMongoRepo.FindTrxOutstandingByID(ID)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, "Sesi Tidak Ditemukan", nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}

		resultTrx = resultDataTrx

		if request.CardType == constans.SETTLEMENT_CODE_QRIS {
			go helperService.CallSyncConfirmTrxToCloud(&ID, *request, resultTrx, svc.Service)

		} else {
			go helperService.CallSyncConfirmTrxCustomCard(*request, resultTrx, svc.Service)
		}
	}

	ipAddr, _ := helpers.GetPrivateIPLocal()

	responseConfirm := &trx.ResponseConfirm{
		DocNo:            resultTrx.DocNo,
		ProductData:      resultTrx.ProductData,
		ProductName:      request.ProductName,
		CardType:         resultTrx.TypeCard,
		CardNumber:       request.CardNumber,
		CheckInDatetime:  resultTrx.CheckInDatetime,
		CheckOutDatetime: request.CheckOutDatetime,
		VehicleNumberIn:  request.VehicleNumber,
		VehicleNumberOut: request.VehicleNumber,
		UuidCard:         request.UUIDCard,
		ShowQRISArea:     constans.EMPTY_VALUE,
		CurrentBalance:   request.CurrentBalance,
		GrandTotal:       request.GrandTotal,
		OuCode:           resultTrx.OuCode,
		OuName:           resultTrx.OuName,
		Address:          config.ADDRESS,
		IpAddr:           ipAddr,
	}
	anyResponseTrx, _ := anypb.New(responseConfirm)

	result = helpers.ResponseJSON(true, constans.SUCCESS_CODE, constans.EMPTY_VALUE, anyResponseTrx)
	return &trx.MyResponse{
		Response: result,
	}, nil
}

func (svc trxService) GetTrxListForDocDate(ctx context.Context, docDate *trx.Param) (*trx.MyResponse, error) {
	var result *trx.Response
	var tempTrxOutstanding *trx.Trx
	var tempTrxOutstandings []*trx.Trx
	var r *http.Request

	docDate.Param = r.URL.Query().Get("docDate")

	resultTrxList, err := svc.Service.TrxMongoRepo.GetTrxListByDocDate(docDate.Param)
	if err != nil {
		result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	var tempTrxInvoiceItems []*trx.TrxInvoiceItem
	for _, x := range resultTrxList {
		for _, v := range x.TrxInvoiceItem {
			tempTrxInvoiceItem := &trx.TrxInvoiceItem{
				DocNo:                  v.DocNo,
				ProductId:              v.ProductId,
				ProductCode:            v.ProductCode,
				ProductName:            v.ProductName,
				IsPctServiceFee:        v.IsPctServiceFee,
				ServiceFee:             v.ServiceFee,
				ServiceFeeMember:       v.ServiceFeeMember,
				Price:                  v.Price,
				BaseTime:               v.BaseTime,
				ProgressiveTime:        v.ProgressiveTime,
				ProgressivePrice:       v.ProgressivePrice,
				IsPct:                  v.IsPct,
				ProgressivePct:         v.ProgressivePct,
				MaxPrice:               v.MaxPrice,
				Is24H:                  v.Is24H,
				OvernightTime:          v.OvernightTime,
				OvernightPrice:         v.OvernightPrice,
				GracePeriod:            v.GracePeriod,
				FlgRepeat:              v.FlgRepeat,
				TotalAmount:            v.TotalAmount,
				TotalProgressiveAmount: v.TotalProgressiveAmount,
			}
			tempTrxInvoiceItems = append(tempTrxInvoiceItems, tempTrxInvoiceItem)
		}
		tempMemberData := &trx.TrxMember{
			DocNo:       x.MemberData.DocNo,
			PartnerCode: x.MemberData.PartnerCode,
			FirstName:   x.MemberData.FirstName,
			LastName:    x.MemberData.LastName,
			RoleType:    x.MemberData.RoleType,
			PhoneNumber: x.MemberData.PhoneNumber,
			Email:       x.MemberData.Email,
			Active:      x.MemberData.Active,
			ActiveAt:    x.MemberData.ActiveAt,
			NonActiveAt: func() *wrapperspb.StringValue {
				if x.MemberData.NonActiveAt != nil {
					return &wrapperspb.StringValue{Value: *x.MemberData.NonActiveAt}
				}
				return nil
			}(),
			OuId:               x.MemberData.OuId,
			TypePartner:        x.MemberData.TypePartner,
			CardNumber:         x.MemberData.CardNumber,
			VehicleNumber:      x.MemberData.VehicleNumber,
			RegisteredDatetime: x.MemberData.RegisteredDatetime,
			DateFrom:           x.MemberData.DateFrom,
			DateTo:             x.MemberData.DateTo,
			ProductId:          x.MemberData.ProductId,
			ProductCode:        x.MemberData.ProductCode,
		}

		convertedTrxAddInfo := make(map[string]*anypb.Any)
		for k, v := range x.TrxAddInfo {
			if msg, ok := v.(goproto.Message); ok {
				anyVal, err := anypb.New(msg)
				if err != nil {
					continue
				}
				convertedTrxAddInfo[k] = anyVal
			}
		}

		tempRequestAddTrxInvoiceDetailItem := &trx.TrxInvoiceDetailItem{
			DocNo:         x.RequestAddTrxInvoiceDetailItem.DocNo,
			ProductCode:   x.RequestAddTrxInvoiceDetailItem.ProductCode,
			InvoiceAmount: x.RequestAddTrxInvoiceDetailItem.InvoiceAmount,
			CreatedAt:     x.RequestAddTrxInvoiceDetailItem.CreatedAt,
			CreatedDate:   x.RequestAddTrxInvoiceDetailItem.CreatedDate,
		}

		tempTrxOutstanding = &trx.Trx{
			DocNo:                          x.DocNo,
			DocDate:                        x.DocDate,
			PaymentRefDocno:                x.PaymentRefDocNo,
			CheckinDateTime:                x.CheckInDatetime,
			CheckoutDateTime:               x.CheckOutDatetime,
			DeviceIdIn:                     x.DeviceIdIn,
			DeviceId:                       x.DeviceId,
			GetIn:                          x.GateIn,
			GetOut:                         x.GateOut,
			CardNumberUUIDIn:               x.CardNumberUUIDIn,
			CardNumberIn:                   x.CardNumberIn,
			CardNumberUUID:                 x.CardNumberUUID,
			CardNumber:                     x.CardNumber,
			TypeCard:                       x.TypeCard,
			BeginningBalance:               x.BeginningBalance,
			ExtLocalDateTime:               x.ExtLocalDatetime,
			ChargeAmount:                   x.ChargeAmount,
			GrandTotal:                     x.GrandTotal,
			ProductCode:                    x.ProductCode,
			ProductName:                    x.ProductName,
			ProductData:                    x.ProductData,
			RequestData:                    x.RequestData,
			RequestOutData:                 x.RequestOutData,
			OuId:                           x.OuId,
			OuName:                         x.OuName,
			OuCode:                         x.OuCode,
			OuSubBranchId:                  x.OuSubBranchId,
			USubBranchName:                 x.OuSubBranchName,
			OuSubBranchCode:                x.OuSubBranchCode,
			MainOuId:                       x.MainOuId,
			MainOuCode:                     x.MainOuCode,
			MainOuName:                     x.MainOuName,
			MemberCode:                     x.MemberCode,
			MemberName:                     x.MemberName,
			MemberType:                     x.MemberType,
			MemberStatus:                   x.MemberStatus,
			MemberExpiredDate:              x.MemberExpiredDate,
			CheckInTime:                    x.CheckInTime,
			CheckOutTime:                   x.CheckOutTime,
			DurationTime:                   x.DurationTime,
			VehicleNumberIn:                x.VehicleNumberIn,
			VehicleNumberOut:               x.VehicleNumberOut,
			LogTrans:                       x.LogTrans,
			MerchantKey:                    x.MerchantKey,
			QrText:                         x.QrText,
			QrA2P:                          x.QrA2P,
			QrTextPaymentOnline:            x.QrTextPaymentOnline,
			TrxInvoiceItem:                 tempTrxInvoiceItems,
			FlagSyncData:                   x.FlagSyncData,
			MemberData:                     tempMemberData,
			TrxAddInfo:                     convertedTrxAddInfo,
			FlagTrxFromCloud:               x.FlagTrxFromCloud,
			IsRsyncDataTrx:                 x.IsRsyncDataTrx,
			ExcludeSf:                      x.ExcludeSf,
			FlagCharge:                     x.FlagCharge,
			ChargeType:                     x.ChargeType,
			RequestAddTrxInvoiceDetailItem: tempRequestAddTrxInvoiceDetailItem,
			LastUpdatedAt:                  x.LastUpdatedAt,
		}
		tempTrxOutstandings = append(tempTrxOutstandings, tempTrxOutstanding)
	}

	for _, data := range tempTrxOutstandings {
		datetimeNow := utils.CurrDatetimeNow()
		checkinDatetimeParse, _ := time.Parse("2006-01-02 15:04", data.CheckinDateTime[:len(data.CheckinDateTime)-3])
		addDayDatetime := checkinDatetimeParse.AddDate(0, 0, 1)

		expiredMinutes := helpers.ConvOvernight24H(checkinDatetimeParse, addDayDatetime, datetimeNow)

		expired24h := fmt.Sprintf("%s-%s-%s-%s", data.OuCode, data.CardNumberUUIDIn, data.ProductCode, data.VehicleNumberOut)
		if expiredMinutes > 0 {
			if data.CardNumberUUIDIn != constans.EMPTY_VALUE {
				log.Println("Basic ExpiredMinutes")
				resultSetRedis := svc.Service.RedisClientLocal.Set(expired24h, data.CheckinDateTime, time.Duration(expiredMinutes)*time.Millisecond)
				log.Println("StatusSetRedis:", resultSetRedis.Val())
			}
		} else {
			if data.CardNumberUUIDIn != constans.EMPTY_VALUE {
				expiredMinutesOvernight := helpers.ConvDifferenceTimeForOvernight(data.CheckinDateTime[:len(data.CheckinDateTime)-3], data.CheckoutDateTime[:len(data.CheckoutDateTime)-3], datetimeNow)
				if expiredMinutesOvernight > 0 {
					log.Println("expired24h:", expired24h)
					log.Println("Overnight ExpiredMinutes")
					resultSetRedis := svc.Service.RedisClientLocal.Set(expired24h, data.CheckinDateTime, time.Duration(expiredMinutesOvernight)*time.Millisecond)
					log.Println("StatusSetRedis:", resultSetRedis.Val())
				}
			}
		}
	}

	result = helpers.ResponseJSON(true, constans.SUCCESS_CODE, constans.EMPTY_VALUE, nil)
	return &trx.MyResponse{
		Response: result,
	}, nil
}

func (svc trxService) UpdateStatusManualTrx(ctx context.Context, status *trx.Param) (*trx.MyResponse, error) {
	var result *trx.Response

	var r *http.Request

	status.Param = r.URL.Query().Get("status")

	redisStatus := svc.Service.RedisClientLocal.Set("P3", status, 0)
	if redisStatus.Err() != nil {
		result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, redisStatus.Err().Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, redisStatus.Err()
	}

	result = helpers.ResponseJSON(true, constans.SUCCESS_CODE, constans.EMPTY_VALUE, nil)
	return &trx.MyResponse{
		Response: result,
	}, nil
}

func (svc trxService) UpdateAutoClearTrx(ctx context.Context, status *trx.Param) (*trx.MyResponse, error) {
	var result *trx.Response

	var r *http.Request

	status.Param = r.URL.Query().Get("status")

	redisStatus := svc.Service.RedisClientLocal.Set("AUTO_CLEAR", status, 0)
	if redisStatus.Err() != nil {
		result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, redisStatus.Err().Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, redisStatus.Err()
	}

	result = helpers.ResponseJSON(true, constans.SUCCESS_CODE, constans.EMPTY_VALUE, nil)
	return &trx.MyResponse{
		Response: result,
	}, nil
}

func (svc trxService) FindTrxOutstandingByIndex(ctx context.Context, param *trx.Param) (*trx.MyResponse, error) {
	var result *trx.Response
	var r *http.Request

	param.Param = r.URL.Query().Get("index")
	var trxOutstanding *models.ResultFindTrxOutstanding
	var err error
	var tempResultTrx *trx.ResultFindTrxOutstanding

	if len(param.Param) >= 68 {
		saltKey := param.Param[len(param.Param)-16:]
		keyword := param.Param[:len(param.Param)-16]

		docNo, err := utils.Decrypt(keyword, saltKey)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)

			return &trx.MyResponse{
				Response: result,
			}, err
		}

		trxOutstanding, err = svc.Service.TrxMongoRepo.FindTrxOutstandingByDocNoCustom(docNo)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}

		tempTrxInvoiceItems := []*trx.TrxInvoiceItem{}
		for _, v := range trxOutstanding.TrxInvoiceItem {
			tempTrxInvoiceItem := &trx.TrxInvoiceItem{
				DocNo:                  v.DocNo,
				ProductId:              v.ProductId,
				ProductCode:            v.ProductCode,
				ProductName:            v.ProductName,
				IsPctServiceFee:        v.IsPctServiceFee,
				ServiceFee:             v.ServiceFee,
				ServiceFeeMember:       v.ServiceFeeMember,
				Price:                  v.Price,
				BaseTime:               v.BaseTime,
				ProgressiveTime:        v.ProgressiveTime,
				ProgressivePrice:       v.ProgressivePrice,
				IsPct:                  v.IsPct,
				ProgressivePct:         v.ProgressivePct,
				MaxPrice:               v.MaxPrice,
				Is24H:                  v.Is24H,
				OvernightTime:          v.OvernightTime,
				OvernightPrice:         v.OvernightPrice,
				GracePeriod:            v.GracePeriod,
				FlgRepeat:              v.FlgRepeat,
				TotalAmount:            v.TotalAmount,
				TotalProgressiveAmount: v.TotalProgressiveAmount,
			}
			tempTrxInvoiceItems = append(tempTrxInvoiceItems, tempTrxInvoiceItem)
		}

		tempResultTrx = &trx.ResultFindTrxOutstanding{
			Id:              trxOutstanding.ID.Hex(),
			DocNo:           trxOutstanding.DocNo,
			GrandTotal:      trxOutstanding.GrandTotal,
			CheckInDatetime: trxOutstanding.CheckInDatetime,
			OverNightPrice:  trxOutstanding.OverNightPrice,
			Is24H:           trxOutstanding.Is24H,
			CardNumber:      trxOutstanding.CardNumber,
			CardNumberUUID:  trxOutstanding.CardNumberUUID,
			OuCode:          trxOutstanding.OuCode,
			VehicleNumberIn: trxOutstanding.VehicleNumberIn,
			TrxInvoiceItem:  tempTrxInvoiceItems,
		}
	} else {
		trxOutstanding, err = svc.Service.TrxMongoRepo.FindTrxOutstandingByUUIDCustom(param.Param)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}
		tempTrxInvoiceItems := []*trx.TrxInvoiceItem{}
		for _, v := range trxOutstanding.TrxInvoiceItem {
			tempTrxInvoiceItem := &trx.TrxInvoiceItem{
				DocNo:                  v.DocNo,
				ProductId:              v.ProductId,
				ProductCode:            v.ProductCode,
				ProductName:            v.ProductName,
				IsPctServiceFee:        v.IsPctServiceFee,
				ServiceFee:             v.ServiceFee,
				ServiceFeeMember:       v.ServiceFeeMember,
				Price:                  v.Price,
				BaseTime:               v.BaseTime,
				ProgressiveTime:        v.ProgressiveTime,
				ProgressivePrice:       v.ProgressivePrice,
				IsPct:                  v.IsPct,
				ProgressivePct:         v.ProgressivePct,
				MaxPrice:               v.MaxPrice,
				Is24H:                  v.Is24H,
				OvernightTime:          v.OvernightTime,
				OvernightPrice:         v.OvernightPrice,
				GracePeriod:            v.GracePeriod,
				FlgRepeat:              v.FlgRepeat,
				TotalAmount:            v.TotalAmount,
				TotalProgressiveAmount: v.TotalProgressiveAmount,
			}
			tempTrxInvoiceItems = append(tempTrxInvoiceItems, tempTrxInvoiceItem)
		}

		tempResultTrx = &trx.ResultFindTrxOutstanding{
			Id:              trxOutstanding.ID.Hex(),
			DocNo:           trxOutstanding.DocNo,
			GrandTotal:      trxOutstanding.GrandTotal,
			CheckInDatetime: trxOutstanding.CheckInDatetime,
			OverNightPrice:  trxOutstanding.OverNightPrice,
			Is24H:           trxOutstanding.Is24H,
			CardNumber:      trxOutstanding.CardNumber,
			CardNumberUUID:  trxOutstanding.CardNumberUUID,
			OuCode:          trxOutstanding.OuCode,
			VehicleNumberIn: trxOutstanding.VehicleNumberIn,
			TrxInvoiceItem:  tempTrxInvoiceItems,
		}
	}
	anyResponseTrx, _ := anypb.New(tempResultTrx)

	result = helpers.ResponseJSON(true, constans.SUCCESS_CODE, constans.EMPTY_VALUE, anyResponseTrx)
	return &trx.MyResponse{
		Response: result,
	}, nil
}

func (svc trxService) UpdateProductPrice(ctx context.Context, input *trx.RequestUpdateProductPrice) (*trx.MyResponse, error) {
	var result *trx.Response
	var r *http.Request
	ouId := r.Context().Value("id").(float64)
	var keyword string

	if err := helpers.BindValidateStruct(input); err != nil {
		result = helpers.ResponseJSON(false, constans.VALIDATE_ERROR_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	dataUser, err := svc.Service.UserMongoRepo.FindUserByIndex(input.Username)
	if err != nil {
		result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, "user not found!", nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}
	var ous []*user.Ou
	for _, v := range dataUser.User.OuList {
		ou := &user.Ou{
			Id:              v.ID,
			OuId:            v.OuId,
			OuName:          v.OuName,
			OuCode:          v.OuCode,
			OuSubBranchId:   v.OuSubBranchId,
			OuSubBranchName: v.OuSubBranchName,
			OuSubBranchCode: v.OuSubBranchCode,
			MainOuId:        v.MainOuId,
			MainOuCode:      v.MainOuCode,
			MainOuName:      v.MainOuName,
		}
		ous = append(ous, ou)
	}

	usersync := &user.UserSync{
		Id:              dataUser.User.Id,
		Name:            dataUser.User.Name,
		Username:        dataUser.User.Username,
		Password:        dataUser.User.Password,
		TypeUser:        dataUser.User.TypeUser,
		Email:           dataUser.User.Email,
		RolesId:         dataUser.User.RolesId,
		RolesName:       dataUser.User.RolesName,
		IsAdmin:         dataUser.User.IsAdmin,
		IsInternal:      dataUser.User.IsInternal,
		Active:          dataUser.User.Active,
		OuList:          ous,
		TaskList:        dataUser.User.TaskList,
		OuDefaultId:     dataUser.User.OuDefaultId,
		PolicyDefaultId: dataUser.User.PolicyDefaultId,
		OuCode:          dataUser.User.OuCode,
		OuName:          dataUser.User.OuName,
		PinUser:         dataUser.User.PinUser,
	}

	var ouSync []*user.OuSync
	for _, x := range dataUser.OuList {
		Ousync := &user.OuSync{
			Id:         x.Id,
			OuCode:     x.OuCode,
			OuName:     x.OuName,
			OuType:     x.OuType,
			OuParentId: x.OuParentId,
		}
		ouSync = append(ouSync, Ousync)
	}
	Devicesync := []*user.DeviceSync{}
	for _, c := range dataUser.DeviceList {
		deviceSync := &user.DeviceSync{
			DeviceId:       c.DeviceId,
			FlgProgressive: c.FlgProgressive,
			MerchantKey:    c.MerchantKey,
		}
		Devicesync = append(Devicesync, deviceSync)
	}

	dataUsers := &user.UserLoginLocal{
		User:               usersync,
		OuList:             ouSync,
		MerchantKeyParking: dataUser.MerchantKeyParking,
		TaskList:           dataUser.TaskList,
		DeviceList:         Devicesync,
	}

	if dataUsers.User.PinUser != input.Pin {
		result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, "pin don't match!", nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	idTrx, _ := primitive.ObjectIDFromHex(input.Id)
	resultTrxOutstanding, err := svc.Service.TrxMongoRepo.FindTrxOutstandingByID(idTrx)
	if err != nil {
		result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	resultProduct, err := svc.Service.ProductRepo.FindPolicyOuProductByAdvance(int64(ouId), input.ProductId)
	if err != nil {
		result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}
	tempProd := &prod.ProductOuWithRules{
		OuId:             resultProduct.ProductOuWithRules.OuId,
		ProductId:        resultProduct.ProductOuWithRules.ProductId,
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
	resultProducts := &prod.PolicyOuProductWithRules{
		OuId:                  resultProduct.OuId,
		OuCode:                resultProduct.OuCode,
		OuName:                resultProduct.OuName,
		ProductId:             resultProduct.ProductId,
		ProductCode:           resultProduct.ProductCode,
		ProductName:           resultProduct.ProductName,
		ServiceFee:            resultProduct.ServiceFee,
		IsPctServiceFee:       resultProduct.IsPctServiceFee,
		IsPctServiceFeeMember: resultProduct.IsPctServiceFeeMember,
		ServiceFeeMember:      resultProduct.ServiceFeeMember,
		ProductRules:          tempProd,
	}

	keyword = resultTrxOutstanding.DocNo
	if resultTrxOutstanding.CardNumberUUIDIn != constans.EMPTY_VALUE {
		keyword = resultTrxOutstanding.CardNumberUUIDIn
	}

	err = svc.Service.TrxMongoRepo.UpdateProductById(keyword, resultProducts.ProductCode, resultProducts.ProductName)
	if err != nil {
		result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}
	boolValue := wrapperspb.Bool(constans.TRUE_VALUE)

	// Convert the wrapper message to an anypb.Any message
	anyValue, _ := anypb.New(boolValue)

	result = helpers.ResponseJSON(true, constans.SUCCESS_CODE, constans.EMPTY_VALUE, anyValue)
	return &trx.MyResponse{
		Response: result,
	}, nil
}

func (svc trxService) RegisterMember(ctx context.Context, input *trx.RequestRegistrationMemberLocal) (*trx.MyResponse, error) {
	var result *trx.Response
	var r *http.Request
	username := r.Context().Value("id").(string)
	ouId := r.Context().Value("id").(float64)
	var trxOutstanding *models.ResultFindTrxOutstanding
	var trxOutstandings *trx.ResultFindTrxOutstanding
	var err error

	if err := helpers.BindValidateStruct(input); err != nil {
		result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	if len(input.Keyword) == 14 {
		trxOutstanding, err = svc.Service.TrxMongoRepo.FindTrxOutstandingByUUIDCustom(input.Keyword)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}
	} else {
		trxOutstanding, err = svc.Service.TrxMongoRepo.FindTrxOutstandingByDocNoCustom(input.Keyword)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}
	}
	tempTrxInvoiceItems := []*trx.TrxInvoiceItem{}
	for _, v := range trxOutstanding.TrxInvoiceItem {
		tempTrxInvoiceItem := &trx.TrxInvoiceItem{
			DocNo:                  v.DocNo,
			ProductId:              v.ProductId,
			ProductCode:            v.ProductCode,
			ProductName:            v.ProductName,
			IsPctServiceFee:        v.IsPctServiceFee,
			ServiceFee:             v.ServiceFee,
			ServiceFeeMember:       v.ServiceFeeMember,
			Price:                  v.Price,
			BaseTime:               v.BaseTime,
			ProgressiveTime:        v.ProgressiveTime,
			ProgressivePrice:       v.ProgressivePrice,
			IsPct:                  v.IsPct,
			ProgressivePct:         v.ProgressivePct,
			MaxPrice:               v.MaxPrice,
			Is24H:                  v.Is24H,
			OvernightTime:          v.OvernightTime,
			OvernightPrice:         v.OvernightPrice,
			GracePeriod:            v.GracePeriod,
			FlgRepeat:              v.FlgRepeat,
			TotalAmount:            v.TotalAmount,
			TotalProgressiveAmount: v.TotalProgressiveAmount,
		}
		tempTrxInvoiceItems = append(tempTrxInvoiceItems, tempTrxInvoiceItem)
	}

	trxOutstandings = &trx.ResultFindTrxOutstanding{
		Id:              trxOutstanding.ID.Hex(),
		DocNo:           trxOutstanding.DocNo,
		GrandTotal:      trxOutstanding.GrandTotal,
		CheckInDatetime: trxOutstanding.CheckInDatetime,
		OverNightPrice:  trxOutstanding.OverNightPrice,
		Is24H:           trxOutstanding.Is24H,
		CardNumber:      trxOutstanding.CardNumber,
		CardNumberUUID:  trxOutstanding.CardNumberUUID,
		OuCode:          trxOutstanding.OuCode,
		VehicleNumberIn: trxOutstanding.VehicleNumberIn,
		TrxInvoiceItem:  tempTrxInvoiceItems,
	}

	requestValidateRegistration := models.ValidateRegis{
		RequestDateFrom: input.StartDate,
		RequestDateTo:   input.EndDate,
	}

	exists, errValidateMember := utils.ValidateBackdate(requestValidateRegistration)
	if !exists {
		result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, errValidateMember, nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	_, exists, err = svc.Service.MemberRepo.IsMemberExistsByPartnerCode(trxOutstandings.DocNo, input.StartDate, input.Keyword)
	if err != nil {
		result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	if exists {
		updateMember := models.UpdateMember{
			PartnerCode: trxOutstandings.DocNo,
			CardNumber:  input.Keyword,
			DateFrom:    input.StartDate,
			DateTo:      input.EndDate,
			UpdatedAt:   utils.Timestamp(),
			UpdatedBy:   username,
		}

		err = svc.Service.MemberRepo.UpdateMemberByPartnerCode(updateMember, nil)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}
	} else {
		addMember := models.Member{
			PartnerCode:        trxOutstandings.DocNo,
			FirstName:          input.FirstName,
			LastName:           input.LastName,
			RoleType:           constans.ROLE_TYPE_GENERAL,
			PhoneNumber:        input.PhoneNumber,
			Email:              constans.EMPTY_VALUE,
			Active:             constans.YES,
			ActiveAt:           utils.Timestamp(),
			NonActiveAt:        nil,
			OuId:               int64(ouId),
			TypePartner:        constans.TYPE_PARTNER_ONE_TIME,
			CardNumber:         input.Keyword,
			VehicleNumber:      input.VehicleNumber,
			RegisteredDatetime: utils.Timestamp(),
			DateFrom:           input.StartDate,
			DateTo:             input.EndDate,
			ProductId:          constans.NULL_LONG_VALUE,
			ProductCode:        constans.EMPTY_VALUE,
			CreatedAt:          utils.Timestamp(),
			CreatedBy:          username,
			UpdatedAt:          utils.Timestamp(),
			UpdatedBy:          username,
		}

		_, err = svc.Service.MemberRepo.AddMember(addMember, nil)
		if err != nil {
			result = helpers.ResponseJSON(false, constans.DATA_ERROR_CODE, err.Error(), nil)
			return &trx.MyResponse{
				Response: result,
			}, err
		}
	}
	boolValue := wrapperspb.Bool(constans.TRUE_VALUE)

	// Convert the wrapper message to an anypb.Any message
	anyValue, err := anypb.New(boolValue)

	result = helpers.ResponseJSON(true, constans.SUCCESS_CODE, constans.EMPTY_VALUE, anyValue)
	return &trx.MyResponse{
		Response: result,
	}, nil
}

func (svc trxService) DecryptMKey(ctx context.Context, input *trx.Decrypt) (*trx.MyResponse, error) {
	var result *trx.Response

	if err := helpers.BindValidateStruct(input); err != nil {
		result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}

	merchantKey, err := utils.DecryptMerchantKey(input.Keyword)
	if err != nil {
		result = helpers.ResponseJSON(false, constans.MALFUNCTION_SYSTEM_CODE, err.Error(), nil)
		return &trx.MyResponse{
			Response: result,
		}, err
	}
	tempMerchantKey := &trx.MerchantKey{
		ID:              merchantKey.ID,
		OuId:            merchantKey.OuId,
		OuName:          merchantKey.OuName,
		OuCode:          merchantKey.OuCode,
		OuSubBranchId:   merchantKey.OuSubBranchId,
		OuSubBranchName: merchantKey.OuSubBranchName,
		OuSubBranchCode: merchantKey.OuSubBranchCode,
		MainOuId:        merchantKey.MainOuId,
		MainOuCode:      merchantKey.MainOuCode,
		MainOuName:      merchantKey.MainOuName,
	}

	anyResponseTrx, _ := anypb.New(tempMerchantKey)

	result = helpers.ResponseJSON(true, constans.SUCCESS_CODE, constans.EMPTY_VALUE, anyResponseTrx)
	return &trx.MyResponse{
		Response: result,
	}, nil
}
