package api

import (
	"github.com/cmkqwerty/movie-ticket-booking-backend/db"
	"github.com/gofiber/fiber/v2"
)

type CinemaHandler struct {
	cinemaStore db.CinemaStore
	movieStore  db.MovieStore
	roomStore   db.HallStore
}

func NewCinemaHandler(cinemaStore db.CinemaStore, movieStore db.MovieStore, roomStore db.HallStore) *CinemaHandler {
	return &CinemaHandler{
		cinemaStore: cinemaStore,
		movieStore:  movieStore,
		roomStore:   roomStore,
	}
}

type CinemaQueryParams struct {
	Halls  bool
	Rating int
}

func (h *CinemaHandler) HandleGetCinemas(c *fiber.Ctx) error {
	var qParams CinemaQueryParams
	if err := c.QueryParser(&qParams); err != nil {
		return err
	}

	cinemas, err := h.cinemaStore.GetCinemas(c.Context(), nil)
	if err != nil {
		return err
	}

	return c.JSON(cinemas)
}
