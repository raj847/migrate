package helperService

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"github.com/raj847/togrpc/config"
	"github.com/raj847/togrpc/constans"
	"github.com/raj847/togrpc/helpers"
	"github.com/raj847/togrpc/models"
	"github.com/raj847/togrpc/proto/trxLocal"
	"github.com/raj847/togrpc/services"
	"github.com/raj847/togrpc/utils"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

func SyncTrxMemberToCloud(Service services.UsecaseService, TrxMember models.TrxMember) error {

	err := Service.MemberMongoRepo.AddTrxMember(TrxMember)
	if err != nil {
		fmt.Println("Error AddTrxMember :", err)
		return err
	}

	trxMemberPub, _ := json.Marshal(TrxMember)
	Service.RedisClient.Publish(constans.CHANNEL_PG_TRX_MEMBER, trxMemberPub)

	return nil
}

func WorkerExtCloud(req interface{}, callbackURL, basicAuth, bearerToken string) (models.WorkerResponseAddMemberExt, error) {
	var result models.WorkerResponseAddMemberExt
	result.SnapshotData = req
	result.URL = callbackURL

	bodyReq, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Err workerAddMemberExtCloud - json.Marshal : ", err.Error())
		result.StatusCode = constans.EMPTY_VALUE
		result.StatusDesc = err.Error()
		return result, err
	}

	bodyReqToStr := string(bodyReq)
	fmt.Println(bodyReqToStr)

	httpReq, err := http.NewRequest("POST", callbackURL, bytes.NewBuffer(bodyReq))
	if err != nil {
		fmt.Println("Err workerAddMemberExtCloud - http.NewRequest ", callbackURL, " :", err.Error())
		result.StatusCode = constans.EMPTY_VALUE
		result.StatusDesc = err.Error()
		return result, err
	}
	defer httpReq.Body.Close()

	httpReq.Close = true
	httpReq.Header.Add("Content-Type", "application/json")
	if basicAuth != constans.EMPTY_VALUE {
		httpReq.Header.Set("Authorization", fmt.Sprintf("%s %s", "Basic", basicAuth))
	} else if bearerToken != constans.EMPTY_VALUE {
		httpReq.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", bearerToken))
	}
	httpReq.Header.Set("Connection", "close")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	//PRINTOUT BODY REQ
	fmt.Println("URL :>", callbackURL)
	fmt.Println("Request :>", bodyReqToStr)

	resp, err := client.Do(httpReq)
	if err != nil {
		fmt.Println("Err workerTrxCallbackURL - client.Do :", err.Error())
		result.StatusCode = constans.EMPTY_VALUE
		result.StatusDesc = err.Error()
		return result, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	result.StatusCode = strconv.Itoa(resp.StatusCode)
	result.StatusDesc = resp.Status
	result.ResponseBody = bodyString

	fmt.Println("ResponseCallback:", utils.ToString(result))

	return result, nil
}

func WorkerExtCloudPayment(req interface{}, callbackURL, basicAuth, bearerToken string) ([]byte, error) {
	var result models.WorkerResponseAddMemberExt
	result.SnapshotData = req
	result.URL = callbackURL

	bodyReq, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Err workerAddMemberExtCloud - json.Marshal : ", err.Error())
		result.StatusCode = constans.EMPTY_VALUE
		result.StatusDesc = err.Error()
		return nil, err
	}

	bodyReqToStr := string(bodyReq)
	fmt.Println(bodyReqToStr)

	httpReq, err := http.NewRequest("POST", callbackURL, bytes.NewBuffer(bodyReq))
	if err != nil {
		fmt.Println("Err workerAddMemberExtCloud - http.NewRequest ", callbackURL, " :", err.Error())
		result.StatusCode = constans.EMPTY_VALUE
		result.StatusDesc = err.Error()
		return nil, err
	}
	defer httpReq.Body.Close()

	httpReq.Close = true
	httpReq.Header.Add("Content-Type", "application/json")
	if basicAuth != constans.EMPTY_VALUE {
		httpReq.Header.Set("Authorization", fmt.Sprintf("%s %s", "Basic", basicAuth))
	} else if bearerToken != constans.EMPTY_VALUE {
		httpReq.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", bearerToken))
	}
	httpReq.Header.Set("Connection", "close")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	//PRINTOUT BODY REQ
	fmt.Println("URL :>", callbackURL)
	fmt.Println("Request :>", bodyReqToStr)

	resp, err := client.Do(httpReq)
	if err != nil {
		fmt.Println("Err workerTrxCallbackURL - client.Do :", err.Error())
		result.StatusCode = constans.EMPTY_VALUE
		result.StatusDesc = err.Error()
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	result.StatusCode = strconv.Itoa(resp.StatusCode)
	result.StatusDesc = resp.Status
	result.ResponseBody = bodyString

	fmt.Println("ResponseCallback:", utils.ToString(result))

	return bodyBytes, nil
}

func CallCheckTrxAlreadyExists(docNo string, svc services.UsecaseService) (errCode string, err error) {
	if utils.IsConnected() {
		responseBodyStr, err := helpers.GetCallAPI(fmt.Sprintf("%s/%s", "/mpos/local/check-trxLocal", docNo))
		if err != nil {
			return constans.MALFUNCTION_SYSTEM_CODE, err
		}

		var responseCheckTrx models.ResponseCheckTrx
		if err = json.Unmarshal([]byte(*responseBodyStr), &responseCheckTrx); err != nil {
			return constans.MALFUNCTION_SYSTEM_CODE, err
		}

		if !responseCheckTrx.Result {
			return constans.CARD_ALREADY_USED_CODE, errors.New("Sesi Kartu Masih Digunakan")
		}

		svc.TrxMongoRepo.RemoveTrxByDocNo(docNo)
	}

	return constans.SUCCESS_CODE, nil
}

