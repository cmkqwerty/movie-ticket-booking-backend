package db

import (
	"context"
	"github.com/cmkqwerty/movie-ticket-booking-backend/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

const hallColl = "halls"

type HallStore interface {
	InsertHall(context.Context, *types.Hall) (*types.Hall, error)
	GetHalls(context.Context, Map) ([]*types.Hall, error)
	GetHallCapacity(context.Context, primitive.ObjectID) (int, error)
}

type MongoHallStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	CinemaStore
}

func NewMongoHallStore(c *mongo.Client, cinemaStore CinemaStore) *MongoHallStore {
	dbname := os.Getenv(MONGO_DB_ENV_NAME)
	return &MongoHallStore{
		client:      c,
		coll:        c.Database(dbname).Collection(hallColl),
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
	filter := Map{"_id": hall.Cinema}
	update := Map{"$push": Map{"halls": hall.ID}}
	if err := s.CinemaStore.UpdateCinema(ctx, filter, update); err != nil {
		return nil, err
	}

	return hall, nil
}

func (s *MongoHallStore) GetHalls(ctx context.Context, filter Map) ([]*types.Hall, error) {
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var halls []*types.Hall
	if err := cur.All(ctx, &halls); err != nil {
		return nil, err
	}

	return halls, nil
}

func (s *MongoHallStore) GetHallCapacity(ctx context.Context, hallID primitive.ObjectID) (int, error) {
	var hall types.Hall
	if err := s.coll.FindOne(ctx, bson.M{"_id": hallID}).Decode(&hall); err != nil {
		return 0, err
	}

	return hall.Capacity, nil
}
