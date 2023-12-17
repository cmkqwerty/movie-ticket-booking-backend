package db

import (
	"context"
	"github.com/cmkqwerty/movie-ticket-booking-backend/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

const bookingColl = "bookings"

type BookingStore interface {
	InsertBooking(context.Context, *types.Booking) (*types.Booking, error)
	GetBookingByID(context.Context, string) (*types.Booking, error)
	GetBookings(context.Context, Map, *Pagination) ([]*types.Booking, error)
	UpdateBooking(context.Context, string, Map) error
	CountBookings(context.Context, Map) (int, error)
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	BookingStore
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore {
	dbname := os.Getenv(MONGO_DB_ENV_NAME)
	return &MongoBookingStore{
		client: client,
		coll:   client.Database(dbname).Collection(bookingColl),
	}
}

func (s *MongoBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	res, err := s.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}

	booking.ID = res.InsertedID.(primitive.ObjectID)
	return booking, nil
}

func (s *MongoBookingStore) GetBookingByID(ctx context.Context, id string) (*types.Booking, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var booking types.Booking
	if err := s.coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&booking); err != nil {
		return nil, err
	}

	return &booking, nil
}

func (s *MongoBookingStore) GetBookings(ctx context.Context, filter Map, pag *Pagination) ([]*types.Booking, error) {
	opts := options.FindOptions{}
	opts.SetSkip((pag.Page - 1) * pag.Limit)
	opts.SetLimit(pag.Limit)

	cur, err := s.coll.Find(ctx, filter, &opts)
	if err != nil {
		return nil, err
	}

	var bookings []*types.Booking
	if err := cur.All(ctx, &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (s *MongoBookingStore) UpdateBooking(ctx context.Context, id string, update Map) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	if _, err := s.coll.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": update}); err != nil {
		return err
	}

	return nil
}

func (s *MongoBookingStore) CountBookings(ctx context.Context, filter Map) (int, error) {
	bookingCount, err := s.coll.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(bookingCount), nil
}
