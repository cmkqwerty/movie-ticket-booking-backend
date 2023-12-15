package main

import (
	"context"
	"fmt"
	"github.com/cmkqwerty/movie-ticket-booking-backend/api"
	"github.com/cmkqwerty/movie-ticket-booking-backend/db"
	"github.com/cmkqwerty/movie-ticket-booking-backend/db/fixtures"
	"github.com/cmkqwerty/movie-ticket-booking-backend/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.URI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.NAME).Drop(ctx); err != nil {
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
}
