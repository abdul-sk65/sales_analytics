package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Customer represents a customer entity
type Customer struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	CustomerID string             `bson:"customer_id" json:"customer_id"`
	Name       string             `bson:"name" json:"name"`
	Email      string             `bson:"email" json:"email"`
	Address    string             `bson:"address" json:"address"`
}

// Product represents a product entity
type Product struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ProductID string             `bson:"product_id" json:"product_id"`
	Name      string             `bson:"name" json:"name"`
	Category  string             `bson:"category" json:"category"`
	UnitPrice float64            `bson:"unit_price" json:"unit_price"`
	Discount  float64            `bson:"discount" json:"discount"`
}

// Order represents an order entity
type Order struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	OrderID       string             `bson:"order_id" json:"order_id"`
	ProductID     string             `bson:"product_id" json:"product_id"`
	CustomerID    string             `bson:"customer_id" json:"customer_id"`
	Region        string             `bson:"region" json:"region"`
	DateOfSale    time.Time          `bson:"date_of_sale" json:"date_of_sale"`
	QuantitySold  int                `bson:"quantity_sold" json:"quantity_sold"`
	ShippingCost  float64            `bson:"shipping_cost" json:"shipping_cost"`
	PaymentMethod string             `bson:"payment_method" json:"payment_method"`
}

// RefreshLog represents a data refresh log entry
type RefreshLog struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	StartTime  time.Time          `bson:"start_time" json:"start_time"`
	EndTime    time.Time          `bson:"end_time" json:"end_time"`
	Status     string             `bson:"status" json:"status"` // success, failed
	RowsLoaded int                `bson:"rows_loaded" json:"rows_loaded"`
	ErrorMsg   string             `bson:"error_msg,omitempty" json:"error_msg,omitempty"`
}

// CSVRecord represents a row from the CSV file
type CSVRecord struct {
	OrderID       string
	ProductID     string
	CustomerID    string
	ProductName   string
	Category      string
	Region        string
	DateOfSale    string
	QuantitySold  string
	UnitPrice     string
	Discount      string
	ShippingCost  string
	PaymentMethod string
	CustomerName  string
	CustomerEmail string
	CustomerAddr  string
}
