package rest

import (
	"fmt"
	"net/http"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/log"
	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	sessionInterceptor interceptors.SessionInterceptor
	logger             log.Logger
}

func NewSessionHandler(sessionInterceptor interceptors.SessionInterceptor, logger log.Logger) *SessionHandler {
	return &SessionHandler{sessionInterceptor: sessionInterceptor, logger: logger}
}

func (h *SessionHandler) Register(router *gin.Engine) {
	group := router.Group("/sessions")
	group.POST("/", h.Create)
	group.GET("/", h.List)
	group.GET("/:id", h.Get)
	group.PATCH("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
}

func (h *SessionHandler) Create(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	create := &models.SessionCreate{}
	if err := ctx.Bind(create); err != nil {
		return
	}
	marks, err := h.sessionInterceptor.Create(ctx.Request.Context(), create, requestUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, marks)
}

func (h *SessionHandler) List(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	filter := &models.SessionFilter{}
	if err := ctx.Bind(filter); err != nil {
		return
	}
	marks, count, err := h.sessionInterceptor.List(ctx.Request.Context(), filter, requestUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.Header("count", fmt.Sprint(count))
	ctx.JSON(http.StatusOK, marks)
}

func (h *SessionHandler) Get(c *gin.Context) {
	requestUser := c.Request.Context().Value(UserContextKey).(*models.User)
	marks, err := h.sessionInterceptor.Get(c.Request.Context(), c.Param("id"), requestUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, marks)
}

func (h *SessionHandler) Update(c *gin.Context) {
	requestUser := c.Request.Context().Value(UserContextKey).(*models.User)
	update := &models.SessionUpdate{}
	if err := c.Bind(update); err != nil {
		return
	}
	update.ID = c.Param("id")
	marks, err := h.sessionInterceptor.Update(c.Request.Context(), update, requestUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, marks)
}

func (h *SessionHandler) Delete(c *gin.Context) {
	requestUser := c.Request.Context().Value(UserContextKey).(*models.User)
	err := h.sessionInterceptor.Delete(c.Request.Context(), c.Param("id"), requestUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.Status(http.StatusNoContent)
}
