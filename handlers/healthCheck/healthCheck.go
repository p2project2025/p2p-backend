package healthcheck

import (
	"net/http"
	"p2p/utils/response"

	"github.com/gin-gonic/gin"
)

type HealthCheckHandler struct{}

func (h *HealthCheckHandler) Ping(c *gin.Context) {
	response.SuccessResponse(c, "Ping successful", "Pong", http.StatusOK)
}

func (h *HealthCheckHandler) Echo(c *gin.Context) {
	var req struct {
		Message string `json:"message"`
	}

	if err := c.BindJSON(&req); err != nil {
		response.HandleError(c, err, "Invalid request format", http.StatusBadRequest)
		return
	}

	res := "Received: " + req.Message + ". Server successfully responded"
	response.SuccessResponse(c, "Echo response", res, http.StatusOK)
}
