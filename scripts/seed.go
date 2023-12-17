package main

import (
	"context"
	"fmt"
	"github.com/cmkqwerty/movie-ticket-booking-backend/api"
	"github.com/cmkqwerty/movie-ticket-booking-backend/db"
	"github.com/cmkqwerty/movie-ticket-booking-backend/db/fixtures"
	"github.com/cmkqwerty/movie-ticket-booking-backend/types"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	var (
		ctx           = context.Background()
		mongoEndpoint = os.Getenv("MONGO_DB_URL")
		mongoDBName   = os.Getenv("MONGO_DB_NAME")
	)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(mongoDBName).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	cinemaStore := db.NewMongoCinemaStore(client)
	store := db.Store{
		User:    db.NewMongoUserStore(client),
		Cinema:  cinemaStore,
		Movie:   db.NewMongoMovieStore(client),
		Hall:    db.NewMongoHallStore(client, cinemaStore),
		Booking: db.NewMongoBookingStore(client),
	}

	user := fixtures.AddUser(&store, "Jimmy", "Scott", false)
	fmt.Println("jimmy ->", api.CreateTokenFromUser(user))
	admin := fixtures.AddUser(&store, "Admin", "Admin", true)
	fmt.Println("admin ->", api.CreateTokenFromUser(admin))
	cinema := fixtures.AddCinema(&store, "CinemaxX", "Berlin", 5, nil)
	movie := fixtures.AddMovie(&store, "The Dark Knight", types.Action)
	hall := fixtures.AddHall(&store, 100, 10.0, cinema.ID, movie.ID)
	booking := fixtures.AddBooking(&store, user.ID, hall.ID, types.Night, time.Now().AddDate(0, 0, 5))
	fmt.Println(booking)

	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("Cinema%d", i)
		location := fmt.Sprintf("Location%d", i)
		fixtures.AddCinema(&store, name, location, rand.Intn(5), nil)
	}
}
