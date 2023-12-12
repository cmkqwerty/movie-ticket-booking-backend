package api

import (
	"fmt"
	"github.com/cmkqwerty/movie-ticket-booking-backend/db"
	"github.com/cmkqwerty/movie-ticket-booking-backend/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type BookHallParams struct {
	Session types.Session `json:"session"`
	Date    time.Time     `json:"date"`
}

type HallHandler struct {
	store *db.Store
}

func NewHallHandler(store *db.Store) *HallHandler {
	return &HallHandler{
		store: store,
	}
}

func (h *HallHandler) HandleBookHall(c *fiber.Ctx) error {
	var params BookHallParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	hallID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}

	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized."})
	}

	booking := types.Booking{
		UserID:  user.ID,
		HallID:  hallID,
		Session: params.Session,
		Date:    params.Date,
	}

	fmt.Printf("%+v\n", booking)
	return nil
}
