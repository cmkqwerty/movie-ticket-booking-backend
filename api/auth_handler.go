package api

import (
	"errors"
	"fmt"
	"github.com/cmkqwerty/movie-ticket-booking-backend/db"
	"github.com/cmkqwerty/movie-ticket-booking-backend/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"os"
	"time"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

func invalidCredentials(c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(map[string]string{"message": "invalid credentials"})
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}

	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println("Error getting user by email:", err)
			return invalidCredentials(c)
		}

		return err
	}

	if !types.IsValidPassword(user.EncryptedPassword, params.Password) {
		return invalidCredentials(c)
	}

	resp := AuthResponse{
		User:  user,
		Token: CreateTokenFromUser(user),
	}

	return c.JSON(resp)
}

func CreateTokenFromUser(user *types.User) string {
	now := time.Now()
	expires := now.Add(time.Hour * 2).Unix()
	claims := jwt.MapClaims{
		"id":      user.ID.Hex(),
		"email":   user.Email,
		"expires": expires,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("Error signing token:", err)
	}

	return tokenString
}
