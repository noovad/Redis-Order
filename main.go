package main

import (
	"log"
	"redis-order/handlers"
	redisdb "redis-order/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	redisdb.InitRedis()

	r := gin.Default()
	r.POST("/order", handlers.AddorderHandler)

	log.Fatal(r.Run(":8080"))
}
