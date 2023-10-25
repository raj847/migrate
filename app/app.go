package app

import (
	"database/sql"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/go-redis/redis"
	"github.com/mkpproduction/mkp-sdk-go/mkp/genautonum"
	"github.com/streadway/amqp"
	"togrpc/repositories"
	"togrpc/repositories/commonRepository"
	"togrpc/repositories/memberRepository"
	"togrpc/repositories/ouRepository"
	"togrpc/repositories/productMembershipRepository"
	"togrpc/repositories/productRepository"
	"togrpc/repositories/trxRepository"
	"togrpc/repositories/userRepository"
	"togrpc/services"
)

func SetupApp(DB *sql.DB, DBCloud *sql.DB, repo repositories.Repository, svc *dynamodb.DynamoDB, rabbitMQ *amqp.Channel, redis *redis.Client, redisLocal *redis.Client, repoGenAutoNum genautonum.Repository) services.UsecaseService {
	// Create table for dynamo db
	//CreateTableMovies(svc)

	// Services
	productRepo := productRepository.NewProductRepository(repo)
	trxMongoRepo := trxRepository.NewTrxMongoRepository(repo)
	genAutoNumRepo := genautonum.NewGenerateAutonumberRepository(repoGenAutoNum)
	memberMongoRepo := memberRepository.NewMemberMongoRepository(repo)
	memberRepo := memberRepository.NewMemberRepository(repo)
	ouMongoRepo := ouRepository.NewOuMongoRepository(repo)
	userMongoRepo := userRepository.NewUserMongoRepository(repo)
	reportTrxRepo := trxRepository.NewReportTrxMongoRepository(repo)
	trxMongoDepositCounterRepo := trxRepository.NewTrxMongoDepositCounterRepository(repo)
	commonRepo := commonRepository.NewCommonRepository(repo)
	productMembershipRepo := productMembershipRepository.NewProductMembershipRepository(repo)

	usecaseSvc := services.NewUsecaseService(DB, DBCloud, svc, rabbitMQ, redis, redisLocal, genAutoNumRepo,
		productRepo, trxMongoRepo, memberMongoRepo, memberRepo, ouMongoRepo, reportTrxRepo, userMongoRepo, trxMongoDepositCounterRepo, productMembershipRepo,
		commonRepo)

	return usecaseSvc
}
