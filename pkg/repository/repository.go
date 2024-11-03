package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"CustomerCRUD/pkg/models"
	"CustomerCRUD/utils"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type CustomerRepository interface {
	GetAllCustomers(ctx context.Context) ([]models.Customer, error)
	GetCustomerByID(ctx context.Context, customerID uuid.UUID) (*models.Customer, error)
	GetCustomerByEmail(ctx context.Context, email string) (*models.Customer, error)
	CreateCustomer(ctx context.Context, customer models.Customer) error
	UpdateCustomer(ctx context.Context, customer models.Customer) error
	DeleteCustomer(ctx context.Context, customerID uuid.UUID) error
}

type customerRepository struct {
	db *sql.DB
}

func (r customerRepository) GetAllCustomers(ctx context.Context) ([]models.Customer, error) {
	rows, err := r.db.Query("SELECT * FROM customers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var c models.Customer
		err := rows.Scan(&c.ID, &c.FirstName, &c.MiddleName, &c.LastName, &c.Email, &c.PhoneNumber)
		if err != nil {
			return nil, fmt.Errorf("error scanning customer rows: %w", err)
		}
		customers = append(customers, c)
	}
	return customers, nil
}

func (r customerRepository) GetCustomerByID(ctx context.Context, customerID uuid.UUID) (*models.Customer, error) {
	var c models.Customer
	query := "SELECT * FROM customers WHERE id = $1"
	err := r.db.QueryRow(query, customerID).Scan(&c.ID, &c.FirstName, &c.MiddleName, &c.LastName, &c.Email, &c.PhoneNumber)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r customerRepository) GetCustomerByEmail(ctx context.Context, email string) (*models.Customer, error) {
	var c models.Customer
	query := "SELECT * FROM customers WHERE email = $1"
	err := r.db.QueryRow(query, email).Scan(&c.ID, &c.FirstName, &c.MiddleName, &c.LastName, &c.Email, &c.PhoneNumber)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r customerRepository) CreateCustomer(ctx context.Context, customer models.Customer) error {
	_, err := r.db.Exec(
		`INSERT INTO customers (id, first_name, middle_name, last_name, email, phone_number)
     VALUES ($1, $2, $3, $4, $5, $6)`,
		customer.ID, customer.FirstName, customer.MiddleName, customer.LastName, customer.Email, customer.PhoneNumber)

	if err != nil {
		return fmt.Errorf("error inserting customer rows: %w", err)
	}
	return nil
}

func (r customerRepository) UpdateCustomer(ctx context.Context, customer models.Customer) error {
	_, err := r.db.Exec(
		`UPDATE customers SET first_name=$1, middle_name=$2, last_name=$3, email=$4, phone_number=$5
         WHERE id=$6`,
		customer.FirstName, customer.MiddleName, customer.LastName, customer.Email, customer.PhoneNumber, customer.ID)
	if err != nil {
		return fmt.Errorf("error updating customer: %w", err)
	}
	return nil
}

func (r customerRepository) DeleteCustomer(ctx context.Context, customerID uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM customers WHERE id=$1", customerID)
	if err != nil {
		return fmt.Errorf("error deleting customer: %w", err)
	}
	return nil
}

func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func GetDB(isLocalDb bool, connStrEnvVar string) (*sql.DB, error) {
	if isLocalDb {
		return utils.GetLocalDB()
	} else {
		connStr := os.Getenv(connStrEnvVar)
		return sql.Open("postgres", connStr)
	}
}
