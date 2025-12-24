package dto

type WalletRequest struct {
	ID            string  `json:"walletId" binding:"required,uuid"`
	OperationType string  `json:"operationType" binding:"required,oneof=DEPOSIT WITHDRAW"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
}

type WalletResponse struct {
	WalletID string  `json:"walletId"`
	Balance  float64 `json:"balance"`
}

type WalletInitRequest struct {
	WalletID string `json:"walletId" binding:"required,uuid"`
}
