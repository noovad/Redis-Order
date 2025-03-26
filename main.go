package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"redis-order/handlers"
	redisdb "redis-order/redis"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	redisdb.InitRedis()

	r := gin.Default()
	r.POST("/order", handlers.AddorderHandler)
	r.GET("/order", handlers.GetOrderDataHandler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go handlers.StartBatchOrderProcessor(ctx)
	go handlers.NotifyOrderProcessed(ctx)

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM) // Menangkap Ctrl+C (SIGINT)
		<-sig                                               // Tunggu sinyal
		fmt.Println("\nReceived shutdown signal. Stopping the server...")

		cancel() // Panggil cancel untuk menghentikan goroutine yang mendengarkan context
	}()

	log.Fatal(r.Run(":8080"))
}
