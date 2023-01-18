package rest

import (
	"fmt"
	"net/http"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/log"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userInterceptor interceptors.UserInterceptor
	logger          log.Logger
}

func NewUserHandler(userInterceptor interceptors.UserInterceptor, logger log.Logger) *UserHandler {
	return &UserHandler{userInterceptor: userInterceptor, logger: logger}
}

func (h *UserHandler) Register(router *gin.Engine) {
	group := router.Group("/users")
	group.POST("/", h.Create)
	group.GET("/", h.List)
	group.GET("/:id", h.Get)
	group.PATCH("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
}

func (h *UserHandler) Create(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	create := &models.UserCreate{}
	if err := ctx.Bind(create); err != nil {
		return
	}
	users, err := h.userInterceptor.Create(ctx.Request.Context(), create, requestUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (h *UserHandler) List(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	filter := &models.UserFilter{}
	if err := ctx.Bind(filter); err != nil {
		return
	}
	users, count, err := h.userInterceptor.List(ctx.Request.Context(), filter, requestUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.Header("count", fmt.Sprint(count))
	ctx.JSON(http.StatusOK, users)
}

func (h *UserHandler) Get(c *gin.Context) {
	requestUser := c.Request.Context().Value(UserContextKey).(*models.User)
	users, err := h.userInterceptor.Get(c.Request.Context(), c.Param("id"), requestUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) Update(c *gin.Context) {
	requestUser := c.Request.Context().Value(UserContextKey).(*models.User)
	update := &models.UserUpdate{}
	if err := c.Bind(update); err != nil {
		return
	}
	update.ID = c.Param("id")
	users, err := h.userInterceptor.Update(c.Request.Context(), update, requestUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) Delete(c *gin.Context) {
	requestUser := c.Request.Context().Value(UserContextKey).(*models.User)
	err := h.userInterceptor.Delete(c.Request.Context(), c.Param("id"), requestUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.Status(http.StatusNoContent)
}
