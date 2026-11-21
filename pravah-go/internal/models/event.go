package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventData struct {
	ID string `bson:"id" json:"id"`
}

type Event struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AppID     string             `bson:"appId" json:"appId"`
	StreamID  string             `bson:"streamId" json:"streamId"`
	Event     string             `bson:"event" json:"event"`
	Data      EventData          `bson:"data" json:"data"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	ExpiresAt time.Time          `bson:"expiresAt" json:"expiresAt"`
}
