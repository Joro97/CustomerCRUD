package repository

import (
	"context"
	"database/sql"

	"CustomerCRUD/pkg/models"
	"CustomerCRUD/utils"

	"github.com/google/uuid"
)

type CustomerRepository interface {
	GetAllCustomers(ctx context.Context) ([]models.Customer, error)
	GetCustomerByID(ctx context.Context, customerID uuid.UUID) (*models.Customer, error)
	GetCustomerByEmail(ctx context.Context, email string) (*models.Customer, error)
	CreateCustomer(ctx context.Context, customer models.Customer) error
	UpdateCustomer(ctx context.Context, customer models.Customer) error
	DeleteCustomer(ctx context.Context, customerID uuid.UUID) error
}

func GetDB(isLocalDb bool, connStrEnvVar string) (*sql.DB, error) {
	if isLocalDb {
		return utils.GetLocalDB()
	} else {
		return nil, nil // TODO:
	}
}
