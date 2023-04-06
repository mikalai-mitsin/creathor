package rest

import (
	"fmt"
	"net/http"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/log"
	"github.com/gin-gonic/gin"
)

type DayHandler struct {
	dayInterceptor interceptors.DayInterceptor
	logger         log.Logger
}

func NewDayHandler(dayInterceptor interceptors.DayInterceptor, logger log.Logger) *DayHandler {
	return &DayHandler{dayInterceptor: dayInterceptor, logger: logger}
}

func (h *DayHandler) Register(router *gin.RouterGroup) {
	group := router.Group("/days")
	group.POST("/", h.Create)
	group.GET("/", h.List)
	group.GET("/:id", h.Get)
	group.PATCH("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
}

// Create        godoc
// @Summary      Store a new Day
// @Description  Takes a Day JSON and store in DB. Return saved JSON.
// @Tags         Day
// @Produce      json
// @Param        Day  body   models.DayCreate  true  "Day JSON"
// @Success      201   {object}  models.Day
// @Failure      400   {object}  errs.Error
// @Failure      401   {object}  errs.Error
// @Failure      403   {object}  errs.Error
// @Failure      404   {object}  errs.Error
// @Failure      405   {object}  errs.Error
// @Failure      500   {object}  errs.Error
// @Failure      503   {object}  errs.Error
// @Router       /days/ [post]
func (h *DayHandler) Create(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	create := &models.DayCreate{}
	_ = ctx.Bind(create)
	day, err := h.dayInterceptor.Create(ctx.Request.Context(), create, requestUser)
	if err != nil {
		decodeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, day)
}

// List          godoc
// @Summary      List Day array
// @Description  Responds with the list of all Day as JSON.
// @Tags         Day
// @Produce      json
// @Param        filter  query   models.DayFilter false "Day filter"
// @Success      200  {array}  models.Day
// @Failure      400   {object}  errs.Error
// @Failure      401   {object}  errs.Error
// @Failure      403   {object}  errs.Error
// @Failure      404   {object}  errs.Error
// @Failure      405   {object}  errs.Error
// @Failure      500   {object}  errs.Error
// @Failure      503   {object}  errs.Error
// @Router       /days [get]
func (h *DayHandler) List(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	filter := &models.DayFilter{}
	_ = ctx.Bind(filter)
	listDays, count, err := h.dayInterceptor.List(ctx.Request.Context(), filter, requestUser)
	if err != nil {
		decodeError(ctx, err)
		return
	}
	ctx.Header("count", fmt.Sprint(count))
	ctx.JSON(http.StatusOK, listDays)
}

// Get           godoc
// @Summary      Get single Day by UUID
// @Description  Returns the Day whose UUID value matches the UUID.
// @Tags         Day
// @Produce      json
// @Param        uuid  path      string  true  "search Day by UUID"
// @Success      200  {object}  models.Day
// @Failure      400   {object}  errs.Error
// @Failure      401   {object}  errs.Error
// @Failure      403   {object}  errs.Error
// @Failure      404   {object}  errs.Error
// @Failure      405   {object}  errs.Error
// @Failure      500   {object}  errs.Error
// @Failure      503   {object}  errs.Error
// @Router       /days/{uuid} [get]
func (h *DayHandler) Get(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	day, err := h.dayInterceptor.Get(
		ctx.Request.Context(),
		models.UUID(ctx.Param("id")),
		requestUser,
	)
	if err != nil {
		decodeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, day)
}

// Update        godoc
// @Summary      Update Day by UUID
// @Description  Returns the updated Day.
// @Tags         Day
// @Produce      json
// @Param        uuid  path      string  true  "update Day by UUID"
// @Param        Day  body   models.DayUpdate  true  "Day JSON"
// @Success      201  {object}  models.Day
// @Failure      400   {object}  errs.Error
// @Failure      401   {object}  errs.Error
// @Failure      403   {object}  errs.Error
// @Failure      404   {object}  errs.Error
// @Failure      405   {object}  errs.Error
// @Failure      500   {object}  errs.Error
// @Failure      503   {object}  errs.Error
// @Router       /days/{uuid} [PATCH]
func (h *DayHandler) Update(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	update := &models.DayUpdate{}
	_ = ctx.Bind(update)
	update.ID = models.UUID(ctx.Param("id"))
	day, err := h.dayInterceptor.Update(ctx.Request.Context(), update, requestUser)
	if err != nil {
		decodeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, day)
}

// Delete        godoc
// @Summary      Delete single Day by UUID
// @Description  Delete the Day whose UUID value matches the UUID.
// @Tags         Day
// @Param        uuid  path      string  true  "delete Day by UUID"
// @Success      204
// @Failure      400   {object}  errs.Error
// @Failure      401   {object}  errs.Error
// @Failure      403   {object}  errs.Error
// @Failure      404   {object}  errs.Error
// @Failure      405   {object}  errs.Error
// @Failure      500   {object}  errs.Error
// @Failure      503   {object}  errs.Error
// @Router       /days/{uuid} [delete]
func (h *DayHandler) Delete(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	err := h.dayInterceptor.Delete(ctx.Request.Context(), models.UUID(ctx.Param("id")), requestUser)
	if err != nil {
		decodeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
