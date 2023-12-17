package api

import (
	"errors"
	"github.com/cmkqwerty/movie-ticket-booking-backend/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type MovieHandler struct {
	store *db.Store
}

func NewMovieHandler(store *db.Store) *MovieHandler {
	return &MovieHandler{
		store: store,
	}
}

func (h *MovieHandler) HandleGetMovie(c *fiber.Ctx) error {
	id := c.Params("id")

	movie, err := h.store.Movie.GetMovieByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrResourceNotFound("movie")
		}
		return err
	}

	return c.JSON(movie)
}

func (h *MovieHandler) HandleGetMovies(c *fiber.Ctx) error {
	movies, err := h.store.Movie.GetMovies(c.Context(), db.Map{})
	if err != nil {
		return ErrResourceNotFound("movie")
	}

	return c.JSON(movies)
}
