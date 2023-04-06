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
// @Failure      400   {object}  errs.Error
// @Failure      401   {object}  errs.Error
// @Failure      403   {object}  errs.Error
// @Failure      404   {object}  errs.Error
// @Failure      405   {object}  errs.Error
// @Failure      500   {object}  errs.Error
// @Failure      503   {object}  errs.Error
// @Router       /arches/ [post]
func (h *ArchHandler) Create(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	create := &models.ArchCreate{}
	_ = ctx.Bind(create)
	arch, err := h.archInterceptor.Create(ctx.Request.Context(), create, requestUser)
	if err != nil {
		decodeError(ctx, err)
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
// @Failure      400   {object}  errs.Error
// @Failure      401   {object}  errs.Error
// @Failure      403   {object}  errs.Error
// @Failure      404   {object}  errs.Error
// @Failure      405   {object}  errs.Error
// @Failure      500   {object}  errs.Error
// @Failure      503   {object}  errs.Error
// @Router       /arches [get]
func (h *ArchHandler) List(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	filter := &models.ArchFilter{}
	_ = ctx.Bind(filter)
	listArches, count, err := h.archInterceptor.List(ctx.Request.Context(), filter, requestUser)
	if err != nil {
		decodeError(ctx, err)
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
// @Failure      400   {object}  errs.Error
// @Failure      401   {object}  errs.Error
// @Failure      403   {object}  errs.Error
// @Failure      404   {object}  errs.Error
// @Failure      405   {object}  errs.Error
// @Failure      500   {object}  errs.Error
// @Failure      503   {object}  errs.Error
// @Router       /arches/{uuid} [get]
func (h *ArchHandler) Get(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	arch, err := h.archInterceptor.Get(
		ctx.Request.Context(),
		models.UUID(ctx.Param("id")),
		requestUser,
	)
	if err != nil {
		decodeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, arch)
}

// Update        godoc
// @Summary      Update Arch by UUID
// @Description  Returns the updated Arch.
// @Tags         Arch
// @Produce      json
// @Param        uuid  path      string  true  "update Arch by UUID"
// @Param        Arch  body   models.ArchUpdate  true  "Arch JSON"
// @Success      201  {object}  models.Arch
// @Failure      400   {object}  errs.Error
// @Failure      401   {object}  errs.Error
// @Failure      403   {object}  errs.Error
// @Failure      404   {object}  errs.Error
// @Failure      405   {object}  errs.Error
// @Failure      500   {object}  errs.Error
// @Failure      503   {object}  errs.Error
// @Router       /arches/{uuid} [PATCH]
func (h *ArchHandler) Update(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	update := &models.ArchUpdate{}
	_ = ctx.Bind(update)
	update.ID = models.UUID(ctx.Param("id"))
	arch, err := h.archInterceptor.Update(ctx.Request.Context(), update, requestUser)
	if err != nil {
		decodeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, arch)
}

// Delete        godoc
// @Summary      Delete single Arch by UUID
// @Description  Delete the Arch whose UUID value matches the UUID.
// @Tags         Arch
// @Param        uuid  path      string  true  "delete Arch by UUID"
// @Success      204
// @Failure      400   {object}  errs.Error
// @Failure      401   {object}  errs.Error
// @Failure      403   {object}  errs.Error
// @Failure      404   {object}  errs.Error
// @Failure      405   {object}  errs.Error
// @Failure      500   {object}  errs.Error
// @Failure      503   {object}  errs.Error
// @Router       /arches/{uuid} [delete]
func (h *ArchHandler) Delete(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	err := h.archInterceptor.Delete(
		ctx.Request.Context(),
		models.UUID(ctx.Param("id")),
		requestUser,
	)
	if err != nil {
		decodeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
