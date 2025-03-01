package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"wgnalvian.com/payment-server/dto"
	"wgnalvian.com/payment-server/service"
)

type TransactionController struct {
	TransactionService *service.TransactionService
}

func (t *TransactionController) TopUp(c *gin.Context) {
	var req dto.TopupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Buat Order ID unik
	orderID := uuid.New().String()

	// Buat transaksi Midtrans
	paymentURL, err := t.TransactionService.CreateMidtransTransaction(orderID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order_id": orderID, "payment_url": paymentURL})
}
