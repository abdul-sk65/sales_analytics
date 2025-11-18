# Sales Analytics System

A high-performance sales analytics system built with Go, Fiber framework, and MongoDB. This system loads large CSV files containing sales data, stores it in a normalized database schema, and provides RESTful APIs for revenue analysis.

## Features

- **Normalized Database Schema**: Separate collections for Customers, Products, and Orders
- **Efficient Data Loading**: Worker pool implementation for parallel CSV processing
- **RESTful API**: Clean API endpoints for data refresh and analytics
- **Revenue Analytics**: Calculate total revenue, revenue by product/category/region
- **Data Refresh Mechanism**: On-demand data refresh with comprehensive logging
- **Robust Error Handling**: Graceful error management throughout the application
- **Performance Optimized**: Database indexes for fast query execution

## Architecture

```
kenshilabs/
├── cmd/
│   └── main.go              # Application entry point
├── config/
│   └── config.go            # Configuration management
├── pkg/
│   └── repository/
│       ├── models.go        # Data models
│       ├── repository.go    # Database operations
│       ├── loader.go        # CSV loading with worker pool
│       └── analytics.go     # Revenue calculations
├── api/
│   ├── routes.go            # Route definitions
│   └── handler.go           # HTTP handlers
├── data/
│   └── sales_data.csv       # Sample CSV data
├── .env                     # Environment variables
├── go.mod                   # Go module definition
└── README.md               # This file
```

## Database Schema

### Collections

1. **customers**: Stores customer information

   - customer_id (unique)
   - name
   - email
   - address

2. **products**: Stores product information

   - product_id (unique)
   - name
   - category
   - unit_price
   - discount

3. **orders**: Stores order transactions

   - order_id (unique)
   - product_id (foreign key)
   - customer_id (foreign key)
   - region
   - date_of_sale
   - quantity_sold
   - shipping_cost
   - payment_method

4. **refresh_logs**: Tracks data refresh operations
   - start_time
   - end_time
   - status
   - rows_loaded
   - error_msg

## Setup

### Prerequisites

- Go 1.21 or higher
- MongoDB 4.4 or higher
- Git

### Installation

1. Clone the repository:

```bash
git clone https://github.com/abdul-sk65/sales_analytics.git
cd sales_analytics
```

2. Install dependencies:

```bash
go mod download
```

3. Create a `.env` file:

```bash
cp .env.example .env
```

4. Configure environment variables in `.env`:

```env
MONGODB_URI=mongodb://localhost:27017
DATABASE_NAME=sales_analytics
CSV_FILE_PATH=./data/sales_data.csv
PORT=8080
WORKER_POOL_SIZE=10
```

5. Create data directory and add CSV file:

```bash
mkdir -p data
# Place your sales_data.csv file in the data directory
```

6. Start MongoDB:

```bash
# If using Docker
docker run -d -p 27017:27017 --name mongodb mongo:latest

# Or start your local MongoDB instance
mongod
```

7. Run the application:

