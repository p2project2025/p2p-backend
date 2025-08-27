package chat

import (
	"net/http"
	"p2p/models"
	"p2p/services/chat"
	"p2p/utils"
	"p2p/utils/response"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatHandler struct {
}

// ✅ Create Chat Handler
func (h *ChatHandler) CreateChat(c *gin.Context) {
	// Get sender_id from context
	userIDVal, exists := c.Get("userID")
	if !exists {
		response.HandleError(c, nil, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var senderObjID primitive.ObjectID
	switch v := userIDVal.(type) {
	case primitive.ObjectID:
		senderObjID = v
	case string:
		objID, err := primitive.ObjectIDFromHex(v)
		if err != nil {
			response.HandleError(c, err, "Invalid user id", http.StatusUnauthorized)
			return
		}
		senderObjID = objID
	default:
		response.HandleError(c, nil, "Invalid user id in context", http.StatusUnauthorized)
		return
	}

	// Receiver & message from form-data
	receiverID := c.PostForm("receiver_id")
	message := c.PostForm("message")

	receiverObjID, err := primitive.ObjectIDFromHex(receiverID)
	if err != nil {
		response.HandleError(c, err, "Invalid receiver_id", http.StatusBadRequest)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		response.HandleError(c, err, "Invalid form data", http.StatusBadRequest)
		return
	}

	files := form.File["files"]
	var fileURLs []string

	// Upload files to Cloudinary
	for _, fileHeader := range files {
		url, err := utils.UploadFormFileToCloudinary(c, fileHeader)
		if err != nil {
			response.HandleError(c, err, "Cloudinary upload failed", http.StatusInternalServerError)
			return
		}
		fileURLs = append(fileURLs, url)
	}

	// Build chat record
	chatData := &models.Chat{
		ID:        primitive.NewObjectID(),
		Sender:    senderObjID,
		Receiver:  receiverObjID,
		Message:   message,
		FilesUrl:  fileURLs,
		Timestamp: time.Now(),
	}

	// Save chat
	s := chat.ChatServiceInterface(&chat.ChatService{})
	if err := s.CreateChat(chatData); err != nil {
		response.HandleError(c, err, "Failed to save chat", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, chatData)
}

// ✅ Fetch Chats Handler
func (h *ChatHandler) FetchChats(c *gin.Context) {
	// Sender from context
	userIDVal, exists := c.Get("userID")
	if !exists {
		response.HandleError(c, nil, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var senderObjID primitive.ObjectID
	switch v := userIDVal.(type) {
	case primitive.ObjectID:
		senderObjID = v
	case string:
		objID, err := primitive.ObjectIDFromHex(v)
		if err != nil {
			response.HandleError(c, err, "Invalid user id", http.StatusUnauthorized)
			return
		}
		senderObjID = objID
	default:
		response.HandleError(c, nil, "Invalid user id type", http.StatusUnauthorized)
		return
	}

	// Receiver from query
	receiverID := c.Query("receiver_id")
	if receiverID == "" {
		response.HandleError(c, nil, "receiver_id is required", http.StatusBadRequest)
		return
	}

	receiverObjID, err := primitive.ObjectIDFromHex(receiverID)
	if err != nil {
		response.HandleError(c, err, "Invalid receiver_id", http.StatusBadRequest)
		return
	}

	// Fetch from service
	service := chat.ChatServiceInterface(&chat.ChatService{})
	chats, err := service.GetChatsBetweenUsers(senderObjID, receiverObjID)
	if err != nil {
		response.HandleError(c, err, "Failed to fetch chats", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, chats)
}

func (h *ChatHandler) GetUniqueChatUsers(c *gin.Context) {
	// Get userID from context
	userIDVal, exists := c.Get("userID")
	if !exists {
		response.HandleError(c, nil, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var userObjID primitive.ObjectID
	switch v := userIDVal.(type) {
	case primitive.ObjectID:
		userObjID = v
	case string:
		objID, err := primitive.ObjectIDFromHex(v)
		if err != nil {
			response.HandleError(c, err, "Invalid user id", http.StatusUnauthorized)
			return
		}
		userObjID = objID
	default:
		response.HandleError(c, nil, "Invalid user id type", http.StatusUnauthorized)
		return
	}

	s := chat.ChatServiceInterface(&chat.ChatService{})
	users, err := s.GetUniqueChatUsers(userObjID)
	if err != nil {
		response.HandleError(c, err, "Failed to fetch unique chat users", http.StatusInternalServerError)
		return
	}

	response.SuccessResponse(c, "Unique chat users fetched successfully", users, http.StatusOK)
}
