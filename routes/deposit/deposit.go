package deposit

import (
	"p2p/handlers/deposit"
	midleware "p2p/utils/midleWare"

	"github.com/gin-gonic/gin"
)

func DepositRoutes(r *gin.Engine) {
	h := deposit.DepositHandler{}
	depositRoutes := r.Group("/deposits")

	depositRoutes.Use(midleware.AuthMiddleware())

	depositRoutes.POST("/", h.CreateDeposit)
	depositRoutes.GET("/", h.ListDeposits)
	depositRoutes.GET("/:id", h.GetDepositByID)
	depositRoutes.GET("/search", h.SearchDepositsByUsername)
	depositRoutes.GET("/user", h.GetUserDeposits)       // GET /deposits/my
	depositRoutes.PUT("/status", h.UpdateDepositStatus) // approve/reject deposit

}
