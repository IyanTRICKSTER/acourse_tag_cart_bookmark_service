package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Subscription struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	UserID    string             `json:"user_id" bson:"user_id"`
	Courses   []Course           `json:"courses" bson:"courses"`
	UpdatedAt *time.Time         `json:"updated_at,omitempty" bson:"updated_at"`
	CreatedAt *time.Time         `json:"created_at,omitempty" bson:"created_at"`
	DeletedAt *time.Time         `json:"deleted_at,omitempty" bson:"deleted_at"`
}
