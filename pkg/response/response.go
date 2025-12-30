package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

func ErrorWithHTTPStatus(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

func BadRequest(c *gin.Context, message string) {
	ErrorWithHTTPStatus(c, http.StatusBadRequest, 400, message)
}

func Unauthorized(c *gin.Context, message string) {
	ErrorWithHTTPStatus(c, http.StatusUnauthorized, 401, message)
}

func Forbidden(c *gin.Context, message string) {
	ErrorWithHTTPStatus(c, http.StatusForbidden, 403, message)
}

func NotFound(c *gin.Context, message string) {
	ErrorWithHTTPStatus(c, http.StatusNotFound, 404, message)
}

func InternalError(c *gin.Context, message string) {
	ErrorWithHTTPStatus(c, http.StatusInternalServerError, 500, message)
}
