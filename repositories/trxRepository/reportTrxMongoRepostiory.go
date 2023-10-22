package trxRepository

import (
	"fmt"
	"ppkgwlocal/constans"
	"ppkgwlocal/models"
	"ppkgwlocal/repositories"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type reportTrxMongoRepository struct {
	RepoDB *repositories.Repository
}

func NewReportTrxMongoRepository(repo *repositories.Repository) *reportTrxMongoRepository {
	return &reportTrxMongoRepository{
		RepoDB: repo,
	}
}

func (ctx reportTrxMongoRepository) GetTrxOmzetByDate(dateFrom string, dateTo string) (models.ResponseTrx, error) {
	var result models.ResponseTrx

	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"docDate": bson.M{"$gte": dateFrom, "$lte": dateTo},
			},
		}, {
			"$group": bson.M{
				"_id":        0,
				"totalOmzet": bson.M{"$sum": "$grandTotal"},
				"totalTrx":   bson.M{"$sum": 1},
			},
		},
	}

	pipe, err := ctx.RepoDB.MongoDB.Collection(constans.TRX_COLLECTION).Aggregate(ctx.RepoDB.Context, pipeline)
	if err != nil {
		return result, err
	}
	defer pipe.Close(ctx.RepoDB.Context)
	for pipe.Next(ctx.RepoDB.Context) {
		pipe.Decode(&result)
	}

	return result, nil
}

func (ctx reportTrxMongoRepository) GetTrxOutstandingByDate(dateFrom string, dateTo string) (models.ResponseTrxOutstanding, error) {
	var result models.ResponseTrxOutstanding

	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"extLocalDatetime": bson.M{"$gte": dateFrom, "$lte": dateTo},
			},
		}, {
			"$group": bson.M{
				"_id":                 0,
				"totalTrxOutstanding": bson.M{"$sum": 1},
			},
		},
	}

	pipe, err := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_COLLECTION).Aggregate(ctx.RepoDB.Context, pipeline)
	if err != nil {
		fmt.Println(err)
		return result, err
	}
	defer pipe.Close(ctx.RepoDB.Context)
	for pipe.Next(ctx.RepoDB.Context) {
		pipe.Decode(&result)
	}

	return result, nil
}

func (ctx reportTrxMongoRepository) GetTrxOvernightByDate() ([]models.TrxInvoiceDetail, error) {
	var result []models.TrxInvoiceDetail

	projections := bson.M{
		"createdDate": 1,
	}

	options := options.FindOptions{
		Projection: projections,
	}

	data, err := ctx.RepoDB.MongoDB.Collection(constans.TRX_INVOICE_DETAIL_ITEM_COLLECTIONS).
		Find(ctx.RepoDB.Context, bson.M{}, &options)
	if err != nil {
		return result, err
	}

	for data.Next(ctx.RepoDB.Context) {
		var val models.TrxInvoiceDetail
		data.Decode(&val)

		result = append(result, val)
	}

	return result, nil
}

func (ctx reportTrxMongoRepository) GetRevenueForIntelligenceTrxList(dateFrom string, dateTo string) ([]models.ChartDailyIntelligenceTrxList, error) {
	var result []models.ChartDailyIntelligenceTrxList

	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{"docDate": bson.M{"$gte": dateFrom, "$lte": dateTo}},
		},
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"Days":  bson.M{"$dayOfMonth": bson.M{"$dateFromString": bson.M{"dateString": "$docDate"}}},
					"Month": bson.M{"$month": bson.M{"$dateFromString": bson.M{"dateString": "$docDate"}}},
					"Year":  bson.M{"$year": bson.M{"$dateFromString": bson.M{"dateString": "$docDate"}}},
				},
				"totalOmzet": bson.M{"$sum": "$grandTotal"},
				"totalTrx":   bson.M{"$sum": 1},
			},
		},
	}

	pipe, err := ctx.RepoDB.MongoDB.Collection(constans.TRX_COLLECTION).Aggregate(ctx.RepoDB.Context, pipeline)
	if err != nil {
		return result, err
	}
	defer pipe.Close(ctx.RepoDB.Context)
	for pipe.Next(ctx.RepoDB.Context) {
		var val models.ChartDailyIntelligenceTrxList
		pipe.Decode(&val)

		result = append(result, val)
	}

	return result, nil
}

