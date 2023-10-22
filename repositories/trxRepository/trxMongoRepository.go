package trxRepository

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
	"ppkgwlocal/constans"
	"ppkgwlocal/models"
	"ppkgwlocal/repositories"
)

type trxMongoRepository struct {
	RepoDB repositories.Repository
}

func NewTrxMongoRepository(repo repositories.Repository) trxMongoRepository {
	return trxMongoRepository{
		RepoDB: repo,
	}
}

func (ctx trxMongoRepository) AddTrxCheckIn(trx models.Trx) (*primitive.ObjectID, error) {
	var ID primitive.ObjectID

	result, err := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_COLLECTION).
		InsertOne(ctx.RepoDB.Context, trx)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		ID = oid
	}

	return &ID, nil
}

func (ctx trxMongoRepository) AddTrx(trx models.Trx) (string, error) {
	_, err := ctx.RepoDB.MongoDB.Collection(constans.TRX_COLLECTION).
		InsertOne(ctx.RepoDB.Context, trx)
	if err != nil {
		return "", err
	}
	return trx.DocNo, nil
}

func (ctx trxMongoRepository) ValTrxExistByUUIDCard(uuidCard string) (string, error) {

	filter := bson.M{
		"cardNumberUuid": uuidCard,
	}

	count, err := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_COLLECTION).
		CountDocuments(ctx.RepoDB.Context, filter)
	if err != nil {
		return constans.MALFUNCTION_SYSTEM_CODE, err
	}

	if count > 0 {
		return constans.CARD_ALREADY_USED_CODE, errors.New("Sesi Kartu Masih Digunakan")
	}

	return constans.SUCCESS_CODE, nil
}

func (ctx trxMongoRepository) IsTrxOutstandingByDocNoForCustom(docNo string) (trx *models.Trx, exists bool, err error) {
	var trxList []models.Trx

	filter := bson.M{
		"docNo": docNo,
	}

	collect := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_COLLECTION)
	data, err := collect.Find(ctx.RepoDB.Context, filter)
	if err != nil {
		return nil, false, err
	}

	for data.Next(ctx.RepoDB.Context) {
		var val models.Trx
		if err = data.Decode(&val); err != nil {
			return nil, false, err
		}

		trxList = append(trxList, val)
	}
	data.Close(ctx.RepoDB.Context)

	if len(trxList) == 0 {
		return nil, false, nil
	}

	return &trxList[0], true, nil
}

func (ctx trxMongoRepository) IsTrxOutstandingByDocNoForLostTicket(docNo string) (trx *models.TrxWithId, exists bool, err error) {
	var trxList []models.TrxWithId

	filter := bson.M{
		"docNo": docNo,
	}

	collect := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_COLLECTION)
	data, err := collect.Find(ctx.RepoDB.Context, filter)
	if err != nil {
		return nil, false, err
	}

	for data.Next(ctx.RepoDB.Context) {
		var val models.TrxWithId
		if err = data.Decode(&val); err != nil {
			return nil, false, err
		}

		trxList = append(trxList, val)
	}
	data.Close(ctx.RepoDB.Context)

	if len(trxList) == 0 {
		return nil, false, nil
	}

	return &trxList[0], true, nil
}