```bash
go run cmd/main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Health Check

**GET** `/health`

Returns the health status of the API.

**Response:**

```json
{
  "status": "healthy",
  "time": "2024-01-15T10:30:00Z"
}
```

### Data Refresh

**POST** `/api/v1/refresh`

Triggers a data refresh from the CSV file. The operation runs asynchronously in the background.

**Response:**

```json
{
  "message": "Data refresh initiated",
  "status": "processing"
}
```

### Get Refresh Logs

**GET** `/api/v1/refresh/logs`

Retrieves the latest 10 data refresh logs.

**Response:**

```json
{
  "logs": [
    {
      "start_time": "2024-01-15T10:00:00Z",
      "end_time": "2024-01-15T10:02:30Z",
      "status": "success",
      "rows_loaded": 1000,
      "error_msg": ""
    }
  ]
}
```

### Revenue Analytics

#### Total Revenue

**GET** `/api/v1/revenue/total?start_date=2023-01-01&end_date=2024-12-31`

Calculates total revenue for the specified date range.

**Query Parameters:**

- `start_date` (required): Start date in YYYY-MM-DD format
- `end_date` (required): End date in YYYY-MM-DD format

**Response:**

```json
{
  "start_date": "2023-01-01",
  "end_date": "2024-12-31",
  "total_revenue": 45678.9
}
```

#### Revenue by Product

**GET** `/api/v1/revenue/by-product?start_date=2023-01-01&end_date=2024-12-31`

Calculates revenue grouped by product, sorted by revenue (descending).

**Response:**

```json
{
  "start_date": "2023-01-01",
  "end_date": "2024-12-31",
  "products": [
    {
      "product_id": "P456",
      "product_name": "iPhone 15 Pro",
      "total_revenue": 12990.0
    },
    {
      "product_id": "P123",
      "product_name": "UltraBoost Running Shoes",
      "total_revenue": 324.0
    }
  ]
}
```

#### Revenue by Category

**GET** `/api/v1/revenue/by-category?start_date=2023-01-01&end_date=2024-12-31`

Calculates revenue grouped by product category.

**Response:**

```json
{
  "start_date": "2023-01-01",
  "end_date": "2024-12-31",
  "categories": [
    {
      "category": "Electronics",
      "total_revenue": 35678.9
    },
    {
      "category": "Shoes",
      "total_revenue": 10000.0
    }
  ]
}
```

#### Revenue by Region

**GET** `/api/v1/revenue/by-region?start_date=2023-01-01&end_date=2024-12-31`

Calculates revenue grouped by region.

**Response:**

```json
{
  "start_date": "2023-01-01",
  "end_date": "2024-12-31",
  "regions": [
    {
      "region": "North America",
      "total_revenue": 20000.0
    },
    {
      "region": "Europe",
      "total_revenue": 15000.0
    }
  ]
}
```

## Testing

### Manual Testing with cURL

1. **Health Check:**

```bash
curl http://localhost:8080/health
```

2. **Trigger Data Refresh:**

```bash
curl -X POST http://localhost:8080/api/v1/refresh
```

3. **Get Refresh Logs:**

```bash
curl http://localhost:8080/api/v1/refresh/logs
```

4. **Get Total Revenue:**

```bash
curl "http://localhost:8080/api/v1/revenue/total?start_date=2023-01-01&end_date=2024-12-31"
```

5. **Get Revenue by Product:**

```bash
curl "http://localhost:8080/api/v1/revenue/by-product?start_date=2023-01-01&end_date=2024-12-31"
```

6. **Get Revenue by Category:**

```bash
curl "http://localhost:8080/api/v1/revenue/by-category?start_date=2023-01-01&end_date=2024-12-31"
```

7. **Get Revenue by Region:**

```bash
curl "http://localhost:8080/api/v1/revenue/by-region?start_date=2023-01-01&end_date=2024-12-31"
```

## Performance Considerations

### Worker Pool

The CSV loader uses a worker pool pattern to process records in parallel:

- Configurable worker pool size (default: 10 workers)
- Buffered channels to prevent blocking
- Concurrent upserts to MongoDB
- Graceful error handling with early termination

### Database Indexes

The system creates the following indexes for optimal query performance:

- Customers: `customer_id` (unique), `email`
- Products: `product_id` (unique), `category`
- Orders: `order_id` (unique), `customer_id`, `product_id`, `date_of_sale`, `region`

### Revenue Calculation Formula

Revenue is calculated as:

```
Revenue = Quantity × (Unit Price - (Unit Price × Discount))
```

This accounts for discounts applied to each order.

## Error Handling

The application implements comprehensive error handling:

- CSV parsing errors are logged and propagated
- Database connection failures are caught at startup
- API errors return appropriate HTTP status codes
- Data refresh failures are logged in the database
- Worker pool errors terminate processing gracefully

## Logging

The application logs:

- Server startup and shutdown events
- Data refresh operations (start, progress, completion)
- Worker pool activities
- Database operations
- API request/response cycles
- Errors and warnings

## Future Enhancements

Potential improvements for production use:

- Add unit and integration tests
- Implement rate limiting
- Add API authentication/authorization
- Support for incremental data updates
- Real-time data streaming
- Caching layer for frequently accessed data
- Scheduled automatic refreshes
- More analytics endpoints (customer lifetime value, top customers, etc.)
- Export results to CSV/Excel
- Dashboard UI

## License

MIT License

## Support

For issues and questions, please open an issue in the repository.