func CallSyncConfirmTrxCustomCard(request models.RequestConfirmTrx, resultTrx models.Trx, svc services.UsecaseService) error {
	requestDataOut, _ := json.Marshal(request)

	cardNumber := constans.EMPTY_VALUE
	logTrans := constans.EMPTY_VALUE
	uuidCard := constans.EMPTY_VALUE
	if request.LogTrans != constans.EMPTY_VALUE {
		cardNumber = request.CardNumber
		logTrans = request.LogTrans
		uuidCard = request.UUIDCard
	}

	merchantKey, err := utils.DecryptMerchantKey(config.MERCHANT_KEY)
	if err != nil {
		log.Println("Error Service DecryptMerchantKey: ", err.Error())
		return errors.New(fmt.Sprintf("%s:%s", "Error Service DecryptMerchantKey", err.Error()))
	}

	resultProduct, err := svc.ProductRepo.FindProductByProductCode(request.ProductCode, merchantKey.OuId)
	if err != nil {
		log.Println("Error Service FindProductByProductCode: ", err.Error())
		return errors.New(fmt.Sprintf("%s:%s", "Error Service FindProductByProductCode", err.Error()))
	}

	checkinDatetimeParse, _ := time.Parse("2006-01-02 15:04:05", request.CheckOutDatetime)
	yearMonth := checkinDatetimeParse.Format("060102")
	HourCheckIn := checkinDatetimeParse.Format("15")
	prefixDocNo := utils.RandStringBytesMaskImprSrcChr(4)
	prefix := fmt.Sprintf("%s%s", yearMonth, HourCheckIn)
	autoNumber, err := svc.GenAutoNumRepo.AutonumberValueWithDatatype(constans.DATATYPE_TRX_LOCAL, prefix, 4)
	if err != nil {
		log.Println("Error Service AutonumberValueWithDatatype: ", err.Error())
		return errors.New(fmt.Sprintf("%s:%s", "Error Service AutonumberValueWithDatatype", err.Error()))
	}

	docNo := fmt.Sprintf("%s%s%d", prefixDocNo, autoNumber, merchantKey.OuId)
	randomKey := utils.RandStringBytesMaskImprSrc(16)
	encQrTxt, _ := utils.Encrypt(docNo, randomKey)
	qrCode := fmt.Sprintf("%s%s", encQrTxt.Result, randomKey)

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
		TotalAmount:      request.GrandTotal,
	}
	trxInvoiceItemList = append(trxInvoiceItemList, trxInvoiceItem)

	trxPrepaid := &models.Trx{
		DocNo:            docNo,
		DocDate:          utils.DateNow(),
		CheckInDatetime:  request.CheckOutDatetime,
		CheckOutDatetime: request.CheckOutDatetime,
		DeviceIdIn:       request.DeviceId,
		DeviceId:         request.DeviceId,
		GateIn:           request.IpTerminal,
		GateOut:          request.IpTerminal,
		CardNumberUUIDIn: uuidCard,
		CardNumberIn:     cardNumber,
		CardNumberUUID:   uuidCard,
		CardNumber:       cardNumber,
		TypeCard:         request.CardType,
		BeginningBalance: float64(request.CurrentBalance),
		ExtLocalDatetime: request.CheckOutDatetime,
		GrandTotal:       request.GrandTotal,
		ProductCode:      resultProduct.ProductCode,
		ProductName:      resultProduct.ProductName,
		ProductData:      utils.ToString(resultProduct),
		RequestData:      utils.ToString(requestDataOut),
		RequestOutData:   utils.ToString(requestDataOut),
		OuId:             merchantKey.OuId,
		OuName:           merchantKey.OuName,
		OuCode:           merchantKey.OuCode,
		OuSubBranchId:    merchantKey.OuSubBranchId,
		OuSubBranchName:  merchantKey.OuSubBranchName,
		OuSubBranchCode:  merchantKey.OuSubBranchCode,
		MainOuId:         merchantKey.MainOuId,
		MainOuCode:       merchantKey.MainOuCode,
		MainOuName:       merchantKey.MainOuName,
		CheckInTime:      checkinDatetimeParse.Unix(),
		CheckOutTime:     checkinDatetimeParse.Unix(),
		DurationTime:     constans.EMPTY_VALUE_INT,
		VehicleNumberIn:  constans.EMPTY_VALUE,
		VehicleNumberOut: constans.EMPTY_VALUE,
		LogTrans:         logTrans,
		MerchantKey:      config.MERCHANT_KEY,
		QrText:           qrCode,
		TrxInvoiceItem:   trxInvoiceItemList,
		FlagSyncData:     constans.TRUE_VALUE,
		MemberData:       nil,
		TrxAddInfo:       nil,
	}

	inquiryDate := request.CheckOutDatetime[:len(request.CheckOutDatetime)-9]
	resultMember, exists, err := svc.MemberRepo.IsMemberByAdvanceIndex(request.UUIDCard, constans.EMPTY_VALUE, inquiryDate, config.MEMBER_BY, false)
	if err != nil {
		log.Println("Error Service IsMemberByAdvanceIndex: ", err.Error())
		return errors.New(fmt.Sprintf("%s:%s", "Error Service IsMemberByAdvanceIndex", err.Error()))
	}

	isFreePass := false
	isSpecialMember := false
	if exists {
		if resultMember.TypePartner == constans.TYPE_PARTNER_FREE_PASS {
			isFreePass = true
		} else if resultMember.TypePartner == constans.TYPE_PARTNER_SPECIAL_MEMBER {
			isSpecialMember = true
		}
	}

	if !isFreePass {
		productData, existsProduct := svc.TrxMongoRepo.IsTrxProductCustomExistsByKeyword(request.UUIDCard)
		if existsProduct {
			request.ProductCode = productData.ProductCode
		}

		resultProduct, err = svc.ProductRepo.FindProductByProductCode(request.ProductCode, merchantKey.OuId)
		if err != nil {
			log.Println("Error Service FindProductByProductCode: ", err.Error())
			return errors.New(fmt.Sprintf("%s:%s", "Error Service FindProductByProductCode", err.Error()))
		}

		if isSpecialMember {
			bodySplit := strings.Split(request.ProductCode, "SPM")
			if len(bodySplit) > 1 {
				if request.ProductCode == resultMember.ProductCode {
					resultProduct, err = svc.ProductRepo.FindProductByProductCode(request.ProductCode, merchantKey.OuId)
					if err != nil {
						log.Println("Error Service FindProductByProductCode: ", err.Error())
						return errors.New(fmt.Sprintf("%s:%s", "Error Service FindProductByProductCode", err.Error()))
					}
				}
			} else {
				request.ProductCode = fmt.Sprintf("%s-%s", "SPM", request.ProductCode)
				if request.ProductCode == resultMember.ProductCode {
					resultProduct, err = svc.ProductRepo.FindProductByProductCode(request.ProductCode, merchantKey.OuId)
					if err != nil {
						log.Println("Error Service FindProductByProductCode: ", err.Error())
						return errors.New(fmt.Sprintf("%s:%s", "Error Service FindProductByProductCode", err.Error()))
					}
				}
			}

		}

		grandTotal := request.GrandTotal
		checkinDate := request.CheckOutDatetime[:len(request.CheckOutDatetime)-9]

		svc.RedisClientLocal.Del(fmt.Sprintf("%s-%s", docNo, constans.MEMBER))
		var resultTrxMemberList []models.TrxMember
		if exists && !isSpecialMember {
			resultActiveMemberList, err := svc.MemberRepo.GetMemberActiveListByPeriod(request.UUIDCard, constans.EMPTY_VALUE, checkinDate, inquiryDate, config.MEMBER_BY, false)
			if err != nil {
				log.Println("Error Service GetMemberActiveListByPeriod: ", err.Error())
				return errors.New(fmt.Sprintf("%s", "Sesi Tidak Ditemukan"))
			}

			for _, rows := range resultActiveMemberList {
				if rows.ProductCode == request.ProductCode {
					trxMemberList := models.TrxMember{
						DocNo:              docNo,
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
		}

		if len(resultTrxMemberList) > 0 {
			memberData := resultTrxMemberList[len(resultTrxMemberList)-1]

			grandTotal = 0
			log.Println("Grand Total After Calc Invoice Amount - Total Progressive Amount:", grandTotal)
			if grandTotal == 0 {
				trxPrepaid.MemberCode = memberData.PartnerCode
				trxPrepaid.MemberName = strings.TrimSpace(fmt.Sprintf("%s %s", memberData.FirstName, memberData.LastName))
				trxPrepaid.MemberType = memberData.TypePartner
				trxPrepaid.MemberData = &memberData
			}
		} else if isSpecialMember && request.ProductCode == resultMember.ProductCode {
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

			trxPrepaid.MemberCode = resultMember.PartnerCode
			trxPrepaid.MemberName = strings.TrimSpace(fmt.Sprintf("%s %s", resultMember.FirstName, resultMember.LastName))
			trxPrepaid.MemberType = resultMember.TypePartner
			trxPrepaid.MemberData = &trxMemberList

		}
		trxPrepaid.GrandTotal = grandTotal
	} else {
		trxPrepaid.MemberCode = resultMember.PartnerCode
		trxPrepaid.MemberName = strings.TrimSpace(fmt.Sprintf("%s %s", resultMember.FirstName, resultMember.LastName))
		trxPrepaid.MemberType = resultMember.TypePartner
	}

	if trxPrepaid.GrandTotal == 0 || request.CardType == constans.SETTLEMENT_CODE_CASH {
		trxPrepaid.TypeCard = constans.SETTLEMENT_CODE_CASH
		trxPrepaid.DeviceId = constans.EMPTY_VALUE
	}

	normalSfee := config.NORMAL_SF
	if request.GrandTotal == constans.EMPTY_VALUE_INT {
		normalSfee = constans.NO
	}

	dataStr, _ := json.Marshal(trxPrepaid)

	if utils.IsConnected() {
		if normalSfee == constans.YES {
			redisStatus := svc.RedisClient.Publish(constans.CHANNEL_REDIS_PG_PARKING_CHECKOUT, dataStr)
			if redisStatus.Err() != nil || redisStatus.Val() == 0 {
				trxPrepaid.FlagSyncData = false
			}
		} else if normalSfee == constans.NO {
			redisStatus := svc.RedisClient.Publish(constans.CHANNEL_REDIS_PG_PARKING_CHECKOUT_EXCLUDE_SF, dataStr)
			if redisStatus.Err() != nil || redisStatus.Val() == 0 {
				trxPrepaid.FlagSyncData = false
			}
		}

	} else {
		trxPrepaid.FlagSyncData = false
	}

	_, err = svc.TrxMongoRepo.AddTrx(*trxPrepaid)
	if err != nil {
		log.Println("ERROR AddTrx : ", err.Error())
		return err
	}

	return nil
}

func CallSyncConfirmTrxToCloud(ID *primitive.ObjectID, request models.RequestConfirmTrx, resultTrx models.Trx, svc services.UsecaseService) error {
	requestDataOut, _ := json.Marshal(request)
	checkoutDatetimeParse, _ := time.Parse("2006-01-02 15:04:05", request.CheckOutDatetime)
	timeCheckOutUnix := checkoutDatetimeParse.Unix()
	durationTime := (timeCheckOutUnix - resultTrx.CheckInTime) / 60

	resultTrx.CheckOutDatetime = request.CheckOutDatetime
	resultTrx.CheckOutTime = timeCheckOutUnix
	resultTrx.DurationTime = durationTime
	resultTrx.CardNumberUUID = request.UUIDCard
	resultTrx.CardNumber = request.CardNumber
	resultTrx.TypeCard = request.CardType
	resultTrx.DeviceId = request.DeviceId
	resultTrx.GateOut = request.IpTerminal
	resultTrx.ProductCode = request.ProductCode
	resultTrx.ProductName = request.ProductName
	resultTrx.RequestOutData = string(requestDataOut)
	resultTrx.LogTrans = request.LogTrans
	resultTrx.GrandTotal = request.GrandTotal
	resultTrx.DocDate = utils.DateNow()
	resultTrx.MerchantKey = config.MERCHANT_KEY

	resultMember := svc.RedisClientLocal.Get(fmt.Sprintf("%s-%s", resultTrx.DocNo, constans.MEMBER))
	log.Println("Result Member:", resultMember)

	if resultMember.Val() != constans.EMPTY_VALUE {
		var trxMember models.TrxMember

		err := json.Unmarshal([]byte(resultMember.Val()), &trxMember)
		if err != nil {
			log.Println("Error Unmarshal Confirm Trx:", err.Error())
		}

		resultTrx.MemberCode = trxMember.PartnerCode
		resultTrx.MemberName = strings.TrimSpace(fmt.Sprintf("%s %s", trxMember.FirstName, trxMember.LastName))
		resultTrx.MemberType = trxMember.TypePartner
	}

	if request.GrandTotal == 0 || request.CardType == constans.SETTLEMENT_CODE_CASH {
		resultTrx.TypeCard = constans.SETTLEMENT_CODE_CASH
		resultTrx.DeviceId = constans.EMPTY_VALUE
	}

	normalSfee := config.NORMAL_SF
	if request.GrandTotal == constans.EMPTY_VALUE_INT {
		normalSfee = constans.NO
	}

	dataStr, _ := json.Marshal(resultTrx)

	log.Println("SETTING SFEE", normalSfee)
	resultTrx.FlagSyncData = true
	if utils.IsConnected() {
		if normalSfee == constans.YES {
			redisStatus := svc.RedisClient.Publish(constans.CHANNEL_REDIS_PG_PARKING_CHECKOUT, dataStr)
			if redisStatus.Err() != nil || redisStatus.Val() == 0 {
				resultTrx.FlagSyncData = false
			}
		} else if normalSfee == constans.NO {
			redisStatus := svc.RedisClient.Publish(constans.CHANNEL_REDIS_PG_PARKING_CHECKOUT_EXCLUDE_SF, dataStr)
			if redisStatus.Err() != nil || redisStatus.Val() == 0 {
				resultTrx.FlagSyncData = false
			}
		}

	} else {
		resultTrx.FlagSyncData = false
	}

	_, err := svc.TrxMongoRepo.AddTrx(resultTrx)
	if err != nil {
		log.Println("ERROR AddTrx : ", err.Error())
		return err
	}

	if ID != nil {
		err = svc.TrxMongoRepo.RemoveTrxByID(*ID)
		if err != nil {
			log.Println("ERROR RemoveTrxByID : ", err.Error())
			return err
		}

		var keyword string = resultTrx.DocNo
		if resultTrx.CardNumberUUIDIn != constans.EMPTY_VALUE {
			keyword = resultTrx.CardNumberUUIDIn
		}
		svc.TrxMongoRepo.RemoveTrxProductCustom(keyword)

	}

	redisStatus := svc.RedisClientLocal.Del(request.ID)
	if redisStatus.Err() != nil {
		return redisStatus.Err()
	}

	redisStatusRemoveP3Qris := svc.RedisClientLocal.Del("PAYMENT-INQUIRY-QRIS-PASS")
	if redisStatusRemoveP3Qris.Err() != nil {
		return redisStatusRemoveP3Qris.Err()
	}

	return nil
}

func CallSyncConfirmTrxForMemberFreePass(request models.RequestConfirmTrx, resultTrx models.Trx, svc services.UsecaseService) error {
	checkoutDatetimeParse, _ := time.Parse("2006-01-02 15:04:05", request.CheckOutDatetime)
	timeCheckOutUnix := checkoutDatetimeParse.Unix()
	durationTime := (timeCheckOutUnix - resultTrx.CheckInTime) / 60

	resultTrx.CheckOutDatetime = request.CheckOutDatetime
	resultTrx.CheckOutTime = timeCheckOutUnix
	resultTrx.DurationTime = durationTime
	resultTrx.DocDate = utils.DateNow()
	resultTrx.DeviceId = request.DeviceId
	resultTrx.TypeCard = request.CardType
	resultTrx.CardNumberUUID = request.UUIDCard
	resultTrx.CardNumber = request.CardNumber
	resultTrx.GateOut = request.IpTerminal
	resultTrx.ProductCode = request.ProductCode
	resultTrx.ProductName = request.ProductName
	resultTrx.GrandTotal = request.GrandTotal
	resultTrx.VehicleNumberOut = request.VehicleNumber
	resultTrx.MerchantKey = config.MERCHANT_KEY

	dataStr, _ := json.Marshal(resultTrx)
	if utils.IsConnected() {
		redisStatus := svc.RedisClient.Publish(constans.CHANNEL_REDIS_PG_PARKING_CHECKOUT, dataStr)
		if redisStatus.Err() != nil || redisStatus.Val() == 0 {
			resultTrx.FlagSyncData = false
		}
	} else {
		resultTrx.FlagSyncData = false
	}

	_, err := svc.TrxMongoRepo.AddTrx(resultTrx)
	if err != nil {
		log.Println("ERROR AddTrx : ", err.Error())
		return err
	}

	err = svc.TrxMongoRepo.RemoveTrxByDocNo(resultTrx.DocNo)
	if err != nil {
		log.Println("ERROR RemoveTrxByDocNo : ", err.Error())
		return err
	}

	redisStatus := svc.RedisClientLocal.Del(request.UUIDCard)
	if redisStatus.Err() != nil {
		return redisStatus.Err()
	}

	return nil
}

func CallSyncConfirmTrxForMemberFreePassCustom(request models.RequestConfirmTrx, resultTrx models.Trx, svc services.UsecaseService, additionalData map[string]interface{}) error {
	checkoutDatetimeParse, _ := time.Parse("2006-01-02 15:04:05", request.CheckOutDatetime)
	timeCheckOutUnix := checkoutDatetimeParse.Unix()
	durationTime := (timeCheckOutUnix - resultTrx.CheckInTime) / 60
	datetimeNow := utils.CurrDatetimeNow()

	resultTrx.CheckOutDatetime = request.CheckOutDatetime
	resultTrx.CheckOutTime = timeCheckOutUnix
	resultTrx.DurationTime = durationTime
	resultTrx.DocDate = utils.DateNow()
	resultTrx.DeviceId = request.DeviceId
	resultTrx.TypeCard = request.CardType

	resultMember := svc.RedisClientLocal.Get(fmt.Sprintf("%s-%s", resultTrx.DocNo, constans.MEMBER))
	log.Println("Result Member:", resultMember)

	if resultMember.Val() != constans.EMPTY_VALUE {
		var trxMember models.TrxMember

		err := json.Unmarshal([]byte(resultMember.Val()), &trxMember)
		if err != nil {
			log.Println("Error Unmarshal Confirm Trx:", err.Error())
		}

		resultTrx.MemberCode = trxMember.PartnerCode
		resultTrx.MemberName = strings.TrimSpace(fmt.Sprintf("%s %s", trxMember.FirstName, trxMember.LastName))
		resultTrx.MemberType = trxMember.TypePartner
	}

	if request.GrandTotal == 0 || request.CardType == constans.SETTLEMENT_CODE_CASH {
		resultTrx.TypeCard = constans.SETTLEMENT_CODE_CASH
		resultTrx.DeviceId = constans.EMPTY_VALUE
	}

	resultTrx.CardNumberUUID = request.UUIDCard
	resultTrx.CardNumber = request.CardNumber
	resultTrx.GateOut = request.IpTerminal
	resultTrx.ProductCode = request.ProductCode
	resultTrx.ProductName = request.ProductName
	resultTrx.GrandTotal = request.GrandTotal
	resultTrx.VehicleNumberOut = request.VehicleNumber
	resultTrx.MerchantKey = config.MERCHANT_KEY

	var redisStatus *redis.IntCmd

	//normalSfee := config.NORMAL_SF
	//if request.GrandTotal == constans.EMPTY_VALUE_INT {
	//	normalSfee = constans.NO
	//}

	filter := bson.M{
		"docNo": resultTrx.DocNo,
	}

	updateSet := bson.M{
		"$set": bson.M{
			"username":  additionalData["username"],
			"shiftCode": request.ShiftCode,
		},
	}

	err := svc.TrxMongoRepo.AddTrxAddInfoInterfaces(filter, updateSet)
	if err != nil {
		log.Println("ERROR AddTrxAddInfoInterfaces : ", err.Error())
		return err
	}

	// Digunakan utk mengambil data terbaru setelah diupdate
	resultTrxAddInfo, exists := svc.TrxMongoRepo.IsTrxAddInfoInterfacesExistsByDocNo(resultTrx.DocNo)
	if exists {
		resultTrx.TrxAddInfo = resultTrxAddInfo
	}

	dataStr, _ := json.Marshal(resultTrx)
	resultTrx.FlagSyncData = true

	if *request.ExcludeSf {
		redisStatus = svc.RedisClient.Publish(constans.CHANNEL_REDIS_PG_PARKING_CHECKOUT_EXCLUDE_SF, dataStr)
	} else {
		redisStatus = svc.RedisClient.Publish(constans.CHANNEL_REDIS_PG_PARKING_CHECKOUT, dataStr)
	}

	resultTrx.FlagSyncData = true
	if redisStatus.Err() != nil || redisStatus.Val() == 0 {
		resultTrx.FlagSyncData = false
	}

	_, err = svc.TrxMongoRepo.AddTrx(resultTrx)
	if err != nil {
		log.Println("ERROR AddTrx : ", err.Error())
		return err
	}

	err = svc.TrxMongoRepo.RemoveTrxByDocNo(resultTrx.DocNo)
	if err != nil {
		log.Println("ERROR RemoveTrxByDocNo : ", err.Error())
		return err
	}

	redisStatus = svc.RedisClientLocal.Del(request.UUIDCard)
	if redisStatus.Err() != nil {
		return redisStatus.Err()
	}

	if resultTrx.MemberCode != constans.EMPTY_VALUE {
		if resultTrx.MemberCode != constans.MANUAL {
			expired24h := fmt.Sprintf("%s-%s-%s-%s", resultTrx.OuCode, resultTrx.MemberCode, request.VehicleNumber, request.ProductCode)
			log.Println("expired24h:", expired24h)
			resultGetRedis := svc.RedisClientLocal.Get(expired24h)
			if resultGetRedis.Val() == constans.EMPTY_VALUE {
				checkinDatetimeParse, _ := time.Parse("2006-01-02 15:04", resultTrx.CheckInDatetime[:len(resultTrx.CheckInDatetime)-3])
				addDayDatetime := checkinDatetimeParse.AddDate(0, 0, 1)
				expiredMinutes := helpers.ConvOvernight24H(checkinDatetimeParse, addDayDatetime, datetimeNow)

				if expiredMinutes > 0 {
					resultSetRedis := svc.RedisClientLocal.Set(expired24h, resultTrx.CheckInDatetime, time.Duration(expiredMinutes)*time.Millisecond)
					log.Println("StatusSetRedis: Basic", resultSetRedis.Val())
				} else {
					checkinDatetimeStr := resultTrx.CheckInDatetime[:len(resultTrx.CheckInDatetime)-3]
					checkoutDatetimeStr := resultTrx.CheckOutDatetime[:len(resultTrx.CheckOutDatetime)-3]

					expiredMinutesOvernight := helpers.ConvDifferenceTimeForOvernight(checkinDatetimeStr, checkoutDatetimeStr, datetimeNow)
					if expiredMinutesOvernight > 0 {
						resultSetRedis := svc.RedisClientLocal.Set(expired24h, resultTrx.CheckInDatetime, time.Duration(expiredMinutesOvernight)*time.Millisecond)
						log.Println("StatusSetRedis: Sisa Overnight", resultSetRedis.Val())
					} else {
						expiredMinutes = helpers.ConvNextInvoiceForTime(checkinDatetimeStr, checkoutDatetimeStr, datetimeNow)
						if expiredMinutes > 0 {
							resultSetRedis := svc.RedisClientLocal.Set(expired24h, resultTrx.CheckInDatetime, time.Duration(expiredMinutes)*time.Millisecond)
							log.Println("StatusSetRedis: NextInvoice Overnight", resultSetRedis.Val())
						}
					}
				}
			}
		}
	}

	return nil
}

func CallSyncConfirmTrxToCloudCustomV2(ID *primitive.ObjectID, request models.RequestConfirmTrx, resultTrx models.Trx, svc services.UsecaseService, datetimeNow time.Time, additionalData map[string]interface{}, normalOut string) error {
	var trxMember *models.TrxMember = nil
	requestDataOut, _ := json.Marshal(request)
	checkoutDatetimeParse, _ := time.Parse("2006-01-02 15:04:05", request.CheckOutDatetime)
	timeCheckOutUnix := checkoutDatetimeParse.Unix()
	durationTime := (timeCheckOutUnix - resultTrx.CheckInTime) / 60

	//isUsedInvoiceCardPayment := resultTrx.TypeCard == constans.EMPTY_VALUE

	resultTrx.CheckOutDatetime = request.CheckOutDatetime
	resultTrx.CheckOutTime = timeCheckOutUnix
	resultTrx.DurationTime = durationTime
	resultTrx.CardNumberUUID = request.UUIDCard
	resultTrx.CardNumber = request.CardNumber
	resultTrx.TypeCard = request.CardType
	resultTrx.DeviceId = request.DeviceId

	if request.GrandTotal == 0 || request.CardType == constans.SETTLEMENT_CODE_CASH {
		resultTrx.TypeCard = constans.SETTLEMENT_CODE_CASH
		resultTrx.DeviceId = constans.EMPTY_VALUE
	}

	resultTrx.GateOut = request.IpTerminal
	resultTrx.RequestOutData = string(requestDataOut)
	resultTrx.LogTrans = request.LogTrans
	resultTrx.ProductCode = request.ProductCode
	resultTrx.ProductName = request.ProductName
	resultTrx.GrandTotal = request.GrandTotal
	resultTrx.DocDate = utils.DateNow()
	resultTrx.MerchantKey = config.MERCHANT_KEY
	resultTrx.VehicleNumberIn = request.VehicleNumber
	resultTrx.VehicleNumberOut = request.VehicleNumber

	resultMember := svc.RedisClientLocal.Get(fmt.Sprintf("%s-%s", resultTrx.DocNo, constans.MEMBER))
	if resultMember.Val() != constans.EMPTY_VALUE {
		err := json.Unmarshal([]byte(resultMember.Val()), &trxMember)
		if err != nil {
			log.Println("Error Unmarshal Confirm Trx:", err.Error())
			return err
		}

		resultTrx.MemberCode = trxMember.PartnerCode
		resultTrx.MemberName = strings.TrimSpace(fmt.Sprintf("%s %s", trxMember.FirstName, trxMember.LastName))
		resultTrx.MemberType = trxMember.TypePartner
	}

	filter := bson.M{
		"docNo": resultTrx.DocNo,
	}

	updateSet := bson.M{
		"$set": bson.M{
			"username":  additionalData["username"],
			"shiftCode": request.ShiftCode,
		},
	}

	err := svc.TrxMongoRepo.AddTrxAddInfoInterfaces(filter, updateSet)
	if err != nil {
		log.Println("ERROR AddTrxAddInfoInterfaces : ", err.Error())
		return err
	}

	// Digunakan utk mengambil data terbaru setelah diupdate
	resultTrxAddInfo, exists := svc.TrxMongoRepo.IsTrxAddInfoInterfacesExistsByDocNo(resultTrx.DocNo)
	if exists {
		resultTrx.TrxAddInfo = resultTrxAddInfo
	}

	dataStr, _ := json.Marshal(resultTrx)
	var redisStatus *redis.IntCmd

	//normalSfee := config.NORMAL_SF
	//if request.GrandTotal == constans.EMPTY_VALUE_INT {
	//	normalSfee = constans.NO
	//}

	if *request.ExcludeSf {
		redisStatus = svc.RedisClient.Publish(constans.CHANNEL_REDIS_PG_PARKING_CHECKOUT_EXCLUDE_SF, dataStr)
	} else {
		redisStatus = svc.RedisClient.Publish(constans.CHANNEL_REDIS_PG_PARKING_CHECKOUT, dataStr)
	}

	resultTrx.FlagSyncData = true
	if redisStatus.Err() != nil || redisStatus.Val() == 0 {
		resultTrx.FlagSyncData = false
	}

	_, err = svc.TrxMongoRepo.AddTrx(resultTrx)
	if err != nil {
		log.Println("ERROR AddTrx : ", err.Error())
		return err
	}

	if ID != nil {
		err = svc.TrxMongoRepo.RemoveTrxByID(*ID)
		if err != nil {
			log.Println("ERROR RemoveTrxByID : ", err.Error())
			return err
		}
	}

	if resultTrx.MemberCode != constans.EMPTY_VALUE {
		if resultTrx.MemberCode != constans.MANUAL {
			expired24h := fmt.Sprintf("%s-%s-%s-%s", resultTrx.OuCode, resultTrx.MemberCode, request.VehicleNumber, request.ProductCode)
			log.Println("expired24h:", expired24h)
			resultGetRedis := svc.RedisClientLocal.Get(expired24h)
			if resultGetRedis.Val() == constans.EMPTY_VALUE {
				checkinDatetimeParse, _ := time.Parse("2006-01-02 15:04", resultTrx.CheckInDatetime[:len(resultTrx.CheckInDatetime)-3])
				addDayDatetime := checkinDatetimeParse.AddDate(0, 0, 1)
				expiredMinutes := helpers.ConvOvernight24H(checkinDatetimeParse, addDayDatetime, datetimeNow)

				if expiredMinutes > 0 {
					resultSetRedis := svc.RedisClientLocal.Set(expired24h, resultTrx.CheckInDatetime, time.Duration(expiredMinutes)*time.Millisecond)
					log.Println("StatusSetRedis: Basic", resultSetRedis.Val())
				} else {
					checkinDatetimeStr := resultTrx.CheckInDatetime[:len(resultTrx.CheckInDatetime)-3]
					checkoutDatetimeStr := resultTrx.CheckOutDatetime[:len(resultTrx.CheckOutDatetime)-3]

					expiredMinutesOvernight := helpers.ConvDifferenceTimeForOvernight(checkinDatetimeStr, checkoutDatetimeStr, datetimeNow)
					if expiredMinutesOvernight > 0 {
						resultSetRedis := svc.RedisClientLocal.Set(expired24h, resultTrx.CheckInDatetime, time.Duration(expiredMinutesOvernight)*time.Millisecond)
						log.Println("StatusSetRedis: Sisa Overnight", resultSetRedis.Val())
					} else {
						expiredMinutes = helpers.ConvNextInvoiceForTime(checkinDatetimeStr, checkoutDatetimeStr, datetimeNow)
						if expiredMinutes > 0 {
							resultSetRedis := svc.RedisClientLocal.Set(expired24h, resultTrx.CheckInDatetime, time.Duration(expiredMinutes)*time.Millisecond)
							log.Println("StatusSetRedis: NextInvoice Overnight", resultSetRedis.Val())
						}
					}
				}
			}
		}
	}

	return nil
}

func CallSyncConfirmTrxToCloudCustom(ID *primitive.ObjectID, request models.RequestConfirmTrx, resultTrx models.Trx, svc services.UsecaseService, datetimeNow time.Time, additionalData map[string]interface{}, normalOut string) error {
	var trxMember *models.TrxMember = nil
	requestDataOut, _ := json.Marshal(request)
	checkoutDatetimeParse, _ := time.Parse("2006-01-02 15:04:05", request.CheckOutDatetime)
	timeCheckOutUnix := checkoutDatetimeParse.Unix()
	durationTime := (timeCheckOutUnix - resultTrx.CheckInTime) / 60

	isUsedInvoiceCardPayment := resultTrx.TypeCard == constans.EMPTY_VALUE

	resultTrx.CheckOutDatetime = request.CheckOutDatetime
	resultTrx.CheckOutTime = timeCheckOutUnix
	resultTrx.DurationTime = durationTime
	resultTrx.CardNumberUUID = request.UUIDCard
	resultTrx.CardNumber = request.CardNumber
	resultTrx.TypeCard = request.CardType
	resultTrx.DeviceId = request.DeviceId

	if request.GrandTotal == 0 || request.CardType == constans.SETTLEMENT_CODE_CASH {
		resultTrx.TypeCard = constans.SETTLEMENT_CODE_CASH
		resultTrx.DeviceId = constans.EMPTY_VALUE
	}

	resultTrx.GateOut = request.IpTerminal
	resultTrx.RequestOutData = string(requestDataOut)
	resultTrx.LogTrans = request.LogTrans
	resultTrx.ProductCode = request.ProductCode
	resultTrx.ProductName = request.ProductName
	resultTrx.GrandTotal = request.GrandTotal
	resultTrx.DocDate = utils.DateNow()
	resultTrx.MerchantKey = config.MERCHANT_KEY
	resultTrx.VehicleNumberIn = request.VehicleNumber
	resultTrx.VehicleNumberOut = request.VehicleNumber

	resultMember := svc.RedisClientLocal.Get(fmt.Sprintf("%s-%s", resultTrx.DocNo, constans.MEMBER))
	if resultMember.Val() != constans.EMPTY_VALUE {
		err := json.Unmarshal([]byte(resultMember.Val()), &trxMember)
		if err != nil {
			log.Println("Error Unmarshal Confirm Trx:", err.Error())
			return err
		}

		resultTrx.MemberCode = trxMember.PartnerCode
		resultTrx.MemberName = strings.TrimSpace(fmt.Sprintf("%s %s", trxMember.FirstName, trxMember.LastName))
		resultTrx.MemberType = trxMember.TypePartner
	}

	filter := bson.M{
		"docNo": resultTrx.DocNo,
	}

	updateSet := bson.M{
		"$set": bson.M{
			"username":  additionalData["username"],
			"shiftCode": request.ShiftCode,
		},
	}

	err := svc.TrxMongoRepo.AddTrxAddInfoInterfaces(filter, updateSet)
	if err != nil {
		log.Println("ERROR AddTrxAddInfoInterfaces : ", err.Error())
		return err
	}

	// Digunakan utk mengambil data terbaru setelah diupdate
	resultTrxAddInfo, exists := svc.TrxMongoRepo.IsTrxAddInfoInterfacesExistsByDocNo(resultTrx.DocNo)
	if exists {
		resultTrx.TrxAddInfo = resultTrxAddInfo
	}

	dataStr, _ := json.Marshal(resultTrx)
	var redisStatus *redis.IntCmd
	normalSfee := config.NORMAL_SF
	if request.GrandTotal == constans.EMPTY_VALUE_INT {
		normalSfee = constans.NO
	}

	if *request.ExcludeSf {
		redisStatus = svc.RedisClient.Publish(constans.CHANNEL_REDIS_PG_PARKING_CHECKOUT_EXCLUDE_SF, dataStr)
	} else {
		if normalSfee == constans.YES {
			redisStatus = svc.RedisClient.Publish(constans.CHANNEL_REDIS_PG_PARKING_CHECKOUT, dataStr)
		} else {
			redisStatus = svc.RedisClient.Publish(constans.CHANNEL_REDIS_PG_PARKING_CHECKOUT_EXCLUDE_SF, dataStr)
		}
	}

	resultTrx.FlagSyncData = true
	if redisStatus.Err() != nil || redisStatus.Val() == 0 {
		resultTrx.FlagSyncData = false
	}

	_, err = svc.TrxMongoRepo.AddTrx(resultTrx)
	if err != nil {
		log.Println("ERROR AddTrx : ", err.Error())
		return err
	}

	if ID != nil {
		err = svc.TrxMongoRepo.RemoveTrxByID(*ID)
		if err != nil {
			log.Println("ERROR RemoveTrxByID : ", err.Error())
			return err
		}
	}

	if normalOut != constans.YES {
		if !(resultTrx.TypeCard == constans.SETTLEMENT_CODE_QRIS || resultTrx.TypeCard == constans.SETTLEMENT_CODE_CASH) && !isUsedInvoiceCardPayment {
			if trxMember == nil {
				expired24h := fmt.Sprintf("%s-%s-%s-%s", resultTrx.OuCode, resultTrx.CardNumberUUIDIn, request.ProductCode, request.VehicleNumber)
				log.Println("expired24h:", expired24h)
				resultGetRedis := svc.RedisClientLocal.Get(expired24h)
				if resultGetRedis.Val() == constans.EMPTY_VALUE {
					checkinDatetimeParse, _ := time.Parse("2006-01-02 15:04", resultTrx.CheckInDatetime[:len(resultTrx.CheckInDatetime)-3])
					addDayDatetime := checkinDatetimeParse.AddDate(0, 0, 1)
					expiredMinutes := helpers.ConvOvernight24H(checkinDatetimeParse, addDayDatetime, datetimeNow)

					if expiredMinutes > 0 {
						resultSetRedis := svc.RedisClientLocal.Set(expired24h, resultTrx.CheckInDatetime, time.Duration(expiredMinutes)*time.Millisecond)
						log.Println("StatusSetRedis: Basic", resultSetRedis.Val())
					} else {
						checkinDatetimeStr := resultTrx.CheckInDatetime[:len(resultTrx.CheckInDatetime)-3]
						checkoutDatetimeStr := resultTrx.CheckOutDatetime[:len(resultTrx.CheckOutDatetime)-3]

						expiredMinutesOvernight := helpers.ConvDifferenceTimeForOvernight(checkinDatetimeStr, checkoutDatetimeStr, datetimeNow)
						if expiredMinutesOvernight > 0 {
							resultSetRedis := svc.RedisClientLocal.Set(expired24h, resultTrx.CheckInDatetime, time.Duration(expiredMinutesOvernight)*time.Millisecond)
							log.Println("StatusSetRedis: Sisa Overnight", resultSetRedis.Val())
						} else {
							expiredMinutes = helpers.ConvNextInvoiceForTime(checkinDatetimeStr, checkoutDatetimeStr, datetimeNow)
							if expiredMinutes > 0 {
								resultSetRedis := svc.RedisClientLocal.Set(expired24h, resultTrx.CheckInDatetime, time.Duration(expiredMinutes)*time.Millisecond)
								log.Println("StatusSetRedis: NextInvoice Overnight", resultSetRedis.Val())
							}
						}
					}
				}
			}
		}
	}

	return nil
}

func CallQRPayment(resultTrx *trxLocal.ResultFindTrxOutstanding, request *trxLocal.RequestInquiryWithoutCard, duration utils.ConvTime, svc services.UsecaseService) error {
	resultInquiryTrxWithCard := make(map[string]interface{})
	resultInquiryTrxWithCard["duration"] = duration

	requestPublish := models.RequestInquiryRedis{
		DocNo:           resultTrx.DocNo,
		ProductCode:     request.ProductCode,
		ProductName:     resultTrx.TrxInvoiceItem[0].ProductName,
		GrandTotal:      resultTrx.TrxInvoiceItem[0].TotalAmount,
		OuCode:          resultTrx.OuCode,
		MKey:            config.MERCHANT_KEY_APPS2PAY,
		PaymentCategory: constans.PAYMENT_CATEGORY,
		DeviceId:        constans.EMPTY_VALUE,
		ChannelCallback: fmt.Sprintf("%s_%s", resultTrx.OuCode, constans.CHANNEL_CALLBACK_REDIS),
	}

	if resultTrx.TrxInvoiceItem[0].TotalAmount > 0 {
		trx, _ := json.Marshal(requestPublish)
		svc.RedisClient.Publish(constans.CHANNEL_PG_INQUIRY_PAYMENT, trx)

		dataTerminal := make(map[string]interface{})
		dataTerminal["display"] = request.TerminalId + "-DISPLAY"
		dataTerminal["confirm"] = request.TerminalId
		dataTerminal["duration"] = resultInquiryTrxWithCard["duration"]

		dataTerminalStr, _ := json.Marshal(dataTerminal)

		redisCmdDisplay := svc.RedisClientLocal.Set(resultTrx.DocNo, dataTerminalStr, 5*time.Minute)
		if redisCmdDisplay.Err() != nil {
			log.Println("Error redisCmdDisplay:", redisCmdDisplay.Err().Error())
			return redisCmdDisplay.Err()
		}

	}

	return nil
}

func CallQRPaymentP3(resultTrx *trxLocal.Trx, duration utils.ConvTime, svc services.UsecaseService) (string, string, error) {
	resultInquiryTrxWithCard := make(map[string]interface{})
	resultInquiryTrxWithCard["duration"] = duration

	requestInquiryQRIS := make(map[string]interface{})
	requestInquiryQRIS["docNo"] = resultTrx.DocNo
	requestInquiryQRIS["paymentMethod"] = constans.PAYMENT_METHOD_QRIS
	requestInquiryQRIS["productCode"] = resultTrx.ProductCode
	requestInquiryQRIS["productName"] = resultTrx.ProductName
	requestInquiryQRIS["grandTotal"] = resultTrx.GrandTotal
	requestInquiryQRIS["channelCallback"] = fmt.Sprintf("%s_%s", resultTrx.OuCode, "CALLBACK_PAYMENT")
	requestInquiryQRIS["mKey"] = config.MERCHANT_KEY_APPS2PAY

	body, err := json.Marshal(requestInquiryQRIS)
	if err != nil {
		log.Println("Error Marshal Request Inquiry QRIS :", err.Error())
		return constans.EMPTY_VALUE, constans.EMPTY_VALUE, errors.New(err.Error())
	}

	responseInquiryPaymentStr, err := helpers.CallHttpCloudServerParking("POST", body, "/mpos/local/inquiry-payment-p3")
	if err != nil {
		log.Println("Error CallHttpCloudServerParking : ", err.Error())
		return constans.EMPTY_VALUE, constans.EMPTY_VALUE, errors.New(err.Error())
	}

	var responseInquiryPayment models.ResponseInquiryPayment
	err = json.Unmarshal([]byte(*responseInquiryPaymentStr), &responseInquiryPayment)
	if err != nil {
		log.Println("Error Unmarshall : ", err.Error())
		return constans.EMPTY_VALUE, constans.EMPTY_VALUE, errors.New(err.Error())
	}

	if !responseInquiryPayment.Success {
		log.Println("Error Inquiry Payment :", responseInquiryPayment.Message)
		return constans.EMPTY_VALUE, constans.EMPTY_VALUE, errors.New(responseInquiryPayment.Message)
	}

	log.Println("RESPONSE A2P : ", utils.ToString(responseInquiryPayment))
	if resultTrx.TrxInvoiceItem[0].TotalAmount > 0 {

		dataTerminal := make(map[string]interface{})
		dataTerminal["display"] = resultTrx.GetIn + "-DISPLAY"
		dataTerminal["confirm"] = resultTrx.GetIn
		dataTerminal["duration"] = resultInquiryTrxWithCard["duration"]

		dataTerminalStr, _ := json.Marshal(dataTerminal)

		redisCmdDisplay := svc.RedisClientLocal.Set(resultTrx.DocNo, dataTerminalStr, constans.EMPTY_VALUE_INT)
		if redisCmdDisplay.Err() != nil {
			log.Println("Error redisCmdDisplay:", redisCmdDisplay.Err().Error())
			return constans.EMPTY_VALUE, constans.EMPTY_VALUE, redisCmdDisplay.Err()
		}

	}

	return responseInquiryPayment.Result.QrCode, responseInquiryPayment.Result.PaymentRefDocNo, nil
}

func CallQRPaymentCustomV2(resultTrx models.ResultFindTrxOutstanding, request models.RequestCustomInquiryWithoutCard, duration utils.ConvTime, svc services.UsecaseService) error {
	resultInquiryTrxWithCard := make(map[string]interface{})
	resultInquiryTrxWithCard["duration"] = duration

	requestPublish := models.RequestInquiryRedis{
		DocNo:           resultTrx.DocNo,
		ProductCode:     request.ProductCode,
		ProductName:     resultTrx.TrxInvoiceItem[0].ProductName,
		GrandTotal:      resultTrx.TrxInvoiceItem[0].TotalAmount,
		OuCode:          resultTrx.OuCode,
		MKey:            config.MERCHANT_KEY_APPS2PAY,
		PaymentCategory: constans.PAYMENT_CATEGORY,
		DeviceId:        constans.EMPTY_VALUE,
		ChannelCallback: fmt.Sprintf("%s_%s", resultTrx.OuCode, constans.CHANNEL_CALLBACK_REDIS),
	}

	if resultTrx.TrxInvoiceItem[0].TotalAmount > 0 {
		trx, _ := json.Marshal(requestPublish)
		svc.RedisClient.Publish(constans.CHANNEL_PG_INQUIRY_PAYMENT, trx)

		dataTerminal := make(map[string]interface{})
		dataTerminal["display"] = request.TerminalId + "-DISPLAY"
		dataTerminal["confirm"] = request.TerminalId
		dataTerminal["duration"] = resultInquiryTrxWithCard["duration"]

		dataTerminalStr, _ := json.Marshal(dataTerminal)

		redisCmdDisplay := svc.RedisClientLocal.Set(resultTrx.DocNo, dataTerminalStr, 5*time.Minute)
		if redisCmdDisplay.Err() != nil {
			log.Println("Error redisCmdDisplay:", redisCmdDisplay.Err().Error())
			return redisCmdDisplay.Err()
		}

	}

	return nil
}

func CheckTrxCloudServer(docNo string) (response *models.ResponseTrxCloud, err error) {
	log.Println("Call CheckTrxCloudServer:", docNo)

	data, err := helpers.GetCallAPI(fmt.Sprintf("%s/%s", "/mpos/local/inquiry", docNo))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(*data), &response)
	if err != nil {
		log.Println("Error Unmarshal CallUpdateServiceFeeA2P:", err.Error())
		return nil, err
	}

	return response, nil
}

func CheckTrxPaymentOnline(docNo string) (response *models.ResponseTrxPaymentOnline, err error) {
	log.Println("Call CheckTrxCloudServer:", docNo)

	data, err := helpers.GetCallAPITrxPaymentOnline(fmt.Sprintf("%s/%s", "/trxLocal/local/check-status", docNo))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(*data), &response)
	if err != nil {
		log.Println("Error Unmarshal CallUpdateServiceFeeA2P:", err.Error())
		return nil, err
	}

	return response, nil
}

