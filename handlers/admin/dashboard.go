package admin

import (
	"net/http"
	"p2p/services/admin"
	"p2p/utils/response"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct{}

func (h *DashboardHandler) GetCounts(c *gin.Context) {
	s := admin.DashboardServiceInterface(&admin.DashboardService{})
	counts, err := s.GetCounts()
	if err != nil {
		response.HandleError(c, err, "Failed to fetch counts", http.StatusInternalServerError)
		return
	}
	response.SuccessResponse(c, "Counts fetched successfully", counts, http.StatusOK)
}