func (ctx reportTrxMongoRepository) GetTrafficParkingFromCheckin(date string, listTime []string) ([]models.TrafficParking, error) {
	var result []models.TrafficParking

	for _, row := range listTime {
		var data models.TrafficParking
		var totalTrx models.TotalTrx
		hoursRange := strings.Split(row, "#")
		data.Hour = fmt.Sprintf("%s-%s", hoursRange[0], hoursRange[1])
		data.TotalTrx = 0

		pipeline := []bson.M{
			bson.M{"$redact": bson.M{
				"$cond": []interface{}{
					bson.M{"$and": []bson.M{
						bson.M{"docDate": date},
						bson.M{"$gte": []interface{}{bson.M{"$substr": []interface{}{"$checkInDatetime", 11, 5}}, hoursRange[0]}},
						bson.M{"$lte": []interface{}{bson.M{"$substr": []interface{}{"$checkInDatetime", 11, 5}}, hoursRange[1]}},
					},
					},
					"$$KEEP",
					"$$PRUNE",
				},
			},
			},
			bson.M{
				"$group": bson.M{
					"_id":      0,
					"totalTrx": bson.M{"$sum": 1},
				},
			},
		}

		pipe, err := ctx.RepoDB.MongoDB.Collection(constans.TRX_COLLECTION).Aggregate(ctx.RepoDB.Context, pipeline)
		if err != nil {
			return result, err
		}
		defer pipe.Close(ctx.RepoDB.Context)
		for pipe.Next(ctx.RepoDB.Context) {
			pipe.Decode(&totalTrx)
		}
		data.TotalTrx = totalTrx.TotalTrx
		result = append(result, data)

	}

	return result, nil
}

func (ctx reportTrxMongoRepository) GetTrxListAdvance(trx models.TrxItems, isExport bool) ([]models.ResultFindTrxForPostBox, error) {
	var result []models.ResultFindTrxForPostBox
	filter := make(map[string]interface{})
	filter["checkOutDatetime"] = bson.M{"$gte": trx.DateFrom, "$lte": trx.DateTo}

	if trx.Keyword != constans.EMPTY_VALUE {
		filter["$or"] = []bson.M{
			bson.M{"docNo": bson.M{
				"$regex": primitive.Regex{
					Pattern: trx.Keyword,
					Options: "i",
				},
			}},
			bson.M{"cardNumberUuidIn": bson.M{
				"$regex": primitive.Regex{
					Pattern: trx.Keyword,
					Options: "i",
				},
			}},
			bson.M{"vehicleNumberIn": bson.M{
				"$regex": primitive.Regex{
					Pattern: trx.Keyword,
					Options: "i",
				},
			}},
			bson.M{"memberName": bson.M{
				"$regex": primitive.Regex{
					Pattern: trx.Keyword,
					Options: "i",
				},
			}},
		}
	}

	if trx.TypeCard != constans.EMPTY_VALUE {
		filter["typeCard"] = trx.TypeCard
	}

	options := options.FindOptions{
		Projection: bson.M{
			"_id":             1,
			"docNo":           1,
			"grandTotal":      1,
			"checkInDatetime": 1,
			"cardNumber":      1,
			"cardNumberUuid":  1,
			"typeCard":        1,
			"ouCode":          1,
			"vehicleNumberIn": 1,
			"gateIn":          1,
			"qrText":          1,
			"trxInvoiceItem":  1,
		},
	}

	if trx.ColumnOrderName != constans.EMPTY_VALUE {
		if trx.AscDesc == constans.ASCENDING {
			options.Sort = bson.M{trx.ColumnOrderName: 1}
		} else if trx.AscDesc == constans.DESCENDING {
			options.Sort = bson.M{trx.ColumnOrderName: -1}
		} else {
			options.Sort = bson.M{"productCode": -1}
		}
	}

	if !isExport {
		options.Limit = &trx.Limit
		options.Skip = &trx.Offset
	}

	collect := ctx.RepoDB.MongoDB.Collection(constans.TRX_COLLECTION)
	data, err := collect.Find(ctx.RepoDB.Context, filter, &options)

	if err != nil {
		return result, err
	}

	defer data.Close(ctx.RepoDB.Context)

	for data.Next(ctx.RepoDB.Context) {
		var val models.ResultFindTrxForPostBox
		data.Decode(&val)

		result = append(result, val)
	}

	return result, nil
}

