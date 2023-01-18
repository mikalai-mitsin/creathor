package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestMiddleware(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), RequestIDContextKey, uuid.New().String())
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}
