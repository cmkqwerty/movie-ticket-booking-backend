package api

import (
	"github.com/cmkqwerty/movie-ticket-booking-backend/db"
	"github.com/gofiber/fiber/v2"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")

	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return ErrResourceNotFound("booking")
	}

	user, err := getAuthUser(c)
	if err != nil {
		return ErrUnauthorized()
	}

	if booking.UserID != user.ID {
		return ErrUnauthorized()
	}

	return c.JSON(booking)
}

type BookingQueryParams struct {
	db.Pagination
	Canceled bool
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	var params BookingQueryParams
	if err := c.QueryParser(&params); err != nil {
		return ErrBadRequest()
	}

	filter := db.Map{
		"canceled": params.Canceled,
	}
	bookings, err := h.store.Booking.GetBookings(c.Context(), filter, &params.Pagination)
	if err != nil {
		return ErrResourceNotFound("booking")
	}

	resp := ResourceResponse{
		Results: len(bookings),
		Data:    bookings,
		Page:    int(params.Page),
	}
	return c.JSON(resp)

}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")

	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return ErrResourceNotFound("booking")
	}

	user, err := getAuthUser(c)
	if err != nil {
		return ErrUnauthorized()
	}

	if booking.UserID != user.ID {
		return ErrUnauthorized()
	}

	if err := h.store.Booking.UpdateBooking(c.Context(), id, db.Map{"canceled": true}); err != nil {
		return err
	}

	return c.JSON(map[string]string{"message": "success"})
}
