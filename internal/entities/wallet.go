package entities

type Wallet struct {
	ID      string  `gorm:"primaryKey;type:uuid"`
	Balance float64 `gorm:"type:numeric(15,2);not null;default:0"`
}

func (*Wallet) TableName() string {
	return "wallets"
}
