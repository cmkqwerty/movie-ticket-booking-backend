package api

import (
	"context"
	"github.com/cmkqwerty/movie-ticket-booking-backend/db"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"reflect"
	"testing"
)

type testDB struct {
	client *mongo.Client
	*db.Store
}

func (tdb *testDB) tearDown(t *testing.T) {
	dbname := os.Getenv(db.MONGO_DB_ENV_NAME)
	if err := tdb.client.Database(dbname).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testDB {
	if err := godotenv.Load("../.env"); err != nil {
		t.Error(err)
	}
	dburi := os.Getenv("MONGO_DB_URL_TEST")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
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
