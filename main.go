package main

import (
	"context"
	"fmt"
	"github.com/labstack/gommon/color"
	"github.com/mkpproduction/mkp-sdk-go/mkp/genautonum"
	"github.com/raj847/togrpc/app"
	"github.com/raj847/togrpc/config"
	"github.com/raj847/togrpc/proto/trx"
	"github.com/raj847/togrpc/repositories"
	"github.com/raj847/togrpc/services/trxService"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	ctx = context.Background()
)

func main() {
	if err := config.OpenConnection(); err != nil {
		panic(fmt.Sprintf("Open Connection Faild: %s", err.Error()))
	}
	defer config.CloseConnectionDB()

	//if err := config.OpenConnectionCloud(); err != nil {
	//	panic(fmt.Sprintf("Open Connection Faild: %s", err.Error()))
	//}
	//defer config.CloseConnectionDBCloud()

	// Mongo DB connection using database default
	mongoDB := config.ConnectMongo(ctx)
	defer config.CloseMongo(ctx)

	// Connection database
	DB := config.DBConnection()
	//DBCloud := config.DBConnectionCloud()

	repoGenAutoNum := genautonum.NewRepository(nil, ctx, mongoDB)

	//Redis Connection
	redis := config.ConnectRedis()
	redisLocal := config.ConnectRedisLocal()

	// Configuration Repository
	repo := repositories.NewRepository(DB, nil, mongoDB, nil, ctx, nil)

	// Configuration Repository and Services
	services := app.SetupApp(DB, nil, repo, nil, nil, redis, redisLocal, repoGenAutoNum)

	listen, err := net.Listen("tcp", fmt.Sprintf("%s%s", ":", config.GetEnv("APP_PORT", "6000")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	trxSvc := trxService.NewTrxService(services)
	trx.RegisterTrxServiceServer(grpcServer, trxSvc)

	boarding()
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func boarding() {
	txt := `
  __  __   _  __  _____      _____    _____     _____   
 |  \/  | | |/ / |  __ \    |  __ \  |  __ \   / ____|  
 | \  / | | ' /  | |__) |   | |__) | | |__) | | |       
 | |\/| | |  <   |  ___/    |  _  /  |  ___/  | |       
 | |  | | | . \  | |        | | \ \  | |      | |____   
 |_|  |_| |_|\_\ |_|        |_|  \_\ |_|       \_____|`
	fmt.Println(txt)

	println()
	println(fmt.Sprintf("%s%s%s", "http server started on ", color.Green("[::]:", ""), color.Green(config.GetEnv("APP_PORT", "6000", ""))))
}
