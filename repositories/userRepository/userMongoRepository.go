package userRepository

import (
	"togrpc/constans"
	"togrpc/models"
	"togrpc/repositories"

	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type userMongoRepository struct {
	RepoDB repositories.Repository
}

func NewUserMongoRepository(repo repositories.Repository) userMongoRepository {
	return userMongoRepository{
		RepoDB: repo,
	}
}

func (ctx userMongoRepository) FindUserByIndex(username string) (*models.UserLoginLocal, error) {
	var result models.UserLoginLocal
	filter := bson.M{
		"user.username": username,
	}

	userLogin := ctx.RepoDB.MongoDB.Collection(constans.TABLE_USER_LOGIN)
	err := userLogin.FindOne(ctx.RepoDB.Context, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ctx userMongoRepository) FindDeviceIdByIndex(username, deviceId string) (*models.DeviceSync, error) {
	var result models.DeviceSync

	filter := bson.M{
		"user.username": username,
		"devicelist": bson.M{
			"$elemMatch": bson.M{
				"deviceid": deviceId,
			},
		},
	}

	projection := bson.M{
		"devicelist.$": 1,
	}

	userLogin := ctx.RepoDB.MongoDB.Collection(constans.TABLE_USER_LOGIN)
	err := userLogin.FindOne(ctx.RepoDB.Context, filter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
