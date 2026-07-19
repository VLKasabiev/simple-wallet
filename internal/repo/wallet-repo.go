package repo

import (
	"fmt"
	"context"
	"github.com/VLKasabiev/simple-wallet/internal/model"
	"github.com/VLKasabiev/simple-wallet/pkg/postgres"
	"time"
)

type WalletRepository struct {
	db *postgres.DB
}

func NewWalletRepositoty(db *postgres.DB) *WalletRepository {
	return &WalletRepository{
		db: db,
	}
}

func (r *WalletRepository) Create(ctx context.Context, wallet *model.Wallet) error {
	query := `
		INSERT INTO wallets (user_id, currency, balance, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	now := time.Now()
	wallet.CreatedAt = now
	wallet.UpdatedAt = now

	err := r.db.QueryRow(ctx, query, wallet.UserID, wallet.Currency, wallet.Balance, wallet.CreatedAt, wallet.UpdatedAt).Scan(&wallet.ID)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

func (r *WalletRepository) GetWalletsByUserID(ctx context.Context, userID int) ([]*model.Wallet, error) {
	query := `
        SELECT id, user_id, currency, balance, created_at, updated_at 
        FROM wallets 
        WHERE user_id = $1
    `
    rows, err := r.db.Query(ctx, query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to query wallets: %w", err)
    }
    defer rows.Close()
    var wallets []*model.Wallet

    for rows.Next() {
        var w model.Wallet
        
        err := rows.Scan(&w.ID, &w.UserID, &w.Currency, &w.Balance, &w.CreatedAt, &w.UpdatedAt)
        if err != nil {
            return nil, fmt.Errorf("failed to scan wallet row: %w", err)
        }
        
        wallets = append(wallets, &w)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error during rows iteration: %w", err)
    }

    return wallets, nil
}