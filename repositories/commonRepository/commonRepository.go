package commonRepository

import (
	"togrpc/constans"
	"togrpc/models"
	"togrpc/repositories"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type commonRepository struct {
	RepoDB repositories.Repository
}

func NewCommonRepository(repo repositories.Repository) commonRepository {
	return commonRepository{
		RepoDB: repo,
	}
}

func (ctx commonRepository) AddComboConstan(comboId, comboCode string, comboConstan models.RequestAddComboConstan) error {
	filter := bson.M{
		"comboId":   comboId,
		"comboCode": comboCode,
	}

	upsert := true
	after := options.After
	opt := options.FindOneAndReplaceOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	collect := ctx.RepoDB.MongoDB.Collection(constans.TABLE_COMBO_CONSTANT)
	result := collect.FindOneAndReplace(ctx.RepoDB.Context, filter, comboConstan, &opt)
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (ctx commonRepository) GetComboConstantListByComboId(comboId string) ([]models.ComboConstant, error) {
	var result []models.ComboConstant

	opt := options.FindOptions{
		Sort: bson.D{{"sort", 1}},
	}

	collect := ctx.RepoDB.MongoDB.Collection(constans.TABLE_COMBO_CONSTANT)
	data, err := collect.Find(ctx.RepoDB.Context, bson.M{"comboId": comboId}, &opt)
	if err != nil {
		return nil, err
	}

	if err := data.All(ctx.RepoDB.Context, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (ctx commonRepository) AddPaymentMethodList(ouId int64, ouCode string, listPaymentMethod models.ListPaymentMethod) error {
	filter := bson.M{
		"ouId":      ouId,
		"comboCode": ouCode,
	}

	upsert := true
	after := options.After
	opt := options.FindOneAndReplaceOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	collect := ctx.RepoDB.MongoDB.Collection(constans.TABLE_LIST_PAYMENT_METHOD)
	result := collect.FindOneAndReplace(ctx.RepoDB.Context, filter, listPaymentMethod, &opt)
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (ctx commonRepository) FindListPaymentMethod(ouId int64) (models.ListPaymentMethod, error) {
	var result models.ListPaymentMethod

	filter := bson.M{"ouid": ouId}
	collect := ctx.RepoDB.MongoDB.Collection(constans.TABLE_LIST_PAYMENT_METHOD)
	collect.FindOne(ctx.RepoDB.Context, filter).Decode(&result)

	return result, nil
}
