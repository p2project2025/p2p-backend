package models

type AdminConfigData struct {
	SecureWalletAddress string `json:"secure_wallet_address" bson:"secure_wallet_address"`
	USDTRate            string `json:"usdt_rate" bson:"usdt_rate"`
	QRCodeURL           string `json:"qr_code_url" bson:"qr_code_url"`
}

type LedgerRes struct {
	TotalDeposits            float64    `json:"total_deposits"`
	TotalWithdrawals         float64    `json:"total_withdrawals"`
	TotalPendingWithdrawals  int64      `json:"total_pending_withdrawals"`
	CurrentTotalBalance      float64    `json:"current_total_balance"`
	PendingWithdrawalsTotal  float64    `json:"pending_withdrawals_total"`
	RejectedWithdrawalsTotal float64    `json:"rejected_withdrawals_total"`
	TodayStats               TodayStats `json:"today_stats"`
}

type TodayStats struct {
	TotalWithdrawals         float64 `json:"total_withdrawals"`
	TotalWithdrawalsPending  float64 `json:"total_withdrawals_pending"`
	TotalWithdrawalsApproved float64 `json:"total_withdrawals_approved"`
	TotalDeposits            float64 `json:"total_deposits"`
	TotalDepositsPending     float64 `json:"total_deposits_pending"`
	TotalDepositsApproved    float64 `json:"total_deposits_approved"`
	NewUsers                 int64   `json:"new_users"`
}