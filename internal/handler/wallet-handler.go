package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/VLKasabiev/simple-wallet/internal/model"
	"github.com/VLKasabiev/simple-wallet/internal/service"
	"github.com/VLKasabiev/simple-wallet/internal/utils/validator"
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
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil || userID <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid user id in url"})
	}

	ctx := c.Request().Context()
	currentUserID, ok := ctx.Value("userID").(int)
	if !ok || currentUserID <= 0 {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	if userID != currentUserID {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "you can not create wallets for other users"})
	}

	var req CreateWalletRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"errors": validator.FormatError(err),
		})
	}

	if !req.Currency.IsValid() {
        return c.JSON(http.StatusBadRequest, echo.Map{
            "error": "unsupported currency, allowed values: RUB, USD, EUR",
        })
    }

	wallet, err := h.walletService.CreateWallet(c.Request().Context(), userID,  req.Currency)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, wallet)
}

func (h *WalletHandler) GetUserWallets(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
    if err != nil || userID <= 0 {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid user id"})
    }

	ctx := c.Request().Context()
	currentUserID, ok := ctx.Value("userID").(int)
	if !ok || currentUserID <= 0 {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	if userID != currentUserID {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "you can't whatch another user's wallets"})
	}

	wallets, err := h.walletService.GetWalletsByUserID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch wallets"})
	}
	
	return c.JSON(http.StatusOK, wallets)
}

func (h *WalletHandler) GetByID(c echo.Context) error {
	walletID, err := strconv.Atoi(c.Param("id"))
	if err != nil || walletID <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid wallet id"})
	}

	ctx := c.Request().Context()
	userID, ok := ctx.Value("userID").(int)
	if !ok || userID <= 0 {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	wallet, err := h.walletService.GetByID(ctx, walletID, userID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrWalletNotFound):
			return c.JSON(http.StatusNotFound, echo.Map{"error": "wallet not found"})

		case errors.Is(err, model.ErrNotWalletOwner):
			return c.JSON(http.StatusForbidden, echo.Map{"error": "you don't have access to this wallet, because it's another user's wallet"})

		default:
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch wallet"})
		}
	}

	return c.JSON(http.StatusOK, wallet)
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

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"errors": validator.FormatError(err),
		})
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

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"errors": validator.FormatError(err),
		})
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


func (h *WalletHandler) Transfer(c echo.Context) error {
	userID, ok := c.Request().Context().Value("userID").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	fromWalletID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid wallet id"})
	}

	var req TransferRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"errors": validator.FormatError(err),
		})
	}

	if req.Amount.LessThanOrEqual(decimal.Zero) {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "amount must be greater than zero"})
	}

	if fromWalletID == req.ToWalletID {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "cannot transfer to the same wallet"})
	}

	ctx := c.Request().Context()
	err = h.walletService.Transfer(ctx, userID, fromWalletID, req.ToWalletID, req.Amount, req.Description)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "transfer successful"})
}