package routes

import (
	"keripik-pangsit/controller"
	"keripik-pangsit/repository"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Init repositories
	productRepo := repository.NewProductRepository()
	customerRepo := repository.NewCustomerRepository()
	transactionRepo := repository.NewTransactionRepository()
	redeemRepo := repository.NewRedeemRepository(customerRepo, productRepo)

	// Init controllers
	productCtrl := controller.NewProductController(productRepo)
	customerCtrl := controller.NewCustomerController(customerRepo)
	transactionCtrl := controller.NewTransactionController(transactionRepo, customerRepo, productRepo)
	redeemCtrl := controller.NewRedeemController(redeemRepo)

	api := r.Group("/api/v1")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok", "service": "Keripik Pangsit API"})
		})

		// Products
		products := api.Group("/products")
		{
			products.POST("", productCtrl.CreateProduct)
			products.GET("", productCtrl.GetProductsByDate)
		}

		// Customers
		customers := api.Group("/customers")
		{
			customers.GET("", customerCtrl.GetAllCustomers)
		}

		// Transactions
		transactions := api.Group("/transactions")
		{
			transactions.POST("", transactionCtrl.CreateTransaction)
			transactions.GET("", transactionCtrl.GetTransactions)
			transactions.GET("/summary", transactionCtrl.GetSummary)
		}

		// Redeem Poin
		redeem := api.Group("/redeem")
		{
			redeem.POST("", redeemCtrl.RedeemPoin)
			redeem.GET("", redeemCtrl.GetRedeemLogs)
		}
	}
}