func (ctx trxMongoRepository) FindTrxOutstandingByDocNo(docNo string, productCode string) (models.ResultFindTrxOutstanding, error) {
	var result models.ResultFindTrxOutstanding
	filter := make(map[string]interface{})

	if productCode != constans.EMPTY_VALUE {
		filter = bson.M{
			"docNo": docNo,
			"trxInvoiceItem": bson.M{
				"$elemMatch": bson.M{
					"docNo":       docNo,
					"productCode": productCode,
				},
			},
		}
	} else {
		filter = bson.M{
			"docNo": docNo,
			"trxInvoiceItem": bson.M{
				"$elemMatch": bson.M{
					"docNo": docNo,
				},
			},
		}
	}

	projection := bson.M{
		"_id":              1,
		"docNo":            1,
		"grandTotal":       1,
		"checkInDatetime":  1,
		"overnightPrice":   1,
		"is24H":            1,
		"ouCode":           1,
		"vehicleNumberIn":  1,
		"trxInvoiceItem.$": 1,
	}

	err := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_COLLECTION).
		FindOne(ctx.RepoDB.Context, filter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (ctx trxMongoRepository) IsTrxOutstandingByDocNo(docNo string, productCode string) (*models.ResultFindTrxOutstanding, bool, error) {
	var result []models.ResultFindTrxOutstanding
	filter := make(map[string]interface{})

	if productCode != constans.EMPTY_VALUE {
		filter = bson.M{
			"docNo": docNo,
			"trxInvoiceItem": bson.M{
				"$elemMatch": bson.M{
					"docNo":       docNo,
					"productCode": productCode,
				},
			},
		}
	} else {
		filter = bson.M{
			"docNo": docNo,
			"trxInvoiceItem": bson.M{
				"$elemMatch": bson.M{
					"docNo": docNo,
				},
			},
		}
	}

	projection := bson.M{
		"_id":              1,
		"docNo":            1,
		"grandTotal":       1,
		"checkInDatetime":  1,
		"overnightPrice":   1,
		"is24H":            1,
		"ouCode":           1,
		"vehicleNumberIn":  1,
		"trxInvoiceItem.$": 1,
	}

	collect := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_COLLECTION)
	data, err := collect.Find(ctx.RepoDB.Context, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, false, err
	}

	for data.Next(ctx.RepoDB.Context) {
		var val models.ResultFindTrxOutstanding
		if err = data.Decode(&val); err != nil {
			return nil, false, err
		}

		result = append(result, val)
	}

	if len(result) == 0 {
		return nil, false, nil
	}

	return &result[0], true, nil
}

func (ctx trxMongoRepository) FindTrxOutstandingByUUID(uuid string, productCode string) (models.ResultFindTrxOutstanding, error) {
	var result models.ResultFindTrxOutstanding
	filter := make(map[string]interface{})

	if productCode != constans.EMPTY_VALUE {
		filter = bson.M{
			"cardNumberUuid": uuid,
			"trxInvoiceItem": bson.M{
				"$elemMatch": bson.M{
					"productCode": productCode,
				},
			},
		}
	} else {
		filter = bson.M{
			"cardNumberUuid": uuid,
		}
	}

	projection := bson.M{
		"_id":              1,
		"docNo":            1,
		"grandTotal":       1,
		"checkInDatetime":  1,
		"overnightPrice":   1,
		"is24H":            1,
		"cardNumber":       1,
		"cardNumberUuid":   1,
		"ouCode":           1,
		"vehicleNumberIn":  1,
		"trxInvoiceItem.$": 1,
	}

	err := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_COLLECTION).
		FindOne(ctx.RepoDB.Context, filter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (ctx trxMongoRepository) IsTrxOutstandingByUUID(uuidCard string) (*models.Trx, bool, error) {
	var result []models.Trx

	filter := bson.M{
		"cardNumberUuid": uuidCard,
	}

	collect := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_COLLECTION)
	data, err := collect.Find(ctx.RepoDB.Context, filter)
	if err != nil {
		return nil, false, err
	}
	defer data.Close(ctx.RepoDB.Context)

	for data.Next(ctx.RepoDB.Context) {
		var val models.Trx
		if err = data.Decode(&val); err != nil {
			return nil, false, err
		}

		result = append(result, val)
	}

	if len(result) == 0 {
		return nil, false, nil
	}

	return &result[0], true, nil
}

