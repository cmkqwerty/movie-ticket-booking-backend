package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/cmkqwerty/movie-ticket-booking-backend/db"
	"github.com/cmkqwerty/movie-ticket-booking-backend/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http/httptest"
	"testing"
)

const (
	testdburi  = "mongodb://localhost:27017"
	testdbname = "movie-ticket-booking-test"
)

type testdb struct {
	db.UserStore
}

func (tdb *testdb) tearDown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testdburi))
	if err != nil {
		log.Fatal(err)
	}

	return &testdb{
		UserStore: db.NewMongoUserStore(client, testdbname),
	}
}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.tearDown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
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
	json.NewDecoder(resp.Body).Decode(&user)
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
	userHandler := NewUserHandler(tdb.UserStore)
	app.Get("/:id", userHandler.HandleGetUser)

	user, err := tdb.UserStore.InsertUser(context.TODO(), &types.User{
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
	json.NewDecoder(resp.Body).Decode(&returnedUser)
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
	userHandler := NewUserHandler(tdb.UserStore)
	app.Get("/", userHandler.HandleGetUsers)

	user1, err := tdb.UserStore.InsertUser(context.TODO(), &types.User{
		FirstName:         "John",
		LastName:          "Doe",
		Email:             "john@doe.com",
		EncryptedPassword: "password12345",
	})
	if err != nil {
		t.Fatal(err)
	}
	user2, err := tdb.UserStore.InsertUser(context.TODO(), &types.User{
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

	var returnedUsers []types.User
	json.NewDecoder(resp.Body).Decode(&returnedUsers)
	if returnedUsers[0].ID != user1.ID {
		t.Errorf("Expected %s, got %s", user1.ID, returnedUsers[0].ID)
	}
	if len(returnedUsers[0].EncryptedPassword) > 0 {
		t.Error("Expected user.EncryptedPassword to be not returned")
	}
	if returnedUsers[0].FirstName != user1.FirstName {
		t.Errorf("Expected %s, got %s", user1.FirstName, returnedUsers[0].FirstName)
	}
	if returnedUsers[0].LastName != user1.LastName {
		t.Errorf("Expected %s, got %s", user1.LastName, returnedUsers[0].LastName)
	}
	if returnedUsers[0].Email != user1.Email {
		t.Errorf("Expected %s, got %s", user1.Email, returnedUsers[0].Email)
	}
	if returnedUsers[1].ID != user2.ID {
		t.Errorf("Expected %s, got %s", user2.ID, returnedUsers[1].ID)
	}
	if len(returnedUsers[1].EncryptedPassword) > 0 {
		t.Error("Expected user.EncryptedPassword to be not returned")
	}
	if returnedUsers[1].FirstName != user2.FirstName {
		t.Errorf("Expected %s, got %s", user2.FirstName, returnedUsers[1].FirstName)
	}
	if returnedUsers[1].LastName != user2.LastName {
		t.Errorf("Expected %s, got %s", user2.LastName, returnedUsers[1].LastName)
	}
	if returnedUsers[1].Email != user2.Email {
		t.Errorf("Expected %s, got %s", user2.Email, returnedUsers[1].Email)
	}
}

func TestUpdateUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.tearDown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Put("/:id", userHandler.HandlePutUser)

	user, err := tdb.UserStore.InsertUser(context.TODO(), &types.User{
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
	userHandler := NewUserHandler(tdb.UserStore)
	app.Delete("/:id", userHandler.HandleDeleteUser)

	user, err := tdb.UserStore.InsertUser(context.TODO(), &types.User{
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

	_, err = tdb.UserStore.GetUserByID(context.TODO(), user.ID.Hex())
	if err == nil {
		t.Error("Expected error, got nil")
	}
}