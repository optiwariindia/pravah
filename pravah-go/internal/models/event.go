package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventData struct {
	ID       string                 `bson:"id" json:"id"`
	Metadata map[string]interface{} `bson:"metadata" json:"metadata"`
}

type Event struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StreamID  string             `bson:"streamId" json:"streamId"`
	Event     string             `bson:"event" json:"event"`
	Data      EventData          `bson:"data" json:"data"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}
