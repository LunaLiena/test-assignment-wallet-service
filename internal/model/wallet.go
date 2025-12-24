package model

import "errors"

type Wallet struct {
	ID      string  `gorm:"primaryKey;type:uuid"`
	Balance float64 `gorm:"type:numeric(15,2);not null;default:0"`
}

func (w *Wallet) Deposit(amount float64) error {
	if amount <= 0 {
		return errors.New("deposit amount must be positive")
	}
	w.Balance += amount
	return nil
}

func (w *Wallet) Withdraw(amount float64) error {
	if amount <= 0 {
		return errors.New("withdrawal amount must be positive")
	}
	if w.Balance < amount {
		return errors.New("insufficient balance")
	}
	w.Balance -= amount
	return nil
}
