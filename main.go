package main

import (
	"context"
	"flag"
	"github.com/cmkqwerty/movie-ticket-booking-backend/api"
	"github.com/cmkqwerty/movie-ticket-booking-backend/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
	}}

func main() {
	listenAddr := flag.String("listenAddr", ":3000", "The listener address of the HTTP API server.")
	flag.Parse()

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	var (
		userHandler   = api.NewUserHandler(db.NewMongoUserStore(client))
		cinemaStore   = db.NewMongoCinemaStore(client)
		hallStore     = db.NewMongoHallStore(client, cinemaStore)
		cinemaHandler = api.NewCinemaHandler(cinemaStore, db.NewMongoMovieStore(client), hallStore)
	)

	app := fiber.New(config)
	apiv1 := app.Group("/api/v1")

	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)

	apiv1.Get("/cinema", cinemaHandler.HandleGetCinemas)

	app.Listen(*listenAddr)
}
