package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Session int

const (
	Morning Session = iota
	Afternoon
	Evening
	Night
)

type Booking struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID   primitive.ObjectID `bson:"userID,omitempty" json:"userID,omitempty"`
	HallID   primitive.ObjectID `bson:"hallID,omitempty" json:"hallID,omitempty"`
	Session  Session            `bson:"session,omitempty" json:"session,omitempty"`
	Date     time.Time          `bson:"date,omitempty" json:"date,omitempty"`
	Canceled bool               `bson:"canceled" json:"canceled"`
}
