package service

import (
	"context"
	"errors"
	"wallet-service/internal/model"
	"wallet-service/internal/repository"
)

type WalletService struct {
	repo repository.WalletRepository
}

func NewWalletService(repo repository.WalletRepository) *WalletService {
	return &WalletService{repo: repo}
}

func (s *WalletService) Deposit(ctx context.Context, walletId string, amount float64) error {

	wallet, err := s.repo.GetByID(ctx, walletId)

	if err != nil {
		if errors.Is(err, repository.ErrWalletNotFound) {
			if err := s.repo.CreateIfNotExists(ctx, walletId); err != nil {
				return err
			}
			wallet = &model.Wallet{ID: walletId, Balance: 0}
		} else {
			return err
		}
	}

	if err := wallet.Deposit(amount); err != nil {
		return err
	}

	return s.repo.Save(ctx, wallet)

}

func (s *WalletService) WithDraw(ctx context.Context, walletID string, amount float64) error {
	wallet, err := s.repo.GetByID(ctx, walletID)
	if err != nil {
		return err
	}

	if err := wallet.Withdraw(amount); err != nil {
		return err
	}
	return s.repo.Save(ctx, wallet)
}

func (s *WalletService) GetBalance(ctx context.Context, walletId string) (float64, error) {

	wallet, err := s.repo.GetByID(ctx, walletId)
	if err != nil {
		return 0, err
	}
	return wallet.Balance, nil
}

func (s *WalletService) CreateWallet(ctx context.Context, walletId string) error {
	return s.repo.CreateIfNotExists(ctx, walletId)
}
