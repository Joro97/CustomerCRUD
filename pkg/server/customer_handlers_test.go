// server/handlers_test.go
package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"CustomerCRUD/pkg/models"
	"CustomerCRUD/pkg/repository"
	"CustomerCRUD/pkg/repository/mocks"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestServer(mockRepo repository.CustomerRepository) *Server {
	return &Server{
		repository: mockRepo,
	}
}

func TestGetAllCustomers(t *testing.T) {
	mockRepo := &mocks.CustomerRepository{}
	s := newTestServer(mockRepo)

	expectedCustomers := []models.Customer{
		{
			ID:        uuid.New(),
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		},
		{
			ID:        uuid.New(),
			FirstName: "Jane",
			LastName:  "Smith",
			Email:     "jane.smith@example.com",
		},
	}

	mockRepo.On("GetAllCustomers", mock.Anything).Return(expectedCustomers, nil)

	req, err := http.NewRequest("GET", "/customers", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.GetAllCustomers)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	var customers []models.Customer
	err = json.Unmarshal(rr.Body.Bytes(), &customers)
	if err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	assert.Equal(t, customers, expectedCustomers)

	mockRepo.AssertExpectations(t)
}

func TestGetAllCustomers_Error(t *testing.T) {
	mockRepo := &mocks.CustomerRepository{}
	s := newTestServer(mockRepo)

	mockRepo.On("GetAllCustomers", mock.Anything).Return(nil, errors.New("database error"))

	req, err := http.NewRequest("GET", "/customers", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.GetAllCustomers)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, status)
	}

	expectedBody := "Problem when retrieving customers, please try again later\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rr.Body.String())
	}

	mockRepo.AssertExpectations(t)
}

func TestGetCustomerByEmail(t *testing.T) {
	mockRepo := &mocks.CustomerRepository{}
	s := newTestServer(mockRepo)

	expectedCustomer := &models.Customer{
		ID:        uuid.New(),
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	mockRepo.On("GetCustomerByEmail", mock.Anything, "john.doe@example.com").Return(expectedCustomer, nil)

	req, err := http.NewRequest("GET", "/customers/email/john.doe@example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"email": "john.doe@example.com"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.GetCustomerByEmail)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	var customer models.Customer
	err = json.Unmarshal(rr.Body.Bytes(), &customer)
	if err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	assert.Equal(t, expectedCustomer.ID, customer.ID)

	mockRepo.AssertExpectations(t)
}

func TestGetCustomerByEmail_NotFound(t *testing.T) {
	mockRepo := &mocks.CustomerRepository{}
	s := newTestServer(mockRepo)

	mockRepo.On("GetCustomerByEmail", mock.Anything, "unknown@example.com").Return(&models.Customer{}, sql.ErrNoRows)

	req, err := http.NewRequest("GET", "/customers/email/unknown@example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"email": "unknown@example.com"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.GetCustomerByEmail)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, status)
	}

	expectedBody := "Customer not found\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rr.Body.String())
	}

	mockRepo.AssertExpectations(t)
}

func TestGetCustomerByEmail_Error(t *testing.T) {
	mockRepo := &mocks.CustomerRepository{}
	s := newTestServer(mockRepo)

	mockRepo.On("GetCustomerByEmail", mock.Anything, "john.doe@example.com").Return(&models.Customer{}, errors.New("database error"))

	req, err := http.NewRequest("GET", "/customers/email/john.doe@example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"email": "john.doe@example.com"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.GetCustomerByEmail)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, status)
	}

	expectedBody := "Failed to retrieve customer\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rr.Body.String())
	}

	mockRepo.AssertExpectations(t)
}

func TestGetCustomerByID(t *testing.T) {
	mockRepo := &mocks.CustomerRepository{}
	s := newTestServer(mockRepo)

	id := uuid.New()
	expectedCustomer := models.Customer{
		ID:        id,
		FirstName: "Alice",
		LastName:  "Wonderland",
		Email:     "alice@example.com",
	}

	mockRepo.On("GetCustomerByID", mock.Anything, id).Return(&expectedCustomer, nil)

	req, err := http.NewRequest("GET", "/customers/"+id.String(), nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": id.String()})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.GetCustomerByID)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	var customer models.Customer
	err = json.Unmarshal(rr.Body.Bytes(), &customer)
	if err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	assert.Equal(t, expectedCustomer, customer)

	mockRepo.AssertExpectations(t)
}

func TestGetCustomerByID_InvalidID(t *testing.T) {
	mockRepo := &mocks.CustomerRepository{}
	s := newTestServer(mockRepo)

	req, err := http.NewRequest("GET", "/customers/invalid-uuid", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "invalid-uuid"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.GetCustomerByID)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
	}

	expectedBody := "Invalid customer ID\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rr.Body.String())
	}

	mockRepo.AssertExpectations(t)
}

