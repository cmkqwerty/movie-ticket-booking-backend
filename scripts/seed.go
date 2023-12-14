package main

import (
	"context"
	"fmt"
	"github.com/cmkqwerty/movie-ticket-booking-backend/api"
	"github.com/cmkqwerty/movie-ticket-booking-backend/db"
	"github.com/cmkqwerty/movie-ticket-booking-backend/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var (
	client       *mongo.Client
	cinemaStore  db.CinemaStore
	movieStore   db.MovieStore
	hallStore    db.HallStore
	userStore    db.UserStore
	bookingStore db.BookingStore
	ctx          = context.Background()
)

func seedUser(isAdmin bool, fName, lName, email, password string) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fName,
		LastName:  lName,
		Email:     email,
		Password:  password,
	})
	if err != nil {
		log.Fatal(err)
	}

	user.IsAdmin = isAdmin

	insertedUser, err := userStore.InsertUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))
	return insertedUser
}

func seedBooking(userID, hallID primitive.ObjectID, session types.Session, date time.Time) {
	booking := &types.Booking{
		UserID:  userID,
		HallID:  hallID,
		Session: session,
		Date:    date,
	}

	if _, err := bookingStore.InsertBooking(ctx, booking); err != nil {
		log.Fatal(err)
	}
}

func seedCinema(name, location string, rating int) *types.Cinema {
	cinema := types.Cinema{
		Name:     name,
		Location: location,
		Halls:    []primitive.ObjectID{},
		Rating:   rating,
	}

	insertedCinema, err := cinemaStore.InsertCinema(ctx, &cinema)
	if err != nil {
		log.Fatal(err)
	}

	return insertedCinema
}

func seedMovie(title string, genre types.Genre) *types.Movie {
	movie := types.Movie{
		Title: title,
		Genre: genre,
	}

	insertedMovie, err := movieStore.InsertMovie(ctx, &movie)
	if err != nil {
		log.Fatal(err)
	}

	return insertedMovie
}

func seedHall(capacity int, price float64, movieID, cinemaID primitive.ObjectID) *types.Hall {
	hall := types.Hall{
		Capacity: capacity,
		Price:    price,
		Movie:    movieID,
		Cinema:   cinemaID,
	}

	insertedHall, err := hallStore.InsertHall(ctx, &hall)
	if err != nil {
		log.Fatal(err)
	}

	return insertedHall
}

func main() {
	jimmy := seedUser(false, "Jimmy", "Scott", "jimmy@scott.com", "12345678")
	seedUser(true, "Admin", "Admin", "admin@admin.com", "secure12345678")
	seedCinema("Mock", "Berlin", 2)
	cinema := seedCinema("CinemaxX", "Berlin", 5)
	movie := seedMovie("The Dark Knight", types.Action)
	hall := seedHall(100, 10.0, movie.ID, cinema.ID)
	seedBooking(jimmy.ID, hall.ID, types.Night, time.Now().AddDate(0, 0, 5))
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.URI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.NAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	cinemaStore = db.NewMongoCinemaStore(client)
	movieStore = db.NewMongoMovieStore(client)
	hallStore = db.NewMongoHallStore(client, cinemaStore)
	userStore = db.NewMongoUserStore(client)
	bookingStore = db.NewMongoBookingStore(client)
}