func (ctx reportTrxMongoRepository) CountGetTrxListAdvance(trx models.TrxItems) (int64, error) {
	filter := make(map[string]interface{})
	filter["docDate"] = bson.M{"$gte": trx.DateFrom, "$lte": trx.DateTo}

	if trx.Keyword != constans.EMPTY_VALUE {
		filter["$or"] = []bson.M{
			bson.M{"docNo": bson.M{
				"$regex": primitive.Regex{
					Pattern: trx.Keyword,
					Options: "i",
				},
			}},
			bson.M{"cardNumberUuid": bson.M{
				"$regex": primitive.Regex{
					Pattern: trx.Keyword,
					Options: "i",
				},
			}},
			bson.M{"vehicleNumberIn": bson.M{
				"$regex": primitive.Regex{
					Pattern: trx.Keyword,
					Options: "i",
				},
			}},
			bson.M{"memberName": bson.M{
				"$regex": primitive.Regex{
					Pattern: trx.Keyword,
					Options: "i",
				},
			}},
		}
	}

	if trx.TypeCard != constans.EMPTY_VALUE {
		filter["typeCard"] = trx.TypeCard
	}

	collect := ctx.RepoDB.MongoDB.Collection(constans.TRX_COLLECTION)
	count, err := collect.CountDocuments(ctx.RepoDB.Context, filter)

	if err != nil {
		return count, err
	}

	return count, nil
}

func (ctx reportTrxMongoRepository) GetTrxVehicleByDocDate(requestTrx models.RequestGetTrxOrderByVehicle) ([]models.GetTrxOrderByVehicle, error) {
	var result []models.GetTrxOrderByVehicle
	var pipeline []interface{}
	start := make(map[string]interface{})
	end := make(map[string]interface{})

	start["$match"] = bson.M{"docDate": requestTrx.Date}

	if requestTrx.ProductCode != constans.EMPTY_VALUE {
		start["$match"] = bson.M{"productCode": requestTrx.ProductCode}
	}

	end["$group"] = bson.M{
		"_id": bson.M{
			"productCode": "$productCode",
			"productName": "$productName",
		},
		"totalOmzet": bson.M{"$sum": "$grandTotal"},
		"totalTrx":   bson.M{"$sum": 1},
	}

	pipeline = append(pipeline, start, end)
	pipe, err := ctx.RepoDB.MongoDB.Collection(constans.TRX_COLLECTION).Aggregate(ctx.RepoDB.Context, pipeline)
	if err != nil {
		return result, err
	}
	defer pipe.Close(ctx.RepoDB.Context)
	for pipe.Next(ctx.RepoDB.Context) {
		var val models.GetTrxOrderByVehicle
		pipe.Decode(&val)

		result = append(result, val)
	}

	return result, nil
}

func (ctx reportTrxMongoRepository) AddTrxLostTicket(trxLostTicket models.LostTicket) (*primitive.ObjectID, error) {
	var ID primitive.ObjectID

	result, err := ctx.RepoDB.MongoDB.Collection(constans.LOST_TICKET_COLLECTION).
		InsertOne(ctx.RepoDB.Context, trxLostTicket)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		ID = oid
	}
	return &ID, nil

}

