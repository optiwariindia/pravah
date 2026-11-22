package db

import (
	"context"
	"log"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/sampledge-sse/pravah-go/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Database *mongo.Database
var healthy atomic.Bool
var indexesReady atomic.Bool

func Connect(uri string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI(uri).
		SetServerSelectionTimeout(utils.DatabaseOperationTimeout()))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	Client = client
	// We'll extract the database name from the URI or use a default one.
	// Typically the URI is mongodb://host:port/dbname
	// For simplicity, let's assume the db name is 'pravah' if not specified.
	Database = client.Database("pravah")
	if Ping(ctx) {
		log.Println("Connected to MongoDB")
		EnsureIndexes()
	} else {
		log.Println("MongoDB is unavailable; starting in degraded live-delivery mode")
	}
}

func Ping(ctx context.Context) bool {
	if Client == nil {
		return false
	}
	err := Client.Ping(ctx, nil)
	healthy.Store(err == nil)
	return err == nil
}

func IsHealthy() bool {
	return healthy.Load()
}

func MarkUnhealthy() {
	healthy.Store(false)
}

func EnsureIndexes() {
	if indexesReady.Load() || Database == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	eventIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "expiresAt", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0).SetName("event_ttl"),
		},
		{
			Keys:    bson.D{{Key: "appId", Value: 1}, {Key: "streamId", Value: 1}, {Key: "_id", Value: 1}},
			Options: options.Index().SetName("tenant_stream_event_id"),
		},
	}
	if _, err := Database.Collection("events").Indexes().CreateMany(ctx, eventIndexes); err != nil {
		slog.Error("failed to create event indexes", "error", err)
		return
	}

	configIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "appId", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("tenant_app_id"),
	}
	if _, err := Database.Collection("app_configs").Indexes().CreateOne(ctx, configIndex); err != nil {
		slog.Error("failed to create app config indexes", "error", err)
		return
	}
	indexesReady.Store(true)
}
