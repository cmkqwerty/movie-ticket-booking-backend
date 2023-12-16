package fixtures

import (
	"context"
	"fmt"
	"github.com/cmkqwerty/movie-ticket-booking-backend/db"
	"github.com/cmkqwerty/movie-ticket-booking-backend/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

func AddUser(store *db.Store, fName, lName string, admin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fName,
		LastName:  lName,
		Email:     fmt.Sprintf("%s@%s.com", fName, lName),
		Password:  fmt.Sprintf("%s_%s", fName, lName),
	})
	if err != nil {
		log.Fatal(err)
	}

	user.IsAdmin = admin

	insertedUser, err := store.User.InsertUser(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	return insertedUser
}

func AddCinema(store *db.Store, name, location string, rating int, hallIDs []primitive.ObjectID) *types.Cinema {
	if hallIDs == nil {
		hallIDs = []primitive.ObjectID{}
	}

	cinema := &types.Cinema{
		Name:     name,
		Location: location,
		Halls:    hallIDs,
		Rating:   rating,
	}

	insertedCinema, err := store.Cinema.InsertCinema(context.Background(), cinema)
	if err != nil {
		log.Fatal(err)
	}

	return insertedCinema
}

func AddMovie(store *db.Store, title string, genre types.Genre) *types.Movie {
	movie := &types.Movie{
		Title: title,
		Genre: genre,
	}

	insertedMovie, err := store.Movie.InsertMovie(context.Background(), movie)
	if err != nil {
		log.Fatal(err)
	}

	return insertedMovie
}

func AddHall(store *db.Store, capacity int, price float64, cinemaID, movieID primitive.ObjectID) *types.Hall {
	hall := &types.Hall{
		Capacity: capacity,
		Price:    price,
		Cinema:   cinemaID,
		Movie:    movieID,
	}

	insertedHall, err := store.Hall.InsertHall(context.Background(), hall)
	if err != nil {
		log.Fatal(err)
	}

	return insertedHall
}

func AddBooking(store *db.Store, userID, hallID primitive.ObjectID, session types.Session, date time.Time) *types.Booking {
	booking := &types.Booking{
		UserID:  userID,
		HallID:  hallID,
		Session: session,
		Date:    date,
	}

	insertedBooking, err := store.Booking.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}

	return insertedBooking
}
