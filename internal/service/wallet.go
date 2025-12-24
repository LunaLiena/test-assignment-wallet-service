package service

import (
	"errors"
	"wallet-service/internal/entities"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletService struct {
	DB *gorm.DB
}

func NewWalletService(db *gorm.DB) *WalletService {
	return &WalletService{DB: db}
}

func (s *WalletService) Deposit(walletId string, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	return s.updateBalance(walletId, amount)
}

func (s *WalletService) WithDraw(walletId string, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	return s.updateBalance(walletId, -amount)
}

func (s *WalletService) GetBalance(walletId string) (float64, error) {

	var wallet entities.Wallet
	err := s.DB.First(&wallet, "id = ?", walletId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("wallet not found")
		}
		return 0, err
	}

	return wallet.Balance, nil
}

func (s *WalletService) CreateWallet(walletId string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var wallet entities.Wallet
		err := tx.First(&wallet, "id = ?", walletId).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				wallet := entities.Wallet{
					ID: walletId,
				}
				return tx.Create(&wallet).Error
			}
			return err
		}
		return nil
	})
}

func (s *WalletService) updateBalance(walletId string, delta float64) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {

		var wallet entities.Wallet
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&wallet, "id = ?", walletId).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Автоматически создаём кошелек при первой операции
				wallet = entities.Wallet{
					ID:      walletId,
					Balance: 0,
				}
				if err := tx.Create(&wallet).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

		newBalance := wallet.Balance + delta
		if newBalance < 0 {
			return errors.New("insufficient balance")
		}

		wallet.Balance = newBalance
		return tx.Save(&wallet).Error
	})
}

func (s *WalletService) ensureWalletExists(tx *gorm.DB, walletId string) error {
	var wallet entities.Wallet
	err := tx.First(&wallet, "id = ?", walletId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("wallet not initialized for user")
		}
		return err
	}
	return nil
}
