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
// @Failure        400   {object}  errs.Error
// @Failure        401   {object}  errs.Error
// @Failure        403   {object}  errs.Error
// @Failure        404   {object}  errs.Error
// @Failure        405   {object}  errs.Error
// @Failure        500   {object}  errs.Error
// @Failure        503   {object}  errs.Error
// @Router         /auth [post]
func (h *AuthHandler) CreateTokenPair(ctx *gin.Context) {
	create := &models.Login{}
	_ = ctx.Bind(create)
	marks, err := h.authInterceptor.CreateToken(ctx.Request.Context(), create)
	if err != nil {
		decodeError(ctx, err)
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
// @Failure        400   {object}  errs.Error
// @Failure        401   {object}  errs.Error
// @Failure        403   {object}  errs.Error
// @Failure        404   {object}  errs.Error
// @Failure        405   {object}  errs.Error
// @Failure        500   {object}  errs.Error
// @Failure        503   {object}  errs.Error
// @Router         /auth [patch]
func (h *AuthHandler) RefreshTokenPair(ctx *gin.Context) {
	form := &Refresh{}
	_ = ctx.Bind(form)
	marks, err := h.authInterceptor.RefreshToken(ctx.Request.Context(), form.Token)
	if err != nil {
		decodeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, marks)
}
