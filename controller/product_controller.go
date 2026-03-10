package controller

import (
	"keripik-pangsit/models"
	"keripik-pangsit/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	productRepo repository.ProductRepository
}

func NewProductController(pr repository.ProductRepository) *ProductController {
	return &ProductController{productRepo: pr}
}

func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var req models.CreateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Validate ukuran
	validUkuran := map[string]bool{"Small": true, "Medium": true, "Large": true}
	if !validUkuran[req.Ukuran] {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Ukuran harus Small, Medium, atau Large",
		})
		return
	}

	// Validate rasa
	validRasa := map[string]bool{
		"Jagung Bakar": true, "Rumput Laut": true, "Original": true,
		"Jagung Manis": true, "Keju Asin": true, "Keju Manis": true, "Pedas": true,
	}
	if !validRasa[req.Rasa] {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Rasa tidak valid",
		})
		return
	}

	product, err := c.productRepo.CreateProduct(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Produk berhasil ditambahkan",
		"data":    product,
	})
}

func (c *ProductController) GetProductsByDate(ctx *gin.Context) {
	tanggal := ctx.Query("tanggal")
	if tanggal == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Parameter tanggal wajib diisi (format: YYYY-MM-DD)",
		})
		return
	}

	products, err := c.productRepo.GetProductsByDate(tanggal)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data produk berhasil diambil",
		"data":    products,
	})
}