func (ctx trxMongoRepository) IsTrxOutstandingByCardNumber(cardNumber string) (*models.Trx, bool, error) {
	var result []models.Trx

	filter := bson.M{
		"cardNumber": cardNumber,
	}

	collect := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_COLLECTION)
	data, err := collect.Find(ctx.RepoDB.Context, filter)
	if err != nil {
		return nil, false, err
	}
	defer data.Close(ctx.RepoDB.Context)

	for data.Next(ctx.RepoDB.Context) {
		var val models.Trx
		if err = data.Decode(&val); err != nil {
			return nil, false, err
		}

		result = append(result, val)
	}

	if len(result) == 0 {
		return nil, false, nil
	}

	return &result[0], true, nil
}

func (ctx trxMongoRepository) FindTrxOutstandingByID(ID primitive.ObjectID) (models.Trx, error) {
	var result models.Trx

	filter := bson.M{
		"_id": ID,
	}

	err := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_COLLECTION).
		FindOne(ctx.RepoDB.Context, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (ctx trxMongoRepository) UpdateProgressivePriceByDocNo(docNo string, progressivePrice float64) error {
	var update interface{}
	filter := bson.M{
		"docNo": docNo,
	}

	update = bson.M{
		"$set": bson.M{"grandTotal": progressivePrice},
	}

	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_COLLECTION).
		UpdateOne(ctx.RepoDB.Context, filter, update, &opt)
	if err != nil {
		return err
	}

	return nil
}

func (ctx trxMongoRepository) RemoveTrxByID(ID primitive.ObjectID) error {

	filter := bson.M{
		"_id": ID,
	}

	data := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_COLLECTION).
		FindOneAndDelete(ctx.RepoDB.Context, filter)
	if data.Err() != nil {
		return data.Err()
	}

	return nil
}

func (ctx trxMongoRepository) RemoveTrxByDocNo(docNo string) error {

	filter := bson.M{
		"docNo": docNo,
	}

	opt := &options.DeleteOptions{}

	_, err := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_COLLECTION).DeleteOne(ctx.RepoDB.Context, filter, opt)
	if err != nil {
		return err
	}

	return nil
}

func (ctx trxMongoRepository) LogActivity(logActivity models.LogActivityTrx) error {
	_, err := ctx.RepoDB.MongoDB.Collection("LogActivity").
		InsertOne(ctx.RepoDB.Context, logActivity)
	if err != nil {
		return err
	}
	return nil
}

func (ctx trxMongoRepository) InvoiceTrx(invoiceTrx models.InvoiceTrx) error {
	_, err := ctx.RepoDB.MongoDB.Collection("InvoiceTrx").
		InsertOne(ctx.RepoDB.Context, invoiceTrx)
	if err != nil {
		return err
	}
	return nil
}

func (ctx trxMongoRepository) UpdateTrxInvoiceItemForTrxOutstanding(docNo string, productCode string, trxInvoiceItems models.TrxInvoiceItem) error {

	filter := bson.M{
		"docNo": docNo,
	}

	update := bson.M{
		"$set": bson.M{
			"trxInvoiceItem.$[elemX].totalAmount": trxInvoiceItems.TotalAmount,
		},
	}

	_, err := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_COLLECTION).
		UpdateOne(ctx.RepoDB.Context, filter, update, options.Update().SetArrayFilters(
			options.ArrayFilters{Filters: []interface{}{
				bson.M{
					"elemX.docNo":       docNo,
					"elemX.productCode": productCode,
				},
			},
			},
		))

	if err != nil {
		return err
	}

	return nil
}

func (ctx trxMongoRepository) IsTrxOutstandingExistByUUIDCard(uuidCard string) (models.Trx, bool) {
	var result models.Trx

	filter := bson.M{
		"cardNumberUuid": uuidCard,
	}

	opt := options.FindOneOptions{
		Sort: bson.M{"dateFrom": 1},
	}

	err := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_COLLECTION).
		FindOne(ctx.RepoDB.Context, filter, &opt).Decode(&result)
	if err != nil {
		return result, false
	}

	return result, true
}

