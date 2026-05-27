package controllers

import (
	"errors"
	"net/http"

	dmerrors "device-management-service/device-management/domain"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func RespondError(c *gin.Context, err error) {
	var appErr *dmerrors.AppError
	if errors.As(err, &appErr) {
		c.JSON(statusFromCode(appErr.Code), ErrorResponse{Code: string(appErr.Code), Message: appErr.Message})
		return
	}
	c.JSON(http.StatusInternalServerError, ErrorResponse{Code: string(dmerrors.ErrInternal), Message: "unexpected server error"})
}

func RespondValidation(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{Code: string(dmerrors.ErrValidation), Message: message})
}

func statusFromCode(code dmerrors.ErrorCode) int {
	switch code {
	case dmerrors.ErrValidation:
		return http.StatusBadRequest
	case dmerrors.ErrNotFound:
		return http.StatusNotFound
	case dmerrors.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
