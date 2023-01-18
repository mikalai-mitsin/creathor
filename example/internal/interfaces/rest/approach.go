package rest

import (
	"fmt"
	"net/http"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/log"
	"github.com/gin-gonic/gin"
)

type ApproachHandler struct {
	approachInterceptor interceptors.ApproachInterceptor
	logger              log.Logger
}

func NewApproachHandler(approachInterceptor interceptors.ApproachInterceptor, logger log.Logger) *ApproachHandler {
	return &ApproachHandler{approachInterceptor: approachInterceptor, logger: logger}
}

func (h *ApproachHandler) Register(router *gin.Engine) {
	group := router.Group("/approaches")
	group.POST("/", h.Create)
	group.GET("/", h.List)
	group.GET("/:id", h.Get)
	group.PATCH("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
}

func (h *ApproachHandler) Create(ctx *gin.Context) {
	create := &models.ApproachCreate{}
	if err := ctx.Bind(create); err != nil {
		return
	}
	marks, err := h.approachInterceptor.Create(ctx.Request.Context(), create)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, marks)
}

func (h *ApproachHandler) List(ctx *gin.Context) {
	filter := &models.ApproachFilter{}
	if err := ctx.Bind(filter); err != nil {
		return
	}
	marks, count, err := h.approachInterceptor.List(ctx.Request.Context(), filter)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.Header("count", fmt.Sprint(count))
	ctx.JSON(http.StatusOK, marks)
}

func (h *ApproachHandler) Get(c *gin.Context) {
	marks, err := h.approachInterceptor.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, marks)
}

func (h *ApproachHandler) Update(c *gin.Context) {
	update := &models.ApproachUpdate{}
	if err := c.Bind(update); err != nil {
		return
	}
	update.ID = c.Param("id")
	marks, err := h.approachInterceptor.Update(c.Request.Context(), update)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, marks)
}

func (h *ApproachHandler) Delete(c *gin.Context) {
	err := h.approachInterceptor.Delete(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.Status(http.StatusNoContent)
}