func (ctx trxMongoRepository) GetTrxMemberByPeriodList(docNo string, inquiryDate string) ([]models.TrxMember, error) {
	var result []models.TrxMember

	filter := bson.M{
		"docNo":    docNo,
		"dateFrom": bson.M{"$lte": inquiryDate},
	}

	data, err := ctx.RepoDB.MongoDB.Collection(constans.TRX_MEMBER_COLLECTIONS, options.Collection()).Find(ctx.RepoDB.Context, filter)
	if err != nil {
		return result, err
	}

	if err := data.All(ctx.RepoDB.Context, &result); err != nil {
		return nil, err
	}

	err = data.Close(ctx.RepoDB.Context)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ctx trxMongoRepository) GetTrxMemberByPeriodListByProductCode(docNo, productCode, inquiryDate string) ([]models.TrxMember, error) {
	var result []models.TrxMember

	filter := bson.M{
		"docNo":       docNo,
		"dateFrom":    bson.M{"$lte": inquiryDate},
		"productCode": productCode,
	}

	data, err := ctx.RepoDB.MongoDB.Collection(constans.TRX_MEMBER_COLLECTIONS, options.Collection()).Find(ctx.RepoDB.Context, filter)
	if err != nil {
		return result, err
	}

	if err := data.All(ctx.RepoDB.Context, &result); err != nil {
		return nil, err
	}

	err = data.Close(ctx.RepoDB.Context)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ctx trxMongoRepository) GetTrxInvoiceDetailsItemsLists(trxMemberList []models.TrxMember) (float64, error) {
	var totalAmount float64

	for _, rows := range trxMemberList {
		if rows.RoleType != constans.TYPE_PARTNER_FREE_PASS {
			var invoiceTrx []models.ResponseTrxInvoiceDetailsItemList

			pipeline := []bson.M{
				bson.M{"$match": bson.M{
					"$and": []bson.M{
						bson.M{
							"createdDate": bson.M{"$gte": rows.DateFrom, "$lte": rows.DateTo},
						},
						bson.M{
							"docNo":       rows.DocNo,
							"productCode": rows.ProductCode,
						},
					},
				},
				}, {
					"$group": bson.M{
						"_id":           nil,
						"invoiceAmount": bson.M{"$sum": "$invoiceAmount"},
					},
				},
			}

			data, err := ctx.RepoDB.MongoDB.Collection(constans.TRX_INVOICE_DETAIL_ITEM_COLLECTIONS).
				Aggregate(ctx.RepoDB.Context, pipeline)
			if err != nil {
				return totalAmount, err
			}

			for data.Next(ctx.RepoDB.Context) {
				var val models.ResponseTrxInvoiceDetailsItemList
				err := data.Decode(&val)
				if err != nil {
					return 0, err
				}

				invoiceTrx = append(invoiceTrx, val)
			}

			err = data.Close(ctx.RepoDB.Context)
			if err != nil {
				return 0, err
			}

			if len(invoiceTrx) > 0 {
				totalAmount += invoiceTrx[0].TotalAmount
			}
		}

	}

	return totalAmount, nil
}

func (ctx trxMongoRepository) GetTrxInvoiceDetailsItemsList(docNo, productCode string) ([]models.TrxInvoiceDetailItem, error) {
	var invoiceTrxDetailItemList []models.TrxInvoiceDetailItem

	filter := bson.M{
		"docNo":       docNo,
		"productCode": productCode,
	}

	data, err := ctx.RepoDB.MongoDB.Collection(constans.TRX_INVOICE_DETAIL_ITEM_COLLECTIONS).Find(ctx.RepoDB.Context, filter)
	if err != nil {
		return nil, err
	}

	if err = data.All(ctx.RepoDB.Context, &invoiceTrxDetailItemList); err != nil {
		return nil, err
	}

	err = data.Close(ctx.RepoDB.Context)
	if err != nil {
		return nil, err
	}

	return invoiceTrxDetailItemList, nil
}

