package response

import (
	"log"
	"p2p/models"

	"github.com/gin-gonic/gin"
)

// HandleError logs the error and sends an error response using Gin.
func HandleError(c *gin.Context, err error, description string, statusCode int) {
	log.Println(description+" : ", err)
	jsonResponse(c, statusCode, "Failed", description, err.Error())
}

// SuccessResponse sends a success response using Gin.
func SuccessResponse(c *gin.Context, description string, responseData interface{}, statusCode int) {
	jsonResponse(c, statusCode, "Success", description, responseData)
}

// jsonResponse is a helper to send JSON response in a common format.
func jsonResponse(c *gin.Context, status int, responseStatus, responseDescription string, responseData interface{}) {
	resp := models.GeneralResponse{
		ResponseStatus:      responseStatus,
		ResponseDescription: responseDescription,
		ResponseData:        responseData,
		StatusCode:          status,
	}

	c.JSON(status, resp)
}
