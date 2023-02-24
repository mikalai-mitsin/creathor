package rest

import (
	"net/http"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/log"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authInterceptor interceptors.AuthInterceptor
	logger          log.Logger
}

func NewAuthHandler(authInterceptor interceptors.AuthInterceptor, logger log.Logger) *AuthHandler {
	return &AuthHandler{authInterceptor: authInterceptor, logger: logger}
}

func (h *AuthHandler) Register(router *gin.RouterGroup) {
	group := router.Group("/auth")
	group.POST("/", h.CreateTokenPair)
	group.PATCH("/", h.RefreshTokenPair)
}

// CreateTokenPair godoc
// @Summary        Create token pair
// @Description    Auth user return access and refresh token.
// @Tags           Auth
// @Produce        json
// @Param          Login  body   models.Login  true  "Login JSON"
// @Success        200   {object}  models.TokenPair
// @Router         /auth [post]
func (h *AuthHandler) CreateTokenPair(ctx *gin.Context) {
	create := &models.Login{}
	if err := ctx.Bind(create); err != nil {
		return
	}
	marks, err := h.authInterceptor.CreateToken(ctx.Request.Context(), create)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, marks)
}

type Refresh struct {
	Token models.Token `json:"token"`
}

// RefreshTokenPair godoc
// @Summary        Refresh token
// @Description    Return new token pair.
// @Tags           Auth
// @Produce        json
// @Param          Refresh  body   Refresh  true  "Refresh token JSON"
// @Success        200   {object}  models.TokenPair
// @Router         /auth [patch]
func (h *AuthHandler) RefreshTokenPair(ctx *gin.Context) {
	form := &Refresh{}
	if err := ctx.Bind(form); err != nil {
		return
	}
	marks, err := h.authInterceptor.RefreshToken(ctx.Request.Context(), form.Token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, marks)
}
