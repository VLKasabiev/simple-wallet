package service

import (
	"context"
	"fmt"

	"github.com/VLKasabiev/simple-wallet/internal/model"
)

type TransactionRepository interface {
	GetTransactions(ctx context.Context, walletID int) ([]*model.Transaction, error)
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

func (s *TransactionService) GetTransactions(ctx context.Context, userID, walletID int) ([]*model.Transaction, error) {

    wallet, err := s.walletProvider.GetWalletByID(ctx, walletID)
    if err != nil {
        return nil, fmt.Errorf("failed to check wallet ownership: %w", err)
    }

    if wallet.UserID != userID {
        return nil, model.ErrNotWalletOwner
    }

    transactions, err := s.repo.GetTransactions(ctx, walletID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch transactions from repository: %w", err)
    }

    return transactions, nil
}