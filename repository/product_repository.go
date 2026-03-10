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

type ProductRepository interface {
	CreateProduct(req models.CreateProductRequest) (*models.Product, error)
	GetProductsByDate(tanggal string) ([]models.Product, error)
	GetProductByID(id string) (*models.Product, error)
	ReduceStock(productID string, qty int) error
}

type productRepository struct{}

func NewProductRepository() ProductRepository {
	return &productRepository{}
}

func (r *productRepository) CreateProduct(req models.CreateProductRequest) (*models.Product, error) {
	tanggal, err := time.Parse("2006-01-02", req.TanggalBuat)
	if err != nil {
		return nil, fmt.Errorf("format tanggal tidak valid, gunakan YYYY-MM-DD")
	}

	product := &models.Product{
		ID:          uuid.New().String(),
		NamaProduk:  req.NamaProduk,
		TipeProduk:  req.TipeProduk,
		Rasa:        req.Rasa,
		Ukuran:      req.Ukuran,
		Harga:       req.Harga,
		Stok:        req.Stok,
		TanggalBuat: tanggal,
	}

	query := `INSERT INTO products (id, nama_produk, tipe_produk, rasa, ukuran, harga, stok, tanggal_buat)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`

	err = config.DB.QueryRowx(query,
		product.ID, product.NamaProduk, product.TipeProduk,
		product.Rasa, product.Ukuran, product.Harga, product.Stok, product.TanggalBuat,
	).StructScan(product)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	ctx := context.Background()
	config.Redis.Del(ctx, fmt.Sprintf("products:date:%s", req.TanggalBuat))

	return product, nil
}

func (r *productRepository) GetProductsByDate(tanggal string) ([]models.Product, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("products:date:%s", tanggal)

	// Try cache first
	cached, err := config.Redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var products []models.Product
		if jsonErr := json.Unmarshal([]byte(cached), &products); jsonErr == nil {
			return products, nil
		}
	}

	// Parse date for range query
	parsedDate, err := time.Parse("2006-01-02", tanggal)
	if err != nil {
		return nil, fmt.Errorf("format tanggal tidak valid")
	}

	query := `SELECT id, nama_produk, tipe_produk, rasa, ukuran, harga, stok, tanggal_buat
	          FROM products
	          WHERE DATE(tanggal_buat) = $1
	          ORDER BY tanggal_buat DESC`

	var products []models.Product
	err = config.DB.Select(&products, query, parsedDate.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}

	// Store in cache for 5 minutes
	if data, jsonErr := json.Marshal(products); jsonErr == nil {
		config.Redis.Set(ctx, cacheKey, string(data), 5*time.Minute)
	}

	return products, nil
}

func (r *productRepository) GetProductByID(id string) (*models.Product, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("product:%s", id)

	cached, err := config.Redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var product models.Product
		if jsonErr := json.Unmarshal([]byte(cached), &product); jsonErr == nil {
			return &product, nil
		}
	}

	var product models.Product
	query := `SELECT * FROM products WHERE id = $1`
	err = config.DB.QueryRowx(query, id).StructScan(&product)
	if err != nil {
		return nil, err
	}

	if data, jsonErr := json.Marshal(product); jsonErr == nil {
		config.Redis.Set(ctx, cacheKey, string(data), 10*time.Minute)
	}

	return &product, nil
}

func (r *productRepository) ReduceStock(productID string, qty int) error {
	query := `UPDATE products SET stok = stok - $1 WHERE id = $2 AND stok >= $1`
	result, err := config.DB.Exec(query, qty, productID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("stok tidak mencukupi atau produk tidak ditemukan")
	}

	// Invalidate product cache
	ctx := context.Background()
	config.Redis.Del(ctx, fmt.Sprintf("product:%s", productID))

	return nil
}
