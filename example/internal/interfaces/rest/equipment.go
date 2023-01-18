package rest

import (
	"fmt"
	"net/http"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/log"
	"github.com/gin-gonic/gin"
)

type EquipmentHandler struct {
	equipmentInterceptor interceptors.EquipmentInterceptor
	logger               log.Logger
}

func NewEquipmentHandler(equipmentInterceptor interceptors.EquipmentInterceptor, logger log.Logger) *EquipmentHandler {
	return &EquipmentHandler{equipmentInterceptor: equipmentInterceptor, logger: logger}
}

func (h *EquipmentHandler) Register(router *gin.Engine) {
	group := router.Group("/equipment")
	group.POST("/", h.Create)
	group.GET("/", h.List)
	group.GET("/:id", h.Get)
	group.PATCH("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
}

func (h *EquipmentHandler) Create(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	create := &models.EquipmentCreate{}
	if err := ctx.Bind(create); err != nil {
		return
	}
	marks, err := h.equipmentInterceptor.Create(ctx.Request.Context(), create, requestUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, marks)
}

func (h *EquipmentHandler) List(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	filter := &models.EquipmentFilter{}
	if err := ctx.Bind(filter); err != nil {
		return
	}
	marks, count, err := h.equipmentInterceptor.List(ctx.Request.Context(), filter, requestUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.Header("count", fmt.Sprint(count))
	ctx.JSON(http.StatusOK, marks)
}

func (h *EquipmentHandler) Get(c *gin.Context) {
	requestUser := c.Request.Context().Value(UserContextKey).(*models.User)
	marks, err := h.equipmentInterceptor.Get(c.Request.Context(), c.Param("id"), requestUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, marks)
}

func (h *EquipmentHandler) Update(c *gin.Context) {
	requestUser := c.Request.Context().Value(UserContextKey).(*models.User)
	update := &models.EquipmentUpdate{}
	if err := c.Bind(update); err != nil {
		return
	}
	update.ID = c.Param("id")
	marks, err := h.equipmentInterceptor.Update(c.Request.Context(), update, requestUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, marks)
}

func (h *EquipmentHandler) Delete(c *gin.Context) {
	requestUser := c.Request.Context().Value(UserContextKey).(*models.User)
	err := h.equipmentInterceptor.Delete(c.Request.Context(), c.Param("id"), requestUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.Status(http.StatusNoContent)
}
