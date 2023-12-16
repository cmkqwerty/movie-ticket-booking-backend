package api

import (
	"fmt"
	"github.com/cmkqwerty/movie-ticket-booking-backend/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"time"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok {
			fmt.Println("No token provided.")
			return ErrUnauthorized()
		}

		claims, err := validateToken(token[0])
		if err != nil {
			return ErrUnauthorized()
		}

		expires := claims["expires"].(float64)
		if time.Now().Unix() > int64(expires) {
			fmt.Println("Token has expired.")
			return NewError(http.StatusUnauthorized, "expired token")
		}

		userID := claims["id"].(string)
		user, err := userStore.GetUserByID(c.Context(), userID)
		if err != nil {
			return ErrUnauthorized()
		}

		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}

func validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Unexpected signing method:", token.Header["alg"])
			return nil, ErrUnauthorized()
		}

		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("Error parsing token:", err)
		return nil, ErrUnauthorized()
	}

	if !token.Valid {
		fmt.Println("Invalid token:")
		return nil, ErrUnauthorized()
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("Error getting claims:", err)
		return nil, ErrUnauthorized()
	}

	return claims, nil
}