func TestGetCustomerByID_NotFound(t *testing.T) {
	mockRepo := &mocks.CustomerRepository{}
	s := newTestServer(mockRepo)

	id := uuid.New()
	mockRepo.On("GetCustomerByID", mock.Anything, id).Return(&models.Customer{}, sql.ErrNoRows)

	req, err := http.NewRequest("GET", "/customers/"+id.String(), nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": id.String()})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.GetCustomerByID)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, status)
	}

	expectedBody := "Customer not found\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rr.Body.String())
	}

	mockRepo.AssertExpectations(t)
}

func TestGetCustomerByID_Error(t *testing.T) {
	mockRepo := &mocks.CustomerRepository{}
	s := newTestServer(mockRepo)

	id := uuid.New()
	mockRepo.On("GetCustomerByID", mock.Anything, id).Return(&models.Customer{}, errors.New("database error"))

	req, err := http.NewRequest("GET", "/customers/"+id.String(), nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": id.String()})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.GetCustomerByID)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, status)
	}

	expectedBody := "Failed to retrieve customer\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rr.Body.String())
	}

	mockRepo.AssertExpectations(t)
}

func TestCreateCustomer(t *testing.T) {
	mockRepo := &mocks.CustomerRepository{}
	s := newTestServer(mockRepo)

	inputCustomer := models.Customer{
		FirstName: "Bob",
		LastName:  "Builder",
		Email:     "bob.builder@example.com",
	}

	mockRepo.On("CreateCustomer", mock.Anything, mock.MatchedBy(func(c models.Customer) bool {
		return c.FirstName == inputCustomer.FirstName &&
			c.LastName == inputCustomer.LastName &&
			c.Email == inputCustomer.Email &&
			c.ID != uuid.Nil
	})).Return(nil)

	body, _ := json.Marshal(inputCustomer)
	req, err := http.NewRequest("POST", "/customers", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.CreateCustomer)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, status)
	}

	var customer models.Customer
	err = json.Unmarshal(rr.Body.Bytes(), &customer)
	if err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	if customer.FirstName != inputCustomer.FirstName || customer.LastName != inputCustomer.LastName || customer.Email != inputCustomer.Email {
		t.Errorf("Expected customer data to match input")
	}

	if customer.ID == uuid.Nil {
		t.Errorf("Expected customer ID to be set")
	}

	mockRepo.AssertExpectations(t)
}

func TestCreateCustomer_InvalidJSON(t *testing.T) {
	mockRepo := &mocks.CustomerRepository{}
	s := newTestServer(mockRepo)

	req, err := http.NewRequest("POST", "/customers", bytes.NewBuffer([]byte("invalid json")))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.CreateCustomer)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
	}

	expectedBody := "Invalid request payload\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rr.Body.String())
	}

	mockRepo.AssertExpectations(t)
}

func TestCreateCustomer_MissingFields(t *testing.T) {
	mockRepo := &mocks.CustomerRepository{}
	s := newTestServer(mockRepo)

	inputCustomer := models.Customer{
		Email: "missing.fields@example.com",
	}

	body, _ := json.Marshal(inputCustomer)
	req, err := http.NewRequest("POST", "/customers", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.CreateCustomer)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
	}

	expectedBody := "First name, last name, and email are required\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rr.Body.String())
	}

	mockRepo.AssertExpectations(t)
}

func TestCreateCustomer_Error(t *testing.T) {
	mockRepo := &mocks.CustomerRepository{}
	s := newTestServer(mockRepo)

	inputCustomer := models.Customer{
		FirstName: "Error",
		LastName:  "Case",
		Email:     "error.case@example.com",
	}

	mockRepo.On("CreateCustomer", mock.Anything, mock.AnythingOfType("models.Customer")).Return(errors.New("database error"))

	body, _ := json.Marshal(inputCustomer)
	req, err := http.NewRequest("POST", "/customers", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.CreateCustomer)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, status)
	}

	expectedBody := "Failed to create customer\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rr.Body.String())
	}

	mockRepo.AssertExpectations(t)
}

