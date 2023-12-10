package db

import (
	"context"
	"github.com/cmkqwerty/movie-ticket-booking-backend/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const hallColl = "halls"

type HallStore interface {
	InsertHall(context.Context, *types.Hall) (*types.Hall, error)
}

type MongoHallStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	CinemaStore
}

func NewMongoHallStore(c *mongo.Client, cinemaStore CinemaStore) *MongoHallStore {
	return &MongoHallStore{
		client:      c,
		coll:        c.Database(DBNAME).Collection(hallColl),
		CinemaStore: cinemaStore,
	}
}

func (s *MongoHallStore) InsertHall(ctx context.Context, hall *types.Hall) (*types.Hall, error) {
	res, err := s.coll.InsertOne(ctx, hall)
	if err != nil {
		return nil, err
	}

	hall.ID = res.InsertedID.(primitive.ObjectID)

	// update cinema with new hall
	filter := bson.M{"_id": hall.Cinema}
	update := bson.M{"$push": bson.M{"halls": hall.ID}}
	if err := s.CinemaStore.Update(ctx, filter, update); err != nil {
		return nil, err
	}

	return hall, nil
}
