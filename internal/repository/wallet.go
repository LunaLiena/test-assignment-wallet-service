package repository

import (
	"context"
	"errors"
	"wallet-service/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository interface {
	GetByID(ctx context.Context, id string) (*model.Wallet, error)
	Save(ctx context.Context, wallet *model.Wallet) error
	CreateIfNotExists(ctx context.Context, id string) error
	AutoMigrate(ctx context.Context) error
}

type GormWalletRepository struct {
	db *gorm.DB
}

func NewGormWalletRepository(db *gorm.DB) *GormWalletRepository {
	return &GormWalletRepository{
		db: db,
	}
}

func (r *GormWalletRepository) GetByID(ctx context.Context, id string) (*model.Wallet, error) {
	var wallet model.Wallet
	if err := r.db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(&wallet, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("wallet not found")
		}
		return nil, err
	}
	return &wallet, nil
}

func (r *GormWalletRepository) Save(ctx context.Context, wallet *model.Wallet) error {
	return r.db.WithContext(ctx).Save(wallet).Error
}

func (r *GormWalletRepository) CreateIfNotExists(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var w model.Wallet
		err := tx.First(&w, "id = ?", id).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				w = model.Wallet{ID: id, Balance: 0}
				return tx.Create(&w).Error
			}
			return err
		}
		return nil
	})
}

func (r *GormWalletRepository) AutoMigrate(ctx context.Context) error {
	return r.db.WithContext(ctx).AutoMigrate(&model.Wallet{})
}
