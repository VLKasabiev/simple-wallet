package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/VLKasabiev/simple-wallet/internal/model"
	"github.com/VLKasabiev/simple-wallet/internal/service"
	"github.com/VLKasabiev/simple-wallet/internal/utils/validator"
	"github.com/labstack/echo/v4"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
}

func NewTransactionHandler(transactionService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}


func (h *TransactionHandler) GetTransactions(c echo.Context) error {
	userID, ok := c.Request().Context().Value("userID").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	walletID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid wallet id"})
    }

	var filter model.TransactionFilter

	if err := c.Bind(&filter); err != nil {
        return c.JSON(http.StatusBadRequest, err.Error())
    }

	if err := c.Validate(&filter); err != nil {
        return c.JSON(http.StatusBadRequest, validator.FormatError(err))
    }

	transactions, err := h.transactionService.GetTransactions(c.Request().Context(), userID, walletID, filter)
	if err != nil {
		slog.Error("failed to get transactions", "wallet_id", walletID, "error", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
    }

    return c.JSON(http.StatusOK, transactions)
}