package withdrawls

import (
	"p2p/handlers/withdrawl"
	midleware "p2p/utils/midleWare"

	"github.com/gin-gonic/gin"
)

func WithdrawlRoutes(r *gin.Engine) {
	h := withdrawl.WithdrawlHandler{}
	withdrawlRoutes := r.Group("/withdrawls")
	withdrawlRoutes.Use(midleware.AuthMiddleware())

	withdrawlRoutes.POST("/", h.CreateWithdrawl)
	withdrawlRoutes.GET("/", h.ListWithdrawls)
	withdrawlRoutes.GET("/:id", h.GetWithdrawlByID)
	withdrawlRoutes.GET("/search", h.SearchWithdrawlsByUsername)
	withdrawlRoutes.GET("/user", h.GetUserWithdrawls)       // GET /withdrawls/my
	withdrawlRoutes.PUT("/status", h.UpdateWithdrawlStatus) // approve/reject deposit

}
