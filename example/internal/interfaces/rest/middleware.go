package rest

import (
	"context"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthMiddleware struct {
	authService interceptors.AuthInterceptor
}

func NewAuthMiddleware(authService interceptors.AuthInterceptor) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func (m *AuthMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		token := c.GetHeader("Authorization")
		user := models.Guest
		if len(token) > 7 {
			token = token[7:]
			authUser, err := m.authService.Auth(ctx, models.Token(token))
			if err == nil {
				user = authUser
			}
		}
		ctx = context.WithValue(ctx, UserContextKey, user)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func RequestMiddleware(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), RequestIDContextKey, uuid.New().String())
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}
