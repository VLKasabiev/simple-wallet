package service

import (
	"context"
	"fmt"
	"github.com/VLKasabiev/simple-wallet/internal/model"
	"github.com/shopspring/decimal"
)


type WalletRepository interface {
	Create(ctx context.Context, wallet *model.Wallet) error
	GetWalletsByUserID(ctx context.Context, userID int) ([]*model.Wallet, error)
}

type WalletService struct {
	repo WalletRepository
}


func NewWalletService(repo WalletRepository) *WalletService {
	return &WalletService{
		repo: repo,
	}
}


func (s *WalletService) CreateWallet(ctx context.Context, userID int, currency string) (*model.Wallet, error) {
	wallet := &model.Wallet{
		UserID: userID,
		Currency: currency,
		Balance: decimal.Zero,
	}

	if err := s.repo.Create(ctx, wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *WalletService) GetWalletsByUserID(ctx context.Context, userID int) ([]*model.Wallet, error) {

	wallets, err := s.repo.GetWalletsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallets from repo: %w", err)
	}

	return wallets, nil
}
