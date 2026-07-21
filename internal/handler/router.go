package handler

import (
    "github.com/labstack/echo/v4"
    "github.com/VLKasabiev/simple-wallet/internal/middleware"
    "github.com/VLKasabiev/simple-wallet/internal/config"
)


func RegisterRoutes(e *echo.Echo, userH *UserHandler, walletH *WalletHandler, transactionH *TransactionHandler, healthH *HealthHandler, jwtCfg *config.JWTConfig) {
    api := e.Group("")

    api.GET("/health", healthH.CheckHealth)

    api.POST("/users", userH.Create)
    api.GET("/users", userH.List)
    api.GET("/users/:id", userH.GetByID)
    api.POST("/users/login", userH.Login)

    protected := e.Group("")
    protected.Use(middleware.AuthMiddleware(jwtCfg.SecretKey))

    // PROTECTED ENDPOINTS (Wallets + Transactions)
    protected.POST("/users/:id/wallets", walletH.Create)
    protected.GET("/users/:id/wallets", walletH.GetUserWallets)
    protected.GET("/wallets/:id", walletH.GetByID)
    protected.GET("/wallets/:id/balance", walletH.GetBalance)
    protected.POST("/wallets/:id/deposit", walletH.Deposit)
    protected.POST("/wallets/:id/withdraw", walletH.Withdraw)
    protected.POST("/wallets/:id/transfer", walletH.Transfer)
    protected.GET("/wallets/:id/transactions", transactionH.GetTransactions)
    
}