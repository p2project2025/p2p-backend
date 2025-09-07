package chats

import (
	"p2p/handlers/chat"
	midleware "p2p/utils/midleWare"

	"github.com/gin-gonic/gin"
)

func RegisterChatRoutes(r *gin.Engine) {
	h := chat.ChatHandler{}
	chatHandler := r.Group("/chat")
	chatHandler.Use(midleware.AuthMiddleware())

	chatHandler.POST("/", h.CreateChat)
	chatHandler.GET("/", h.FetchChats)
	chatHandler.GET("/users", h.GetUniqueChatUsers)
	chatHandler.PUT("/read", h.UpdateChatsReadStatus)

}
