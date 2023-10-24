package services

import (
	"database/sql"
	"togrpc/repositories"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/go-redis/redis"
	"github.com/mkpproduction/mkp-sdk-go/mkp/genautonum"
	"github.com/streadway/amqp"
)

type UsecaseService struct {
	RepoDB                     *sql.DB
	repoDBCloud                *sql.DB
	DynamoDB                   *dynamodb.DynamoDB
	RabbitMQ                   *amqp.Channel
	RedisClient                *redis.Client
	RedisClientLocal           *redis.Client
	GenAutoNumRepo             genautonum.GenerateAutonumberRepository
	ProductRepo                repositories.ProductMongoRepository
	TrxMongoRepo               repositories.TrxMongoRepository
	MemberMongoRepo            repositories.MemberMongoRepository
	MemberRepo                 repositories.MemberRepository
	OuMongoRepo                repositories.OuMongoRepository
	ReportTrxMongoRepo         repositories.ReportTrxMongoRepository
	UserMongoRepo              repositories.UserMongoRepository
	TrxMongoDepositCounterRepo repositories.TrxMongoDepositCounterRepository
	ProductMembershipRepo      repositories.ProductMembershipRepository
	CommonRepo                 repositories.CommonRepository
}

func NewUsecaseService(
	repoDB *sql.DB,
	repoDBCloud *sql.DB,
	dynamodb *dynamodb.DynamoDB,
	rabbitMQ *amqp.Channel,
	Redis *redis.Client,
	RedisClientLocal *redis.Client,
	GenAutoNumRepo genautonum.GenerateAutonumberRepository,
	ProductRepo repositories.ProductMongoRepository,
	TrxMongoRepo repositories.TrxMongoRepository,
	MemberMongoRepo repositories.MemberMongoRepository,
	MemberRepo repositories.MemberRepository,
	OuMongoRepo repositories.OuMongoRepository,
	ReportTrxMongoRepo repositories.ReportTrxMongoRepository,
	UserMongoRepo repositories.UserMongoRepository,
	TrxMongoDepositCounterRepo repositories.TrxMongoDepositCounterRepository,
	ProductMembershipRepo repositories.ProductMembershipRepository,
	CommonRepo repositories.CommonRepository,

) UsecaseService {
	return UsecaseService{
		RepoDB:                     repoDB,
		repoDBCloud:                repoDBCloud,
		DynamoDB:                   dynamodb,
		RabbitMQ:                   rabbitMQ,
		RedisClient:                Redis,
		RedisClientLocal:           RedisClientLocal,
		GenAutoNumRepo:             GenAutoNumRepo,
		ProductRepo:                ProductRepo,
		TrxMongoRepo:               TrxMongoRepo,
		MemberMongoRepo:            MemberMongoRepo,
		MemberRepo:                 MemberRepo,
		OuMongoRepo:                OuMongoRepo,
		ReportTrxMongoRepo:         ReportTrxMongoRepo,
		UserMongoRepo:              UserMongoRepo,
		TrxMongoDepositCounterRepo: TrxMongoDepositCounterRepo,
		ProductMembershipRepo:      ProductMembershipRepo,
		CommonRepo:                 CommonRepo,
	}
}
