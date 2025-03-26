package redis

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client
var Ctx = context.Background()

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: "",
		DB:       0,
	})

	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis!")

	// err = Rdb.FlushDB(Ctx).Err()
	// if err != nil {
	// 	fmt.Println("Failed to flush Redis DB:", err)
	// } else {
	// 	fmt.Println("Redis DB flushed successfully")
	// }
}
