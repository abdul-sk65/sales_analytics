package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RevenueResult revenue calculation result
type RevenueResult struct {
	TotalRevenue float64 `bson:"total_revenue" json:"total_revenue"`
}

// ProductRevenueResult revenue by product
type ProductRevenueResult struct {
	ProductID    string  `bson:"_id" json:"product_id"`
	ProductName  string  `bson:"product_name" json:"product_name"`
	TotalRevenue float64 `bson:"total_revenue" json:"total_revenue"`
}

// CategoryRevenueResult revenue by category
type CategoryRevenueResult struct {
	Category     string  `bson:"_id" json:"category"`
	TotalRevenue float64 `bson:"total_revenue" json:"total_revenue"`
}

// RegionRevenueResult revenue by region
type RegionRevenueResult struct {
	Region       string  `bson:"_id" json:"region"`
	TotalRevenue float64 `bson:"total_revenue" json:"total_revenue"`
}

// CalculateTotalRevenue total revenue for a date range
func (r *MongoRepository) CalculateTotalRevenue(ctx context.Context, startDate, endDate time.Time) (float64, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"date_of_sale": bson.M{
				"$gte": startDate,
				"$lte": endDate,
			},
		}}},
		// JOIN with products to get price and discount
		{{Key: "$lookup", Value: bson.M{
			"from":         "products",
			"localField":   "product_id",
			"foreignField": "product_id",
			"as":           "product",
		}}},
		{{Key: "$unwind", Value: "$product"}},
		{{Key: "$group", Value: bson.M{
			"_id": nil,
			"total_revenue": bson.M{
				"$sum": bson.M{
					"$multiply": bson.A{
						"$quantity_sold",
						bson.M{"$subtract": bson.A{
							"$product.unit_price",
							bson.M{"$multiply": bson.A{
								"$product.unit_price",
								"$product.discount",
							}},
						}},
					},
				},
			},
		}}},
	}

	cursor, err := r.GetCollection("orders").Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var results []RevenueResult
	if err := cursor.All(ctx, &results); err != nil {
		return 0, err
	}

	if len(results) > 0 {
		return results[0].TotalRevenue, nil
	}
	return 0, nil
}

// CalculateRevenueByProduct revenue grouped by product
func (r *MongoRepository) CalculateRevenueByProduct(ctx context.Context, startDate, endDate time.Time) ([]ProductRevenueResult, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"date_of_sale": bson.M{
				"$gte": startDate,
				"$lte": endDate,
			},
		}}},
		// JOIN with products to get price, discount, and name
		{{Key: "$lookup", Value: bson.M{
			"from":         "products",
			"localField":   "product_id",
			"foreignField": "product_id",
			"as":           "product",
		}}},
		{{Key: "$unwind", Value: "$product"}},
		{{Key: "$group", Value: bson.M{
			"_id":          "$product_id",
			"product_name": bson.M{"$first": "$product.name"},
			"total_revenue": bson.M{
				"$sum": bson.M{
					"$multiply": bson.A{
						"$quantity_sold",
						bson.M{"$subtract": bson.A{
							"$product.unit_price",
							bson.M{"$multiply": bson.A{
								"$product.unit_price",
								"$product.discount",
							}},
						}},
					},
				},
			},
		}}},
		{{Key: "$sort", Value: bson.M{"total_revenue": -1}}},
	}

	cursor, err := r.GetCollection("orders").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []ProductRevenueResult
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// CalculateRevenueByCategory revenue grouped by category
func (r *MongoRepository) CalculateRevenueByCategory(ctx context.Context, startDate, endDate time.Time) ([]CategoryRevenueResult, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"date_of_sale": bson.M{
				"$gte": startDate,
				"$lte": endDate,
			},
		}}},
		// JOIN with products to get category, price, and discount
		{{Key: "$lookup", Value: bson.M{
			"from":         "products",
			"localField":   "product_id",
			"foreignField": "product_id",
			"as":           "product",
		}}},
		{{Key: "$unwind", Value: "$product"}},
		{{Key: "$group", Value: bson.M{
			"_id": "$product.category",
			"total_revenue": bson.M{
				"$sum": bson.M{
					"$multiply": bson.A{
						"$quantity_sold",
						bson.M{"$subtract": bson.A{
							"$product.unit_price",
							bson.M{"$multiply": bson.A{
								"$product.unit_price",
								"$product.discount",
							}},
						}},
					},
				},
			},
		}}},
		{{Key: "$sort", Value: bson.M{"total_revenue": -1}}},
	}

	cursor, err := r.GetCollection("orders").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []CategoryRevenueResult
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// CalculateRevenueByRegion revenue grouped by region
func (r *MongoRepository) CalculateRevenueByRegion(ctx context.Context, startDate, endDate time.Time) ([]RegionRevenueResult, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"date_of_sale": bson.M{
				"$gte": startDate,
				"$lte": endDate,
			},
		}}},
		// JOIN with products to get price and discount
		{{Key: "$lookup", Value: bson.M{
			"from":         "products",
			"localField":   "product_id",
			"foreignField": "product_id",
			"as":           "product",
		}}},
		{{Key: "$unwind", Value: "$product"}},
		{{Key: "$group", Value: bson.M{
			"_id": "$region",
			"total_revenue": bson.M{
				"$sum": bson.M{
					"$multiply": bson.A{
						"$quantity_sold",
						bson.M{"$subtract": bson.A{
							"$product.unit_price",
							bson.M{"$multiply": bson.A{
								"$product.unit_price",
								"$product.discount",
							}},
						}},
					},
				},
			},
		}}},
		{{Key: "$sort", Value: bson.M{"total_revenue": -1}}},
	}

	cursor, err := r.GetCollection("orders").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []RegionRevenueResult
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// GetRefreshLogs retrieves the latest refresh logs
func (r *MongoRepository) GetRefreshLogs(ctx context.Context, limit int) ([]RefreshLog, error) {
	opts := options.Find().SetSort(bson.M{"start_time": -1}).SetLimit(int64(limit))

	cursor, err := r.GetCollection("refresh_logs").Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []RefreshLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, err
	}

	return logs, nil
}
