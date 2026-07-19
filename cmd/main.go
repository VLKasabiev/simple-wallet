package main

import (
	"context"
	"log/slog"
	"os"
	"github.com/VLKasabiev/simple-wallet/pkg/log"
	"github.com/VLKasabiev/simple-wallet/internal/config"
	"github.com/VLKasabiev/simple-wallet/internal/handler"
	"github.com/VLKasabiev/simple-wallet/internal/service"
	"github.com/VLKasabiev/simple-wallet/pkg/postgres"
	"github.com/VLKasabiev/simple-wallet/internal/repo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	
	ctx := context.Background()

	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	if logFile := log.InitLogger(cfg.IsProd); logFile != nil {
		defer logFile.Close()
	}

	slog.Info("configuration loaded successfully")

	db, err := postgres.Connect(ctx, cfg.Postgres)
	if err != nil {
		slog.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	slog.Info("Postgres successfully connected")

	userRepo := repo.NewUserRepository(db)
	userService := service.NewUserService(userRepo, cfg)
	userHandler := handler.NewUserHandler(userService)

	walletRepo := repo.NewWalletRepositoty(db)
	walletsService := service.NewWalletService(walletRepo)
	walletHandler := handler.NewWalletHandler(walletsService)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	handler.RegisterRoutes(e, userHandler, walletHandler, cfg.JWT)

	slog.Info("starting HTTP server", "port", cfg.GetWebPort())
	if err := e.Start(":" + cfg.GetWebPort()); err != nil {
		slog.Error("server crashed", "error", err)
		os.Exit(1)
	}
}