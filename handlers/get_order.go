package handlers

import (
	"fmt"
	"net/http"
	redisdb "redis-order/redis"
	"time"

	"github.com/gin-gonic/gin"
)

// Simulate fetching order data from a database (dummy).
func fetchOrderFromDatabase() (string, string, error) {
	time.Sleep(3 * time.Second)

	product := "Dummy Product"
	price := "100"

	return product, price, nil
}

// GetOrderDataHandler retrieves the order data from Redis or database if not found.
func GetOrderDataHandler(c *gin.Context) {
	orderData, err := redisdb.Rdb.HGetAll(redisdb.Ctx, "order-99").Result()

	if len(orderData) == 0 || err != nil {
		fmt.Printf("No data found for order %s. Fetching from database...\n", "order-99")
		product, price, dbErr := fetchOrderFromDatabase()
		if dbErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching data from database"})
			return
		}

		err := redisdb.Rdb.HSet(redisdb.Ctx, "order-99", "product", product, "price", price).Err()
		if err != nil {
			fmt.Println("Failed to cache order data:", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"order_id": "order-99",
			"product":  product,
			"price":    price,
		})

	} else {
		product, productExists := orderData["product"]
		price, priceExists := orderData["price"]

		if !productExists || !priceExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Incomplete data for order"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"order_id": "order-99",
			"product":  product,
			"price":    price,
		})
	}
}
