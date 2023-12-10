package db

import (
	"context"
	"github.com/cmkqwerty/movie-ticket-booking-backend/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const movieColl = "movies"

type MovieStore interface {
	InsertMovie(context.Context, *types.Movie) (*types.Movie, error)
}

type MongoMovieStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoMovieStore(c *mongo.Client) *MongoMovieStore {
	return &MongoMovieStore{
		client: c,
		coll:   c.Database(DBNAME).Collection(movieColl),
	}
}

func (s *MongoMovieStore) InsertMovie(ctx context.Context, movie *types.Movie) (*types.Movie, error) {
	res, err := s.coll.InsertOne(ctx, movie)
	if err != nil {
		return nil, err
	}

	movie.ID = res.InsertedID.(primitive.ObjectID)

	return movie, nil
}
