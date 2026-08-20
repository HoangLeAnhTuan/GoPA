package http

import (
	"github.com/gin-gonic/gin"
	"gopa/internal/core/domain"
	"gopa/internal/services"
	"gopa/pkg/utils/response"
)

type PomodoroHandler struct {
	pomodoro *services.PomodoroService
}

func NewPomodoroHandler(pomodoro *services.PomodoroService) *PomodoroHandler {
	return &PomodoroHandler{pomodoro: pomodoro}
}

func (h *PomodoroHandler) RegisterRoutes(api *gin.RouterGroup, authenticate gin.HandlerFunc) {
	p := api.Group("/pomodoro", authenticate)
	p.GET("", h.get)
	p.PUT("", h.set)
	p.POST("/stop", h.stop)
	p.GET("/history", h.history)
}

func (h *PomodoroHandler) get(c *gin.Context) {
	state, err := h.pomodoro.Get(c.Request.Context(), userID(c))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, state)
}

func (h *PomodoroHandler) set(c *gin.Context) {
	var state domain.PomodoroState
	if !bindJSON(c, &state) {
		return
	}
	if err := h.pomodoro.Set(c.Request.Context(), userID(c), state); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"status": "ok"})
}

func (h *PomodoroHandler) stop(c *gin.Context) {
	if err := h.pomodoro.Stop(c.Request.Context(), userID(c)); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"stopped": true})
}

func (h *PomodoroHandler) history(c *gin.Context) {
	limit := queryLimit(c, 20)
	cursor := c.Query("cursor")
	items, err := h.pomodoro.ListHistory(c.Request.Context(), userID(c), limit, cursor)
	if err != nil {
		response.Fail(c, err)
		return
	}
	var nextCursor string
	if len(items) == limit {
		nextCursor = items[len(items)-1].StartedAt.Format("2006-01-02T15:04:05.999999Z07:00")
	}
	response.Paginated(c, items, nextCursor, nil)
}