func ConsumeTrxForScheduling(svc services.UsecaseService, trx string) error {
	log.Println("ConsumeTrxForScheduling with data:", trx)
	//log.Println("ConsumeTrxForScheduling IsConnected:", utils.IsConnected())

	//if utils.IsConnected() {
	statusPublish := svc.RedisClient.Publish(config.DEFAULTChannelRedisParking, trx)
	log.Println("Publish Redis Cloud:", config.DEFAULTChannelRedisParking, "Status:", statusPublish.Val())

	if statusPublish.Val() == 0 {
		statusPublish = svc.RedisClientLocal.Publish(config.DEFAULTChannelRedisParking, trx)
		log.Println("Publish Redis Local (Redis cloud not connected):", config.DEFAULTChannelRedisParking, "Status:", statusPublish.Val())
	}

	//} else {
	//statusPublish := svc.RedisClientLocal.Publish(config.DEFAULTChannelRedisParking, trxLocal)
	//log.Println("Publish Redis Local (No internet connection):", config.DEFAULTChannelRedisParking, "Status:", statusPublish.Val())
	//}

	return nil
}

func CheckInTrxDeposit(svc services.UsecaseService, request models.RequestTrxDepositCounter, docNo, qrCode string, merchantKey models.MerchantKey, outputProduct models.PolicyOuProductDepositCounterWithRules) error {
	log.Println("CheckInTrxDeposit with data:", constans.EMPTY_VALUE)

	var trxInvoiceItemList []models.TrxInvoiceDepositCounterItem
	trxItems := models.TrxInvoiceDepositCounterItem{
		DocNo:            docNo,
		ProductId:        outputProduct.ProductId,
		ProductCode:      outputProduct.ProductCode,
		ProductName:      outputProduct.ProductName,
		IsPctServiceFee:  outputProduct.IsPctServiceFee,
		ServiceFee:       outputProduct.ServiceFee,
		ServiceFeeMember: outputProduct.ServiceFeeMember,
		Price:            outputProduct.ProductOuDepositCounterWithRules.Price,
	}
	trxInvoiceItemList = append(trxInvoiceItemList, trxItems)
	trx := models.TrxCustom{
		DocNo:                        docNo,
		CheckInDatetime:              request.CheckInDatetime,
		ExtLocalDatetime:             utils.Timestamp(),
		OuId:                         merchantKey.OuId,
		OuName:                       merchantKey.OuName,
		OuCode:                       merchantKey.OuCode,
		ProductCode:                  outputProduct.ProductCode,
		OuSubBranchId:                merchantKey.OuSubBranchId,
		OuSubBranchName:              merchantKey.OuSubBranchName,
		OuSubBranchCode:              merchantKey.OuSubBranchCode,
		MainOuId:                     merchantKey.MainOuId,
		MainOuName:                   merchantKey.MainOuName,
		MainOuCode:                   merchantKey.MainOuCode,
		QrText:                       qrCode,
		TrxInvoiceItemDepositCounter: trxInvoiceItemList,
		MerchantKey:                  config.MERCHANT_KEY,
	}

	trxHeader, _ := json.Marshal(trx)
	//if utils.IsConnected() {
	statusPublish := svc.RedisClient.Publish(config.DEFAULTChannelRedisParkingDeposit, string(trxHeader))
	log.Println("Publish Redis Cloud:", config.DEFAULTChannelRedisParkingDeposit, "Status:", statusPublish.Val())

	if statusPublish.Val() == 0 {
		statusPublish = svc.RedisClientLocal.Publish(config.DEFAULTChannelRedisParkingDeposit, trx)
		log.Println("Publish Redis Local (Redis cloud not connected):", config.DEFAULTChannelRedisParkingDeposit, "Status:", statusPublish.Val())
	}

	//} else {
	//statusPublish := svc.RedisClientLocal.Publish(config.DEFAULTChannelRedisParking, trxLocal)
	//log.Println("Publish Redis Local (No internet connection):", config.DEFAULTChannelRedisParking, "Status:", statusPublish.Val())
	//}

	return nil
}

