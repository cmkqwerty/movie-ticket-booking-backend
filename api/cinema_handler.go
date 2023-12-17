package api

import (
	"github.com/cmkqwerty/movie-ticket-booking-backend/db"
	"github.com/gofiber/fiber/v2"
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

	cinema, err := h.store.Cinema.GetCinemaByID(c.Context(), id)
	if err != nil {
		return ErrResourceNotFound("cinema")
	}

	return c.JSON(cinema)
}

type CinemaQueryParams struct {
	db.Pagination
	Rating int
}

func (h *CinemaHandler) HandleGetCinemas(c *fiber.Ctx) error {
	var params CinemaQueryParams
	if err := c.QueryParser(&params); err != nil {
		return ErrBadRequest()
	}

	filter := db.Map{
		"rating": params.Rating,
	}
	cinemas, err := h.store.Cinema.GetCinemas(c.Context(), filter, &params.Pagination)
	if err != nil {
		return ErrResourceNotFound("cinema")
	}

	resp := ResourceResponse{
		Results: len(cinemas),
		Data:    cinemas,
		Page:    int(params.Page),
	}
	return c.JSON(resp)
}

func (h *CinemaHandler) HandleGetHalls(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID()
	}

	filter := db.Map{"cinema": objID}
	halls, err := h.store.Hall.GetHalls(c.Context(), filter)
	if err != nil {
		return ErrResourceNotFound("hall")
	}

	return c.JSON(halls)
}
