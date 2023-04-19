package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AbortWithMessage[T any](c *gin.Context, message T) {
	c.AbortWithStatusJSON(http.StatusBadRequest, message)
}
