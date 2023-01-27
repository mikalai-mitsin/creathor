package rest

import (
	"fmt"
	"net/http"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/log"
	"github.com/gin-gonic/gin"
)

type PlanHandler struct {
	planInterceptor interceptors.PlanInterceptor
	logger          log.Logger
}

func NewPlanHandler(planInterceptor interceptors.PlanInterceptor, logger log.Logger) *PlanHandler {
	return &PlanHandler{planInterceptor: planInterceptor, logger: logger}
}

func (h *PlanHandler) Register(router *gin.RouterGroup) {
	group := router.Group("/plans")
	group.POST("/", h.Create)
	group.GET("/", h.List)
	group.GET("/:id", h.Get)
	group.PATCH("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
}

// Create        godoc
// @Summary      Store a new Plan
// @Description  Takes a Plan JSON and store in DB. Return saved JSON.
// @Tags         Plan
// @Produce      json
// @Param        Plan  body   models.PlanCreate  true  "Plan JSON"
// @Success      201   {object}  models.Plan
// @Router       /plans [post]
func (h *PlanHandler) Create(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	create := &models.PlanCreate{}
	if err := ctx.Bind(create); err != nil {
		return
	}
	plan, err := h.planInterceptor.Create(ctx.Request.Context(), create, requestUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusCreated, plan)
}

// List          godoc
// @Summary      List Plan array
// @Description  Responds with the list of all Plan as JSON.
// @Tags         Plan
// @Produce      json
// @Param        filter  query   models.PlanFilter false "Plan filter"
// @Success      200  {array}  models.Plan
// @Router       /plans [get]
func (h *PlanHandler) List(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	filter := &models.PlanFilter{}
	if err := ctx.Bind(filter); err != nil {
		return
	}
	listPlans, count, err := h.planInterceptor.List(ctx.Request.Context(), filter, requestUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.Header("count", fmt.Sprint(count))
	ctx.JSON(http.StatusOK, listPlans)
}

// Get           godoc
// @Summary      Get single Plan by UUID
// @Description  Returns the Plan whose UUID value matches the UUID.
// @Tags         Plan
// @Produce      json
// @Param        uuid  path      string  true  "search Plan by UUID"
// @Success      200  {object}  models.Plan
// @Router       /plans/{uuid} [get]
func (h *PlanHandler) Get(c *gin.Context) {
	requestUser := c.Request.Context().Value(UserContextKey).(*models.User)
	plan, err := h.planInterceptor.Get(c.Request.Context(), models.UUID(c.Param("id")), requestUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, plan)
}

// Update        godoc
// @Summary      Update Plan by UUID
// @Description  Returns the updated Plan.
// @Tags         Plan
// @Produce      json
// @Param        uuid  path      string  true  "update Plan by UUID"
// @Param        Plan  body   models.PlanUpdate  true  "Plan JSON"
// @Success      201  {object}  models.Plan
// @Router       /plans/{uuid} [PATCH]
func (h *PlanHandler) Update(c *gin.Context) {
	requestUser := c.Request.Context().Value(UserContextKey).(*models.User)
	update := &models.PlanUpdate{}
	if err := c.Bind(update); err != nil {
		return
	}
	update.ID = models.UUID(c.Param("id"))
	plan, err := h.planInterceptor.Update(c.Request.Context(), update, requestUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, plan)
}

// Delete        godoc
// @Summary      Delete single Plan by UUID
// @Description  Delete the Plan whose UUID value matches the UUID.
// @Tags         Plan
// @Param        uuid  path      string  true  "delete Plan by UUID"
// @Success      204
// @Router       /plans/{uuid} [delete]
func (h *PlanHandler) Delete(c *gin.Context) {
	requestUser := c.Request.Context().Value(UserContextKey).(*models.User)
	err := h.planInterceptor.Delete(c.Request.Context(), models.UUID(c.Param("id")), requestUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.Status(http.StatusNoContent)
}
