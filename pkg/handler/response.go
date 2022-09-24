package handler

import (
	"hacknu/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}

// serverErrorMessage creator
func serverErrorMessage(StatusCode int, Message string) *model.DefaultResponse {
	response := &model.DefaultResponse{}
	response.StatusCode = StatusCode
	response.Data = Message
	return response
}

// defaultErrorHandler error only with status code
func defaultErrorHandler(c *gin.Context, Err error) {

	fullError := Err.Error()

	parts := strings.Split(fullError, "|")
	mainMessage := strings.TrimSpace(parts[0])

	switch mainMessage {

	case "bad request":
		c.AbortWithStatusJSON(http.StatusOK, serverErrorMessage(2001, fullError))
	case "not found":
		c.AbortWithStatusJSON(http.StatusOK, serverErrorMessage(2002, fullError))
	case "incorrect boyman":
		c.AbortWithStatusJSON(http.StatusOK, serverErrorMessage(2003, fullError))
	case "incorrect token":
		c.AbortWithStatusJSON(http.StatusOK, serverErrorMessage(2004, fullError))
	case "username is already taken":
		c.AbortWithStatusJSON(http.StatusOK, serverErrorMessage(2005, fullError))

	case "file system":
		c.AbortWithStatusJSON(http.StatusOK, serverErrorMessage(8001, fullError))
	case "server error":
		c.AbortWithStatusJSON(http.StatusOK, serverErrorMessage(8002, fullError))
	default:
		c.AbortWithStatusJSON(http.StatusOK, serverErrorMessage(8000, fullError))
	}
}

// sendGeneral sends general data
func sendGeneral(data interface{}, c *gin.Context) {
	gr := model.SuccessResponse()
	gr.Data = data

	c.JSON(http.StatusOK, gr)
}

// sendSuccess sends response success
func sendSuccess(c *gin.Context) {
	gr := &model.DefaultResponse{StatusCode: 1000}

	c.JSON(http.StatusOK, gr)
}
