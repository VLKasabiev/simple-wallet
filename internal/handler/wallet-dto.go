package handler

import "github.com/shopspring/decimal"

type CreateWalletRequest struct {
	UserID   int    `json:"user_id"`
	Currency string `json:"currency"`
}

type UpdateBalanceRequest struct {
	Amount decimal.Decimal `json:"amount"`
}

type TransferRequest struct {
	ToWalletID  int             `json:"to_wallet_id"`
	Amount      decimal.Decimal `json:"amount"`
	Description string          `json:"description"`
}