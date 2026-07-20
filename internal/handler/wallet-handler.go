package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/VLKasabiev/simple-wallet/internal/model"
	"github.com/VLKasabiev/simple-wallet/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)

type WalletHandler struct {
	walletService *service.WalletService
}


func NewWalletHandler(ws *service.WalletService) *WalletHandler {
	return &WalletHandler{
		walletService: ws,
	}
}


func (h *WalletHandler) Create(c echo.Context) error {
	var req CreateWalletRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	wallet, err := h.walletService.CreateWallet(c.Request().Context(), req.UserID,  req.Currency)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, wallet)
}

func (h *WalletHandler) GetUserWallets(c echo.Context) error {
	userID, ok := c.Request().Context().Value("userID").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	wallets, err := h.walletService.GetWalletsByUserID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch wallets"})
	}
	
	return c.JSON(http.StatusOK, wallets)
}

func (h *WalletHandler) GetBalance(c echo.Context) error {
	walletIDStr := c.Param("id")
	walletID, err := strconv.Atoi(walletIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid wallet id format"})
	}

	userID, ok := c.Request().Context().Value("userID").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	balance, err := h.walletService.GetBalance(c.Request().Context(), walletID, userID)
	if err != nil {

		if errors.Is(err, model.ErrNotWalletOwner) {
			return c.JSON(http.StatusForbidden, echo.Map{
				"error": "Access denied. This wallet belongs to another user.",
			})
		}

		if errors.Is(err, model.ErrWalletNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "wallet not found"})
		}

		slog.Error("failed to get wallet balance", 
        "error", err, 
        "wallet_id", walletID, 
        "user_id", userID,
		)

		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"wallet_id": walletID,
		"balance":   balance,
		"currency":  "RUB",
	})
}

func (h *WalletHandler) Deposit(c echo.Context) error {
	walletIDStr := c.Param("id")
	walletID, err := strconv.Atoi(walletIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid wallet id format"})
	}

	var req UpdateBalanceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	userID, ok := c.Request().Context().Value("userID").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	updatedWallet, err := h.walletService.Deposit(c.Request().Context(), walletID, userID, req.Amount)
	if err != nil {
		if errors.Is(err, model.ErrNotWalletOwner) {
			return c.JSON(http.StatusForbidden, echo.Map{
				"error": "Access denied. You can only deposit money into your own wallet.",
			})
		}

		if errors.Is(err, model.ErrWalletNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "wallet not found"})
		}
		
		slog.Error("failed to deposit to wallet", 
			"error", err, 
			"wallet_id", walletID, 
			"user_id", userID,
		)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"wallet_id": updatedWallet.ID,
		"balance":   updatedWallet.Balance,
		"currency":  updatedWallet.Currency,
	})
}


func (h *WalletHandler) Withdraw(c echo.Context) error {
	walletIDStr := c.Param("id")
	walletID, err := strconv.Atoi(walletIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid wallet id format"})
	}

	var req UpdateBalanceRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	if req.Amount.LessThanOrEqual(decimal.Zero) {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "amount must be greater than zero"})
	}

	userID, ok := c.Request().Context().Value("userID").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	updatedWallet, err := h.walletService.Withdraw(c.Request().Context(), walletID, userID, req.Amount)
	if err != nil {
		if errors.Is(err, model.ErrNotWalletOwner) {
			return c.JSON(http.StatusForbidden, echo.Map{"error": "Access denied. You can only withdraw from your own wallet."})
		}
		if errors.Is(err, model.ErrWalletNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "wallet not found"})
		}
		
		if errors.Is(err, model.ErrInsufficientBalance) {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{"error": "insufficient funds on wallet balance"})
		}
		
		slog.Error("failed to withdraw from wallet", "error", err, "wallet_id", walletID, "user_id", userID)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"wallet_id": updatedWallet.ID,
		"balance":   updatedWallet.Balance,
		"currency":  updatedWallet.Currency,
	})
}