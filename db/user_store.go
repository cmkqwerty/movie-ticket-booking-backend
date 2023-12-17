package db

import (
	"context"
	"fmt"
	"github.com/cmkqwerty/movie-ticket-booking-backend/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

const userColl = "users"

type Dropper interface {
	Drop(context.Context) error
}

type UserStore interface {
	InsertUser(context.Context, *types.User) (*types.User, error)
	GetUserByID(context.Context, string) (*types.User, error)
	GetUserByEmail(context.Context, string) (*types.User, error)
	GetUsers(context.Context, *Pagination) ([]*types.User, error)
	UpdateUser(ctx context.Context, filter Map, params types.UpdateUserParams) error
	DeleteUser(context.Context, string) error

	Dropper
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(c *mongo.Client) *MongoUserStore {
	dbname := os.Getenv(MONGO_DB_ENV_NAME)
	return &MongoUserStore{
		client: c,
		coll:   c.Database(dbname).Collection(userColl),
	}
}

func (s *MongoUserStore) Drop(ctx context.Context) error {
	fmt.Println("--- dropping users collection ---")
	return s.coll.Drop(ctx)
}

func (s *MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = res.InsertedID.(primitive.ObjectID)

	return user, nil
}

func (s *MongoUserStore) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user types.User
	if err := s.coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	var user types.User
	if err := s.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context, pag *Pagination) ([]*types.User, error) {
	opts := options.FindOptions{}
	opts.SetSkip((pag.Page - 1) * pag.Limit)
	opts.SetLimit(pag.Limit)

	cur, err := s.coll.Find(ctx, bson.M{}, &opts)
	if err != nil {
		return nil, err
	}

	var users []*types.User
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// TODO: Handle if user doesn't exist
	if _, err := s.coll.DeleteOne(ctx, bson.M{"_id": objID}); err != nil {
		return err
	}

	return nil
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, filter Map, params types.UpdateUserParams) error {
	objID, err := primitive.ObjectIDFromHex(filter["_id"].(string))
	if err != nil {
		return err
	}

	filter["_id"] = objID
	update := bson.D{
		{
			"$set", params.ToBSON(),
		},
	}
	_, err = s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
