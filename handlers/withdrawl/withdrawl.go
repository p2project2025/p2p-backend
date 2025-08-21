package withdrawl

import (
	"net/http"
	"p2p/models"
	"p2p/services/withdrawl"
	"p2p/utils/response"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WithdrawlHandler struct{}

// Create withdrawl request
func (h *WithdrawlHandler) CreateWithdrawl(c *gin.Context) {
	var req models.WithdrawlRequest

	userIDVal, exists := c.Get("userID")
	if !exists {
		response.HandleError(c, nil, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	oid, _ := primitive.ObjectIDFromHex(userIDVal.(string))
	req.UserId = oid
	req.Status = "Pending"

	if err := c.BindJSON(&req); err != nil {
		response.HandleError(c, err, "Invalid request format", http.StatusBadRequest)
		return
	}

	s := withdrawl.WithdrawlServiceInterface(&withdrawl.WithdrawlService{})
	if err := s.CreateWithdrawl(req); err != nil {
		response.HandleError(c, err, "Failed to create withdrawl request", http.StatusInternalServerError)
		return
	}

	response.SuccessResponse(c, "Withdrawl request created successfully", nil, http.StatusCreated)
}

func (h *WithdrawlHandler) UpdateWithdrawlStatus(c *gin.Context) {
	var req models.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err, "Invalid request format", http.StatusBadRequest)
		return
	}

	approve := strings.EqualFold(req.Status, "Approved")

	s := withdrawl.WithdrawlServiceInterface(&withdrawl.WithdrawlService{})
	if err := s.UpdateWithdrawStatus(req.ID, approve); err != nil {
		response.HandleError(c, err, "Failed to update withdrawl status", http.StatusInternalServerError)
		return
	}

	response.SuccessResponse(c, "Withdrawl status updated successfully", nil, http.StatusOK)
}

// List withdrawls with pagination
func (h *WithdrawlHandler) ListWithdrawls(c *gin.Context) {

	s := withdrawl.WithdrawlServiceInterface(&withdrawl.WithdrawlService{})
	results, err := s.ListWithdrawls()
	if err != nil {
		response.HandleError(c, err, "Failed to fetch withdrawls", http.StatusInternalServerError)
		return
	}

	response.SuccessResponse(c, "Withdrawls fetched successfully", results, http.StatusOK)
}

// Get withdrawl by ID
func (h *WithdrawlHandler) GetWithdrawlByID(c *gin.Context) {
	id := c.Param("id")
	s := withdrawl.WithdrawlServiceInterface(&withdrawl.WithdrawlService{})
	result, err := s.GetWithdrawlByID(id)
	if err != nil {
		response.HandleError(c, err, "Withdrawl not found", http.StatusNotFound)
		return
	}

	response.SuccessResponse(c, "Withdrawl fetched successfully", result, http.StatusOK)
}

// Search withdrawls by username
func (h *WithdrawlHandler) SearchWithdrawlsByUsername(c *gin.Context) {
	username := c.Query("username")

	s := withdrawl.WithdrawlServiceInterface(&withdrawl.WithdrawlService{})
	results, err := s.SearchWithdrawlsByUsername(username)
	if err != nil {
		response.HandleError(c, err, "Failed to search withdrawls", http.StatusInternalServerError)
		return
	}

	response.SuccessResponse(c, "Withdrawls fetched successfully", results, http.StatusOK)
}

func (h *WithdrawlHandler) GetUserWithdrawls(c *gin.Context) {
	// Get userID from context
	userIDVal, exists := c.Get("userID")
	if !exists {
		response.HandleError(c, nil, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	userID := userIDVal.(string)

	s := withdrawl.WithdrawlServiceInterface(&withdrawl.WithdrawlService{})
	data, err := s.GetWithdrawlsByUserID(userID)
	if err != nil {
		response.HandleError(c, err, "Failed to fetch withdrawls", http.StatusInternalServerError)
		return
	}

	response.SuccessResponse(c, "Withdrawls fetched successfully", data, http.StatusOK)
}
