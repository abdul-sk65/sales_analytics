package repository

import (
	"context"
	"fmt"

	"sales_analytics/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoRepository handles all database operations
type MongoRepository struct {
	client *mongo.Client
	db     *mongo.Database
	config *config.Config
}

// NewMongoRepository creates a new MongoDB repository
func NewMongoRepository(ctx context.Context, cfg *config.Config) (*MongoRepository, error) {
	clientOptions := options.Client().ApplyURI(cfg.MongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database(cfg.DatabaseName)

	repo := &MongoRepository{
		client: client,
		db:     db,
		config: cfg,
	}
	// Create indexes
	if err := repo.createIndexes(ctx); err != nil {
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}

	return repo, nil
}

// createIndexes creates necessary indexes for optimal query performance
func (r *MongoRepository) createIndexes(ctx context.Context) error {
	// Customer indexes
	customerIndexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "customer_id", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "email", Value: 1}}},
	}
	if _, err := r.db.Collection("customers").Indexes().CreateMany(ctx, customerIndexes); err != nil {
		return err
	}

	// Product indexes
	productIndexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "product_id", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "category", Value: 1}}},
	}
	if _, err := r.db.Collection("products").Indexes().CreateMany(ctx, productIndexes); err != nil {
		return err
	}

	// Order indexes
	orderIndexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "order_id", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "customer_id", Value: 1}}},
		{Keys: bson.D{{Key: "product_id", Value: 1}}},
		{Keys: bson.D{{Key: "date_of_sale", Value: 1}}},
		{Keys: bson.D{{Key: "region", Value: 1}}},
	}
	if _, err := r.db.Collection("orders").Indexes().CreateMany(ctx, orderIndexes); err != nil {
		return err
	}

	return nil
}

// Disconnect closes the MongoDB connection
func (r *MongoRepository) Disconnect(ctx context.Context) error {
	return r.client.Disconnect(ctx)
}

// GetCollection returns a MongoDB collection
func (r *MongoRepository) GetCollection(name string) *mongo.Collection {
	return r.db.Collection(name)
}