func CalculateAfterMember(resultTrxOutstanding models.ResultFindTrxOutstanding, endDateOfMember, inquirtDatetime, typePartner string) float64 {
	var grandTotal float64 = constans.EMPTY_VALUE_INT
	var overNightGrandTotal float64

	dateOffMember, _ := time.Parse("2006-01-02", endDateOfMember)
	checkouDate, _ := time.Parse("2006-01-02 15:04:05", inquirtDatetime)
	diffTime := utils.ConvDiffTime(dateOffMember, checkouDate)

	if typePartner == constans.TYPE_PARTNER_ONE_TIME {
		endDateParse, _ := time.Parse("2006-01-02 15:04", endDateOfMember)
		endDateFormat := endDateParse.Format("2006-01-02 15:04:05")
		dateOffMember, _ = time.Parse("2006-01-02 15:04:05", endDateFormat)
		diffTime = utils.ConvDiffTime(dateOffMember, checkouDate)
	}

	log.Println(utils.ToString(diffTime))
	if diffTime.Days >= constans.EMPTY_VALUE_INT {
		if diffTime.Hours >= constans.EMPTY_VALUE_INT {
			if diffTime.Minutes >= constans.EMPTY_VALUE_INT {
				grandTotal = resultTrxOutstanding.TrxInvoiceItem[0].Price
				dayDuration := (diffTime.Days * 24)
				d, _ := time.ParseDuration(fmt.Sprintf("%d%s", dayDuration, "h"))
				h, _ := time.ParseDuration(fmt.Sprintf("%d%s", diffTime.Hours, "h"))
				minute := h.Minutes() + float64(diffTime.Minutes) + d.Minutes()

				if resultTrxOutstanding.TrxInvoiceItem[0].ProgressivePrice != constans.EMPTY_VALUE_INT {
					duration := int64(minute) / resultTrxOutstanding.TrxInvoiceItem[0].ProgressiveTime
					log.Println(duration)
					grandTotal += (float64(duration) * resultTrxOutstanding.TrxInvoiceItem[0].ProgressivePrice)
				}

				if resultTrxOutstanding.TrxInvoiceItem[0].Is24H == constans.YES {
					if diffTime.Days >= constans.ONE_DAY_OVERNIGHT {
						overNightGrandTotal += (resultTrxOutstanding.TrxInvoiceItem[0].OvernightPrice * float64(diffTime.Days))
					}

					if diffTime.Hours >= constans.TWELVE_HOUR_OVERNIGHT {
						overNightGrandTotal += resultTrxOutstanding.TrxInvoiceItem[0].OvernightPrice
					}

					grandTotal += overNightGrandTotal
				} else if resultTrxOutstanding.TrxInvoiceItem[0].Is24H == constans.NO {
					if diffTime.Days >= constans.ONE_DAY_OVERNIGHT {
						overNightCount := math.Round(float64(dayDuration) / float64(12))
						overNightGrandTotal += (resultTrxOutstanding.TrxInvoiceItem[0].OvernightPrice * float64(overNightCount))
					}

					if diffTime.Hours >= constans.TWELVE_HOUR_OVERNIGHT {
						overNightGrandTotal += resultTrxOutstanding.TrxInvoiceItem[0].OvernightPrice
					}

					grandTotal += overNightGrandTotal
				} else {
					ovrNightTimeInt, _ := strconv.Atoi(resultTrxOutstanding.TrxInvoiceItem[0].OvernightTime)
					overnighTime := (1440 + ovrNightTimeInt) / 60
					log.Println("overnight:", overnighTime)
					overNightCount := math.Round(float64(dayDuration) / float64(overnighTime))
					overNightGrandTotal += (resultTrxOutstanding.TrxInvoiceItem[0].OvernightPrice * float64(overNightCount))

				}

				if resultTrxOutstanding.TrxInvoiceItem[0].MaxPrice > 0 && (grandTotal) >= resultTrxOutstanding.TrxInvoiceItem[0].MaxPrice {
					grandTotal = resultTrxOutstanding.TrxInvoiceItem[0].MaxPrice
				}

			}
		}
	}

	return grandTotal
}