func (ctx reportTrxMongoRepository) GetListTrxLostTicket(lostTicket models.RequestFindLostTicket) (*[]models.FindLostTicket, error) {
	var result []models.FindLostTicket
	var options = options.FindOptions{}
	filter := make(map[string]interface{})
	filter["createdAt"] = bson.M{"$gte": lostTicket.DateFrom, "$lte": lostTicket.DateTo}

	if lostTicket.Keyword != constans.EMPTY_VALUE {
		filter["$or"] = []bson.M{
			bson.M{"custIdentifier": bson.M{
				"$regex": primitive.Regex{
					Pattern: lostTicket.Keyword,
					Options: "i",
				},
			}},
			bson.M{"vehicleNumber": bson.M{
				"$regex": primitive.Regex{
					Pattern: lostTicket.Keyword,
					Options: "i",
				},
			}},
			bson.M{"vehicleRegistration": bson.M{
				"$regex": primitive.Regex{
					Pattern: lostTicket.Keyword,
					Options: "i",
				},
			}},
			bson.M{"officerName": bson.M{
				"$regex": primitive.Regex{
					Pattern: lostTicket.Keyword,
					Options: "i",
				},
			}},
		}
	}

	filter["type"] = "CONFIRM"
	if lostTicket.ColumnOrderName != constans.EMPTY_VALUE {
		if lostTicket.AscDesc == constans.ASCENDING {
			options.Sort = bson.M{lostTicket.ColumnOrderName: 1}
		} else if lostTicket.AscDesc == constans.DESCENDING {
			options.Sort = bson.M{lostTicket.ColumnOrderName: -1}
		} else {
			options.Sort = bson.M{"productCode": -1}
		}
	}

	options.Limit = &lostTicket.Limit
	options.Skip = &lostTicket.Offset

	collect := ctx.RepoDB.MongoDB.Collection(constans.LOST_TICKET_COLLECTION)
	data, err := collect.Find(ctx.RepoDB.Context, filter, &options)

	if err != nil {
		return nil, err
	}

	defer data.Close(ctx.RepoDB.Context)

	for data.Next(ctx.RepoDB.Context) {
		var val models.FindLostTicket
		data.Decode(&val)

		result = append(result, val)
	}

	return &result, nil
}

func (ctx reportTrxMongoRepository) CountListTrxLostTicket(lostTicket models.RequestFindLostTicket) (*int64, error) {
	filter := make(map[string]interface{})
	filter["createdAt"] = bson.M{"$gte": lostTicket.DateFrom, "$lte": lostTicket.DateTo}

	if lostTicket.Keyword != constans.EMPTY_VALUE {
		filter["$or"] = []bson.M{
			bson.M{"custIdentifier": bson.M{
				"$regex": primitive.Regex{
					Pattern: lostTicket.Keyword,
					Options: "i",
				},
			}},
			bson.M{"vehicleNumber": bson.M{
				"$regex": primitive.Regex{
					Pattern: lostTicket.Keyword,
					Options: "i",
				},
			}},
			bson.M{"vehicleRegistration": bson.M{
				"$regex": primitive.Regex{
					Pattern: lostTicket.Keyword,
					Options: "i",
				},
			}},
			bson.M{"officerName": bson.M{
				"$regex": primitive.Regex{
					Pattern: lostTicket.Keyword,
					Options: "i",
				},
			}},
		}
	}

	filter["type"] = "CONFIRM"
	collect := ctx.RepoDB.MongoDB.Collection(constans.LOST_TICKET_COLLECTION)
	resultCountData, err := collect.CountDocuments(ctx.RepoDB.Context, filter)
	if err != nil {
		return nil, err
	}

	return &resultCountData, nil
}

