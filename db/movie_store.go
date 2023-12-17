package db

import (
	"context"
	"github.com/cmkqwerty/movie-ticket-booking-backend/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const movieColl = "movies"

type MovieStore interface {
	InsertMovie(context.Context, *types.Movie) (*types.Movie, error)
	GetMovieByID(context.Context, string) (*types.Movie, error)
	GetMovies(context.Context, Map) ([]*types.Movie, error)
}

type MongoMovieStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoMovieStore(c *mongo.Client) *MongoMovieStore {
	return &MongoMovieStore{
		client: c,
		coll:   c.Database(NAME).Collection(movieColl),
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

func (s *MongoMovieStore) GetMovieByID(ctx context.Context, id string) (*types.Movie, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var movie types.Movie
	if err := s.coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&movie); err != nil {
		return nil, err
	}

	return &movie, nil
}

func (s *MongoMovieStore) GetMovies(ctx context.Context, filter Map) ([]*types.Movie, error) {
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var movies []*types.Movie
	if err := cur.All(ctx, &movies); err != nil {
		return nil, err
	}

	return movies, nil
}
