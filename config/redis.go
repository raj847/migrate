package config

import (
	"flag"
	"fmt"

	"github.com/go-redis/redis"
)

var connection redis.Client

func ConnectRedis() *redis.Client {
	connStr := fmt.Sprintf("%s:%s", REDISHost, REDISPort)
	fmt.Println("Connection Redis:", connStr)
	var addr = flag.String("Server", connStr, "Redis server address")

	rdb := redis.NewClient(&redis.Options{
		Addr:     *addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return rdb
}

func ConnectRedisLocal() *redis.Client {
	connStr := fmt.Sprintf("%s:%s", REDISHostLocal, REDISPortLocal)
	fmt.Println("Connection Redis Local:", connStr)
	var addr = flag.String("Server Local", connStr, "Redis server address")

	rdb := redis.NewClient(&redis.Options{
		Addr:     *addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return rdb
}

func CloseRedis() {
	connection.Close()
}