func TestUpdateCustomer(t *testing.T) {
	mockRepo := &mocks.CustomerRepository{}
	s := newTestServer(mockRepo)

	id := uuid.New()
	updatedCustomer := models.Customer{
		ID:        id,
		FirstName: "Updated",
		LastName:  "User",
		Email:     "updated.user@example.com",
	}

	mockRepo.On("UpdateCustomer", mock.Anything, updatedCustomer).Return(nil)

	body, _ := json.Marshal(updatedCustomer)
	req, err := http.NewRequest("PUT", "/customers/"+id.String(), bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": id.String()})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.UpdateCustomer)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	var customer models.Customer
	err = json.Unmarshal(rr.Body.Bytes(), &customer)
	if err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	assert.Equal(t, customer, updatedCustomer)

	mockRepo.AssertExpectations(t)
}

func TestUpdateCustomer_InvalidID(t *testing.T) {
	mockRepo := &mocks.CustomerRepository{}
	s := newTestServer(mockRepo)

	req, err := http.NewRequest("PUT", "/customers/invalid-uuid", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "invalid-uuid"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.UpdateCustomer)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
	}

	expectedBody := "invalid customer ID\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rr.Body.String())
	}

	mockRepo.AssertExpectations(t)
}

func TestUpdateCustomer_InvalidJSON(t *testing.T) {
	mockRepo := &mocks.CustomerRepository{}
	s := newTestServer(mockRepo)

	id := uuid.New()
	req, err := http.NewRequest("PUT", "/customers/"+id.String(), bytes.NewBuffer([]byte("invalid json")))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": id.String()})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.UpdateCustomer)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
	}

	expectedBody := "invalid request payload\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rr.Body.String())
	}

	mockRepo.AssertExpectations(t)
}

func TestUpdateCustomer_Error(t *testing.T) {
	mockRepo := &mocks.CustomerRepository{}
	s := newTestServer(mockRepo)

	id := uuid.New()
	updatedCustomer := models.Customer{
		ID:        id,
		FirstName: "Error",
		LastName:  "Case",
		Email:     "error.case@example.com",
	}

	mockRepo.On("UpdateCustomer", mock.Anything, updatedCustomer).Return(errors.New("database error"))

	body, _ := json.Marshal(updatedCustomer)
	req, err := http.NewRequest("PUT", "/customers/"+id.String(), bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": id.String()})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.UpdateCustomer)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, status)
	}

	expectedBody := "Failed to update customer\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rr.Body.String())
	}

	mockRepo.AssertExpectations(t)
}

func TestDeleteCustomer(t *testing.T) {
	mockRepo := &mocks.CustomerRepository{}
	s := newTestServer(mockRepo)

	id := uuid.New()
	mockRepo.On("DeleteCustomer", mock.Anything, id).Return(nil)

	req, err := http.NewRequest("DELETE", "/customers/"+id.String(), nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": id.String()})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.DeleteCustomer)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, status)
	}

	mockRepo.AssertExpectations(t)
}

func TestDeleteCustomer_InvalidID(t *testing.T) {
	mockRepo := &mocks.CustomerRepository{}
	s := newTestServer(mockRepo)

	req, err := http.NewRequest("DELETE", "/customers/invalid-uuid", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "invalid-uuid"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.DeleteCustomer)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
	}

	expectedBody := "Invalid customer ID\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rr.Body.String())
	}

	mockRepo.AssertExpectations(t)
}

func TestDeleteCustomer_Error(t *testing.T) {
	mockRepo := &mocks.CustomerRepository{}
	s := newTestServer(mockRepo)

	id := uuid.New()
	mockRepo.On("DeleteCustomer", mock.Anything, id).Return(errors.New("database error"))

	req, err := http.NewRequest("DELETE", "/customers/"+id.String(), nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": id.String()})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.DeleteCustomer)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, status)
	}

	expectedBody := "Failed to delete customer\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rr.Body.String())
	}

	mockRepo.AssertExpectations(t)
}
