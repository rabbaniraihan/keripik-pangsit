package models

import "time"

// ===================== PRODUCT =====================

type Product struct {
	ID          string    `db:"id" json:"id"`
	NamaProduk  string    `db:"nama_produk" json:"nama_produk"`
	TipeProduk  string    `db:"tipe_produk" json:"tipe_produk"`
	Rasa        string    `db:"rasa" json:"rasa"`
	Ukuran      string    `db:"ukuran" json:"ukuran"`
	Harga       float64   `db:"harga" json:"harga"`
	Stok        int       `db:"stok" json:"stok"`
	TanggalBuat time.Time `db:"tanggal_buat" json:"tanggal_buat"`
}

type CreateProductRequest struct {
	NamaProduk  string  `json:"nama_produk" binding:"required"`
	TipeProduk  string  `json:"tipe_produk" binding:"required"`
	Rasa        string  `json:"rasa" binding:"required"`
	Ukuran      string  `json:"ukuran" binding:"required"`
	Harga       float64 `json:"harga" binding:"required"`
	Stok        int     `json:"stok" binding:"required"`
	TanggalBuat string  `json:"tanggal_buat" binding:"required"` // format: "2006-01-02"
}

// ===================== CUSTOMER =====================

type Customer struct {
	ID          string    `db:"id" json:"id"`
	NamaCustomer string   `db:"nama_customer" json:"nama_customer"`
	Poin        int       `db:"poin" json:"poin"`
	TanggalDaftar time.Time `db:"tanggal_daftar" json:"tanggal_daftar"`
}

// ===================== TRANSACTION =====================

type Transaction struct {
	ID              string    `db:"id" json:"id"`
	CustomerID      string    `db:"customer_id" json:"customer_id"`
	NamaCustomer    string    `db:"nama_customer" json:"nama_customer"`
	ProductID       string    `db:"product_id" json:"product_id"`
	NamaProduk      string    `db:"nama_produk" json:"nama_produk"`
	UkuranProduk    string    `db:"ukuran_produk" json:"ukuran_produk"`
	Rasa            string    `db:"rasa" json:"rasa"`
	Quantity        int       `db:"quantity" json:"quantity"`
	Harga           float64   `db:"harga" json:"harga"`
	TotalHarga      float64   `db:"total_harga" json:"total_harga"`
	TanggalTransaksi time.Time `db:"tanggal_transaksi" json:"tanggal_transaksi"`
	IsCustomerBaru  bool      `db:"is_customer_baru" json:"is_customer_baru"`
}

type CreateTransactionRequest struct {
	NamaCustomer string `json:"nama_customer" binding:"required"`
	ProductID    string `json:"product_id" binding:"required"`
	Quantity     int    `json:"quantity" binding:"required,min=1"`
}

// ===================== REDEEM POIN =====================

type RedeemPoinRequest struct {
	NamaCustomer string `json:"nama_customer" binding:"required"`
	Ukuran       string `json:"ukuran" binding:"required"`
	ProductID    string `json:"product_id" binding:"required"`
}

// ===================== SUMMARY =====================

type DailySummary struct {
	TanggalMulai    string        `json:"tanggal_mulai"`
	TanggalAkhir    string        `json:"tanggal_akhir"`
	TotalCustomer   int           `json:"total_customer"`
	AdaCustomerBaru bool          `json:"ada_customer_baru"`
	TotalIncome     float64       `json:"total_income"`
	BestSeller      string        `json:"best_seller"`
	TotalProdukTerjual int        `json:"total_produk_terjual"`
	Transaksi       []Transaction `json:"transaksi"`
}

// ===================== POINT REDEMPTION LOG =====================

type RedeemLog struct {
	ID           string    `db:"id" json:"id"`
	CustomerID   string    `db:"customer_id" json:"customer_id"`
	NamaCustomer string    `db:"nama_customer" json:"nama_customer"`
	ProductID    string    `db:"product_id" json:"product_id"`
	NamaProduk   string    `db:"nama_produk" json:"nama_produk"`
	UkuranProduk string    `db:"ukuran_produk" json:"ukuran_produk"`
	PoinDigunakan int      `db:"poin_digunakan" json:"poin_digunakan"`
	TanggalRedeem time.Time `db:"tanggal_redeem" json:"tanggal_redeem"`
}
