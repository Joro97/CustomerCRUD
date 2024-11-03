package models

import "github.com/google/uuid"

type Customer struct {
	ID          uuid.UUID `json:"id"`
	FirstName   string    `json:"first_name"`
	MiddleName  string    `json:"middle_name,omitempty"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number,omitempty"`
}
