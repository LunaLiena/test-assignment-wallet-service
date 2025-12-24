package repository

import (
	"errors"
	"wallet-service/internal/entities"
	"wallet-service/pkg/database"

	"gorm.io/gorm"
)

func InitSchema() {
	if err := database.DB.AutoMigrate(&entities.Wallet{}); err != nil {
		panic("Failed to migrate database")
	}
}

func GetWalletForUpdate(id string) (*entities.Wallet, error) {
	var wallet entities.Wallet
	if err := database.DB.Set("gorm:query_option", "FOR UPDATE").First(&wallet, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("wallet not found")
		}
		return nil, err
	}
	return &wallet, nil
}

func GetOrCreateWallet(id string) (*entities.Wallet, error) {
	var wallet entities.Wallet
	if err := database.DB.First(&wallet, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			wallet.ID = id
			wallet.Balance = 0

			if err := database.DB.Create(&wallet).Error; err != nil {
				return nil, err
			}
			return &wallet, nil
		}
		return nil, err
	}
	return &wallet, nil
}

func UpdateWallet(wallet *entities.Wallet) error {
	if err := database.DB.Save(wallet).Error; err != nil {
		return err
	}
	return nil
}
