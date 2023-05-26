package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/suhel-kap/rss-aggregator/internal/database"
)

func TestHandlerReadiness(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerReadiness)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func GetApiCfg() *apiConfig {
	godotenv.Load()
	conn, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		fmt.Printf(err.Error())
	}

	// Create a new API server with a mock database
	api := &apiConfig{
		DB: database.New(conn),
	}

	return api
}

func TestDBCreate(t *testing.T) {
	godotenv.Load()
	conn, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close()

	err = conn.Ping()
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	t.Logf("Successfully connected to database")
}

func TestCreateUser(t *testing.T) {
	godotenv.Load()
	conn, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		t.Errorf(err.Error())
	}

	// Create a new API server with a mock database
	api := &apiConfig{
		DB: database.New(conn),
	}

	// Create a new user
	user := struct {
		Name string `json:"name"`
	}{
		Name: "Testuser",
	}

	// Encode the user as JSON
	jsonUser, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	// Send a POST request to create the user
	req, err := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(jsonUser))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.handlerCreateUser)

	handler.ServeHTTP(rr, req)

	// Check that the response status code is 201 Created
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Decode the response body into a User struct
	var createdUser User
	err = json.NewDecoder(rr.Body).Decode(&createdUser)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the created user has the correct name
	if createdUser.Name != user.Name {
		t.Errorf("handler returned unexpected user name: got %v want %v",
			createdUser.Name, user.Name)
	}

	// Check that the created user has a non-zero ID
	if createdUser.ID == uuid.Nil {
		t.Errorf("handler returned user with zero ID")
	}

	// Check that the created user has a non-zero API key
	if createdUser.ApiKey == "" {
		t.Errorf("handler returned user with empty API key")
	}
}

type mockDB struct{}

func GetUUID() uuid.UUID {
	id, err := uuid.Parse("bccc704a-a573-4b19-af21-ff6ef51bb5bf")
	if err != nil {
		fmt.Println(err.Error())
	}
	return id
}

func (m *mockDB) CreateUser(user *User) error {
	user.ID = GetUUID()
	return nil
}

func (m *mockDB) GetUser(username string) (*User, error) {
	return &User{
		ID:   GetUUID(),
		Name: "testuser",
	}, nil
}
