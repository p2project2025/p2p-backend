package deposit

import (
	"log"
	"net/http"
	"p2p/models"
	"p2p/services/deposit"
	"p2p/utils/response"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DepositHandler struct{}

// Create deposit request
func (h *DepositHandler) CreateDeposit(c *gin.Context) {
	var req models.DepositRequest

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

	s := deposit.DepositServiceInterface(&deposit.DepositService{})
	if err := s.CreateDeposit(req); err != nil {
		response.HandleError(c, err, "Failed to create deposit request", http.StatusInternalServerError)
		return
	}

	response.SuccessResponse(c, "Deposit request created successfully", nil, http.StatusCreated)
}

func (h *DepositHandler) UpdateDepositStatus(c *gin.Context) {
	log.Println("UpdateDepositStatus called")
	var req models.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err, "Invalid request format", http.StatusBadRequest)
		return
	}

	approve := strings.EqualFold(req.Status, "Approved")

	s := deposit.DepositServiceInterface(&deposit.DepositService{})
	if err := s.UpdateDepositStatus(req.ID, approve); err != nil {
		response.HandleError(c, err, "Failed to update deposit status", http.StatusInternalServerError)
		return
	}

	response.SuccessResponse(c, "Deposit status updated successfully", nil, http.StatusOK)
}

func (h *DepositHandler) GetUserDeposits(c *gin.Context) {
	// Get userID from context
	userIDVal, exists := c.Get("userID")
	if !exists {
		response.HandleError(c, nil, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	userID := userIDVal.(string)

	s := deposit.DepositServiceInterface(&deposit.DepositService{})
	data, err := s.GetDepositsByUserID(userID)
	if err != nil {
		response.HandleError(c, err, "Failed to fetch deposits", http.StatusInternalServerError)
		return
	}

	response.SuccessResponse(c, "Deposits fetched successfully", data, http.StatusOK)
}

// List deposits with pagination
func (h *DepositHandler) ListDeposits(c *gin.Context) {

	s := deposit.DepositServiceInterface(&deposit.DepositService{})
	results, err := s.ListDeposits()
	if err != nil {
		response.HandleError(c, err, "Failed to fetch deposits", http.StatusInternalServerError)
		return
	}

	response.SuccessResponse(c, "Deposits fetched successfully", results, http.StatusOK)
}

// Get deposit by ID
func (h *DepositHandler) GetDepositByID(c *gin.Context) {
	id := c.Param("id")
	s := deposit.DepositServiceInterface(&deposit.DepositService{})
	result, err := s.GetDepositByID(id)
	if err != nil {
		response.HandleError(c, err, "Deposit not found", http.StatusNotFound)
		return
	}

	response.SuccessResponse(c, "Deposit fetched successfully", result, http.StatusOK)
}

// Search deposits by username
func (h *DepositHandler) SearchDepositsByUsername(c *gin.Context) {
	username := c.Query("username")

	s := deposit.DepositServiceInterface(&deposit.DepositService{})
	results, err := s.SearchDepositsByUsername(username)
	if err != nil {
		response.HandleError(c, err, "Failed to search deposits", http.StatusInternalServerError)
		return
	}

	response.SuccessResponse(c, "Deposits fetched successfully", results, http.StatusOK)
}
