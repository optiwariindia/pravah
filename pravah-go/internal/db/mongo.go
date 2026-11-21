package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Database *mongo.Database

func Connect(uri string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	Client = client
	// We'll extract the database name from the URI or use a default one.
	// Typically the URI is mongodb://host:port/dbname
	// For simplicity, let's assume the db name is 'pravah' if not specified.
	Database = client.Database("pravah")
	EnsureIndexes()
	log.Println("Connected to MongoDB")
}

func EnsureIndexes() {
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
		log.Fatalf("Failed to create event indexes: %v", err)
	}

	configIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "appId", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("tenant_app_id"),
	}
	if _, err := Database.Collection("app_configs").Indexes().CreateOne(ctx, configIndex); err != nil {
		log.Fatalf("Failed to create app config indexes: %v", err)
	}
}