func CallSyncConfirmTrxDepositCounterToCloud(ID *primitive.ObjectID, request models.RequestConfirmTrxDepositCounter, resultTrx models.TrxDepositCounter, svc services.UsecaseService, additionalData map[string]interface{}) error {
	requestDataOut, _ := json.Marshal(request)
	var trxInvoiceItemList []models.TrxInvoiceDepositCounterItem

	trxInvoiceItem := models.TrxInvoiceDepositCounterItem{
		DocNo:            resultTrx.TrxInvoiceItem[0].DocNo,
		ProductId:        resultTrx.TrxInvoiceItem[0].ProductId,
		ProductCode:      resultTrx.TrxInvoiceItem[0].ProductCode,
		ProductName:      resultTrx.TrxInvoiceItem[0].ProductName,
		IsPctServiceFee:  resultTrx.TrxInvoiceItem[0].IsPctServiceFee,
		ServiceFee:       resultTrx.TrxInvoiceItem[0].ServiceFee,
		ServiceFeeMember: resultTrx.TrxInvoiceItem[0].ServiceFeeMember,
		Price:            resultTrx.TrxInvoiceItem[0].Price,
		TotalAmount:      resultTrx.TrxInvoiceItem[0].TotalAmount,
	}

	trxInvoiceItemList = append(trxInvoiceItemList, trxInvoiceItem)
	trx := models.TrxCustom{
		DocNo:                        resultTrx.DocNoDepo,
		DocDate:                      utils.DateNow(),
		CheckInDatetime:              resultTrx.CheckInDatetime,
		CheckOutDatetime:             request.CheckOutDatetime,
		DeviceIdIn:                   request.DeviceId,
		DeviceId:                     request.DeviceId,
		GateIn:                       constans.EMPTY_VALUE,
		GateOut:                      constans.EMPTY_VALUE,
		CardNumberUUIDIn:             request.UUIDCard,
		CardNumberIn:                 request.CardNumber,
		CardNumberUUID:               request.UUIDCard,
		CardNumber:                   request.CardNumber,
		TypeCard:                     request.CardType,
		BeginningBalance:             constans.EMPTY_VALUE_INT,
		ExtLocalDatetime:             resultTrx.ExtLocalDatetime,
		GrandTotal:                   request.GrandTotal,
		ProductCode:                  resultTrx.ProductCode,
		ProductName:                  resultTrx.ProductName,
		ProductData:                  resultTrx.ProductData,
		RequestData:                  resultTrx.RequestData,
		RequestOutData:               string(requestDataOut),
		OuId:                         resultTrx.OuId,
		OuName:                       resultTrx.OuName,
		OuCode:                       resultTrx.OuCode,
		OuSubBranchId:                resultTrx.OuSubBranchId,
		OuSubBranchName:              resultTrx.OuSubBranchName,
		OuSubBranchCode:              resultTrx.OuSubBranchCode,
		MainOuId:                     resultTrx.MainOuId,
		MainOuCode:                   resultTrx.MainOuCode,
		MainOuName:                   resultTrx.MainOuName,
		LogTrans:                     request.LogTrans,
		MerchantKey:                  config.MERCHANT_KEY,
		QrText:                       resultTrx.QrText,
		TrxInvoiceItemDepositCounter: trxInvoiceItemList,
		FlgDepositCounter:            true,
	}

	resultTrx.CheckOutDatetime = request.CheckOutDatetime
	resultTrx.CardNumberUUID = request.UUIDCard
	resultTrx.CardNumber = request.CardNumber
	resultTrx.TypeCard = request.CardType
	resultTrx.DeviceId = request.DeviceId
	resultTrx.ProductCode = request.ProductCode
	resultTrx.ProductName = request.ProductName
	resultTrx.RequestOutData = string(requestDataOut)
	resultTrx.LogTrans = request.LogTrans
	resultTrx.GrandTotal = request.GrandTotal
	resultTrx.DocDateDepo = utils.DateNow()
	resultTrx.MerchantKey = config.MERCHANT_KEY
	resultTrx.UsernameOut = request.Username
	resultTrx.ShiftCodeOut = request.ShiftCode

	if request.GrandTotal == 0 || request.CardType == constans.SETTLEMENT_CODE_CASH {
		resultTrx.TypeCard = constans.SETTLEMENT_CODE_CASH
		resultTrx.DeviceId = constans.EMPTY_VALUE
		trx.DeviceId = constans.EMPTY_VALUE
		trx.TypeCard = constans.SETTLEMENT_CODE_CASH
	}

	config.NORMAL_SF = constans.YES
	if request.GrandTotal == constans.EMPTY_VALUE_INT {
		config.NORMAL_SF = constans.NO
	}

	filter := bson.M{
		"docNo": resultTrx.DocNoDepo,
	}

	updateSet := bson.M{
		"$set": bson.M{
			"username":  additionalData["username"],
			"shiftCode": request.ShiftCode,
		},
	}

	err := svc.TrxMongoRepo.AddTrxAddInfoInterfaces(filter, updateSet)
	if err != nil {
		log.Println("ERROR AddTrxAddInfoInterfaces : ", err.Error())
		return err
	}

	// Digunakan utk mengambil data terbaru setelah diupdate
	resultTrxAddInfo, exists := svc.TrxMongoRepo.IsTrxAddInfoInterfacesExistsByDocNo(resultTrx.DocNoDepo)
	if exists {
		trx.TrxAddInfo = resultTrxAddInfo
	}

	dataStr, _ := json.Marshal(trx)
	log.Println("SETTING SFEE", config.NORMAL_SF)
	trx.FlagSyncData = true
	resultTrx.FlagSyncData = true
	if utils.IsConnected() {
		if config.NORMAL_SF == constans.YES {
			redisStatus := svc.RedisClient.Publish(constans.PG_PARKING_CHECKOUT_DEPOSIT_COUNTER, dataStr)
			if redisStatus.Err() != nil || redisStatus.Val() == 0 {
				trx.FlagSyncData = false
				resultTrx.FlagSyncData = false
			}
		} else if config.NORMAL_SF == constans.NO {
			redisStatus := svc.RedisClient.Publish(constans.PG_PARKING_CHECKOUT_DEPOSIT_COUNTER_EXCLUDE_SF, dataStr)
			if redisStatus.Err() != nil || redisStatus.Val() == 0 {
				trx.FlagSyncData = false
				resultTrx.FlagSyncData = false
			}
		}

	} else {
		trx.FlagSyncData = false
		resultTrx.FlagSyncData = false
	}

	_, err = svc.TrxMongoDepositCounterRepo.AddTrxDepositCounter(resultTrx)
	if err != nil {
		log.Println("ERROR AddTrxDepositCounter : ", err.Error())
		return err
	}

	if ID != nil {
		err = svc.TrxMongoDepositCounterRepo.RemoveTrxDepositCounterByID(*ID)
		if err != nil {
			log.Println("ERROR RemoveTrxDepositCounterByID : ", err.Error())
			return err
		}
	}

	redisStatus := svc.RedisClientLocal.Del(request.ID)
	if redisStatus.Err() != nil {
		return redisStatus.Err()
	}

	return nil
}

