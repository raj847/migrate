package memberRepository

import (
	"togrpc/constans"
	"togrpc/models"
	"togrpc/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type memberMongoRepository struct {
	RepoDB repositories.Repository
}

func NewMemberMongoRepository(repo repositories.Repository) memberMongoRepository {
	return memberMongoRepository{
		RepoDB: repo,
	}
}

func (ctx memberMongoRepository) AddTrxMember(trxMember models.TrxMember) error {
	_, err := ctx.RepoDB.MongoDB.Collection(constans.TRX_MEMBER_COLLECTIONS).InsertOne(ctx.RepoDB.Context, trxMember)
	if err != nil {
		return err
	}

	return nil
}

func (ctx memberMongoRepository) AddMemberTemp(memberData models.Member) (*primitive.ObjectID, error) {
	var ID primitive.ObjectID

	result, err := ctx.RepoDB.MongoDB.Collection(constans.TRX_MEMBER_TEMPORARY).
		InsertOne(ctx.RepoDB.Context, memberData)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		ID = oid
	}

	return &ID, nil
}

// func (ctx memberMongoRepository) AddMemberList(member interface{}) error {
// 	var requestMember []interface{}

// 	requestMember = append(requestMember, member)

// 	_, err := ctx.RepoDB.MongoDB.Collection(constans.MEMBER_COLLECTIONS).InsertMany(ctx.RepoDB.Context, requestMember)
// 	if err != nil {
// 		return err
// 	}

// 	return nil

// }

// func (ctx memberMongoRepository) IsMemberActiveExistsByUUIDCard(uuidCard string, dateNow string) (models.Member, bool) {
// 	var result models.Member

// 	fmt.Println("DateNow", dateNow)

// 	filter := bson.M{
// 		"cardNumber": uuidCard,
// 		"dateFrom":   bson.M{"$gte": dateNow},
// 		"dateTo":     bson.M{"$lte": dateNow},
// 		"active":     constans.YES,
// 	}

// 	err := ctx.RepoDB.MongoDB.Collection(constans.MEMBER_COLLECTIONS).FindOne(ctx.RepoDB.Context, filter).Decode(&result)
// 	if err != nil {
// 		return result, false
// 	}

// 	return result, true
// }

// func (ctx memberMongoRepository) IsPartnerExistsByPartnerCode(partnerCode string) ([]models.ResponseFindPartner, bool) {
// 	var result []models.ResponseFindPartner

// 	filter := bson.M{
// 		"partnerCode": partnerCode,
// 	}

// 	projections := bson.M{
// 		"partnerId":   1,
// 		"dateFrom":    1,
// 		"dateTo":      1,
// 		"cardNumber":  1,
// 		"ouId":        1,
// 		"productId":   1,
// 		"partnerCode": 1,
// 	}

// 	options := options.FindOptions{
// 		Projection: projections,
// 	}

// 	data, err := ctx.RepoDB.MongoDB.Collection(constans.MEMBER_COLLECTIONS).
// 		Find(ctx.RepoDB.Context, filter, &options)
// 	if err != nil {
// 		return result, false
// 	}

// 	for data.Next(ctx.RepoDB.Context) {
// 		var val models.ResponseFindPartner
// 		data.Decode(&val)

// 		result = append(result, val)
// 	}

// 	if len(result) == 0 {
// 		return result, false
// 	}

// 	return result, true
// }

// func (ctx memberMongoRepository) IsPartnerExistsByID(ID primitive.ObjectID) (models.ResponseIsPartnerExistsByID, bool) {
// 	var result models.ResponseIsPartnerExistsByID

// 	filter := bson.M{
// 		"_id": ID,
// 	}

// 	projections := bson.M{
// 		"partnerCode": 1,
// 		"dateFrom":    1,
// 		"dateTo":      1,
// 		"activeAt":    1,
// 		"ouId":        1,
// 	}

// 	err := ctx.RepoDB.MongoDB.Collection(constans.MEMBER_COLLECTIONS).FindOne(ctx.RepoDB.Context, filter,
// 		options.FindOne().SetProjection(projections)).Decode(&result)
// 	if err != nil {
// 		return result, false
// 	}

// 	return result, true
// }

// func (ctx memberMongoRepository) ActivationPartner(ID primitive.ObjectID, memberNew models.EditPartner) error {
// 	filter := bson.M{
// 		"_id": ID,
// 	}

// 	update := bson.M{
// 		"$set": bson.M{
// 			"active":      memberNew.Active,
// 			"activeAt":    memberNew.ActiveAt,
// 			"nonActiveAt": memberNew.NonActiveAt,
// 		},
// 	}

// 	err := ctx.RepoDB.MongoDB.Collection(constans.MEMBER_COLLECTIONS).
// 		FindOneAndUpdate(ctx.RepoDB.Context, filter, update)
// 	if err.Err() != nil {
// 		return err.Err()
// 	}

// 	return nil
// }

// func (ctx memberMongoRepository) RemoveMemberByID(ID primitive.ObjectID) error {
// 	filter := bson.M{
// 		"_id": ID,
// 	}

// 	err := ctx.RepoDB.MongoDB.Collection(constans.MEMBER_COLLECTIONS).
// 		FindOneAndDelete(ctx.RepoDB.Context, filter)
// 	if err.Err() != nil {
// 		return err.Err()
// 	}

// 	return nil
// }

// func (ctx memberMongoRepository) FindMemberAdvance(keyword string, limit int64, offset int64) ([]models.FindMember, error) {
// 	var result []models.FindMember

// 	filter := make(map[string]interface{})

// 	if keyword != constans.EMPTY_VALUE {
// 		filter = bson.M{
// 			"$or": []bson.M{
// 				bson.M{"partnerCode": bson.M{
// 					"$regex": primitive.Regex{
// 						Pattern: keyword,
// 						Options: "i",
// 					},
// 				}},
// 				bson.M{"firstName": bson.M{
// 					"$regex": primitive.Regex{
// 						Pattern: keyword,
// 						Options: "i",
// 					},
// 				}},
// 				bson.M{"lastName": bson.M{
// 					"$regex": primitive.Regex{
// 						Pattern: keyword,
// 						Options: "i",
// 					},
// 				}},
// 			},
// 		}
// 	}

// 	options := options.FindOptions{
// 		Limit: &limit,
// 		Skip:  &offset,
// 	}

// 	data, err := ctx.RepoDB.MongoDB.Collection(constans.MEMBER_COLLECTIONS).
// 		Find(ctx.RepoDB.Context, filter, &options)
// 	if err != nil {
// 		return result, err
// 	}

// 	if err := data.All(ctx.RepoDB.Context, &result); err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }

// func (ctx memberMongoRepository) GetMemberForTrxOutstandingByUUID(uuidCard string) ([]models.TrxMember, bool) {
// 	var result []models.TrxMember

// 	filter := bson.M{
// 		"cardNumber": uuidCard,
// 	}

// 	data, err := ctx.RepoDB.MongoDB.Collection(constans.TRX_MEMBER_COLLECTIONS).Find(ctx.RepoDB.Context, filter)
// 	if err != nil {
// 		return result, false
// 	}

// 	for data.Next(ctx.RepoDB.Context) {
// 		var val models.TrxMember
// 		data.Decode(val)

// 		result = append(result, val)
// 	}

// 	return result, true
// }