func (ctx trxMongoRepository) AddTrxAddInfoInterfaces(filter interface{}, updateSet interface{}) error {
	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := ctx.RepoDB.MongoDB.Collection(constans.TRX_ADD_INFO_COLLECTION).
		UpdateOne(ctx.RepoDB.Context, filter, updateSet, &opt)
	if err != nil {
		return err
	}

	return nil
}

func (ctx trxMongoRepository) IsTrxAddInfoInterfacesExistsByDocNo(docNo string) (map[string]interface{}, bool) {
	var result map[string]interface{}
	filter := bson.M{
		"docNo": docNo,
	}

	err := ctx.RepoDB.MongoDB.Collection(constans.TRX_ADD_INFO_COLLECTION).
		FindOne(ctx.RepoDB.Context, filter).Decode(&result)
	if err != nil {
		return nil, false
	}

	return result, true
}

func (ctx trxMongoRepository) GetTrxListForSyncDataFailed() (trxList []models.Trx, err error) {
	limit := int64(10)

	filter := bson.M{
		"flagSyncData": false,
	}

	opt := &options.FindOptions{
		Sort:  bson.M{"docDate": 1},
		Limit: &limit,
	}

	collect := ctx.RepoDB.MongoDB.Collection(constans.TRX_COLLECTION)
	data, err := collect.Find(ctx.RepoDB.Context, filter, opt)
	if err != nil {
		return nil, err
	}

	for data.Next(ctx.RepoDB.Context) {
		var val models.Trx
		if err = data.Decode(&val); err != nil {
			return nil, err
		}

		trxList = append(trxList, val)
	}
	data.Close(ctx.RepoDB.Context)

	return trxList, nil
}

func (ctx trxMongoRepository) IsHistoryTrxListByCheckInAndCheckoutDate(cardNumberUuid string, productCode, checkInDate string, vehicleNumber string) (trx *models.Trx, exists bool, err error) {
	var result []models.Trx

	filter := bson.M{
		"cardNumberUuidIn": cardNumberUuid,
		"productCode":      productCode,
		"checkInDatetime": bson.M{"$regex": primitive.Regex{
			Pattern: checkInDate,
			Options: "i",
		}},
		"checkOutDatetime": bson.M{"$regex": primitive.Regex{
			Pattern: checkInDate,
			Options: "i",
		}},
		"grandTotal":       bson.M{"$gt": 0},
		"vehicleNumberOut": vehicleNumber,
	}

	opt := &options.FindOptions{}

	collect := ctx.RepoDB.MongoDB.Collection(constans.TRX_COLLECTION)
	cursor, err := collect.Find(ctx.RepoDB.Context, filter, opt)
	if err != nil {
		return nil, false, err
	}

	if err = cursor.All(ctx.RepoDB.Context, &result); err != nil {
		return nil, false, err
	}
	cursor.Close(ctx.RepoDB.Context)

	if len(result) == 0 {
		return nil, false, nil
	}

	return &result[0], true, nil
}

func (ctx trxMongoRepository) GetTrxListByDocDate(docDate string) (trxList []models.Trx, err error) {

	filter := bson.M{
		"docDate": docDate,
	}

	opt := &options.FindOptions{
		Sort: bson.M{"docDate": 1},
	}

	collect := ctx.RepoDB.MongoDB.Collection(constans.TRX_COLLECTION)
	data, err := collect.Find(ctx.RepoDB.Context, filter, opt)
	if err != nil {
		return nil, err
	}

	for data.Next(ctx.RepoDB.Context) {
		var val models.Trx
		if err = data.Decode(&val); err != nil {
			return nil, err
		}

		trxList = append(trxList, val)
	}
	data.Close(ctx.RepoDB.Context)

	return trxList, nil
}

