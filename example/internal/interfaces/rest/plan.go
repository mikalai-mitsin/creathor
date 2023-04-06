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
// @Failure      400   {object}  errs.Error
// @Failure      401   {object}  errs.Error
// @Failure      403   {object}  errs.Error
// @Failure      404   {object}  errs.Error
// @Failure      405   {object}  errs.Error
// @Failure      500   {object}  errs.Error
// @Failure      503   {object}  errs.Error
// @Router       /plans/ [post]
func (h *PlanHandler) Create(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	create := &models.PlanCreate{}
	_ = ctx.Bind(create)
	plan, err := h.planInterceptor.Create(ctx.Request.Context(), create, requestUser)
	if err != nil {
		decodeError(ctx, err)
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
// @Failure      400   {object}  errs.Error
// @Failure      401   {object}  errs.Error
// @Failure      403   {object}  errs.Error
// @Failure      404   {object}  errs.Error
// @Failure      405   {object}  errs.Error
// @Failure      500   {object}  errs.Error
// @Failure      503   {object}  errs.Error
// @Router       /plans [get]
func (h *PlanHandler) List(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	filter := &models.PlanFilter{}
	_ = ctx.Bind(filter)
	listPlans, count, err := h.planInterceptor.List(ctx.Request.Context(), filter, requestUser)
	if err != nil {
		decodeError(ctx, err)
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
// @Failure      400   {object}  errs.Error
// @Failure      401   {object}  errs.Error
// @Failure      403   {object}  errs.Error
// @Failure      404   {object}  errs.Error
// @Failure      405   {object}  errs.Error
// @Failure      500   {object}  errs.Error
// @Failure      503   {object}  errs.Error
// @Router       /plans/{uuid} [get]
func (h *PlanHandler) Get(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	plan, err := h.planInterceptor.Get(
		ctx.Request.Context(),
		models.UUID(ctx.Param("id")),
		requestUser,
	)
	if err != nil {
		decodeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, plan)
}

// Update        godoc
// @Summary      Update Plan by UUID
// @Description  Returns the updated Plan.
// @Tags         Plan
// @Produce      json
// @Param        uuid  path      string  true  "update Plan by UUID"
// @Param        Plan  body   models.PlanUpdate  true  "Plan JSON"
// @Success      201  {object}  models.Plan
// @Failure      400   {object}  errs.Error
// @Failure      401   {object}  errs.Error
// @Failure      403   {object}  errs.Error
// @Failure      404   {object}  errs.Error
// @Failure      405   {object}  errs.Error
// @Failure      500   {object}  errs.Error
// @Failure      503   {object}  errs.Error
// @Router       /plans/{uuid} [PATCH]
func (h *PlanHandler) Update(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	update := &models.PlanUpdate{}
	_ = ctx.Bind(update)
	update.ID = models.UUID(ctx.Param("id"))
	plan, err := h.planInterceptor.Update(ctx.Request.Context(), update, requestUser)
	if err != nil {
		decodeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, plan)
}

// Delete        godoc
// @Summary      Delete single Plan by UUID
// @Description  Delete the Plan whose UUID value matches the UUID.
// @Tags         Plan
// @Param        uuid  path      string  true  "delete Plan by UUID"
// @Success      204
// @Failure      400   {object}  errs.Error
// @Failure      401   {object}  errs.Error
// @Failure      403   {object}  errs.Error
// @Failure      404   {object}  errs.Error
// @Failure      405   {object}  errs.Error
// @Failure      500   {object}  errs.Error
// @Failure      503   {object}  errs.Error
// @Router       /plans/{uuid} [delete]
func (h *PlanHandler) Delete(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	err := h.planInterceptor.Delete(
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
