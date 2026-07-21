package utils

import (
	"fmt"

	"github.com/VLKasabiev/simple-wallet/internal/model"
	"github.com/shopspring/decimal"
)

var exchangeRatesToRUB = map[model.Currency]decimal.Decimal{
	model.CurrencyRUB: decimal.NewFromInt(1),
	model.CurrencyUSD: decimal.NewFromInt(90),
	model.CurrencyEUR: decimal.NewFromInt(100),
}

func Convert(amount decimal.Decimal, from, to model.Currency) (decimal.Decimal, error) {
	if from == to {
		return amount, nil
	}

	rateFrom, okFrom := exchangeRatesToRUB[from]
	rateTo, okTo := exchangeRatesToRUB[to]

	if !okFrom || !okTo {
		return decimal.Zero, fmt.Errorf("unsupported currency conversion: %s to %s", from, to)
	}

	amountInRUB := amount.Mul(rateFrom)

	convertedAmount := amountInRUB.DivRound(rateTo, 2)

	return convertedAmount, nil
}