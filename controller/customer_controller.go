package controller

import (
	"keripik-pangsit/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerController struct {
	customerRepo repository.CustomerRepository
}

func NewCustomerController(cr repository.CustomerRepository) *CustomerController {
	return &CustomerController{customerRepo: cr}
}

func (c *CustomerController) GetAllCustomers(ctx *gin.Context) {
	customers, err := c.customerRepo.GetAllCustomers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data customer berhasil diambil",
		"data":    customers,
	})
}
