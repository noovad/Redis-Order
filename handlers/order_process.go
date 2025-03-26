package handlers

import (
	"context"
	"fmt"
	"net/http"
	redisdb "redis-order/redis"
	"time"

	"github.com/gin-gonic/gin"
)

func AddorderHandler(c *gin.Context) {
	id, err := redisdb.Rdb.Incr(redisdb.Ctx, "order_id_counter").Result()
	if err != nil {
		fmt.Println("Error incrementing ID:", err)
		return
	}

	key := fmt.Sprintf("order-%d", id)
	product := c.PostForm("product")
	price := c.PostForm("price")

	if product == "" || price == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product and price are required"})
		return
	}

	pipe := redisdb.Rdb.TxPipeline()

	pipe.RPush(redisdb.Ctx, "order_queue", key)
	pipe.HSet(redisdb.Ctx, key, "product", product, "price", price)
	pipe.Expire(redisdb.Ctx, "order_queue", 10*time.Minute)

	_, err = pipe.Exec(redisdb.Ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed adding order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order successfully added!", "order_id": id})
}

func StartBatchOrderProcessor(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Exiting batch order processor")
			return
		default:
			orderBatch, err := redisdb.Rdb.LRange(redisdb.Ctx, "order_queue", 0, -1).Result()
			if err != nil {
				fmt.Println("Failed to fetch orders:", err)
				return
			}

			if len(orderBatch) > 0 {
				pipe := redisdb.Rdb.TxPipeline()

				for _, orderKey := range orderBatch {
					pipe.HSet(redisdb.Ctx, orderKey, "status", "processed")
					pipe.Publish(redisdb.Ctx, "order_processed", orderKey)
				}

				_, err = pipe.Exec(redisdb.Ctx)
				if err != nil {
					fmt.Println("Error executing transaction:", err)
					return
				} else {
					_, err = redisdb.Rdb.LTrim(redisdb.Ctx, "order_queue", int64(len(orderBatch)), -1).Result()
					if err != nil {
						fmt.Println("Error removing processed orders from queue:", err)
						return
					}
				}
			}

			time.Sleep(5 * time.Minute)
		}
	}
}

// Simulate sending a notification, e.g., to a user or another system
func NotifyOrderProcessed(ctx context.Context) {
	subscriber := redisdb.Rdb.Subscribe(redisdb.Ctx, "order_processed")
	defer subscriber.Close()

	for {
		select {
		case msg := <-subscriber.Channel():
			orderKey := msg.Payload
			orderData, err := redisdb.Rdb.HGetAll(redisdb.Ctx, orderKey).Result()
			if err != nil {
				fmt.Printf("Error fetching data for order %s: %v\n", orderKey, err)
				continue
			}

			product, productExists := orderData["product"]
			price, priceExists := orderData["price"]

			if !productExists || !priceExists {
				fmt.Printf("Incomplete data for order %s. Product or price missing.\n", orderKey)
				continue
			}

			fmt.Printf("Order %s has been processed. Product: %s, Price: %s. Sending notification...\n", orderKey, product, price)

		case <-ctx.Done():
			fmt.Println("Exiting notification handler")
			return
		}
	}
}