func (ctx reportTrxMongoRepository) GetListTrxOutstanding(trx models.RequestFindTrxOutstanding) (*[]models.ResultFindTrxForPostBox, error) {
	var result []models.ResultFindTrxForPostBox

	filter := make(map[string]interface{})
	filter["checkInDatetime"] = bson.M{"$gte": trx.DateFrom, "$lte": trx.DateTo}

	if trx.Keyword != constans.EMPTY_VALUE {
		filter["$or"] = []bson.M{
			bson.M{"cardNumberUuid": bson.M{
				"$regex": primitive.Regex{
					Pattern: trx.Keyword,
					Options: "i",
				},
			}},
			bson.M{"cardNumber": bson.M{
				"$regex": primitive.Regex{
					Pattern: trx.Keyword,
					Options: "i",
				},
			}},
			bson.M{"docNo": bson.M{
				"$regex": primitive.Regex{
					Pattern: trx.Keyword,
					Options: "i",
				},
			}},
			bson.M{"deviceIdIn": bson.M{
				"$regex": primitive.Regex{
					Pattern: trx.Keyword,
					Options: "i",
				},
			}},
			bson.M{"gateIn": bson.M{
				"$regex": primitive.Regex{
					Pattern: trx.Keyword,
					Options: "i",
				},
			}},
			bson.M{"vehicleNumberIn": bson.M{
				"$regex": primitive.Regex{
					Pattern: trx.Keyword,
					Options: "i",
				},
			}},
		}
	}

	options := options.FindOptions{
		Projection: bson.M{
			"_id":             1,
			"docNo":           1,
			"grandTotal":      1,
			"checkInDatetime": 1,
			"cardNumber":      1,
			"cardNumberUuid":  1,
			"typeCard":        1,
			"ouCode":          1,
			"vehicleNumberIn": 1,
			"gateIn":          1,
			"qrText":          1,
			"trxInvoiceItem":  1,
		},
	}

	if trx.ColumnOrderName != constans.EMPTY_VALUE {
		if trx.AscDesc == constans.ASCENDING {
			options.Sort = bson.M{trx.ColumnOrderName: 1}
		} else if trx.AscDesc == constans.DESCENDING {
			options.Sort = bson.M{trx.ColumnOrderName: -1}
		} else {
			options.Sort = bson.M{"checkInDatetime": -1}
		}
	}

	options.Limit = &trx.Limit
	options.Skip = &trx.Offset

	collect := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_COLLECTION)
	data, err := collect.Find(ctx.RepoDB.Context, filter, &options)

	if err != nil {
		return nil, err
	}

	defer data.Close(ctx.RepoDB.Context)

	for data.Next(ctx.RepoDB.Context) {
		var val models.ResultFindTrxForPostBox
		if err = data.Decode(&val); err != nil {
			return nil, err
		}
		result = append(result, val)
	}

	return &result, nil
}

func (ctx reportTrxMongoRepository) CountTrxOutstanding(dateFrom string, dateTo string, keyword string) (models.ResponseTrxOutstanding, error) {
	var result models.ResponseTrxOutstanding

	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"checkInDatetime": bson.M{"$gte": dateFrom, "$lte": dateTo},
				"$or": []bson.M{
					bson.M{"cardNumberUuid": bson.M{
						"$regex": primitive.Regex{
							Pattern: keyword,
							Options: "i",
						},
					}},
					bson.M{"cardNumber": bson.M{
						"$regex": primitive.Regex{
							Pattern: keyword,
							Options: "i",
						},
					}},
				},
			},
		}, {
			"$group": bson.M{
				"_id":                 0,
				"totalTrxOutstanding": bson.M{"$sum": 1},
			},
		},
	}

	pipe, err := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_COLLECTION).Aggregate(ctx.RepoDB.Context, pipeline)
	if err != nil {
		fmt.Println(err)
		return result, err
	}
	defer pipe.Close(ctx.RepoDB.Context)
	for pipe.Next(ctx.RepoDB.Context) {
		pipe.Decode(&result)
	}

	return result, nil
}

