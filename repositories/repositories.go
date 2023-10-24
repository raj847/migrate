package repositories

import (
	"context"
	"database/sql"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	DB        *sql.DB
	DBCloud   *sql.DB
	MongoDB   *mongo.Database
	MongoDBPg *mongo.Database
	Context   context.Context
	DynamoDB  *dynamodb.DynamoDB
}

func NewRepository(conn *sql.DB, connCloud *sql.DB, MongoDB *mongo.Database, MongoDBPg *mongo.Database, ctx context.Context, dynamodb *dynamodb.DynamoDB) Repository {
	return Repository{
		DB:        conn,
		DBCloud:   connCloud,
		MongoDB:   MongoDB,
		MongoDBPg: MongoDBPg,
		Context:   ctx,
		DynamoDB:  dynamodb,
	}
}
