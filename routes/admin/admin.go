package admin

import (
	"p2p/handlers/admin"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine) {
	h := admin.AdminHandler{}
	d := admin.DashboardHandler{}
	adminRoutes := r.Group("/admin")
	adminRoutes.POST("/register", h.RegisterAdmin)
	adminRoutes.POST("/login", h.SignInAdmin)

	adminRoutes.GET("/dashboard/counts", d.GetCounts)
	// Admin Config
	adminRoutes.POST("/config/wallet", h.UpsertSecureWalletAddress) // update wallet address
	adminRoutes.POST("/config/usdt", h.UpsertUSDTRate)              // update USDT rate
	adminRoutes.POST("/config/qrcode", h.UpsertQRCode)              // update QR code via upload
	adminRoutes.GET("/config", h.FetchAdminConfig)                  // fetch current config

}
