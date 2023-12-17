package api

import (
	"context"
	"fmt"
	"github.com/cmkqwerty/movie-ticket-booking-backend/db"
	"github.com/cmkqwerty/movie-ticket-booking-backend/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type BookHallParams struct {
	Session types.Session `json:"session"`
	Date    time.Time     `json:"date"`
}

func (p BookHallParams) validate() error {
	now := time.Now()
	if now.After(p.Date) {
		return NewError(http.StatusBadRequest, "Invalid date.")
	}

	if p.Session < types.Morning || p.Session > types.Night {
		return NewError(http.StatusBadRequest, "Invalid session.")
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
	halls, err := h.store.Hall.GetHalls(c.Context(), db.Map{})
	if err != nil {
		return ErrResourceNotFound("hall")
	}

	return c.JSON(halls)
}

func (h *HallHandler) HandleBookHall(c *fiber.Ctx) error {
	var params BookHallParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}

	hallID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return ErrInvalidID()
	}

	if err := params.validate(); err != nil {
		return err
	}

	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return ErrUnauthorized()
	}

	// check if there is enough capacity for the session
	ok, err = h.isHallAvailableForBooking(c.Context(), hallID, params.Session)
	if err != nil {
		return ErrResourceNotFound("hall")
	}
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"message": fmt.Sprintf("Hall %s is full.", hallID.Hex())})
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
	where := db.Map{
		"hallID":   hallID,
		"session":  session,
		"canceled": false,
	}
	bookingsCount, err := h.store.Booking.CountBookings(ctx, where)
	if err != nil {
		return false, ErrResourceNotFound("booking")
	}
	capacity, err := h.store.Hall.GetHallCapacity(ctx, hallID)
	if err != nil {
		return false, ErrResourceNotFound("hall")
	}

	if bookingsCount >= capacity {
		return false, nil
	}

	return true, nil
}
