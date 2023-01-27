package rest

import (
	"fmt"
	"net/http"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/log"
	"github.com/gin-gonic/gin"
)

type ArchHandler struct {
	archInterceptor interceptors.ArchInterceptor
	logger          log.Logger
}

func NewArchHandler(archInterceptor interceptors.ArchInterceptor, logger log.Logger) *ArchHandler {
	return &ArchHandler{archInterceptor: archInterceptor, logger: logger}
}

func (h *ArchHandler) Register(router *gin.RouterGroup) {
	group := router.Group("/arches")
	group.POST("/", h.Create)
	group.GET("/", h.List)
	group.GET("/:id", h.Get)
	group.PATCH("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
}

// Create        godoc
// @Summary      Store a new Arch
// @Description  Takes a Arch JSON and store in DB. Return saved JSON.
// @Tags         Arch
// @Produce      json
// @Param        Arch  body   models.ArchCreate  true  "Arch JSON"
// @Success      201   {object}  models.Arch
// @Router       /arches [post]
func (h *ArchHandler) Create(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	create := &models.ArchCreate{}
	if err := ctx.Bind(create); err != nil {
		return
	}
	arch, err := h.archInterceptor.Create(ctx.Request.Context(), create, requestUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusCreated, arch)
}

// List          godoc
// @Summary      List Arch array
// @Description  Responds with the list of all Arch as JSON.
// @Tags         Arch
// @Produce      json
// @Param        filter  query   models.ArchFilter false "Arch filter"
// @Success      200  {array}  models.Arch
// @Router       /arches [get]
func (h *ArchHandler) List(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	filter := &models.ArchFilter{}
	if err := ctx.Bind(filter); err != nil {
		return
	}
	listArches, count, err := h.archInterceptor.List(ctx.Request.Context(), filter, requestUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.Header("count", fmt.Sprint(count))
	ctx.JSON(http.StatusOK, listArches)
}

// Get           godoc
// @Summary      Get single Arch by UUID
// @Description  Returns the Arch whose UUID value matches the UUID.
// @Tags         Arch
// @Produce      json
// @Param        uuid  path      string  true  "search Arch by UUID"
// @Success      200  {object}  models.Arch
// @Router       /arches/{uuid} [get]
func (h *ArchHandler) Get(c *gin.Context) {
	requestUser := c.Request.Context().Value(UserContextKey).(*models.User)
	arch, err := h.archInterceptor.Get(c.Request.Context(), models.UUID(c.Param("id")), requestUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, arch)
}

// Update        godoc
// @Summary      Update Arch by UUID
// @Description  Returns the updated Arch.
// @Tags         Arch
// @Produce      json
// @Param        uuid  path      string  true  "update Arch by UUID"
// @Param        Arch  body   models.ArchUpdate  true  "Arch JSON"
// @Success      201  {object}  models.Arch
// @Router       /arches/{uuid} [PATCH]
func (h *ArchHandler) Update(c *gin.Context) {
	requestUser := c.Request.Context().Value(UserContextKey).(*models.User)
	update := &models.ArchUpdate{}
	if err := c.Bind(update); err != nil {
		return
	}
	update.ID = models.UUID(c.Param("id"))
	arch, err := h.archInterceptor.Update(c.Request.Context(), update, requestUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, arch)
}

// Delete        godoc
// @Summary      Delete single Arch by UUID
// @Description  Delete the Arch whose UUID value matches the UUID.
// @Tags         Arch
// @Param        uuid  path      string  true  "delete Arch by UUID"
// @Success      204
// @Router       /arches/{uuid} [delete]
func (h *ArchHandler) Delete(c *gin.Context) {
	requestUser := c.Request.Context().Value(UserContextKey).(*models.User)
	err := h.archInterceptor.Delete(c.Request.Context(), models.UUID(c.Param("id")), requestUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.Status(http.StatusNoContent)
}
