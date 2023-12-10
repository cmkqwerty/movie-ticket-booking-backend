package db

import (
	"context"
	"github.com/cmkqwerty/movie-ticket-booking-backend/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const cinemaColl = "cinemas"

type CinemaStore interface {
	Insert(context.Context, *types.Cinema) (*types.Cinema, error)
	Update(context.Context, bson.M, bson.M) error
	GetCinemas(context.Context, bson.M) ([]*types.Cinema, error)
}

type MongoCinemaStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoCinemaStore(c *mongo.Client) *MongoCinemaStore {
	return &MongoCinemaStore{
		client: c,
		coll:   c.Database(DBNAME).Collection(cinemaColl),
	}
}

func (s *MongoCinemaStore) Insert(ctx context.Context, cinema *types.Cinema) (*types.Cinema, error) {
	res, err := s.coll.InsertOne(ctx, cinema)
	if err != nil {
		return nil, err
	}

	cinema.ID = res.InsertedID.(primitive.ObjectID)

	return cinema, nil
}

func (s *MongoCinemaStore) Update(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (s *MongoCinemaStore) GetCinemas(ctx context.Context, filter bson.M) ([]*types.Cinema, error) {
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var cinemas []*types.Cinema
	if err := cur.All(ctx, &cinemas); err != nil {
		return nil, err
	}

	return cinemas, nil
}
