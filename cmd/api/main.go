package main

import (
	"log"
	
	"github.com/anti-fraud-golang/internal/handlers"
	"github.com/anti-fraud-golang/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	// Inicializa stores
	profileStore := services.NewInMemoryProfileStore()
	blacklistStore := services.NewInMemoryBlacklistStore()
	
	// Adiciona dados de exemplo
	profileStore.CreateSampleProfile("USER456")
	blacklistStore.AddSampleBlacklist()
	
	// Inicializa serviço de detecção de fraude
	fraudService := services.NewFraudDetectionService(profileStore, blacklistStore)
	
	// Inicializa handlers
	fraudHandler := handlers.NewFraudHandler(fraudService)
	
	// Configura router
	router := gin.Default()
	
	// Middleware para CORS
	router.Use(corsMiddleware())
	
	// Rotas da API
	api := router.Group("/api/v1")
	{
		// Health check
		api.GET("/health", fraudHandler.HealthCheck)
		
		// Transações
		transactions := api.Group("/transaction")
		{
			transactions.POST("/analyze", fraudHandler.AnalyzeTransaction)
		}
		
		// Analytics
		analytics := api.Group("/analytics")
		{
			analytics.GET("/:user_id", fraudHandler.GetAnalytics)
		}
	}
	
	// Rota raiz
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "Anti-Fraud API",
			"version": "1.0.0",
			"status":  "running",
			"endpoints": []string{
				"GET  /api/v1/health",
				"POST /api/v1/transaction/analyze",
				"GET  /api/v1/analytics/:user_id",
			},
		})
	})
	
	// Inicia servidor
	port := ":8080"
	log.Printf("🚀 Anti-Fraud API iniciada em http://localhost%s", port)
	log.Printf("📚 Documentação disponível em http://localhost%s/", port)
	
	if err := router.Run(port); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}

// corsMiddleware adiciona headers CORS
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}
