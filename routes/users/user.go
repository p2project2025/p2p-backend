package users

import (
	"p2p/handlers/users"
	midleware "p2p/utils/midleWare"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	h := users.UserHandler{}
	d := users.DashboardHandler{}
	userRoutes := r.Group("/users")
	userRoutes.POST("/register", h.RegisterUser)
	userRoutes.POST("/login", h.SignInUser)
	userRoutes.PUT("/:id/block", h.BlockUser)
	userRoutes.PUT("/forgot", h.UpdatePassword)

	authUserRoutes := userRoutes.Group("/auth")
	authUserRoutes.Use(midleware.AuthMiddleware())
	authUserRoutes.GET("/all", h.GetAllUsers)
	authUserRoutes.GET("/dashboard", d.GetUserDashboard)

}
