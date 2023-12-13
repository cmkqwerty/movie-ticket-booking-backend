package api

import (
	"context"
	"fmt"
	"github.com/cmkqwerty/movie-ticket-booking-backend/db"
	"github.com/cmkqwerty/movie-ticket-booking-backend/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type BookHallParams struct {
	Session types.Session `json:"session"`
	Date    time.Time     `json:"date"`
}

func (p BookHallParams) validate() error {
	now := time.Now()
	if now.After(p.Date) {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid date.")
	}

	if p.Session < types.Morning || p.Session > types.Night {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid session.")
	}

	return nil
}

type HallHandler struct {
	store *db.Store
}

func NewHallHandler(store *db.Store) *HallHandler {
	return &HallHandler{
		store: store,
	}
}

func (h *HallHandler) HandleGetHalls(c *fiber.Ctx) error {
	halls, err := h.store.Hall.GetHalls(c.Context(), bson.M{})
	if err != nil {
		return err
	}

	return c.JSON(halls)
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

	if err := params.validate(); err != nil {
		return err
	}

	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized."})
	}

	// check if there is enough capacity for the session
	ok, err = h.isHallAvailableForBooking(c.Context(), hallID, params.Session)
	if err != nil {
		return err
	}
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(genericResp{
			Type: "error",
			Msg:  fmt.Sprintf("Hall %s is full.", hallID.Hex()),
		})
	}

	booking := types.Booking{
		UserID:  user.ID,
		HallID:  hallID,
		Session: params.Session,
		Date:    params.Date,
	}

	inserted, err := h.store.Booking.InsertBooking(c.Context(), &booking)
	if err != nil {
		return err
	}

	return c.JSON(inserted)
}

func (h *HallHandler) isHallAvailableForBooking(ctx context.Context, hallID primitive.ObjectID, session types.Session) (bool, error) {
	where := bson.M{
		"hallID":  hallID,
		"session": session,
	}
	bookings, err := h.store.Booking.GetBookings(ctx, where)
	if err != nil {
		return false, err
	}
	capacity, err := h.store.Hall.GetHallCapacity(ctx, hallID)
	if err != nil {
		return false, err
	}

	if len(bookings) >= capacity {
		return false, nil
	}

	return true, nil
}