// server/integration_test.go
package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"CustomerCRUD/pkg/models"
	"CustomerCRUD/pkg/repository"
	"CustomerCRUD/pkg/server"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var serverAddress = "localhost:8081"

func TestMain(m *testing.M) {
	// Set up the test database connection
	testDBURL := os.Getenv("TEST_DATABASE_URL")
	if testDBURL == "" {
		log.Fatal("TEST_DATABASE_URL environment variable not set")
	}

	// Initialize the repository
	db, err := repository.GetDB(false, "TEST_DATABASE_URL")
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := runMigrations(testDBURL); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	repo := repository.NewCustomerRepository(db)

	// Initialize the server
	srv := server.NewServer(repo)
	srv.SetupRoutes()

	go func() {
		if err := http.ListenAndServe(":8081", srv.Router); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Wait for the server to start
	time.Sleep(1 * time.Second)

	// Run tests
	code := m.Run()

	os.Exit(code)
}

func runMigrations(dbURL string) error {
	// Get the absolute path to the migrations directory
	dir, err := filepath.Abs("../../migrations")
	if err != nil {
		return err
	}

	m, err := migrate.New(
		"file://"+dir,
		dbURL)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func TestIntegration_CreateAndGetCustomer(t *testing.T) {
	baseURL := "http://" + serverAddress

	// Create a new customer
	customer := models.Customer{
		FirstName: "Integration",
		LastName:  "Test",
		Email:     "integration.test@example.com",
	}

	createBody, _ := json.Marshal(customer)
	createResp, err := http.Post(baseURL+"/customers", "application/json", bytes.NewBuffer(createBody))
	if err != nil {
		t.Fatal(err)
	}
	defer createResp.Body.Close()

	if createResp.StatusCode != http.StatusCreated {
		fmt.Println(createResp.Body)
		t.Fatalf("Expected status 201 Created, got %d", createResp.StatusCode)
	}

	var createdCustomer models.Customer
	err = json.NewDecoder(createResp.Body).Decode(&createdCustomer)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Get the customer by ID
	getResp, err := http.Get(baseURL + "/customers/" + createdCustomer.ID.String())
	if err != nil {
		t.Fatal(err)
	}
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %d", getResp.StatusCode)
	}

	var fetchedCustomer models.Customer
	err = json.NewDecoder(getResp.Body).Decode(&fetchedCustomer)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if fetchedCustomer != createdCustomer {
		t.Fatalf("Expected fetched customer to be %v, got %v", createdCustomer, fetchedCustomer)
	}

	// Clean up: delete the customer
	req, err := http.NewRequest("DELETE", baseURL+"/customers/"+createdCustomer.ID.String(), nil)
	if err != nil {
		t.Fatal(err)
	}
	delResp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer delResp.Body.Close()

	if delResp.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status 204 No Content, got %d", delResp.StatusCode)
	}
}
