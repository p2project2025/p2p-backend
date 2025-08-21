package users

import (
	"net/http"
	"p2p/models"
	"p2p/services/users"
	midleware "p2p/utils/midleWare"
	"p2p/utils/response"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct{}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req models.User

	if err := c.BindJSON(&req); err != nil {
		response.HandleError(c, err, "Invalid request format", http.StatusBadRequest)
		return
	}
	s := users.UserServiceInterface(&users.UserService{})
	userId, err := s.RegisterUser(req)
	if err != nil {
		response.HandleError(c, err, "Failed to register user", http.StatusInternalServerError)
		return
	}
	response.SuccessResponse(c, "User registered successfully", gin.H{"user_id": userId}, http.StatusCreated)
}

func (h *UserHandler) SignInUser(c *gin.Context) {
	var req models.Login

	if err := c.BindJSON(&req); err != nil {
		response.HandleError(c, err, "Invalid request format", http.StatusBadRequest)
		return
	}

	s := users.UserServiceInterface(&users.UserService{})
	user, err := s.SignInUser(req)
	if err != nil {
		response.HandleError(c, err, "Failed to sign in", http.StatusForbidden)
		return
	}

	token, err := midleware.GenerateJWT(user.Email, user.ID.Hex(), user.Role)
	if err != nil {
		response.HandleError(c, err, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	c.Header("Authorization", "Bearer "+token)

	c.SetCookie("token", token, 3600, "/", "", false, true) // 1 hour, HttpOnly

	response.SuccessResponse(c, "User signed in successfully", gin.H{
		"user":  user,
		"token": token}, http.StatusOK)
}

func (h *UserHandler) BlockUser(c *gin.Context) {
	// Get user_id from URL param
	userIDParam := c.Param("id")
	userID, err := primitive.ObjectIDFromHex(userIDParam)
	if err != nil {
		response.HandleError(c, err, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Parse request body (expects { "is_blocked": true/false })
	var req struct {
		IsBlocked bool `json:"is_blocked"`
	}
	if err := c.BindJSON(&req); err != nil {
		response.HandleError(c, err, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Call service
	s := users.UserServiceInterface(&users.UserService{})
	err = s.BlockUser(userID, req.IsBlocked)
	if err != nil {
		response.HandleError(c, err, "Failed to update block status", http.StatusInternalServerError)
		return
	}

	status := "unblocked"
	if req.IsBlocked {
		status = "blocked"
	}

	response.SuccessResponse(c, "User "+status+" successfully", gin.H{"user_id": userID.Hex()}, http.StatusOK)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	s := users.UserServiceInterface(&users.UserService{})
	userList, err := s.GetAllUsers()
	if err != nil {
		response.HandleError(c, err, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	response.SuccessResponse(c, "Users fetched successfully", userList, http.StatusOK)
}
