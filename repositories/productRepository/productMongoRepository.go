package productRepository

import (
	"errors"
	"fmt"
	"togrpc/constans"
	"togrpc/models"
	"togrpc/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type productRepository struct {
	RepoDB repositories.Repository
}

func NewProductRepository(repo repositories.Repository) productRepository {
	return productRepository{
		RepoDB: repo,
	}
}

func (ctx productRepository) FindProductByProductCode(productCode string, ouId int64) (models.PolicyOuProductWithRules, error) {
	var result models.PolicyOuProductWithRules

	filter := bson.M{
		"ouId":        ouId,
		"productCode": productCode,
	}

	data := ctx.RepoDB.MongoDB.Collection(constans.TABLE_POLICY_OU_PRODUCT).FindOne(ctx.RepoDB.Context, filter)
	if data.Err() != nil {
		fmt.Println(data.Err())
		return result, errors.New("Policy Ou Product Not Found")
	}

	err := data.Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (ctx productRepository) FindPolicyOuProductByAdvance(OuID int64, ProductId int64) (models.PolicyOuProductWithRules, error) {
	var result models.PolicyOuProductWithRules

	filter := bson.M{
		"ouId":      OuID,
		"productId": ProductId,
	}

	collect := ctx.RepoDB.MongoDB.Collection(constans.TABLE_POLICY_OU_PRODUCT)
	data := collect.FindOne(ctx.RepoDB.Context, filter)
	if data.Err() != nil {
		return result, errors.New("policy ou product not found")
	}

	err := data.Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (ctx productRepository) GetPolicyOuProductList() ([]models.PolicyOuProductWithRules, error) {
	var result []models.PolicyOuProductWithRules

	collect := ctx.RepoDB.MongoDB.Collection(constans.TABLE_POLICY_OU_PRODUCT)
	data, err := collect.Find(ctx.RepoDB.Context, bson.M{})
	if err != nil {
		return result, errors.New("policy ou product not found")
	}

	for data.Next(ctx.RepoDB.Context) {
		var val models.PolicyOuProductWithRules
		data.Decode(&val)

		result = append(result, val)
	}

	return result, nil
}

func (ctx productRepository) GetPolicyOuProductListAdvance(product models.RequestGetPolicyOuProductList) ([]models.PolicyOuProductWithRules, error) {
	var result []models.PolicyOuProductWithRules
	filter := make(map[string]interface{})

	if product.Keyword != constans.EMPTY_VALUE {
		filter["$or"] = []bson.M{
			bson.M{"productCode": bson.M{
				"$regex": primitive.Regex{
					Pattern: product.Keyword,
					Options: "i",
				},
			}},
			bson.M{"productName": bson.M{
				"$regex": primitive.Regex{
					Pattern: product.Keyword,
					Options: "i",
				},
			}},
		}
	}

	options := options.FindOptions{}
	if product.ColumnOrderName != constans.EMPTY_VALUE {
		if product.AscDesc == constans.ASCENDING {
			options.Sort = bson.M{product.ColumnOrderName: 1}
		} else if product.AscDesc == constans.DESCENDING {
			options.Sort = bson.M{product.ColumnOrderName: -1}
		} else {
			options.Sort = bson.M{"productRules.price": -1}
		}
	} else {
		options.Sort = bson.M{"productRules.price": 1}

	}

	collect := ctx.RepoDB.MongoDB.Collection(constans.TABLE_POLICY_OU_PRODUCT)
	data, err := collect.Find(ctx.RepoDB.Context, filter, &options)

	if err != nil {
		return nil, err
	}

	defer data.Close(ctx.RepoDB.Context)

	for data.Next(ctx.RepoDB.Context) {
		var val models.PolicyOuProductWithRules
		data.Decode(&val)

		result = append(result, val)
	}

	return result, nil
}

func (ctx productRepository) CountGetPolicyOuProductListAdvance(keyword string) (int64, error) {
	filter := make(map[string]interface{})

	if keyword != constans.EMPTY_VALUE {
		filter["$or"] = []bson.M{
			bson.M{"productCode": bson.M{
				"$regex": primitive.Regex{
					Pattern: keyword,
					Options: "i",
				},
			}},
			bson.M{"productName": bson.M{
				"$regex": primitive.Regex{
					Pattern: keyword,
					Options: "i",
				},
			}},
		}
	}

	collect := ctx.RepoDB.MongoDB.Collection(constans.TABLE_POLICY_OU_PRODUCT)
	count, err := collect.CountDocuments(ctx.RepoDB.Context, filter)

	if err != nil {
		return count, err
	}

	return count, nil
}

func (ctx productRepository) GetPolicyOuProductByOuId(OuID int64) (*[]models.PolicyOuProductWithRules, error) {
	var result []models.PolicyOuProductWithRules

	filter := bson.M{
		"ouId": OuID,
	}

	collect := ctx.RepoDB.MongoDB.Collection(constans.TABLE_POLICY_OU_PRODUCT)
	data, err := collect.Find(ctx.RepoDB.Context, filter)
	if err != nil {
		return nil, errors.New("policy ou product not found")
	}

	defer data.Close(ctx.RepoDB.Context)
	for data.Next(ctx.RepoDB.Context) {
		var val models.PolicyOuProductWithRules
		data.Decode(&val)

		result = append(result, val)
	}

	return &result, nil
}

func (ctx productRepository) FindProductDepositByProductCode(productCode string, ouId int64) (models.PolicyOuProductDepositCounterWithRules, error) {
	var result models.PolicyOuProductDepositCounterWithRules

	filter := bson.M{
		"ouId":        ouId,
		"productCode": productCode,
	}

	data := ctx.RepoDB.MongoDB.Collection(constans.TABLE_POLICY_OU_PRODUCT_DEPOSIT_COUNTER).FindOne(ctx.RepoDB.Context, filter)
	if data.Err() != nil {
		fmt.Println(data.Err())
		return result, errors.New("Policy Ou Product Deposit Counter Not Found")
	}

	err := data.Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (ctx productRepository) GetPolicyOuProductDepositByOuId(OuID int64) (*[]models.PolicyOuProductDepositCounterWithRules, error) {
	var result []models.PolicyOuProductDepositCounterWithRules

	filter := bson.M{
		"ouId": OuID,
	}

	collect := ctx.RepoDB.MongoDB.Collection(constans.TABLE_POLICY_OU_PRODUCT_DEPOSIT_COUNTER)
	data, err := collect.Find(ctx.RepoDB.Context, filter)
	if err != nil {
		return nil, errors.New("policy ou product deposit counter not found")
	}

	defer data.Close(ctx.RepoDB.Context)
	for data.Next(ctx.RepoDB.Context) {
		var val models.PolicyOuProductDepositCounterWithRules
		data.Decode(&val)

		result = append(result, val)
	}

	return &result, nil
}

func (ctx productRepository) GetPolicyOuProductDepositListAdvance(product models.RequestGetPolicyOuProductList) ([]models.PolicyOuProductDepositCounterWithRules, error) {
	var result []models.PolicyOuProductDepositCounterWithRules
	filter := make(map[string]interface{})

	if product.Keyword != constans.EMPTY_VALUE {
		filter["$or"] = []bson.M{
			bson.M{"productCode": bson.M{
				"$regex": primitive.Regex{
					Pattern: product.Keyword,
					Options: "i",
				},
			}},
			bson.M{"productName": bson.M{
				"$regex": primitive.Regex{
					Pattern: product.Keyword,
					Options: "i",
				},
			}},
		}
	}

	options := options.FindOptions{}
	if product.ColumnOrderName != constans.EMPTY_VALUE {
		if product.AscDesc == constans.ASCENDING {
			options.Sort = bson.M{product.ColumnOrderName: 1}
		} else if product.AscDesc == constans.DESCENDING {
			options.Sort = bson.M{product.ColumnOrderName: -1}
		} else {
			options.Sort = bson.M{"productCode": -1}
		}
	}

	collect := ctx.RepoDB.MongoDB.Collection(constans.TABLE_POLICY_OU_PRODUCT_DEPOSIT_COUNTER)
	data, err := collect.Find(ctx.RepoDB.Context, filter, &options)

	if err != nil {
		return nil, err
	}

	defer data.Close(ctx.RepoDB.Context)

	for data.Next(ctx.RepoDB.Context) {
		var val models.PolicyOuProductDepositCounterWithRules
		data.Decode(&val)

		result = append(result, val)
	}

	return result, nil
}
