package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
