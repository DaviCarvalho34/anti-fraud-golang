package handlers

import (
	"net/http"
	"time"
	
	"github.com/anti-fraud-golang/internal/models"
	"github.com/anti-fraud-golang/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// FraudHandler handler para endpoints de fraude
type FraudHandler struct {
	fraudService *services.FraudDetectionService
}

// NewFraudHandler cria uma nova instância do handler
func NewFraudHandler(fraudService *services.FraudDetectionService) *FraudHandler {
	return &FraudHandler{
		fraudService: fraudService,
	}
}

// AnalyzeTransactionRequest request para análise de transação
type AnalyzeTransactionRequest struct {
	TransactionID string              `json:"transaction_id,omitempty"`
	UserID        string              `json:"user_id" binding:"required"`
	Amount        float64             `json:"amount" binding:"required,gt=0"`
	Currency      string              `json:"currency" binding:"required"`
	Merchant      string              `json:"merchant" binding:"required"`
	Location      models.Location     `json:"location" binding:"required"`
	DeviceInfo    *models.DeviceInfo  `json:"device_info,omitempty"`
	CardLast4     string              `json:"card_last4,omitempty"`
	CardType      string              `json:"card_type,omitempty"`
	Description   string              `json:"description,omitempty"`
}

// AnalyzeTransaction analisa uma transação
// @Summary Analisa uma transação para detectar fraude
// @Description Recebe os dados de uma transação e retorna a análise de risco
// @Tags fraud
// @Accept json
// @Produce json
// @Param transaction body AnalyzeTransactionRequest true "Dados da transação"
// @Success 200 {object} models.FraudAnalysisResult
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/transaction/analyze [post]
func (h *FraudHandler) AnalyzeTransaction(c *gin.Context) {
	var req AnalyzeTransactionRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}
	
	// Gera ID da transação se não fornecido
	transactionID := req.TransactionID
	if transactionID == "" {
		transactionID = "TXN-" + uuid.New().String()
	}
	
	// Cria objeto de transação
	transaction := &models.Transaction{
		ID:          transactionID,
		UserID:      req.UserID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Merchant:    req.Merchant,
		Location:    req.Location,
		Timestamp:   time.Now(),
		CardLast4:   req.CardLast4,
		CardType:    req.CardType,
		Description: req.Description,
	}
	
	if req.DeviceInfo != nil {
		transaction.DeviceInfo = *req.DeviceInfo
	}
	
	// Analisa a transação
	result, err := h.fraudService.AnalyzeTransaction(transaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Analysis failed",
			Message: err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, result)
}

// GetAnalytics retorna analytics de um usuário
// @Summary Retorna estatísticas de transações de um usuário
// @Description Obtém analytics e histórico de fraude de um usuário
// @Tags analytics
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} services.TransactionAnalytics
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/analytics/{user_id} [get]
func (h *FraudHandler) GetAnalytics(c *gin.Context) {
	userID := c.Param("user_id")
	
	analytics, err := h.fraudService.GetTransactionAnalytics(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get analytics",
			Message: err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, analytics)
}

// HealthCheck verifica o status da API
// @Summary Health check
// @Description Verifica se a API está funcionando
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /api/v1/health [get]
func (h *FraudHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:    "healthy",
		Service:   "anti-fraud-api",
		Timestamp: time.Now(),
		Version:   "1.0.0",
	})
}

// ErrorResponse resposta de erro padrão
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// HealthResponse resposta de health check
type HealthResponse struct {
	Status    string    `json:"status"`
	Service   string    `json:"service"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}
