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

type TransactionRepository interface {
	CreateTransaction(req models.CreateTransactionRequest, customer *models.Customer, product *models.Product, isNew bool) (*models.Transaction, error)
	GetTransactionsByDateRange(startDate, endDate string) ([]models.Transaction, error)
	GetDailySummary(startDate, endDate string) (*models.DailySummary, error)
}

type transactionRepository struct{}

func NewTransactionRepository() TransactionRepository {
	return &transactionRepository{}
}

func (r *transactionRepository) CreateTransaction(
	req models.CreateTransactionRequest,
	customer *models.Customer,
	product *models.Product,
	isNew bool,
) (*models.Transaction, error) {
	totalHarga := product.Harga * float64(req.Quantity)

	trx := &models.Transaction{
		ID:               uuid.New().String(),
		CustomerID:       customer.ID,
		NamaCustomer:     customer.NamaCustomer,
		ProductID:        product.ID,
		NamaProduk:       product.NamaProduk,
		UkuranProduk:     product.Ukuran,
		Rasa:             product.Rasa,
		Quantity:         req.Quantity,
		Harga:            product.Harga,
		TotalHarga:       totalHarga,
		TanggalTransaksi: time.Now(),
		IsCustomerBaru:   isNew,
	}

	query := `INSERT INTO transactions 
	          (id, customer_id, nama_customer, product_id, nama_produk, ukuran_produk, rasa, quantity, harga, total_harga, tanggal_transaksi, is_customer_baru)
	          VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING *`

	err := config.DB.QueryRowx(query,
		trx.ID, trx.CustomerID, trx.NamaCustomer,
		trx.ProductID, trx.NamaProduk, trx.UkuranProduk,
		trx.Rasa, trx.Quantity, trx.Harga, trx.TotalHarga,
		trx.TanggalTransaksi, trx.IsCustomerBaru,
	).StructScan(trx)
	if err != nil {
		return nil, err
	}

	// Invalidate transaction cache
	ctx := context.Background()
	config.Redis.Del(ctx, "transactions:*")

	return trx, nil
}

func (r *transactionRepository) GetTransactionsByDateRange(startDate, endDate string) ([]models.Transaction, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("transactions:range:%s:%s", startDate, endDate)

	cached, err := config.Redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var transactions []models.Transaction
		if jsonErr := json.Unmarshal([]byte(cached), &transactions); jsonErr == nil {
			return transactions, nil
		}
	}

	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("format tanggal mulai tidak valid")
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("format tanggal akhir tidak valid")
	}
	// Include the entire end day
	endOfDay := end.Add(24*time.Hour - time.Second)

	query := `SELECT * FROM transactions 
	          WHERE tanggal_transaksi >= $1 AND tanggal_transaksi <= $2
	          ORDER BY tanggal_transaksi DESC`

	var transactions []models.Transaction
	err = config.DB.Select(&transactions, query, start, endOfDay)
	if err != nil {
		return nil, err
	}

	if data, jsonErr := json.Marshal(transactions); jsonErr == nil {
		config.Redis.Set(ctx, cacheKey, string(data), 2*time.Minute)
	}

	return transactions, nil
}

func (r *transactionRepository) GetDailySummary(startDate, endDate string) (*models.DailySummary, error) {
	transactions, err := r.GetTransactionsByDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Count unique customers
	customerSet := make(map[string]bool)
	adaCustomerBaru := false
	var totalIncome float64
	productSales := make(map[string]int)
	totalProdukTerjual := 0

	for _, t := range transactions {
		customerSet[t.CustomerID] = true
		if t.IsCustomerBaru {
			adaCustomerBaru = true
		}
		totalIncome += t.TotalHarga
		key := fmt.Sprintf("%s - %s", t.NamaProduk, t.Rasa)
		productSales[key] += t.Quantity
		totalProdukTerjual += t.Quantity
	}

	// Find best seller
	bestSeller := ""
	maxQty := 0
	for prodKey, qty := range productSales {
		if qty > maxQty {
			maxQty = qty
			bestSeller = prodKey
		}
	}

	summary := &models.DailySummary{
		TanggalMulai:       startDate,
		TanggalAkhir:       endDate,
		TotalCustomer:      len(customerSet),
		AdaCustomerBaru:    adaCustomerBaru,
		TotalIncome:        totalIncome,
		BestSeller:         bestSeller,
		TotalProdukTerjual: totalProdukTerjual,
		Transaksi:          transactions,
	}

	return summary, nil
}