func (ctx reportTrxMongoRepository) GetListTrxOneTimeMember(trx models.RequestFindOneTimeMember, isExport bool) (*[]models.TrxOneTimeMember, error) {
	var result []models.TrxOneTimeMember
	filter := make(map[string]interface{})
	filter["checkOutDatetime"] = bson.M{"$gte": trx.DateFrom, "$lte": trx.DateTo}
	filter["memberCode"] = bson.M{"$ne": constans.EMPTY_VALUE}
	filter["memberType"] = constans.TYPE_PARTNER_ONE_TIME

	if trx.Keyword != constans.EMPTY_VALUE {
		filter["$or"] = []bson.M{
			bson.M{"docNo": bson.M{
				"$regex": primitive.Regex{
					Pattern: trx.Keyword,
					Options: "i",
				},
			}},
			bson.M{"cardNumberUuidIn": bson.M{
				"$regex": primitive.Regex{
					Pattern: trx.Keyword,
					Options: "i",
				},
			}},
			bson.M{"vehicleNumberIn": bson.M{
				"$regex": primitive.Regex{
					Pattern: trx.Keyword,
					Options: "i",
				},
			}},
			bson.M{"memberName": bson.M{
				"$regex": primitive.Regex{
					Pattern: trx.Keyword,
					Options: "i",
				},
			}},
		}
	}

	options := options.FindOptions{
		Projection: bson.M{
			"_id":              1,
			"docNo":            1,
			"docDate":          1,
			"productName":      1,
			"memberName":       1,
			"cardNumberUuidIn": 1,
			"checkInDatetime":  1,
			"checkOutDatetime": 1,
			"vehicleNumberIn":  1,
		},
	}

	if trx.ColumnOrderName != constans.EMPTY_VALUE {
		if trx.AscDesc == constans.ASCENDING {
			options.Sort = bson.M{trx.ColumnOrderName: 1}
		} else if trx.AscDesc == constans.DESCENDING {
			options.Sort = bson.M{trx.ColumnOrderName: -1}
		} else {
			options.Sort = bson.M{"productCode": -1}
		}
	}

	if !isExport {
		options.Limit = &trx.Limit
		options.Skip = &trx.Offset
	}

	collect := ctx.RepoDB.MongoDB.Collection(constans.TRX_COLLECTION)
	data, err := collect.Find(ctx.RepoDB.Context, filter, &options)

	if err != nil {
		return nil, err
	}

	defer data.Close(ctx.RepoDB.Context)

	for data.Next(ctx.RepoDB.Context) {
		var val models.TrxOneTimeMember
		data.Decode(&val)

		result = append(result, val)
	}

	return &result, nil
}

func (ctx reportTrxMongoRepository) CountGetListTrxOneTimeMember(trx models.RequestFindOneTimeMember) (int64, error) {
	filter := make(map[string]interface{})
	filter["docDate"] = bson.M{"$gte": trx.DateFrom, "$lte": trx.DateTo}

	if trx.Keyword != constans.EMPTY_VALUE {
		filter["$or"] = []bson.M{
			bson.M{"docNo": bson.M{
				"$regex": primitive.Regex{
					Pattern: trx.Keyword,
					Options: "i",
				},
			}},
			bson.M{"cardNumberUuidIn": bson.M{
				"$regex": primitive.Regex{
					Pattern: trx.Keyword,
					Options: "i",
				},
			}},
			bson.M{"vehicleNumberIn": bson.M{
				"$regex": primitive.Regex{
					Pattern: trx.Keyword,
					Options: "i",
				},
			}},
			bson.M{"memberName": bson.M{
				"$regex": primitive.Regex{
					Pattern: trx.Keyword,
					Options: "i",
				},
			}},
		}
	}

	filter["memberCode"] = bson.M{"$ne": constans.EMPTY_VALUE}
	collect := ctx.RepoDB.MongoDB.Collection(constans.TRX_COLLECTION)
	count, err := collect.CountDocuments(ctx.RepoDB.Context, filter)

	if err != nil {
		return count, err
	}

	return count, nil
}

func (ctx reportTrxMongoRepository) FindQrTextByDocNo(docNo string) (*string, error) {
	var qrText string
	filter := bson.M{
		"docNo": docNo,
	}

	err := ctx.RepoDB.MongoDB.Collection(constans.TRX_PRODUCT_CUSTOM_COLLECTION).
		FindOne(ctx.RepoDB.Context, filter).Decode(&qrText)
	if err != nil {
		return nil, err
	}

	return &qrText, nil
}

