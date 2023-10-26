package trxRepository

import (
	"github.com/raj847/togrpc/constans"
	"github.com/raj847/togrpc/models"
	"github.com/raj847/togrpc/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type trxMongoDepositCounterRepository struct {
	RepoDB repositories.Repository
}

func NewTrxMongoDepositCounterRepository(repo repositories.Repository) trxMongoDepositCounterRepository {
	return trxMongoDepositCounterRepository{
		RepoDB: repo,
	}
}

func (ctx trxMongoDepositCounterRepository) AddTrxCheckInDepositCenter(trx models.TrxDepositCounter) (*primitive.ObjectID, error) {
	var ID primitive.ObjectID

	result, err := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_DEPOSIT_COUNTER_COLLECTION).
		InsertOne(ctx.RepoDB.Context, trx)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		ID = oid
	}

	return &ID, nil
}

func (ctx trxMongoDepositCounterRepository) IsTrxDepositCenterOutstandingByDocNo(docNo, productCode string) (*models.ResultFindTrxDepositCounterOutstanding, error) {

	var result models.ResultFindTrxDepositCounterOutstanding
	filter := make(map[string]interface{})

	if productCode != constans.EMPTY_VALUE {
		filter = bson.M{
			"docNoDepo": docNo,
			"trxInvoiceItem": bson.M{
				"$elemMatch": bson.M{
					"docNoDepo":   docNo,
					"productCode": productCode,
				},
			},
		}
	} else {
		filter = bson.M{
			"docNoDepo": docNo,
			"trxInvoiceItem": bson.M{
				"$elemMatch": bson.M{
					"docNoDepo": docNo,
				},
			},
		}
	}

	projection := bson.M{
		"_id":              1,
		"docNoDepo":        1,
		"grandTotal":       1,
		"merk":             1,
		"depositorName":    1,
		"productName":      1,
		"checkInDatetime":  1,
		"ouCode":           1,
		"trxInvoiceItem.$": 1,
	}

	err := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_DEPOSIT_COUNTER_COLLECTION).
		FindOne(ctx.RepoDB.Context, filter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil

}

func (ctx trxMongoDepositCounterRepository) FindTrxDepositCenterOutstandingByID(ID primitive.ObjectID) (models.TrxDepositCounter, error) {
	var result models.TrxDepositCounter

	filter := bson.M{
		"_id": ID,
	}

	err := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_DEPOSIT_COUNTER_COLLECTION).
		FindOne(ctx.RepoDB.Context, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (ctx trxMongoDepositCounterRepository) AddTrxDepositCounter(trx models.TrxDepositCounter) (string, error) {
	_, err := ctx.RepoDB.MongoDB.Collection(constans.TRX_DEPOSIT_COUNTER_COLLECTION).
		InsertOne(ctx.RepoDB.Context, trx)
	if err != nil {
		return "", err
	}
	return trx.DocNoDepo, nil
}

func (ctx trxMongoDepositCounterRepository) RemoveTrxDepositCounterByID(ID primitive.ObjectID) error {

	filter := bson.M{
		"_id": ID,
	}
	data := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_DEPOSIT_COUNTER_COLLECTION).
		FindOneAndDelete(ctx.RepoDB.Context, filter)
	if data.Err() != nil {
		return data.Err()
	}

	return nil
}

func (ctx trxMongoDepositCounterRepository) GetListTrxDepositCounter(trxDepo models.TrxItemsDepositCounter) (*[]models.FindTrxDepositCounter, error) {
	var result []models.FindTrxDepositCounter
	var options = options.FindOptions{}
	filter := make(map[string]interface{})
	filter["checkInDatetime"] = bson.M{"$gte": trxDepo.DateFrom, "$lte": trxDepo.DateTo}

	if trxDepo.Keyword != constans.EMPTY_VALUE {
		filter["$or"] = []bson.M{
			bson.M{"depositorName": bson.M{
				"$regex": primitive.Regex{
					Pattern: trxDepo.Keyword,
					Options: "i",
				},
			}},
			bson.M{"merk": bson.M{
				"$regex": primitive.Regex{
					Pattern: trxDepo.Keyword,
					Options: "i",
				},
			}},
		}
	}

	if trxDepo.ColumnOrderName != constans.EMPTY_VALUE {
		if trxDepo.AscDesc == constans.ASCENDING {
			options.Sort = bson.M{trxDepo.ColumnOrderName: 1}
		} else if trxDepo.AscDesc == constans.DESCENDING {
			options.Sort = bson.M{trxDepo.ColumnOrderName: -1}
		} else {
			options.Sort = bson.M{"checkInDatetime": -1}
		}
	}

	options.Limit = &trxDepo.Limit
	options.Skip = &trxDepo.Offset

	collect := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_DEPOSIT_COUNTER_COLLECTION)
	data, err := collect.Find(ctx.RepoDB.Context, filter, &options)

	if err != nil {
		return nil, err
	}

	defer data.Close(ctx.RepoDB.Context)

	for data.Next(ctx.RepoDB.Context) {
		var val models.FindTrxDepositCounter
		data.Decode(&val)

		result = append(result, val)
	}

	return &result, nil
}

func (ctx trxMongoDepositCounterRepository) CountGetTrxDepositCounter(trxDepo models.TrxItemsDepositCounter) (int64, error) {
	filter := make(map[string]interface{})
	filter["checkInDatetime"] = bson.M{"$gte": trxDepo.DateFrom, "$lte": trxDepo.DateTo}

	if trxDepo.Keyword != constans.EMPTY_VALUE {
		filter["$or"] = []bson.M{
			bson.M{"depositorName": bson.M{
				"$regex": primitive.Regex{
					Pattern: trxDepo.Keyword,
					Options: "i",
				},
			}},
			bson.M{"merk": bson.M{
				"$regex": primitive.Regex{
					Pattern: trxDepo.Keyword,
					Options: "i",
				},
			}},
		}
	}

	collect := ctx.RepoDB.MongoDB.Collection(constans.TRX_OUTSTANDING_DEPOSIT_COUNTER_COLLECTION)
	count, err := collect.CountDocuments(ctx.RepoDB.Context, filter)

	if err != nil {
		return count, err
	}

	return count, nil
}
