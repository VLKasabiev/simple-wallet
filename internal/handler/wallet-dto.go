package handler

import (
	"github.com/VLKasabiev/simple-wallet/internal/model"
	"github.com/shopspring/decimal"
)	

type CreateWalletRequest struct {
	Currency model.Currency `json:"currency"`
}

type UpdateBalanceRequest struct {
	Amount decimal.Decimal `json:"amount"`
}

type TransferRequest struct {
	ToWalletID  int             `json:"to_wallet_id"`
	Amount      decimal.Decimal `json:"amount"`
	Description string          `json:"description"`
}