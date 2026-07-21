package model

import (
	"time"
	"github.com/shopspring/decimal"
)


type Currency string

const (
	CurrencyRUB Currency = "RUB"
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
)

type Wallet struct {
	ID        int             `json:"id" db:"id"`
	UserID    int             `json:"user_id" db:"user_id"`
	Currency  Currency        `json:"currency" db:"currency"`
	Balance   decimal.Decimal `json:"balance" db:"balance"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

func (c Currency) IsValid() bool {
	switch c {
	case CurrencyRUB, CurrencyUSD, CurrencyEUR:
		return true
	default:
		return false
	}
}

