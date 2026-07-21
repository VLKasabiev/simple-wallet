package service

import (
	"context"
	"fmt"

	"github.com/VLKasabiev/simple-wallet/internal/model"
)

type TransactionRepository interface {
	GetTransactions(ctx context.Context, walletID int, filter model.TransactionFilter) ([]*model.Transaction, error)
}

type WalletProvider interface {
    GetWalletByID(ctx context.Context, walletID int) (*model.Wallet, error)
}

type TransactionService struct {
	repo TransactionRepository
	walletProvider WalletProvider
}

func NewTransactionService(repo TransactionRepository, walletProvider WalletProvider) *TransactionService {
	return &TransactionService{
		repo: repo,
		walletProvider: walletProvider,
	}
}

func (s *TransactionService) GetTransactions(ctx context.Context, userID, walletID int, filter model.TransactionFilter) ([]*model.Transaction, error) {

    wallet, err := s.walletProvider.GetWalletByID(ctx, walletID)
    if err != nil {
        return nil, fmt.Errorf("failed to check wallet ownership: %w", err)
    }

    if wallet.UserID != userID {
        return nil, model.ErrNotWalletOwner
    }


	if filter.Type != "" {
        if filter.Type != "withdraw" && filter.Type != "deposit" {
            return nil, fmt.Errorf("invalid transaction type: %s", filter.Type)
        }
    }

    if filter.Status != "" {
        if filter.Status != "success" && filter.Status != "failed" {
            return nil, fmt.Errorf("invalid transaction status: %s", filter.Status)
        }
    }

    if filter.Sort == "" {
        filter.Sort = "created_at_desc" 
    } else {
        if filter.Sort != "created_at_desc" && filter.Sort != "created_at_asc" {
            return nil, fmt.Errorf("invalid sort parameter: %s", filter.Sort)
        }
    }

    transactions, err := s.repo.GetTransactions(ctx, walletID, filter)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch transactions from repository: %w", err)
    }

    return transactions, nil
}