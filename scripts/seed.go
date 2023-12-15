package main

import (
	"context"
	"github.com/cmkqwerty/movie-ticket-booking-backend/db"
	"github.com/cmkqwerty/movie-ticket-booking-backend/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var (
	client      *mongo.Client
	cinemaStore db.CinemaStore
	movieStore  db.MovieStore
	hallStore   db.HallStore
	userStore   db.UserStore
	ctx         = context.Background()
)

func seedUser(fName, lName, email string) {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fName,
		LastName:  lName,
		Email:     email,
		Password:  "securePassword12345",
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = userStore.InsertUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
}

func seedCinema(name, location string, rating int) {
	cinema := types.Cinema{
		Name:     name,
		Location: location,
		Halls:    []primitive.ObjectID{},
		Rating:   rating,
	}
	movie := types.Movie{
		Title: "Aftersun",
		Genre: types.Drama,
	}
	halls := []types.Hall{
		{
			Capacity: 100,
			Price:    20.0,
		},
		{
			Capacity: 200,
			Price:    10.0,
		},
	}

	insertedCinema, err := cinemaStore.InsertCinema(ctx, &cinema)
	if err != nil {
		log.Fatal(err)
	}
	insertedMovie, err := movieStore.InsertMovie(ctx, &movie)
	if err != nil {
		log.Fatal(err)
	}

	for _, hall := range halls {
		hall.Cinema = insertedCinema.ID
		hall.Movie = insertedMovie.ID
		_, err := hallStore.InsertHall(ctx, &hall)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	seedCinema("CinemaxX", "Berlin", 5)
	seedCinema("Mock", "Berlin", 2)
	seedUser("Jimmy", "Scott", "jimmy@scott.com")
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
}
