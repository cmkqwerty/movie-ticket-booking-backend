package api

import (
	"context"
	"github.com/cmkqwerty/movie-ticket-booking-backend/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"testing"
)

type testDB struct {
	client *mongo.Client
	*db.Store
}

func (tdb *testDB) tearDown(t *testing.T) {
	if err := tdb.client.Database(db.NAME).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testDB {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.URI))
	if err != nil {
		t.Fatal(err)
	}

	cinemaStore := db.NewMongoCinemaStore(client)
	return &testDB{
		client: client,
		Store: &db.Store{
			User:    db.NewMongoUserStore(client),
			Cinema:  cinemaStore,
			Movie:   db.NewMongoMovieStore(client),
			Hall:    db.NewMongoHallStore(client, cinemaStore),
			Booking: db.NewMongoBookingStore(client),
		},
	}
}

func contains[T comparable](s []T, elem T) bool {
	for _, e := range s {
		if reflect.DeepEqual(e, elem) {
			return true
		}
	}

	return false
}
