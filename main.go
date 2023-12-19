package main

import (
	"context"
	"github.com/cmkqwerty/movie-ticket-booking-backend/api"
	"github.com/cmkqwerty/movie-ticket-booking-backend/db"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {
	mongoEndpoint := os.Getenv("MONGO_DB_URL")
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	}

	var (
		userStore    = db.NewMongoUserStore(client)
		cinemaStore  = db.NewMongoCinemaStore(client)
		movieStore   = db.NewMongoMovieStore(client)
		hallStore    = db.NewMongoHallStore(client, cinemaStore)
		bookingStore = db.NewMongoBookingStore(client)
		store        = &db.Store{
			User:    userStore,
			Cinema:  cinemaStore,
			Movie:   movieStore,
			Hall:    hallStore,
			Booking: bookingStore,
		}
		userHandler    = api.NewUserHandler(store)
		cinemaHandler  = api.NewCinemaHandler(store)
		movieHandler   = api.NewMovieHandler(store)
		hallHandler    = api.NewHallHandler(store)
		authHandler    = api.NewAuthHandler(userStore)
		bookingHandler = api.NewBookingHandler(store)

		app   = fiber.New(config)
		auth  = app.Group("/api")
		apiV1 = app.Group("/api/v1", api.JWTAuthentication(userStore))
		admin = apiV1.Group("/admin", api.AdminAuth)
	)

	// Auth routes
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// Versioned API routes
	// User routes
	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)
	apiV1.Get("/user", userHandler.HandleGetUsers)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiV1.Put("/user/:id", userHandler.HandlePutUser)

	// Cinema routes
	apiV1.Get("/cinema", cinemaHandler.HandleGetCinemas)
	apiV1.Get("/cinema/:id", cinemaHandler.HandleGetCinema)
	apiV1.Get("/cinema/:id/halls", cinemaHandler.HandleGetHalls)

	apiV1.Get("/hall", hallHandler.HandleGetHalls)
	apiV1.Post("/hall/:id/book", hallHandler.HandleBookHall)

	apiV1.Get("/movie/:id", movieHandler.HandleGetMovie)
	apiV1.Get("/movie", movieHandler.HandleGetMovies)

	// Booking routes
	admin.Get("/booking", bookingHandler.HandleGetBookings)
	apiV1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiV1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

	listenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")
	app.Listen(listenAddr)
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
