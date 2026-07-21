package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/VLKasabiev/simple-wallet/internal/model"
	"github.com/VLKasabiev/simple-wallet/pkg/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
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

func (r *WalletRepository) GetWalletByID(ctx context.Context, walletID int) (*model.Wallet, error) {
	var wallet model.Wallet

	query := `
		SELECT id, user_id, balance, currency, created_at 
		FROM wallets 
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, walletID).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.Balance,  
		&wallet.Currency,
		&wallet.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrWalletNotFound
		}
		
		return nil, fmt.Errorf("failed to execute select query for wallet: %w", err)
	}

	return &wallet, nil
}

func (r *WalletRepository) Deposit(ctx context.Context, walletID int, amount decimal.Decimal) (*model.Wallet, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to start transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	var wallet model.Wallet

	updateQuery := `
		UPDATE wallets 
		SET balance = balance + $1 
		WHERE id = $2 
		RETURNING id, user_id, balance, currency, created_at
	`

	err = tx.QueryRow(ctx, updateQuery, amount, walletID).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.Balance,
		&wallet.Currency,
		&wallet.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrWalletNotFound
		}
		return nil, fmt.Errorf("failed to execute deposit update query: %w", err)
	}

	transactionQuery := `
		INSERT INTO transactions (wallet_id, type, amount, status, created_at)
		VALUES ($1, 'deposit', $2, 'success', NOW())
	`
	_, err = tx.Exec(ctx, transactionQuery, walletID, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to log deposit transaction: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit deposit transaction: %w", err)
	}

	return &wallet, nil
}


func (r *WalletRepository) Withdraw(ctx context.Context, walletID int, amount decimal.Decimal) (*model.Wallet, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to start tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var wallet model.Wallet
	selectQuery := `
		SELECT id, user_id, balance, currency, created_at 
		FROM wallets 
		WHERE id = $1 FOR UPDATE`
	
	err = tx.QueryRow(ctx, selectQuery, walletID).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.Balance,
		&wallet.Currency,
		&wallet.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrWalletNotFound
		}
		return nil, fmt.Errorf("failed to lock wallet: %w", err)
	}

	if wallet.Balance.LessThan(amount) {
		return nil, model.ErrInsufficientBalance
	}

	updateQuery := `UPDATE wallets SET balance = balance - $1 WHERE id = $2`
	_, err = tx.Exec(ctx, updateQuery, amount, walletID)
	if err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	
	transactionQuery := `
		INSERT INTO transactions (wallet_id, type, amount, status, created_at)
		VALUES ($1, 'withdraw', $2, 'success', NOW())`
	_, err = tx.Exec(ctx, transactionQuery, walletID, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to log transaction: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}

	wallet.Balance = wallet.Balance.Sub(amount)
	return &wallet, nil
}


func (r *WalletRepository) Transfer(ctx context.Context, fromWalletID, toWalletID int, amount decimal.Decimal, desc string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("unable to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Это для случая когда 2 кошелька паралельно отправляют друг другу деньги, чтобы избежать deadlock
	firstID, secondID := fromWalletID, toWalletID
	if fromWalletID > toWalletID {
		firstID, secondID = toWalletID, fromWalletID
	}

	var dummyID int

	err = tx.QueryRow(ctx, `SELECT id FROM wallets WHERE id = $1 FOR UPDATE`, firstID).Scan(&dummyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ErrWalletNotFound
		}
		return fmt.Errorf("failed to lock first wallet: %w", err)
	}

	err = tx.QueryRow(ctx, `SELECT id FROM wallets WHERE id = $1 FOR UPDATE`, secondID).Scan(&dummyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ErrWalletNotFound
		}
		return fmt.Errorf("failed to lock second wallet: %w", err)
	}

	withdrawQuery := `UPDATE wallets SET balance = balance - $1 WHERE id = $2 AND balance >= $1`
	cmdTag, err := tx.Exec(ctx, withdrawQuery, amount, fromWalletID)
	if err != nil {
		return fmt.Errorf("failed to withdraw funds: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return model.ErrInsufficientBalance
	}

	depositQuery := `UPDATE wallets SET balance = balance + $1 WHERE id = $2`
	_, err = tx.Exec(ctx, depositQuery, amount, toWalletID)
	if err != nil {
		return fmt.Errorf("failed to deposit funds: %w", err)
	}

	transactionQuery := `
        INSERT INTO transactions (wallet_id, type, amount, status, description, created_at)
        VALUES ($1, 'transfer_out', $2, 'success', $3, NOW())`
	_, err = tx.Exec(ctx, transactionQuery, fromWalletID, amount, fmt.Sprintf("Перевод на кошелек #%d: %s", toWalletID, desc))
	if err != nil {
		return fmt.Errorf("failed to log transfer_out: %w", err)
	}

	_, err = tx.Exec(ctx, transactionQuery, toWalletID, amount, fmt.Sprintf("Перевод от кошелька #%d: %s", fromWalletID, desc))
	if err != nil {
		return fmt.Errorf("failed to log transfer_in: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}