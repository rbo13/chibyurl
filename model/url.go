package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// URL represents a url to be
// shortened.
type URL struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Alias     string             `bson:"alias" json:"alias"`
	URL       string             `bson:"url" json:"url"`
	Click     int32              `bson:"click" json:"click"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
