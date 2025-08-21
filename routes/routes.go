package routes

import (
	healthcheck "p2p/handlers/healthCheck"
	"p2p/routes/admin"
	"p2p/routes/chats"
	"p2p/routes/deposit"
	"p2p/routes/users"
	"p2p/routes/withdrawls"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	healthcheckRoutes(r)
	admin.AdminRoutes(r)
	users.UserRoutes(r)
	deposit.DepositRoutes(r)
	withdrawls.WithdrawlRoutes(r)
	chats.RegisterChatRoutes(r)
}

func healthcheckRoutes(r *gin.Engine) {
	h := healthcheck.HealthCheckHandler{}
	r.GET("/ping", h.Ping)
	r.POST("/echo", h.Echo)
}
