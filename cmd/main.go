package main

import (
	"log"
	"os"
	"wallet-service/internal/handler"
	"wallet-service/internal/repository"
	"wallet-service/internal/service"
	"wallet-service/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	database.ConnectDB()
	repository.InitSchema()

	walletService := service.NewWalletService(database.DB)
	walletHandler := handler.NewWalletHandler(walletService)

	r := gin.Default()

	r.POST("/api/v1/wallet", walletHandler.HandleWalletOperation)
	r.GET("/api/v1/wallets/:id", walletHandler.GetWalletBalance)
	r.POST("/api/v1/wallets/init", walletHandler.InitWallet)
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}

}
