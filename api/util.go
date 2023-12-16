package api

import (
	"github.com/cmkqwerty/movie-ticket-booking-backend/types"
	"github.com/gofiber/fiber/v2"
)

func getAuthUser(c *fiber.Ctx) (*types.User, error) {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return nil, ErrUnauthorized()
	}

	return user, nil
}
