package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type App struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	AppID     string             `bson:"appId" json:"appId"`
	APIKey    string             `bson:"apiKey" json:"apiKey"`
	Status    string             `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}
