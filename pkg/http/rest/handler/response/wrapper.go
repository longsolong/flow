package response

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"net/http"
)

// APIException ...
type APIException struct {
	Ctx struct {
		RequestID string `json:"requestID"`
	} `json:"ctx"`
	Status      int    `json:"-"`
	Code        int    `json:"code"`
	Message     string `json:"message"`
	UserMessage string `json:"userMessage"`
	Data        gin.H  `json:"data"`
}

// Error ...
func (e *APIException) Error() string {
	return e.Message
}

func newAPIException(status, code int, message, userMessage string, data gin.H) *APIException {
	return &APIException{
		Status:      status,
		Code:        code,
		Message:     message,
		UserMessage: userMessage,
		Data:        data,
	}
}

// ServerError ...
func ServerError(message string) *APIException {
	return newAPIException(http.StatusInternalServerError, -1, message, http.StatusText(http.StatusInternalServerError), gin.H{})
}

// UnknownError ...
func UnknownError() *APIException {
	return newAPIException(http.StatusInternalServerError, -1, "unknown error", http.StatusText(http.StatusInternalServerError), gin.H{})
}

// NotFound ...
func NotFound(message string) *APIException {
	return newAPIException(http.StatusNotFound, -1, message, http.StatusText(http.StatusNotFound), gin.H{})
}

// ParameterError ...
func ParameterError(message string) *APIException {
	return newAPIException(http.StatusBadRequest, -1, message, http.StatusText(http.StatusBadRequest), gin.H{})
}

// BadRequestError ...
func BadRequestError(message string) *APIException {
	return newAPIException(http.StatusBadRequest, -1, message, http.StatusText(http.StatusBadRequest), gin.H{})
}

// HandlerFunc ...
type HandlerFunc func(c *gin.Context) (gin.H, error)

// Wrapper ...
func Wrapper(handler HandlerFunc) func(c *gin.Context) {
	return func(c *gin.Context) {
		data, err := handler(c)
		if err != nil {
			var apiException *APIException
			if h, ok := err.(*APIException); ok {
				apiException = h
			} else if e, ok := err.(error); ok {
				apiException = ServerError(e.Error())
			} else {
				apiException = UnknownError()
			}
			apiException.Ctx.RequestID = requestid.Get(c)
			c.JSON(apiException.Status, apiException)
		} else {
			successResponse := newAPIException(http.StatusOK, 0, "success", http.StatusText(http.StatusOK), data)
			successResponse.Ctx.RequestID = requestid.Get(c)
			c.JSON(http.StatusOK, successResponse)
		}
	}
}
