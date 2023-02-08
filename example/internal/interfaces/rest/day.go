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
// @Router       /days [post]
func (h *DayHandler) Create(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	create := &models.DayCreate{}
	if err := ctx.Bind(create); err != nil {
		return
	}
	day, err := h.dayInterceptor.Create(ctx.Request.Context(), create, requestUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
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
// @Router       /days [get]
func (h *DayHandler) List(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	filter := &models.DayFilter{}
	if err := ctx.Bind(filter); err != nil {
		return
	}
	listDays, count, err := h.dayInterceptor.List(ctx.Request.Context(), filter, requestUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
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
// @Router       /days/{uuid} [get]
func (h *DayHandler) Get(c *gin.Context) {
	requestUser := c.Request.Context().Value(UserContextKey).(*models.User)
	day, err := h.dayInterceptor.Get(c.Request.Context(), models.UUID(c.Param("id")), requestUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, day)
}

// Update        godoc
// @Summary      Update Day by UUID
// @Description  Returns the updated Day.
// @Tags         Day
// @Produce      json
// @Param        uuid  path      string  true  "update Day by UUID"
// @Param        Day  body   models.DayUpdate  true  "Day JSON"
// @Success      201  {object}  models.Day
// @Router       /days/{uuid} [PATCH]
func (h *DayHandler) Update(c *gin.Context) {
	requestUser := c.Request.Context().Value(UserContextKey).(*models.User)
	update := &models.DayUpdate{}
	if err := c.Bind(update); err != nil {
		return
	}
	update.ID = models.UUID(c.Param("id"))
	day, err := h.dayInterceptor.Update(c.Request.Context(), update, requestUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, day)
}

// Delete        godoc
// @Summary      Delete single Day by UUID
// @Description  Delete the Day whose UUID value matches the UUID.
// @Tags         Day
// @Param        uuid  path      string  true  "delete Day by UUID"
// @Success      204
// @Router       /days/{uuid} [delete]
func (h *DayHandler) Delete(c *gin.Context) {
	requestUser := c.Request.Context().Value(UserContextKey).(*models.User)
	err := h.dayInterceptor.Delete(c.Request.Context(), models.UUID(c.Param("id")), requestUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.Status(http.StatusNoContent)
}