func CallSyncConfirmTrxLostTicket(ID *primitive.ObjectID, request models.RequestConfirmLostTicket, lostTicket models.LostTicket, svc services.UsecaseService) error {
	requestDataOut, _ := json.Marshal(request)

	lostTicket.LostTicketOutDatetime = utils.Timestamp()
	lostTicket.TrxOutstanding.CardNumberUUID = request.UUIDCard
	lostTicket.TrxOutstanding.CardNumber = request.CardNumber
	lostTicket.TrxOutstanding.TypeCard = request.CardType
	lostTicket.TrxOutstanding.DeviceId = request.DeviceId
	lostTicket.TrxOutstanding.GateOut = request.IpTerminal
	lostTicket.TrxOutstanding.ProductCode = request.ProductCode
	lostTicket.TrxOutstanding.ProductName = request.ProductName
	lostTicket.TrxOutstanding.RequestOutData = string(requestDataOut)
	lostTicket.TrxOutstanding.LogTrans = request.LogTrans
	lostTicket.TrxOutstanding.GrandTotal = request.GrandTotal
	lostTicket.TrxOutstanding.DocDate = utils.DateNow()
	lostTicket.TrxOutstanding.MerchantKey = config.MERCHANT_KEY
	lostTicket.TrxOutstanding.ChargeType = lostTicket.Type

	if request.GrandTotal == 0 || request.CardType == constans.SETTLEMENT_CODE_CASH {
		lostTicket.TrxOutstanding.TypeCard = constans.SETTLEMENT_CODE_CASH
		lostTicket.TrxOutstanding.DeviceId = constans.EMPTY_VALUE
	}

	filter := bson.M{
		"docNo": lostTicket.DocNo,
	}

	updateSet := bson.M{
		"$set": bson.M{
			"username":  request.Username,
			"shiftCode": request.ShiftCode,
		},
	}

	err := svc.TrxMongoRepo.AddTrxAddInfoInterfaces(filter, updateSet)
	if err != nil {
		log.Println("ERROR AddTrxAddInfoInterfaces : ", err.Error())
		return err
	}

	// Digunakan utk mengambil data terbaru setelah diupdate
	resultTrxAddInfo, exists := svc.TrxMongoRepo.IsTrxAddInfoInterfacesExistsByDocNo(lostTicket.DocNo)
	if exists {
		lostTicket.TrxOutstanding.TrxAddInfo = resultTrxAddInfo
	}

	err = svc.ReportTrxMongoRepo.UpdateTypeTrxLostTicket(*ID, "CONFIRM", lostTicket)
	if err != nil {
		log.Println("ERROR UpdateTypeTrxLostTicket : ", err.Error())
		return err
	}

	dataTrx := utils.ToString(lostTicket.TrxOutstanding)
	redisStatus := svc.RedisClientLocal.Set(fmt.Sprintf("%s-%s", lostTicket.DocNo, "LOST_TICKET"), dataTrx, 15*time.Minute)
	if redisStatus.Err() != nil {
		log.Println("ERROR Set Redis Client : ", redisStatus.Err().Error())
		return err
	}

	return nil
}

