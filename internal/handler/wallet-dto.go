package handler

import (
	"github.com/VLKasabiev/simple-wallet/internal/model"
	"github.com/shopspring/decimal"
)	

type CreateWalletRequest struct {
	Currency model.Currency `json:"currency" validate:"required,oneof=RUB USD EUR"`
}

type UpdateBalanceRequest struct {
	Amount decimal.Decimal `json:"amount" validate:"required,decimal_gt_zero"`
}

type TransferRequest struct {
	ToWalletID int `json:"to_wallet_id" validate:"required,gt=0"`
	Amount decimal.Decimal `json:"amount" validate:"required,decimal_gt_zero"`
	Description string `json:"description" validate:"max=255"`
}