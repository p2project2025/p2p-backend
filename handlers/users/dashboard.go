package users

import (
	"net/http"
	"p2p/services/users"
	"p2p/utils/response"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DashboardHandler struct{}

func (h *DashboardHandler) GetUserDashboard(c *gin.Context) {
	s := users.DashboardServiceInterface(&users.DashboardService{})
	userIDVal, exists := c.Get("userID")
	if !exists {
		response.HandleError(c, nil, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	oid, _ := primitive.ObjectIDFromHex(userIDVal.(string))

	dashInfo, err := s.GetUserDashboard(oid)
	if err != nil {
		response.HandleError(c, err, "Failed to fetch user dashboard", http.StatusInternalServerError)
		return
	}

	response.SuccessResponse(c, "Counts fetched successfully", dashInfo, http.StatusOK)
}
