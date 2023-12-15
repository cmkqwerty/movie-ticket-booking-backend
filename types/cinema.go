package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Cinema struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string               `bson:"name" json:"name"`
	Location string               `bson:"location" json:"location"`
	Halls    []primitive.ObjectID `bson:"halls" json:"halls"`
	Rating   int                  `bson:"rating" json:"rating"`
}

type Hall struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Capacity int                `bson:"capacity" json:"capacity"`
	Price    float64            `bson:"price" json:"price"`
	Movie    primitive.ObjectID `bson:"movie" json:"movie"`
	Cinema   primitive.ObjectID `bson:"cinema" json:"cinema"`
}

type Genre int

const (
	Action Genre = iota
	Comedy
	Drama
	Horror
	Thriller
)

type Movie struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title string             `bson:"title" json:"title"`
	Genre Genre              `bson:"genre" json:"genre"`
}
