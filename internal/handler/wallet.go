package handler

import (
	"net/http"
	"wallet-service/internal/dto"
	"wallet-service/internal/service"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	service *service.WalletService
}

func NewWalletHandler(service *service.WalletService) *WalletHandler {
	return &WalletHandler{
		service: service,
	}
}

func (h *WalletHandler) InitWallet(c *gin.Context) {
	var req dto.WalletInitRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateWallet(req.WalletID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.Status(http.StatusCreated)
}

// POST /api/v1/wallet
func (h *WalletHandler) HandleWalletOperation(c *gin.Context) {
	var req dto.WalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var err error

	switch req.OperationType {
	case "DEPOSIT":
		err = h.service.Deposit(req.ID, req.Amount)
	case "WITHDRAW":
		err = h.service.WithDraw(req.ID, req.Amount)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid operation type",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// GET /api/v1/wallets/:id
func (h *WalletHandler) GetWalletBalance(c *gin.Context) {
	id := c.Param("id")

	balance, err := h.service.GetBalance(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := dto.WalletResponse{
		WalletID: id,
		Balance:  balance,
	}

	c.JSON(http.StatusOK, resp)
}
