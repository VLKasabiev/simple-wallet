package repo

import (
	"context"
	"fmt"

	"github.com/VLKasabiev/simple-wallet/internal/model"
	"github.com/VLKasabiev/simple-wallet/pkg/postgres"
)

type TransactionRepository struct {
	db *postgres.DB
}


func NewTransactionRepository(db *postgres.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}


func (r *TransactionRepository) GetTransactions(ctx context.Context, walletID int, filter model.TransactionFilter) ([]*model.Transaction, error) {
	query := `
        SELECT id, wallet_id, type, amount, status, COALESCE(description, ''), created_at
        FROM transactions 
        WHERE wallet_id = $1
    `
	args := []interface{}{walletID}
	counter := 2

	if filter.Type != "" {
		query += fmt.Sprintf(" AND type = $%d", counter)
		args = append(args, filter.Type)
		counter++
	}

	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", counter)
		args = append(args, filter.Status)
		counter++
	}

	switch filter.Sort {
	case "created_at_asc":
		query += " ORDER BY created_at ASC"
	case "created_at_desc":
		query += " ORDER BY created_at DESC"
	default:
		query += " ORDER BY created_at DESC"
	}
    
    rows, err := r.db.Query(ctx, query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to query wallet's transactions: %w", err)
    }
    defer rows.Close()
    var transactions []*model.Transaction

    for rows.Next() {
        var t model.Transaction
        
        err := rows.Scan(&t.ID, &t.WalletID, &t.Type, &t.Amount, &t.Status, &t.Description, &t.CreatedAt)
        if err != nil {
            return nil, fmt.Errorf("failed to scan wallet row: %w", err)
        }
        
        transactions = append(transactions, &t)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error during rows iteration: %w", err)
    }

    return transactions, nil

}