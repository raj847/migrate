package ouRepository

import (
	"togrpc/constans"
	"togrpc/models"
	"togrpc/repositories"

	"gopkg.in/mgo.v2/bson"
)

type ouMongoRepository struct {
	RepoDB repositories.Repository
}

func NewOuMongoRepository(repo repositories.Repository) ouMongoRepository {
	return ouMongoRepository{
		RepoDB: repo,
	}
}

func (ctx ouMongoRepository) FindOuByOuId(OuId int64) (models.Ou, error) {
	var result models.Ou

	filter := bson.M{
		"ouId": OuId,
	}

	err := ctx.RepoDB.MongoDB.Collection(constans.OU_STRUCTURE_COLLECTIONS).
		FindOne(ctx.RepoDB.Context, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}
