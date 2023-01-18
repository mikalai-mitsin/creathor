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
	create := &models.UserCreate{}
	if err := ctx.Bind(create); err != nil {
		return
	}
	marks, err := h.userInterceptor.Create(ctx.Request.Context(), create)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, marks)
}

func (h *UserHandler) List(ctx *gin.Context) {
	filter := &models.UserFilter{}
	if err := ctx.Bind(filter); err != nil {
		return
	}
	marks, count, err := h.userInterceptor.List(ctx.Request.Context(), filter)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.Header("count", fmt.Sprint(count))
	ctx.JSON(http.StatusOK, marks)
}

func (h *UserHandler) Get(c *gin.Context) {
	marks, err := h.userInterceptor.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, marks)
}

func (h *UserHandler) Update(c *gin.Context) {
	update := &models.UserUpdate{}
	if err := c.Bind(update); err != nil {
		return
	}
	update.ID = c.Param("id")
	marks, err := h.userInterceptor.Update(c.Request.Context(), update)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, marks)
}

func (h *UserHandler) Delete(c *gin.Context) {
	err := h.userInterceptor.Delete(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.Status(http.StatusNoContent)
}
