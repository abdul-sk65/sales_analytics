package repository

import (
	"context"
	"fmt"

	"sales_analytics/config"

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

	return repo, nil
}

// Disconnect closes the MongoDB connection
func (r *MongoRepository) Disconnect(ctx context.Context) error {
	return r.client.Disconnect(ctx)
}

// GetCollection returns a MongoDB collection
func (r *MongoRepository) GetCollection(name string) *mongo.Collection {
	return r.db.Collection(name)
}
