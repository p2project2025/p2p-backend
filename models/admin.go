package models

type AdminConfigData struct {
	SecureWalletAddress string `json:"secure_wallet_address" bson:"secure_wallet_address"`
	USDTRate            string `json:"usdt_rate" bson:"usdt_rate"`
	QRCodeURL           string `json:"qr_code_url" bson:"qr_code_url"`
}
