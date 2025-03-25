package handlers

import (
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
	status := "pending"

	redisdb.Rdb.HSet(redisdb.Ctx, key, "prodcut", product, "price", price, "status", status)
	redisdb.Rdb.Expire(redisdb.Ctx, key, time.Minute*15)

	redisdb.Rdb.Publish(redisdb.Ctx, "order_channel", "New order: "+key)

	c.JSON(http.StatusOK, gin.H{"message": "Order successfully added!"})
}
