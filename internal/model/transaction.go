package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type TransactionType string
type TransactionStatus string

const (
	TypeDeposit TransactionType = "deposit"
	TypeWithdraw TransactionType = "withdraw"

	StatusSuccess TransactionStatus= "success"
	StatusFailed TransactionStatus = "failed"
)

type Transaction struct {
	ID          int64             `json:"id" db:"id"`
	WalletID    int64             `json:"wallet_id" db:"wallet_id"`
	Type        TransactionType   `json:"type" db:"type"`
	Amount      decimal.Decimal   `json:"amount" db:"amount"`
	Status      TransactionStatus `json:"status" db:"status"`
	Description string            `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time         `json:"created_at" db:"created_at"`
}

type TransactionFilter struct {
	Type   string `query:"type" validate:"omitempty,oneof=deposit withdraw transfer_in transfer_out"`
	Status string `query:"status" validate:"omitempty,oneof=success failed"`
	Sort   string `query:"sort" validate:"omitempty,oneof=created_at_desc created_at_asc"`
}