func CallSyncConfirmTrxForLostTicket(request models.RequestConfirmTrx, resultTrx models.Trx, svc services.UsecaseService) error {
	requestDataOut, _ := json.Marshal(request)

	checkoutDatetimeParse, _ := time.Parse("2006-01-02 15:04:05", request.CheckOutDatetime)
	timeCheckOutUnix := checkoutDatetimeParse.Unix()
	durationTime := (timeCheckOutUnix - resultTrx.CheckInTime) / 60

	resultTrx.CheckOutDatetime = request.CheckOutDatetime
	resultTrx.CheckOutTime = timeCheckOutUnix
	resultTrx.DurationTime = durationTime
	resultTrx.CardNumberUUID = request.UUIDCard
	resultTrx.CardNumber = request.CardNumber
	resultTrx.TypeCard = request.CardType
	resultTrx.DeviceId = request.DeviceId
	resultTrx.GateOut = request.IpTerminal
	resultTrx.ProductCode = request.ProductCode
	resultTrx.ProductName = request.ProductName
	resultTrx.RequestOutData = string(requestDataOut)
	resultTrx.LogTrans = request.LogTrans
	resultTrx.GrandTotal += request.GrandTotal
	resultTrx.DocDate = utils.DateNow()
	resultTrx.MerchantKey = config.MERCHANT_KEY
	resultTrx.FlagCharge = constans.TRUE_VALUE

	config.NORMAL_SF = constans.YES
	if request.GrandTotal == constans.EMPTY_VALUE_INT {
		config.NORMAL_SF = constans.NO
	}

	dataStr, _ := json.Marshal(resultTrx)
	log.Println("SETTING SFEE", config.NORMAL_SF)
	resultTrx.FlagSyncData = true
	log.Println(utils.ToString(resultTrx))
	if utils.IsConnected() {
		if config.NORMAL_SF == constans.YES {
			redisStatus := svc.RedisClient.Publish(constans.CHANNEL_REDIS_PG_PARKING_CHECKOUT, dataStr)
			if redisStatus.Err() != nil || redisStatus.Val() == 0 {
				resultTrx.FlagSyncData = false
			}
		} else if config.NORMAL_SF == constans.NO {
			redisStatus := svc.RedisClient.Publish(constans.CHANNEL_REDIS_PG_PARKING_CHECKOUT_EXCLUDE_SF, dataStr)
			if redisStatus.Err() != nil || redisStatus.Val() == 0 {
				resultTrx.FlagSyncData = false
			}
		}

	} else {
		resultTrx.FlagSyncData = false
	}

	_, err := svc.TrxMongoRepo.AddTrx(resultTrx)
	if err != nil {
		log.Println("ERROR AddTrx : ", err.Error())
		return err
	}

	err = svc.TrxMongoRepo.RemoveTrxByDocNo(resultTrx.DocNo)
	if err != nil {
		log.Println("ERROR RemoveTrxByDocNo : ", err.Error())
		return err
	}

	redisStatus := svc.RedisClientLocal.Del(fmt.Sprintf("%s#%s", "LOST-TICKET", resultTrx.DocNo))
	if redisStatus.Err() != nil {
		return redisStatus.Err()
	}

	return nil
}

