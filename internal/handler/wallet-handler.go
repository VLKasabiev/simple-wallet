package handler

import (
	"net/http"

	"github.com/VLKasabiev/simple-wallet/internal/service"
	"github.com/labstack/echo/v4"
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