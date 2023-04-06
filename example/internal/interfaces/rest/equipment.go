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

func NewEquipmentHandler(
	equipmentInterceptor interceptors.EquipmentInterceptor,
	logger log.Logger,
) *EquipmentHandler {
	return &EquipmentHandler{equipmentInterceptor: equipmentInterceptor, logger: logger}
}

func (h *EquipmentHandler) Register(router *gin.RouterGroup) {
	group := router.Group("/equipment")
	group.POST("/", h.Create)
	group.GET("/", h.List)
	group.GET("/:id", h.Get)
	group.PATCH("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
}

// Create        godoc
// @Summary      Store a new Equipment
// @Description  Takes a Equipment JSON and store in DB. Return saved JSON.
// @Tags         Equipment
// @Produce      json
// @Param        Equipment  body   models.EquipmentCreate  true  "Equipment JSON"
// @Success      201   {object}  models.Equipment
// @Failure      400   {object}  errs.Error
// @Failure      401   {object}  errs.Error
// @Failure      403   {object}  errs.Error
// @Failure      404   {object}  errs.Error
// @Failure      405   {object}  errs.Error
// @Failure      500   {object}  errs.Error
// @Failure      503   {object}  errs.Error
// @Router       /equipment/ [post]
func (h *EquipmentHandler) Create(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	create := &models.EquipmentCreate{}
	_ = ctx.Bind(create)
	equipment, err := h.equipmentInterceptor.Create(ctx.Request.Context(), create, requestUser)
	if err != nil {
		decodeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, equipment)
}

// List          godoc
// @Summary      List Equipment array
// @Description  Responds with the list of all Equipment as JSON.
// @Tags         Equipment
// @Produce      json
// @Param        filter  query   models.EquipmentFilter false "Equipment filter"
// @Success      200  {array}  models.Equipment
// @Failure      400   {object}  errs.Error
// @Failure      401   {object}  errs.Error
// @Failure      403   {object}  errs.Error
// @Failure      404   {object}  errs.Error
// @Failure      405   {object}  errs.Error
// @Failure      500   {object}  errs.Error
// @Failure      503   {object}  errs.Error
// @Router       /equipment [get]
func (h *EquipmentHandler) List(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	filter := &models.EquipmentFilter{}
	_ = ctx.Bind(filter)
	listEquipment, count, err := h.equipmentInterceptor.List(
		ctx.Request.Context(),
		filter,
		requestUser,
	)
	if err != nil {
		decodeError(ctx, err)
		return
	}
	ctx.Header("count", fmt.Sprint(count))
	ctx.JSON(http.StatusOK, listEquipment)
}

// Get           godoc
// @Summary      Get single Equipment by UUID
// @Description  Returns the Equipment whose UUID value matches the UUID.
// @Tags         Equipment
// @Produce      json
// @Param        uuid  path      string  true  "search Equipment by UUID"
// @Success      200  {object}  models.Equipment
// @Failure      400   {object}  errs.Error
// @Failure      401   {object}  errs.Error
// @Failure      403   {object}  errs.Error
// @Failure      404   {object}  errs.Error
// @Failure      405   {object}  errs.Error
// @Failure      500   {object}  errs.Error
// @Failure      503   {object}  errs.Error
// @Router       /equipment/{uuid} [get]
func (h *EquipmentHandler) Get(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	equipment, err := h.equipmentInterceptor.Get(
		ctx.Request.Context(),
		models.UUID(ctx.Param("id")),
		requestUser,
	)
	if err != nil {
		decodeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, equipment)
}

// Update        godoc
// @Summary      Update Equipment by UUID
// @Description  Returns the updated Equipment.
// @Tags         Equipment
// @Produce      json
// @Param        uuid  path      string  true  "update Equipment by UUID"
// @Param        Equipment  body   models.EquipmentUpdate  true  "Equipment JSON"
// @Success      201  {object}  models.Equipment
// @Failure      400   {object}  errs.Error
// @Failure      401   {object}  errs.Error
// @Failure      403   {object}  errs.Error
// @Failure      404   {object}  errs.Error
// @Failure      405   {object}  errs.Error
// @Failure      500   {object}  errs.Error
// @Failure      503   {object}  errs.Error
// @Router       /equipment/{uuid} [PATCH]
func (h *EquipmentHandler) Update(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	update := &models.EquipmentUpdate{}
	_ = ctx.Bind(update)
	update.ID = models.UUID(ctx.Param("id"))
	equipment, err := h.equipmentInterceptor.Update(ctx.Request.Context(), update, requestUser)
	if err != nil {
		decodeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, equipment)
}

// Delete        godoc
// @Summary      Delete single Equipment by UUID
// @Description  Delete the Equipment whose UUID value matches the UUID.
// @Tags         Equipment
// @Param        uuid  path      string  true  "delete Equipment by UUID"
// @Success      204
// @Failure      400   {object}  errs.Error
// @Failure      401   {object}  errs.Error
// @Failure      403   {object}  errs.Error
// @Failure      404   {object}  errs.Error
// @Failure      405   {object}  errs.Error
// @Failure      500   {object}  errs.Error
// @Failure      503   {object}  errs.Error
// @Router       /equipment/{uuid} [delete]
func (h *EquipmentHandler) Delete(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	err := h.equipmentInterceptor.Delete(
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