func (ctx reportTrxMongoRepository) FindTrxLostTicketByDocNo(id primitive.ObjectID) (*models.FindLostTicket, error) {
	var result models.FindLostTicket
	filter := bson.M{
		"_id": id,
	}

	err := ctx.RepoDB.MongoDB.Collection(constans.LOST_TICKET_COLLECTION).
		FindOne(ctx.RepoDB.Context, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (ctx reportTrxMongoRepository) FindTrxLostTicketByID(ID primitive.ObjectID) (*models.LostTicket, error) {
	var result models.LostTicket

	filter := bson.M{
		"_id": ID,
	}

	err := ctx.RepoDB.MongoDB.Collection(constans.LOST_TICKET_COLLECTION).
		FindOne(ctx.RepoDB.Context, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ctx reportTrxMongoRepository) UpdateTypeTrxLostTicket(ID primitive.ObjectID, typeTrx string, trx models.LostTicket) error {
	filter := bson.M{
		"_id": ID,
	}

	update := bson.M{
		"$set": bson.M{
			"lostTicketOutDatetime":         trx.LostTicketOutDatetime,
			"type":                          typeTrx,
			"trxOutstanding.cardNumberUuid": trx.TrxOutstanding.CardNumberUUID,
			"trxOutstanding.cardNumber":     trx.TrxOutstanding.CardNumber,
			"trxOutstanding.typeCard":       trx.TrxOutstanding.TypeCard,
			"trxOutstanding.deviceId":       trx.TrxOutstanding.DeviceId,
			"trxOutstanding.gateOut":        trx.TrxOutstanding.GateOut,
			"trxOutstanding.productCode":    trx.TrxOutstanding.ProductCode,
			"trxOutstanding.productName":    trx.TrxOutstanding.ProductName,
			"trxOutstanding.requestOutData": trx.TrxOutstanding.RequestOutData,
			"trxOutstanding.logTrans":       trx.TrxOutstanding.LogTrans,
			"trxOutstanding.grandTotal":     trx.TrxOutstanding.GrandTotal,
			"trxOutstanding.chargeAmount":   trx.ChargeAmount,
		},
	}

	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := ctx.RepoDB.MongoDB.Collection(constans.LOST_TICKET_COLLECTION).
		UpdateOne(ctx.RepoDB.Context, filter, update, &opt)
	if err != nil {
		return err
	}

	return nil
}

func (ctx reportTrxMongoRepository) UpdateQrCodeTrxLostTicket(ID primitive.ObjectID, qrCode string) error {
	filter := bson.M{
		"_id": ID,
	}

	update := bson.M{
		"$set": bson.M{
			"qrCodeLostTicket": qrCode,
		},
	}

	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := ctx.RepoDB.MongoDB.Collection(constans.LOST_TICKET_COLLECTION).
		UpdateOne(ctx.RepoDB.Context, filter, update, &opt)
	if err != nil {
		return err
	}

	return nil
}

func (ctx reportTrxMongoRepository) UpdatePaymentRefDocNoCodeTrxLostTicket(ID primitive.ObjectID, inquiryTrx models.RequestUpdateTrxInquiryQris) error {
	filter := bson.M{
		"_id": ID,
	}

	update := bson.M{
		"$set": bson.M{
			"paymentRefDocNo": inquiryTrx.PaymentRefDocNo,
			"grandTotal":      inquiryTrx.GrandTotal,
			"qrCodePayment":   inquiryTrx.QrPayment,
		},
	}

	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := ctx.RepoDB.MongoDB.Collection(constans.LOST_TICKET_COLLECTION).
		UpdateOne(ctx.RepoDB.Context, filter, update, &opt)
	if err != nil {
		return err
	}

	return nil
}

func (ctx reportTrxMongoRepository) FindQrLostTicketByDocNo(docNo string) (result *models.ResponseFindQrLostTicket, exists bool) {
	filter := bson.M{
		"docNo":                 docNo,
		"type":                  "CONFIRM",
		"lostTicketOutDatetime": bson.M{"$ne": constans.EMPTY_VALUE},
	}

	err := ctx.RepoDB.MongoDB.Collection(constans.LOST_TICKET_COLLECTION).FindOne(ctx.RepoDB.Context, filter).Decode(&result)
	if err != nil {
		return nil, false
	}

	return result, true
}
