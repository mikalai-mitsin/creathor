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

func NewSessionHandler(
	sessionInterceptor interceptors.SessionInterceptor,
	logger log.Logger,
) *SessionHandler {
	return &SessionHandler{sessionInterceptor: sessionInterceptor, logger: logger}
}

func (h *SessionHandler) Register(router *gin.RouterGroup) {
	group := router.Group("/sessions")
	group.POST("/", h.Create)
	group.GET("/", h.List)
	group.GET("/:id", h.Get)
	group.PATCH("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
}

// Create        godoc
// @Summary      Store a new Session
// @Description  Takes a Session JSON and store in DB. Return saved JSON.
// @Tags         Session
// @Produce      json
// @Param        Session  body   models.SessionCreate  true  "Session JSON"
// @Success      201   {object}  models.Session
// @Router       /sessions [post]
func (h *SessionHandler) Create(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	create := &models.SessionCreate{}
	if err := ctx.Bind(create); err != nil {
		return
	}
	session, err := h.sessionInterceptor.Create(ctx.Request.Context(), create, requestUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusCreated, session)
}

// List          godoc
// @Summary      List Session array
// @Description  Responds with the list of all Session as JSON.
// @Tags         Session
// @Produce      json
// @Param        filter  query   models.SessionFilter false "Session filter"
// @Success      200  {array}  models.Session
// @Router       /sessions [get]
func (h *SessionHandler) List(ctx *gin.Context) {
	requestUser := ctx.Request.Context().Value(UserContextKey).(*models.User)
	filter := &models.SessionFilter{}
	if err := ctx.Bind(filter); err != nil {
		return
	}
	listSessions, count, err := h.sessionInterceptor.List(
		ctx.Request.Context(),
		filter,
		requestUser,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.Header("count", fmt.Sprint(count))
	ctx.JSON(http.StatusOK, listSessions)
}

// Get           godoc
// @Summary      Get single Session by UUID
// @Description  Returns the Session whose UUID value matches the UUID.
// @Tags         Session
// @Produce      json
// @Param        uuid  path      string  true  "search Session by UUID"
// @Success      200  {object}  models.Session
// @Router       /sessions/{uuid} [get]
func (h *SessionHandler) Get(c *gin.Context) {
	requestUser := c.Request.Context().Value(UserContextKey).(*models.User)
	session, err := h.sessionInterceptor.Get(
		c.Request.Context(),
		models.UUID(c.Param("id")),
		requestUser,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, session)
}

// Update        godoc
// @Summary      Update Session by UUID
// @Description  Returns the updated Session.
// @Tags         Session
// @Produce      json
// @Param        uuid  path      string  true  "update Session by UUID"
// @Param        Session  body   models.SessionUpdate  true  "Session JSON"
// @Success      201  {object}  models.Session
// @Router       /sessions/{uuid} [PATCH]
func (h *SessionHandler) Update(c *gin.Context) {
	requestUser := c.Request.Context().Value(UserContextKey).(*models.User)
	update := &models.SessionUpdate{}
	if err := c.Bind(update); err != nil {
		return
	}
	update.ID = models.UUID(c.Param("id"))
	session, err := h.sessionInterceptor.Update(c.Request.Context(), update, requestUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, session)
}

// Delete        godoc
// @Summary      Delete single Session by UUID
// @Description  Delete the Session whose UUID value matches the UUID.
// @Tags         Session
// @Param        uuid  path      string  true  "delete Session by UUID"
// @Success      204
// @Router       /sessions/{uuid} [delete]
func (h *SessionHandler) Delete(c *gin.Context) {
	requestUser := c.Request.Context().Value(UserContextKey).(*models.User)
	err := h.sessionInterceptor.Delete(c.Request.Context(), models.UUID(c.Param("id")), requestUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.Status(http.StatusNoContent)
}
