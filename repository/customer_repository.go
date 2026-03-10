package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"keripik-pangsit/config"
	"keripik-pangsit/models"
	"time"

	"github.com/google/uuid"
)

type CustomerRepository interface {
	GetOrCreateCustomer(namaCustomer string) (*models.Customer, bool, error)
	GetAllCustomers() ([]models.Customer, error)
	GetCustomerByName(nama string) (*models.Customer, error)
	AddPoin(customerID string, poin int) error
	DeductPoin(customerID string, poin int) error
}

type customerRepository struct{}

func NewCustomerRepository() CustomerRepository {
	return &customerRepository{}
}

// Returns customer, isNew, error
func (r *customerRepository) GetOrCreateCustomer(namaCustomer string) (*models.Customer, bool, error) {
	ctx := context.Background()

	// Check if customer exists
	var customer models.Customer
	query := `SELECT * FROM customers WHERE LOWER(nama_customer) = LOWER($1)`
	err := config.DB.QueryRowx(query, namaCustomer).StructScan(&customer)
	if err == nil {
		return &customer, false, nil
	}

	// Create new customer
	newCustomer := &models.Customer{
		ID:            uuid.New().String(),
		NamaCustomer:  namaCustomer,
		Poin:          0,
		TanggalDaftar: time.Now(),
	}

	insertQuery := `INSERT INTO customers (id, nama_customer, poin, tanggal_daftar)
	                VALUES ($1, $2, $3, $4) RETURNING *`

	err = config.DB.QueryRowx(insertQuery,
		newCustomer.ID, newCustomer.NamaCustomer, newCustomer.Poin, newCustomer.TanggalDaftar,
	).StructScan(newCustomer)
	if err != nil {
		return nil, false, err
	}

	// Invalidate cache
	config.Redis.Del(ctx, "customers:all")

	return newCustomer, true, nil
}

func (r *customerRepository) GetAllCustomers() ([]models.Customer, error) {
	ctx := context.Background()
	cacheKey := "customers:all"

	cached, err := config.Redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var customers []models.Customer
		if jsonErr := json.Unmarshal([]byte(cached), &customers); jsonErr == nil {
			return customers, nil
		}
	}

	var customers []models.Customer
	query := `SELECT * FROM customers ORDER BY tanggal_daftar DESC`
	err = config.DB.Select(&customers, query)
	if err != nil {
		return nil, err
	}

	if data, jsonErr := json.Marshal(customers); jsonErr == nil {
		config.Redis.Set(ctx, cacheKey, string(data), 5*time.Minute)
	}

	return customers, nil
}

func (r *customerRepository) GetCustomerByName(nama string) (*models.Customer, error) {
	var customer models.Customer
	query := `SELECT * FROM customers WHERE LOWER(nama_customer) = LOWER($1)`
	err := config.DB.QueryRowx(query, nama).StructScan(&customer)
	if err != nil {
		return nil, fmt.Errorf("customer tidak ditemukan")
	}
	return &customer, nil
}

func (r *customerRepository) AddPoin(customerID string, poin int) error {
	query := `UPDATE customers SET poin = poin + $1 WHERE id = $2`
	_, err := config.DB.Exec(query, poin, customerID)
	if err != nil {
		return err
	}
	ctx := context.Background()
	config.Redis.Del(ctx, "customers:all")
	return nil
}

func (r *customerRepository) DeductPoin(customerID string, poin int) error {
	query := `UPDATE customers SET poin = poin - $1 WHERE id = $2 AND poin >= $1`
	result, err := config.DB.Exec(query, poin, customerID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("poin tidak mencukupi")
	}
	ctx := context.Background()
	config.Redis.Del(ctx, "customers:all")
	return nil
}
