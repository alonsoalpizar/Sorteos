package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sorteos-platform/backend/internal/usecase/admin/wallet"
	"github.com/sorteos-platform/backend/pkg/logger"
	"gorm.io/gorm"
)

// WalletHandler maneja las peticiones HTTP relacionadas con administración de billeteras
type WalletHandler struct {
	listWalletsUC           *wallet.ListWalletsUseCase
	viewWalletDetailsUC     *wallet.ViewWalletDetailsUseCase
	listWalletTransactionsUC *wallet.ListWalletTransactionsUseCase
	freezeWalletUC          *wallet.FreezeWalletUseCase
	unfreezeWalletUC        *wallet.UnfreezeWalletUseCase
	log                     *logger.Logger
}

// NewWalletHandler crea una nueva instancia del handler
func NewWalletHandler(db *gorm.DB, log *logger.Logger) *WalletHandler {
	return &WalletHandler{
		listWalletsUC:           wallet.NewListWalletsUseCase(db, log),
		viewWalletDetailsUC:     wallet.NewViewWalletDetailsUseCase(db, log),
		listWalletTransactionsUC: wallet.NewListWalletTransactionsUseCase(db, log),
		freezeWalletUC:          wallet.NewFreezeWalletUseCase(db, log),
		unfreezeWalletUC:        wallet.NewUnfreezeWalletUseCase(db, log),
		log:                     log,
	}
}

// List lista billeteras con filtros y paginación
// GET /api/v1/admin/wallets
func (h *WalletHandler) List(c *gin.Context) {
	// Construir input desde query params
	// Frontend envía page 0-indexed, convertimos a 1-indexed
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	page = page + 1 // Convertir de 0-indexed a 1-indexed

	// Aceptar tanto "limit" como "page_size"
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if pageSize == 0 {
		pageSize, _ = strconv.Atoi(c.DefaultQuery("page_size", "20"))
	}

	input := &wallet.ListWalletsInput{
		Page:     page,
		PageSize: pageSize,
		Search:   c.Query("search"),
		Status:   c.Query("status"),
		OrderBy:  c.Query("order_by"),
	}

	// Ejecutar caso de uso
	output, err := h.listWalletsUC.Execute(c.Request.Context(), input)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}

// ViewDetails obtiene detalles de una billetera
// GET /api/v1/admin/wallets/:id
func (h *WalletHandler) ViewDetails(c *gin.Context) {
	walletID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "ID de billetera inválido",
		})
		return
	}

	input := &wallet.ViewWalletDetailsInput{
		WalletID: walletID,
	}

	output, err := h.viewWalletDetailsUC.Execute(c.Request.Context(), input)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}

// ListTransactions lista transacciones de una billetera
// GET /api/v1/admin/wallets/:id/transactions
func (h *WalletHandler) ListTransactions(c *gin.Context) {
	walletID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "ID de billetera inválido",
		})
		return
	}

	// Frontend envía page 0-indexed, convertimos a 1-indexed
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	page = page + 1 // Convertir de 0-indexed a 1-indexed

	// Aceptar tanto "limit" como "page_size"
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if pageSize == 0 {
		pageSize, _ = strconv.Atoi(c.DefaultQuery("page_size", "20"))
	}

	input := &wallet.ListWalletTransactionsInput{
		WalletID: walletID,
		Page:     page,
		PageSize: pageSize,
		Type:     c.Query("type"),
		Status:   c.Query("status"),
	}

	output, err := h.listWalletTransactionsUC.Execute(c.Request.Context(), input)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}

// Freeze congela una billetera
// PUT /api/v1/admin/wallets/:id/freeze
func (h *WalletHandler) Freeze(c *gin.Context) {
	walletID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "ID de billetera inválido",
		})
		return
	}

	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "La razón es requerida",
		})
		return
	}

	input := &wallet.FreezeWalletInput{
		WalletID: walletID,
		AdminID:  adminID,
		Reason:   req.Reason,
	}

	output, err := h.freezeWalletUC.Execute(c.Request.Context(), input)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}

// Unfreeze descongela una billetera
// PUT /api/v1/admin/wallets/:id/unfreeze
func (h *WalletHandler) Unfreeze(c *gin.Context) {
	walletID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "ID de billetera inválido",
		})
		return
	}

	adminID, err := getAdminIDFromContext(c)
	if err != nil {
		handleError(c, err)
		return
	}

	input := &wallet.UnfreezeWalletInput{
		WalletID: walletID,
		AdminID:  adminID,
	}

	output, err := h.unfreezeWalletUC.Execute(c.Request.Context(), input)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    output,
	})
}
