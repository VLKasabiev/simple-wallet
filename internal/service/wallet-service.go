package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/VLKasabiev/simple-wallet/internal/model"
	"github.com/shopspring/decimal"
)


type WalletRepository interface {
	Create(ctx context.Context, wallet *model.Wallet) error
	GetWalletsByUserID(ctx context.Context, userID int) ([]*model.Wallet, error)
	GetWalletByID(ctx context.Context, wallerID int) (*model.Wallet, error)
	Deposit(ctx context.Context, walletID int, amount decimal.Decimal) (*model.Wallet, error)
	Withdraw(ctx context.Context, walletID int, amount decimal.Decimal) (*model.Wallet, error)
	Transfer(ctx context.Context, fromWalletID, toWalletID int, amount decimal.Decimal, desc string) error
}

type WalletService struct {
	repo WalletRepository
}


func NewWalletService(repo WalletRepository) *WalletService {
	return &WalletService{
		repo: repo,
	}
}


func (s *WalletService) CreateWallet(ctx context.Context, userID int, currency model.Currency) (*model.Wallet, error) {
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

func (s *WalletService) GetByID(ctx context.Context, walletID int, userID int) (*model.Wallet, error) {
	wallet, err := s.repo.GetWalletByID(ctx, walletID)
	if err != nil {
		if errors.Is(err, model.ErrWalletNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	if wallet.UserID != userID {
		return nil, model.ErrNotWalletOwner
	}

	return wallet, nil
}

func (s *WalletService) GetBalance(ctx context.Context, walletID, userID int) (decimal.Decimal, error) {

	wallet, err := s.GetByID(ctx, walletID, userID)
	if err != nil {
		return decimal.Zero, err
	}

	return wallet.Balance, nil
}


func (s *WalletService) Deposit(ctx context.Context, walletID, userID int, amount decimal.Decimal) (*model.Wallet, error) {
	wallet, err := s.repo.GetWalletByID(ctx, walletID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify wallet before deposit: %w", err)
	}

	if wallet.UserID != userID {
		return nil, model.ErrNotWalletOwner
	}

	updatedWallet, err := s.repo.Deposit(ctx, walletID, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to apply deposit in repository: %w", err)
	}

	return updatedWallet, nil
}


func (s *WalletService) Withdraw(ctx context.Context, walletID, userID int, amount decimal.Decimal) (*model.Wallet, error) {
	wallet, err := s.repo.GetWalletByID(ctx, walletID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify wallet before withdraw: %w", err)
	}

	if wallet.UserID != userID {
		return nil, model.ErrNotWalletOwner
	}

	updatedWallet, err := s.repo.Withdraw(ctx, walletID, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to apply withdraw in repository: %w", err)
	}

	return updatedWallet, nil
}

func (s *WalletService) Transfer(ctx context.Context, userID, fromWalletID, toWalletID int, amount decimal.Decimal, desc string) error {
	fromWallet, err := s.repo.GetWalletByID(ctx, fromWalletID)
	if err != nil {
		return fmt.Errorf("failed to get sender wallet: %w", err)
	}

	if fromWallet.UserID != userID {
		return model.ErrNotWalletOwner
	}

	err = s.repo.Transfer(ctx, fromWalletID, toWalletID, amount, desc)
    if err != nil {
        return err 
    }

	return nil
}
