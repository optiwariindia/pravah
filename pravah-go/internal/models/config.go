package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppConfig struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	AppID          string             `bson:"appId" json:"appId"`
	RetentionHours int                `bson:"retentionHours" json:"retentionHours"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
}