func (ctx trxMongoRepository) UpdateTrxByInterface(filter interface{}, updateSet interface{}) error {

	upsert := true
	opt := &options.UpdateOptions{
		Upsert: &upsert,
	}

	collect := ctx.RepoDB.MongoDB.Collection(constans.TRX_COLLECTION)
	_, err := collect.UpdateOne(ctx.RepoDB.Context, filter, updateSet, opt)
	if err != nil {
		return err
	}

	return nil
}

func (ctx trxMongoRepository) AddTrxOutstandingForClearSession(trx models.TrxOutstandingForClearSession) (ID *primitive.ObjectID, err error) {

	result, err := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_FOR_CLEAR_SESSION_COLLECTION).InsertOne(ctx.RepoDB.Context, trx)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		ID = &oid
	}

	return ID, nil
}

func (ctx trxMongoRepository) IsTrxOutstandingForClearSession(docNo string) (trx *models.TrxOutstandingForClearSession, exists bool, err error) {
	var trxList []models.TrxOutstandingForClearSession

	filter := bson.M{
		"refDocNo": docNo,
	}

	opt := &options.FindOptions{}

	collect := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_FOR_CLEAR_SESSION_COLLECTION)
	result, err := collect.Find(ctx.RepoDB.Context, filter, opt)
	if err != nil {
		return nil, false, err
	}

	for result.Next(ctx.RepoDB.Context) {
		var val models.TrxOutstandingForClearSession
		if err := result.Decode(&val); err != nil {
			return nil, false, err
		}

		trxList = append(trxList, val)
	}
	result.Close(ctx.RepoDB.Context)

	if len(trxList) == 0 {
		return nil, false, nil
	}

	return &trxList[0], true, nil
}

func (ctx trxMongoRepository) GetListTrxByRangeDate(dateFrom, dateTo, ouCode string, limit int64) (*[]models.Trx, error) {
	var result []models.Trx

	filter := make(map[string]interface{})
	filter["docDate"] = bson.M{"$gte": dateFrom, "$lte": dateTo}
	filter["ouCode"] = ouCode
	//filter["docNo"] = "OOG922102714000453"

	option := options.FindOptions{}
	if limit != constans.EMPTY_VALUE_INT {
		option = options.FindOptions{
			Limit: &limit,
		}
	}

	collect := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_COLLECTION)
	data, err := collect.Find(ctx.RepoDB.Context, filter, &option)

	if err != nil {
		return nil, err
	}
	defer data.Close(ctx.RepoDB.Context)

	for data.Next(ctx.RepoDB.Context) {
		var logtrx models.Trx
		data.Decode(&logtrx)
		result = append(result, logtrx)
	}

	return &result, nil
}

func (ctx trxMongoRepository) FindTrxOutstandingByDocNoCustom(docNo string) (*models.ResultFindTrxOutstanding, error) {
	var result models.ResultFindTrxOutstanding

	filter := bson.M{
		"docNo": docNo,
	}

	err := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_COLLECTION).FindOne(ctx.RepoDB.Context, filter).
		Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ctx trxMongoRepository) FindTrxOutstandingByUUIDCustom(uuidCard string) (*models.ResultFindTrxOutstanding, error) {
	var result models.ResultFindTrxOutstanding

	filter := bson.M{
		"cardNumberUuidIn": uuidCard,
	}

	err := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_COLLECTION).FindOne(ctx.RepoDB.Context, filter).
		Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ctx trxMongoRepository) UpdateProductById(keyword, productCode, productName string) error {

	filter := bson.M{
		"keyword": keyword,
	}

	trxProduct := models.TrxProductCustom{
		Keyword:     keyword,
		ProductCode: productCode,
		ProductName: productName,
	}

	after := options.After
	upsert := true
	opt := options.FindOneAndReplaceOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	err := ctx.RepoDB.MongoDB.Collection(constans.TRX_PRODUCT_CUSTOM_COLLECTION).
		FindOneAndReplace(ctx.RepoDB.Context, filter, trxProduct, &opt)
	if err.Err() != nil {
		return err.Err()
	}

	return nil
}

