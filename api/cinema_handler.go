package api

import (
	"github.com/cmkqwerty/movie-ticket-booking-backend/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CinemaHandler struct {
	store *db.Store
}

func NewCinemaHandler(store *db.Store) *CinemaHandler {
	return &CinemaHandler{
		store: store,
	}
}

func (h *CinemaHandler) HandleGetCinema(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID()
	}

	cinema, err := h.store.Cinema.GetCinemaByID(c.Context(), objID)
	if err != nil {
		return ErrResourceNotFound("cinema")
	}

	return c.JSON(cinema)
}

func (h *CinemaHandler) HandleGetCinemas(c *fiber.Ctx) error {
	cinemas, err := h.store.Cinema.GetCinemas(c.Context(), nil)
	if err != nil {
		return ErrResourceNotFound("cinema")
	}

	return c.JSON(cinemas)
}

func (h *CinemaHandler) HandleGetHalls(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID()
	}

	filter := bson.M{"cinema": objID}
	halls, err := h.store.Hall.GetHalls(c.Context(), filter)
	if err != nil {
		return ErrResourceNotFound("hall")
	}

	return c.JSON(halls)
}
