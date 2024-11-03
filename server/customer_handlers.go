package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"CustomerCRUD/pkg/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func (s *Server) GetAllCustomers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	customers, err := s.repository.GetAllCustomers(ctx)
	if err != nil {
		log.Errorf("error getting customers: %v", err)
		http.Error(w, "Problem when retrieving customers, please try again later", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

func (s *Server) GetCustomerByEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	email := mux.Vars(r)["email"]
	customer, err := s.repository.GetCustomerByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Customer not found", http.StatusNotFound)
		} else {
			log.Errorf("error getting customer by email: %v", err)
			http.Error(w, "Failed to retrieve customer", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

func (s *Server) GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Errorf("failed to parse customer ID: %s", err)
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	customer, err := s.repository.GetCustomerByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Customer not found", http.StatusNotFound)
		} else {
			log.Errorf("error getting customer by ID: %v", err)
			http.Error(w, "Failed to retrieve customer", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

func (s *Server) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var c models.Customer
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if c.Email == "" || c.FirstName == "" || c.LastName == "" {
		http.Error(w, "First name, last name, and email are required", http.StatusBadRequest)
		return
	}

	c.ID = uuid.New()
	err := s.repository.CreateCustomer(ctx, c)
	if err != nil {
		log.Errorf("Failed to create customer: %s", err)
		http.Error(w, "Failed to create customer", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func (s *Server) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid customer ID", http.StatusBadRequest)
		return
	}

	var c models.Customer
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		log.Errorf("failed to parse customer update: %v", err)
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	c.ID = id

	err = s.repository.UpdateCustomer(ctx, c)
	if err != nil {
		log.Errorf("failed to update customer: %v", err)
		http.Error(w, "Failed to update customer", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func (s *Server) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Errorf("failed to parse customer ID: %v", err)
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	err = s.repository.DeleteCustomer(ctx, id)
	if err != nil {
		log.Errorf("failed to delete customer: %v", err)
		http.Error(w, "Failed to delete customer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