func (ctx trxMongoRepository) IsTrxProductCustomExistsByKeyword(keyword string) (*models.TrxProductCustom, bool) {
	var result models.TrxProductCustom

	filter := bson.M{
		"keyword": keyword,
	}
	err := ctx.RepoDB.MongoDB.Collection(constans.TRX_PRODUCT_CUSTOM_COLLECTION).
		FindOne(ctx.RepoDB.Context, filter).Decode(&result)
	if err != nil {
		return nil, false
	}
	return &result, true
}

func (ctx trxMongoRepository) RemoveTrxProductCustom(keyword string) error {
	filter := bson.M{
		"keyword": keyword,
	}
	_, err := ctx.RepoDB.MongoDB.Collection(constans.TRX_PRODUCT_CUSTOM_COLLECTION).
		DeleteOne(ctx.RepoDB.Context, filter)
	if err != nil {
		return err
	}

	return nil
}

func (ctx trxMongoRepository) UpdateProgressivePriceAndChargeAmountByDocNo(docNo string, progressivePrice, chargeAmount float64) error {
	var update interface{}
	filter := bson.M{
		"docNo": docNo,
	}

	update = bson.M{
		"$set": bson.M{
			"grandTotal":   progressivePrice,
			"chargeAmount": chargeAmount,
		},
	}

	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_COLLECTION).
		UpdateOne(ctx.RepoDB.Context, filter, update, &opt)
	if err != nil {
		return err
	}

	return nil
}

func (ctx trxMongoRepository) FindTrxFlgSyncFalse(flgSyncFalse bool, checkInDateFrom, checkInDateTo string) (*[]models.Trx, error) {
	var result []models.Trx
	filter := make(map[string]interface{})
	//filter["flagSyncData"] = flgSyncFalse
	filter["checkInDatetime"] = bson.M{"$gte": checkInDateFrom, "$lte": checkInDateTo}
	filter["trxInvoiceItem.totalAmount"] = 0

	collect := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_COLLECTION)
	data, err := collect.Find(ctx.RepoDB.Context, filter)
	if err != nil {
		return nil, err
	}
	defer data.Close(ctx.RepoDB.Context)

	for data.Next(ctx.RepoDB.Context) {
		var val models.Trx
		if err = data.Decode(&val); err != nil {
			return nil, err
		}

		result = append(result, val)
	}

	if len(result) == 0 {
		return nil, nil
	}

	return &result, nil
}

func (ctx trxMongoRepository) IsTrxOutstandingByDocNoNew(DocNo string) (models.Trx, bool) {
	var result models.Trx

	filter := bson.M{
		"docNo": DocNo,
	}

	err := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_COLLECTION).
		FindOne(ctx.RepoDB.Context, filter).Decode(&result)
	if err != nil {
		return result, false
	}

	return result, true
}

func (ctx trxMongoRepository) AddTrxCheckInv2(docNo string, trx models.Trx) error {
	filter := bson.M{
		"docNo": docNo,
	}
	data := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_COLLECTION).FindOneAndReplace(ctx.RepoDB.Context, filter, trx)
	if data.Err() != nil {
		return data.Err()
	}

	return nil
}

func (ctx trxMongoRepository) AddTrxInvoiceDetailItem(trxInvoice models.TrxInvoiceDetailItem) error {

	_, err := ctx.RepoDB.MongoDB.Collection(constans.TRX_INVOICE_DETAIL_ITEM_COLLECTIONS).
		InsertOne(ctx.RepoDB.Context, trxInvoice)
	if err != nil {
		return err
	}

	return nil
}

func (ctx trxMongoRepository) FindIDTrxOutstandingByCard(uuidCard string) (models.TrxWithId, error) {
	var result models.TrxWithId

	filter := bson.M{
		"cardNumberUuidIn": uuidCard,
	}

	err := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_COLLECTION).
		FindOne(ctx.RepoDB.Context, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil

}
