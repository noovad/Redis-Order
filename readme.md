# Redis Order Processing API

This is a **learning project** aimed at understanding **Redis** and its application in **order processing**. The project is built using **Go (Golang)** and the **Gin** web framework. It demonstrates how to use Redis for **caching, queueing, and pub/sub messaging** in an order processing system.

## Features

- **Add Orders**: Orders are stored in Redis with a unique ID and placed in a queue.
- **Fetch Orders**: Retrieves order data from Redis or simulates fetching from a database if not available.
- **Batch Processing**: Periodically processes orders from the queue and updates their status.
- **Order Notifications**: Uses Redis Pub/Sub to notify when an order is processed.

## Project Structure

```
redis-order/
│── handlers/        # API handlers for order operations
│── redis/           # Redis initialization and connection
│── main.go          # Main application entry point
│── .env             # Environment variables (e.g., Redis connection)
│── go.mod, go.sum   # Go module files
```

## Installation & Setup

1. **Install Redis** (if not already installed)
2. **Run Redis** locally:
3. **Clone the Repository**
4. **Set up Environment Variables**
   Create a `.env` file and set the Redis host:
   ```
   REDIS_HOST=localhost:6379
   ```
5. **Run the Project**
   ```
   go mod tidy
   go run main.go
   ```

## API Endpoints

### 1. **Add Order**

**`POST /order`**

- Adds an order to the Redis queue
- **Request Body** (form-data):
  ```
  product: "Laptop"
  price: "1500"
  ```
- **Response:**
  ```json
  {
    "message": "Order successfully added!",
    "order_id": "order-1"
  }
  ```

### 2. **Get Order Data**

**`GET /order`**

- Retrieves order details from Redis or fetches from the database if missing.
- **Response:**
  ```json
  {
    "orderKey": "order-99",
    "product": "Dummy Product",
    "price": "100"
  }
  ```