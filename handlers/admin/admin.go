package admin

import (
	"net/http"
	"p2p/models"
	"p2p/services/admin"
	"p2p/utils"
	midleware "p2p/utils/midleWare"
	"p2p/utils/response"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct{}

func (h *AdminHandler) RegisterAdmin(c *gin.Context) {
	var req models.User

	if err := c.BindJSON(&req); err != nil {
		response.HandleError(c, err, "Invalid request format", http.StatusBadRequest)
		return
	}

	s := admin.AdminServiceInterface(&admin.AdminService{})
	userId, err := s.RegisterAdmin(req)
	if err != nil {
		response.HandleError(c, err, "Failed to register admin", http.StatusInternalServerError)
		return
	}
	response.SuccessResponse(c, "admin registered successfully", gin.H{"admin_id": userId}, http.StatusCreated)
}

func (h *AdminHandler) SignInAdmin(c *gin.Context) {
	var req models.Login

	if err := c.BindJSON(&req); err != nil {
		response.HandleError(c, err, "Invalid request format", http.StatusBadRequest)
		return
	}

	s := admin.AdminServiceInterface(&admin.AdminService{})
	user, err := s.SignInAdmin(req)
	if err != nil {
		response.HandleError(c, err, "Failed to sign in", http.StatusForbidden)
		return
	}

	midleware.GenerateJWT(user.Email, user.ID.Hex(), user.Role)

	response.SuccessResponse(c, "Admin signed in successfully", gin.H{"admin": user}, http.StatusOK)
}

func (h *AdminHandler) FetchAdminConfig(c *gin.Context) {
	s := admin.AdminServiceInterface(&admin.AdminService{})
	config, err := s.FetchAdminConfig()
	if err != nil {
		response.HandleError(c, err, "Failed to fetch admin config", http.StatusInternalServerError)
		return
	}

	response.SuccessResponse(c, "Admin config fetched successfully", gin.H{"config": config}, http.StatusOK)
}

func (h *AdminHandler) UpsertSecureWalletAddress(c *gin.Context) {
	var req struct {
		SecureWalletAddress string `json:"secure_wallet_address"`
	}

	if err := c.BindJSON(&req); err != nil {
		response.HandleError(c, err, "Invalid request format", http.StatusBadRequest)
		return
	}

	s := admin.AdminServiceInterface(&admin.AdminService{})
	id, err := s.UpsertAdminConfig(models.AdminConfigData{SecureWalletAddress: req.SecureWalletAddress})
	if err != nil {
		response.HandleError(c, err, "Failed to update wallet address", http.StatusInternalServerError)
		return
	}

	response.SuccessResponse(c, "Wallet address updated successfully", gin.H{"config_id": id}, http.StatusOK)
}

func (h *AdminHandler) UpsertUSDTRate(c *gin.Context) {
	var req struct {
		USDTRate string `json:"usdt_rate"`
	}

	if err := c.BindJSON(&req); err != nil {
		response.HandleError(c, err, "Invalid request format", http.StatusBadRequest)
		return
	}

	s := admin.AdminServiceInterface(&admin.AdminService{})
	id, err := s.UpsertAdminConfig(models.AdminConfigData{USDTRate: req.USDTRate})
	if err != nil {
		response.HandleError(c, err, "Failed to update USDT rate", http.StatusInternalServerError)
		return
	}

	response.SuccessResponse(c, "USDT rate updated successfully", gin.H{"config_id": id}, http.StatusOK)
}

func (h *AdminHandler) UpsertQRCode(c *gin.Context) {
	// Expect multipart/form-data with "file"
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.HandleError(c, err, "File not provided", http.StatusBadRequest)
		return
	}

	url, err := utils.UploadFormFileToCloudinary(c, fileHeader)
	if err != nil {
		response.HandleError(c, err, "Cloudinary upload failed", http.StatusInternalServerError)
		return
	}

	s := admin.AdminServiceInterface(&admin.AdminService{})
	id, err := s.UpsertAdminConfig(models.AdminConfigData{QRCodeURL: url})
	if err != nil {
		response.HandleError(c, err, "Failed to update QR code", http.StatusInternalServerError)
		return
	}

	response.SuccessResponse(c, "QR code updated successfully", gin.H{
		"config_id": id,
		"qr_code":   url,
	}, http.StatusOK)
}

func (h *AdminHandler) GetLedgerStats(c *gin.Context) {
	s := admin.AdminServiceInterface(&admin.AdminService{})
	stats, err := s.GetLedgerStats()
	if err != nil {
		response.HandleError(c, err, "Failed to fetch ledger stats", http.StatusInternalServerError)
		return
	}

	response.SuccessResponse(c, "Ledger stats fetched successfully", stats, http.StatusOK)
}
