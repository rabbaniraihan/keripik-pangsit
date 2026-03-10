package controller

import (
	"keripik-pangsit/models"
	"keripik-pangsit/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RedeemController struct {
	redeemRepo repository.RedeemRepository
}

func NewRedeemController(rr repository.RedeemRepository) *RedeemController {
	return &RedeemController{redeemRepo: rr}
}

func (c *RedeemController) RedeemPoin(ctx *gin.Context) {
	var req models.RedeemPoinRequest
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

	log, err := c.redeemRepo.RedeemPoin(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Penukaran poin berhasil",
		"data":    log,
	})
}

func (c *RedeemController) GetRedeemLogs(ctx *gin.Context) {
	logs, err := c.redeemRepo.GetRedeemLogs()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Log penukaran poin berhasil diambil",
		"data":    logs,
	})
}
