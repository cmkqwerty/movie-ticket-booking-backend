package api

import (
	"encoding/json"
	"github.com/cmkqwerty/movie-ticket-booking-backend/api/middleware"
	"github.com/cmkqwerty/movie-ticket-booking-backend/db/fixtures"
	"github.com/cmkqwerty/movie-ticket-booking-backend/types"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestUserGetBooking(t *testing.T) {
	db := setup(t)
	defer db.tearDown(t)

	var (
		nonAuthUser    = fixtures.AddUser(db.Store, "heron", "preston", false)
		user           = fixtures.AddUser(db.Store, "heron", "preston", false)
		cinema         = fixtures.AddCinema(db.Store, "babylon", "berlin", 4, nil)
		movie          = fixtures.AddMovie(db.Store, "the lighthouse", types.Horror)
		hall           = fixtures.AddHall(db.Store, 100, 10.0, cinema.ID, movie.ID)
		booking        = fixtures.AddBooking(db.Store, user.ID, hall.ID, types.Morning, time.Now().AddDate(0, 0, 1))
		app            = fiber.New()
		route          = app.Group("/", middleware.JWTAuthentication(db.User))
		bookingHandler = NewBookingHandler(db.Store)
	)

	route.Get("/:id", bookingHandler.HandleGetBooking)

	req := httptest.NewRequest("GET", "/"+booking.ID.Hex(), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	var returnedBooking *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&returnedBooking); err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if reflect.DeepEqual(returnedBooking, booking) {
		t.Fatalf("expected booking %v, got %v", booking, returnedBooking)
	}

	// Test that non-authenticated users cannot access the route
	req = httptest.NewRequest("GET", "/"+booking.ID.Hex(), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(nonAuthUser))

	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestAdminGetBookings(t *testing.T) {
	db := setup(t)
	defer db.tearDown(t)

	var (
		adminUser      = fixtures.AddUser(db.Store, "admin", "admin", true)
		user           = fixtures.AddUser(db.Store, "heron", "preston", false)
		cinema         = fixtures.AddCinema(db.Store, "babylon", "berlin", 4, nil)
		movie          = fixtures.AddMovie(db.Store, "the lighthouse", types.Horror)
		hall           = fixtures.AddHall(db.Store, 100, 10.0, cinema.ID, movie.ID)
		booking        = fixtures.AddBooking(db.Store, user.ID, hall.ID, types.Morning, time.Now().AddDate(0, 0, 1))
		app            = fiber.New()
		admin          = app.Group("/", middleware.JWTAuthentication(db.User), middleware.AdminAuth)
		bookingHandler = NewBookingHandler(db.Store)
	)
	booking.Date = time.Time{}
	admin.Get("/", bookingHandler.HandleGetBookings)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	for _, b := range bookings {
		b.Date = time.Time{}
	}
	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking, got %d", len(bookings))
	}
	if !contains(bookings, booking) {
		t.Fatalf("expected booking %v to be present in response", booking)
	}

	// test unauthorized access
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))

	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatal("expected status code is non 200")
	}
}
