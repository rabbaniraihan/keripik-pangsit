package repository

import (
	"fmt"
	"keripik-pangsit/config"
	"keripik-pangsit/models"
	"time"

	"github.com/google/uuid"
)

// Poin yang dibutuhkan berdasarkan ukuran
var PoinPerUkuran = map[string]int{
	"Small":  200,
	"Medium": 300,
	"Large":  500,
}

type RedeemRepository interface {
	RedeemPoin(req models.RedeemPoinRequest) (*models.RedeemLog, error)
	GetRedeemLogs() ([]models.RedeemLog, error)
}

type redeemRepository struct {
	customerRepo CustomerRepository
	productRepo  ProductRepository
}

func NewRedeemRepository(cr CustomerRepository, pr ProductRepository) RedeemRepository {
	return &redeemRepository{customerRepo: cr, productRepo: pr}
}

func (r *redeemRepository) RedeemPoin(req models.RedeemPoinRequest) (*models.RedeemLog, error) {
	// Get customer
	customer, err := r.customerRepo.GetCustomerByName(req.NamaCustomer)
	if err != nil {
		return nil, err
	}

	// Get product
	product, err := r.productRepo.GetProductByID(req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("produk tidak ditemukan")
	}

	// Validate ukuran
	poinDibutuhkan, ok := PoinPerUkuran[req.Ukuran]
	if !ok {
		return nil, fmt.Errorf("ukuran tidak valid: Small, Medium, atau Large")
	}

	// Check product ukuran matches
	if product.Ukuran != req.Ukuran {
		return nil, fmt.Errorf("ukuran produk tidak sesuai dengan pilihan redeem")
	}

	// Check customer poin
	if customer.Poin < poinDibutuhkan {
		return nil, fmt.Errorf("poin tidak cukup. dibutuhkan %d poin, Anda memiliki %d poin", poinDibutuhkan, customer.Poin)
	}

	// Check stok produk
	if product.Stok < 1 {
		return nil, fmt.Errorf("stok produk habis")
	}

	// Begin transaction
	tx, err := config.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Deduct poin
	_, err = tx.Exec(`UPDATE customers SET poin = poin - $1 WHERE id = $2`, poinDibutuhkan, customer.ID)
	if err != nil {
		return nil, err
	}

	// Reduce stok
	result, err := tx.Exec(`UPDATE products SET stok = stok - 1 WHERE id = $1 AND stok >= 1`, product.ID)
	if err != nil {
		return nil, err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return nil, fmt.Errorf("stok produk habis")
	}

	// Create redeem log
	log := &models.RedeemLog{
		ID:            uuid.New().String(),
		CustomerID:    customer.ID,
		NamaCustomer:  customer.NamaCustomer,
		ProductID:     product.ID,
		NamaProduk:    product.NamaProduk,
		UkuranProduk:  product.Ukuran,
		PoinDigunakan: poinDibutuhkan,
		TanggalRedeem: time.Now(),
	}

	_, err = tx.Exec(`INSERT INTO redeem_logs (id, customer_id, nama_customer, product_id, nama_produk, ukuran_produk, poin_digunakan, tanggal_redeem)
	                  VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		log.ID, log.CustomerID, log.NamaCustomer,
		log.ProductID, log.NamaProduk, log.UkuranProduk,
		log.PoinDigunakan, log.TanggalRedeem,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Invalidate caches
	// (customer and product cache cleared separately)

	return log, nil
}

func (r *redeemRepository) GetRedeemLogs() ([]models.RedeemLog, error) {
	var logs []models.RedeemLog
	query := `SELECT * FROM redeem_logs ORDER BY tanggal_redeem DESC`
	err := config.DB.Select(&logs, query)
	return logs, err
}