func CalculateFee(resultTrxOutstanding models.ResultFindTrxOutstanding, checkInDate, inquirtDatetime string) float64 {
	var grandTotal float64 = constans.EMPTY_VALUE_INT
	var overNightGrandTotal float64
	var baseTime int64
	var gracePeriod int64
	var overnight float64

	checkinDate, _ := time.Parse("2006-01-02 15:04:05", checkInDate)
	checkouDate, _ := time.Parse("2006-01-02 15:04:05", inquirtDatetime)
	diffTime := utils.ConvDiffTime(checkinDate, checkouDate)

	if diffTime.Days >= constans.EMPTY_VALUE_INT {
		if diffTime.Hours >= constans.EMPTY_VALUE_INT {
			if diffTime.Minutes >= constans.EMPTY_VALUE_INT {
				grandTotal = resultTrxOutstanding.TrxInvoiceItem[0].Price
				dayDuration := (diffTime.Days * 24)

				d, _ := time.ParseDuration(fmt.Sprintf("%d%s", dayDuration, "h"))
				h, _ := time.ParseDuration(fmt.Sprintf("%d%s", diffTime.Hours, "h"))
				minute := h.Minutes() + float64(diffTime.Minutes) + d.Minutes()

				if resultTrxOutstanding.TrxInvoiceItem[0].ProgressivePrice != constans.EMPTY_VALUE_INT {

					dayDuration += diffTime.Hours
					if resultTrxOutstanding.TrxInvoiceItem[0].GracePeriod != constans.EMPTY_VALUE_INT {
						gracePeriod = resultTrxOutstanding.TrxInvoiceItem[0].GracePeriod
					}

					if resultTrxOutstanding.TrxInvoiceItem[0].BaseTime != constans.EMPTY_VALUE_INT {
						baseTime = resultTrxOutstanding.TrxInvoiceItem[0].BaseTime
					}

					if resultTrxOutstanding.TrxInvoiceItem[0].Is24H == "C" {
						ovrNightTimeInt, _ := strconv.Atoi(resultTrxOutstanding.TrxInvoiceItem[0].OvernightTime)
						overnighTime := (1440 + ovrNightTimeInt) / 60
						overNightCount := math.Floor(float64(dayDuration) / float64(overnighTime))
						overnight = overNightCount

						if overNightCount >= 1 {
							minutePlus := float64(diffTime.Minutes) + d.Minutes()
							if resultTrxOutstanding.TrxInvoiceItem[0].GracePeriod != constans.EMPTY_VALUE_INT {
								if int64(minutePlus) > resultTrxOutstanding.TrxInvoiceItem[0].GracePeriod {
									gracePeriod = resultTrxOutstanding.TrxInvoiceItem[0].GracePeriod * int64(overNightCount)
								}
							}

							if resultTrxOutstanding.TrxInvoiceItem[0].BaseTime != constans.EMPTY_VALUE_INT {
								if int64(minutePlus) > resultTrxOutstanding.TrxInvoiceItem[0].BaseTime {
									baseTime = resultTrxOutstanding.TrxInvoiceItem[0].BaseTime * int64(overNightCount)
								}
							}

							grandTotal = overNightCount * resultTrxOutstanding.TrxInvoiceItem[0].Price
						}

					}

					minute -= float64(baseTime)
					minute -= float64(gracePeriod)
					duration := int64(minute) / resultTrxOutstanding.TrxInvoiceItem[0].ProgressiveTime
					grandTotal += (float64(duration) * resultTrxOutstanding.TrxInvoiceItem[0].ProgressivePrice)
				}

				if resultTrxOutstanding.TrxInvoiceItem[0].Is24H != constans.EMPTY_VALUE {
					if resultTrxOutstanding.TrxInvoiceItem[0].Is24H == constans.YES {
						overNightGrandTotal += (resultTrxOutstanding.TrxInvoiceItem[0].OvernightPrice * float64(diffTime.Days))
						grandTotal += overNightGrandTotal
					} else if resultTrxOutstanding.TrxInvoiceItem[0].Is24H == constans.NO {
						overNightCount := math.Floor(float64(dayDuration) / float64(12))
						overNightGrandTotal += (resultTrxOutstanding.TrxInvoiceItem[0].OvernightPrice * float64(overNightCount))
						grandTotal += overNightGrandTotal
					} else {
						overnightPrice := resultTrxOutstanding.TrxInvoiceItem[0].OvernightPrice
						ovrNightTimeInt, _ := strconv.Atoi(resultTrxOutstanding.TrxInvoiceItem[0].OvernightTime)
						overnighTime := (1440 + ovrNightTimeInt) / 60
						overNightCount := math.Floor(float64(dayDuration) / float64(overnighTime))
						overNightGrandTotal += (overnightPrice * float64(overNightCount))
						grandTotal += overNightGrandTotal
					}
				}

				if resultTrxOutstanding.TrxInvoiceItem[0].MaxPrice > 0 && (grandTotal) >= resultTrxOutstanding.TrxInvoiceItem[0].MaxPrice {
					grandTotal = resultTrxOutstanding.TrxInvoiceItem[0].MaxPrice
				}

			}
		}
	}

	if overnight < 1 {
		grandTotal = resultTrxOutstanding.TrxInvoiceItem[0].TotalAmount
	}

	return grandTotal
}

func CalculateInvoiceDetail(resultTrxOutstanding models.ResultFindTrxOutstanding, checkInDate, inquirtDatetime string) (float64, float64) {
	var grandTotal float64 = constans.EMPTY_VALUE_INT
	var overNightGrandTotal float64
	var baseTime int64
	var gracePeriod int64

	checkinDate, _ := time.Parse("2006-01-02 15:04:05", checkInDate)
	checkouDate, _ := time.Parse("2006-01-02 15:04:05", inquirtDatetime)
	diffTime := utils.ConvDiffTime(checkinDate, checkouDate)

	if diffTime.Days >= constans.EMPTY_VALUE_INT {
		if diffTime.Hours >= constans.EMPTY_VALUE_INT {
			if diffTime.Minutes >= constans.EMPTY_VALUE_INT {
				grandTotal = resultTrxOutstanding.TrxInvoiceItem[0].Price
				dayDuration := (diffTime.Days * 24)

				d, _ := time.ParseDuration(fmt.Sprintf("%d%s", dayDuration, "h"))
				h, _ := time.ParseDuration(fmt.Sprintf("%d%s", diffTime.Hours, "h"))
				minute := h.Minutes() + float64(diffTime.Minutes) + d.Minutes()

				dayDuration += diffTime.Hours
				if resultTrxOutstanding.TrxInvoiceItem[0].ProgressivePrice != constans.EMPTY_VALUE_INT {

					if resultTrxOutstanding.TrxInvoiceItem[0].GracePeriod != constans.EMPTY_VALUE_INT {
						gracePeriod = resultTrxOutstanding.TrxInvoiceItem[0].GracePeriod
					}

					if resultTrxOutstanding.TrxInvoiceItem[0].BaseTime != constans.EMPTY_VALUE_INT {
						baseTime = resultTrxOutstanding.TrxInvoiceItem[0].BaseTime
					}

					if resultTrxOutstanding.TrxInvoiceItem[0].Is24H == "C" {
						ovrNightTimeInt, _ := strconv.Atoi(resultTrxOutstanding.TrxInvoiceItem[0].OvernightTime)
						overnighTime := (1440 + ovrNightTimeInt) / 60
						overNightCount := math.Floor(float64(dayDuration) / float64(overnighTime))

						if overNightCount >= 1 {
							minutePlus := float64(diffTime.Minutes) + d.Minutes()
							if resultTrxOutstanding.TrxInvoiceItem[0].GracePeriod != constans.EMPTY_VALUE_INT {
								if int64(minutePlus) > resultTrxOutstanding.TrxInvoiceItem[0].GracePeriod {
									gracePeriod = resultTrxOutstanding.TrxInvoiceItem[0].GracePeriod * int64(overNightCount)
								}
							}

							if resultTrxOutstanding.TrxInvoiceItem[0].BaseTime != constans.EMPTY_VALUE_INT {
								if int64(minutePlus) > resultTrxOutstanding.TrxInvoiceItem[0].BaseTime {
									baseTime = resultTrxOutstanding.TrxInvoiceItem[0].BaseTime * int64(overNightCount)
								}
							}

							grandTotal = overNightCount * resultTrxOutstanding.TrxInvoiceItem[0].Price
						}

					}

					minute -= float64(baseTime)
					minute -= float64(gracePeriod)
					duration := int64(minute) / resultTrxOutstanding.TrxInvoiceItem[0].ProgressiveTime
					grandTotal += (float64(duration) * resultTrxOutstanding.TrxInvoiceItem[0].ProgressivePrice)
				}

				if resultTrxOutstanding.TrxInvoiceItem[0].Is24H != constans.EMPTY_VALUE {
					if resultTrxOutstanding.TrxInvoiceItem[0].Is24H == constans.YES {
						overNightGrandTotal += (resultTrxOutstanding.TrxInvoiceItem[0].OvernightPrice * float64(diffTime.Days))
					} else if resultTrxOutstanding.TrxInvoiceItem[0].Is24H == constans.NO {
						overNightCount := math.Floor(float64(dayDuration) / float64(12))
						overNightGrandTotal += (resultTrxOutstanding.TrxInvoiceItem[0].OvernightPrice * float64(overNightCount))
					} else {
						overnightPrice := resultTrxOutstanding.TrxInvoiceItem[0].OvernightPrice
						ovrNightTimeInt, _ := strconv.Atoi(resultTrxOutstanding.TrxInvoiceItem[0].OvernightTime)
						overnighTime := (1440 + ovrNightTimeInt) / 60
						log.Println("OVERNIGHT TIME:", overnighTime)
						log.Println("Day Duration: ", dayDuration)
						overNightCount := math.Floor(float64(dayDuration) / float64(overnighTime))
						log.Println("Overnight Count :", overNightCount)
						overNightGrandTotal += (overnightPrice * float64(overNightCount))
					}
				}

				if resultTrxOutstanding.TrxInvoiceItem[0].MaxPrice > 0 && (grandTotal) >= resultTrxOutstanding.TrxInvoiceItem[0].MaxPrice {
					grandTotal = resultTrxOutstanding.TrxInvoiceItem[0].MaxPrice
					overNightGrandTotal = constans.EMPTY_VALUE_INT
					if resultTrxOutstanding.TrxInvoiceItem[0].OvernightTime > constans.EMPTY_VALUE {
						if resultTrxOutstanding.TrxInvoiceItem[0].OvernightPrice > constans.EMPTY_VALUE_INT {
							overNightGrandTotal = resultTrxOutstanding.TrxInvoiceItem[0].OvernightPrice
						}
					}
				}

			}
		}
	}

	return grandTotal, overNightGrandTotal
}
