package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/cmkqwerty/movie-ticket-booking-backend/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http/httptest"
	"testing"
)

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.tearDown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.Store)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@doe.com",
		Password:  "password12345",
	}
	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	var user types.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		t.Fatal(err)
	}
	if user.ID.IsZero() {
		t.Error("Expected user.ID to be not empty")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Error("Expected user.EncryptedPassword to be not returned")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("Expected %s, got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("Expected %s, got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("Expected %s, got %s", params.Email, user.Email)
	}
}

func TestGetUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.tearDown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.Store)
	app.Get("/:id", userHandler.HandleGetUser)

	user, err := tdb.Store.User.InsertUser(context.TODO(), &types.User{
		FirstName:         "John",
		LastName:          "Doe",
		Email:             "john@doe.com",
		EncryptedPassword: "password12345",
	})
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/"+user.ID.Hex(), nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	var returnedUser types.User
	err = json.NewDecoder(resp.Body).Decode(&returnedUser)
	if err != nil {
		t.Fatal(err)
	}
	if returnedUser.ID != user.ID {
		t.Errorf("Expected %s, got %s", user.ID, returnedUser.ID)
	}
	if len(returnedUser.EncryptedPassword) > 0 {
		t.Error("Expected user.EncryptedPassword to be not returned")
	}
	if returnedUser.FirstName != user.FirstName {
		t.Errorf("Expected %s, got %s", user.FirstName, returnedUser.FirstName)
	}
	if returnedUser.LastName != user.LastName {
		t.Errorf("Expected %s, got %s", user.LastName, returnedUser.LastName)
	}
	if returnedUser.Email != user.Email {
		t.Errorf("Expected %s, got %s", user.Email, returnedUser.Email)
	}
}

func TestGetUsers(t *testing.T) {
	tdb := setup(t)
	defer tdb.tearDown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.Store)
	app.Get("/", userHandler.HandleGetUsers)

	user1, err := tdb.Store.User.InsertUser(context.TODO(), &types.User{
		FirstName:         "John",
		LastName:          "Doe",
		Email:             "john@doe.com",
		EncryptedPassword: "password12345",
	})
	if err != nil {
		t.Fatal(err)
	}
	user2, err := tdb.Store.User.InsertUser(context.TODO(), &types.User{
		FirstName:         "Jane",
		LastName:          "Doe",
		Email:             "jane@doe.com",
		EncryptedPassword: "password12345",
	})
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	var returnedUsers []*types.User
	var response *ResourceResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	dataSlice := response.Data.([]interface{})

	for _, item := range dataSlice {
		dataMap, ok := item.(map[string]interface{})
		if !ok {
			t.Fatal("unexpected type in dataMap")
		}

		userID, _ := primitive.ObjectIDFromHex(dataMap["id"].(string))
		user := &types.User{
			ID:        userID,
			FirstName: dataMap["firstName"].(string),
			LastName:  dataMap["lastName"].(string),
			Email:     dataMap["email"].(string),
			IsAdmin:   dataMap["isAdmin"].(bool),
		}

		returnedUsers = append(returnedUsers, user)
	}

	user1.EncryptedPassword = ""
	user2.EncryptedPassword = ""
	if !contains(returnedUsers, user1) || !contains(returnedUsers, user2) {
		t.Error("Expected both users to be returned")
	}
}

func TestUpdateUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.tearDown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.Store)
	app.Put("/:id", userHandler.HandlePutUser)

	user, err := tdb.Store.User.InsertUser(context.TODO(), &types.User{
		FirstName:         "John",
		LastName:          "Doe",
		Email:             "john@doe.com",
		EncryptedPassword: "password12345",
	})
	if err != nil {
		t.Fatal(err)
	}

	params := types.UpdateUserParams{
		FirstName: "Jane",
		LastName:  "",
	}
	b, _ := json.Marshal(params)

	req := httptest.NewRequest("PUT", "/"+user.ID.Hex(), bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
}

func TestDeleteUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.tearDown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.Store)
	app.Delete("/:id", userHandler.HandleDeleteUser)

	user, err := tdb.Store.User.InsertUser(context.TODO(), &types.User{
		FirstName:         "John",
		LastName:          "Doe",
		Email:             "john@doe.com",
		EncryptedPassword: "password12345",
	})
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("DELETE", "/"+user.ID.Hex(), nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	_, err = tdb.Store.User.GetUserByID(context.TODO(), user.ID.Hex())
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
