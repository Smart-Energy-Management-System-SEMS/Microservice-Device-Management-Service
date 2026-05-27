package sharedinterfaces

import (
	"errors"
	"net/http"

	shared "device-management-service/shared/domain"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func RespondError(c *gin.Context, err error) {
	var appErr *shared.AppError
	if errors.As(err, &appErr) {
		c.JSON(statusFromCode(appErr.Code), ErrorResponse{Code: string(appErr.Code), Message: appErr.Message})
		return
	}
	c.JSON(http.StatusInternalServerError, ErrorResponse{Code: string(shared.ErrInternal), Message: "unexpected server error"})
}

func RespondValidation(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{Code: string(shared.ErrValidation), Message: message})
}

func statusFromCode(code shared.ErrorCode) int {
	switch code {
	case shared.ErrValidation:
		return http.StatusBadRequest
	case shared.ErrNotFound:
		return http.StatusNotFound
	case shared.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
