package controller

import (
	"keripik-pangsit/models"
	"keripik-pangsit/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	transactionRepo repository.TransactionRepository
	customerRepo    repository.CustomerRepository
	productRepo     repository.ProductRepository
}

func NewTransactionController(
	tr repository.TransactionRepository,
	cr repository.CustomerRepository,
	pr repository.ProductRepository,
) *TransactionController {
	return &TransactionController{
		transactionRepo: tr,
		customerRepo:    cr,
		productRepo:     pr,
	}
}

func (c *TransactionController) CreateTransaction(ctx *gin.Context) {
	var req models.CreateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Get or create customer
	customer, isNew, err := c.customerRepo.GetOrCreateCustomer(req.NamaCustomer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memproses customer: " + err.Error(),
		})
		return
	}

	// Get product
	product, err := c.productRepo.GetProductByID(req.ProductID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Produk tidak ditemukan",
		})
		return
	}

	// Check stok
	if product.Stok < req.Quantity {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Stok produk tidak mencukupi",
		})
		return
	}

	// Create transaction
	trx, err := c.transactionRepo.CreateTransaction(req, customer, product, isNew)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal membuat transaksi: " + err.Error(),
		})
		return
	}

	// Reduce stok
	if err := c.productRepo.ReduceStock(product.ID, req.Quantity); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengurangi stok: " + err.Error(),
		})
		return
	}

	// Add poin: setiap kelipatan Rp 1.000 = 1 poin
	totalHarga := product.Harga * float64(req.Quantity)
	poinDidapat := int(totalHarga / 1000)
	if poinDidapat > 0 {
		if err := c.customerRepo.AddPoin(customer.ID, poinDidapat); err != nil {
			// Log error but don't fail transaction
			_ = err
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success":       true,
		"message":       "Transaksi berhasil",
		"data":          trx,
		"customer_baru": isNew,
		"poin_didapat":  poinDidapat,
	})
}

func (c *TransactionController) GetTransactions(ctx *gin.Context) {
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	if startDate == "" || endDate == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Parameter start_date dan end_date wajib diisi (format: YYYY-MM-DD)",
		})
		return
	}

	transactions, err := c.transactionRepo.GetTransactionsByDateRange(startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data transaksi berhasil diambil",
		"data":    transactions,
	})
}

func (c *TransactionController) GetSummary(ctx *gin.Context) {
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	if startDate == "" || endDate == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Parameter start_date dan end_date wajib diisi (format: YYYY-MM-DD)",
		})
		return
	}

	summary, err := c.transactionRepo.GetDailySummary(startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Summary transaksi berhasil diambil",
		"data":    summary,
	})
}
