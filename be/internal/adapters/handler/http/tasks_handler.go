package http

import (
	"github.com/gin-gonic/gin"
	"gopa/internal/core/domain"
	"gopa/internal/services"
	"gopa/pkg/utils/response"
)

type TasksHandler struct {
	tasks *services.TaskService
}

func NewTasksHandler(tasks *services.TaskService) *TasksHandler {
	return &TasksHandler{tasks: tasks}
}

func (h *TasksHandler) RegisterRoutes(api *gin.RouterGroup, authenticate gin.HandlerFunc) {
	tasks := api.Group("/tasks", authenticate)
	tasks.GET("", h.list)
	tasks.POST("", h.create)
	tasks.GET("/:id", h.get)
	tasks.PATCH("/:id", h.update)
	tasks.DELETE("/:id", h.delete)
	tasks.PATCH("/:id/status", h.updateStatus)
}

func (h *TasksHandler) list(c *gin.Context) {
	tasks, err := h.tasks.List(c.Request.Context(), userID(c))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, tasks)
}

func (h *TasksHandler) get(c *gin.Context) {
	id, ok := parseParamUUID(c, "id")
	if !ok {
		return
	}
	task, err := h.tasks.Get(c.Request.Context(), userID(c), id)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, task)
}

type createTaskRequest struct {
	Title              string              `json:"title" binding:"required,max=200"`
	Description        string              `json:"description"`
	Status             domain.TaskStatus   `json:"status"`
	Priority           domain.TaskPriority `json:"priority"`
	Category           domain.TaskCategory `json:"category"`
	DueDate            *FlexibleDate       `json:"due_date"`
	EstimatedPomodoros int                 `json:"estimated_pomodoros"`
}

func (h *TasksHandler) create(c *gin.Context) {
	var req createTaskRequest
	if !bindJSON(c, &req) {
		return
	}
	task, err := h.tasks.Create(c.Request.Context(), userID(c), services.CreateTaskInput{
		Title:              req.Title,
		Description:        req.Description,
		Status:             req.Status,
		Priority:           req.Priority,
		Category:           req.Category,
		DueDate:            req.DueDate.Time(),
		EstimatedPomodoros: req.EstimatedPomodoros,
	})
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, task)
}

type updateTaskRequest struct {
	Title              string              `json:"title" binding:"required,max=200"`
	Description        string              `json:"description"`
	Status             domain.TaskStatus   `json:"status"`
	Priority           domain.TaskPriority `json:"priority"`
	Category           domain.TaskCategory `json:"category"`
	DueDate            *FlexibleDate       `json:"due_date"`
	Position           *int64              `json:"position"`
	EstimatedPomodoros *int                `json:"estimated_pomodoros"`
	CompletedPomodoros *int                `json:"completed_pomodoros"`
}

func (h *TasksHandler) update(c *gin.Context) {
	id, ok := parseParamUUID(c, "id")
	if !ok {
		return
	}
	var req updateTaskRequest
	if !bindJSON(c, &req) {
		return
	}
	task, err := h.tasks.Update(c.Request.Context(), userID(c), id, services.UpdateTaskInput{
		Title:              req.Title,
		Description:        req.Description,
		Status:             req.Status,
		Priority:           req.Priority,
		Category:           req.Category,
		DueDate:            req.DueDate.Time(),
		Position:           req.Position,
		EstimatedPomodoros: req.EstimatedPomodoros,
		CompletedPomodoros: req.CompletedPomodoros,
	})
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, task)
}

type updateTaskStatusRequest struct {
	Status   domain.TaskStatus `json:"status" binding:"required"`
	Position *int64            `json:"position"`
}

func (h *TasksHandler) updateStatus(c *gin.Context) {
	id, ok := parseParamUUID(c, "id")
	if !ok {
		return
	}
	var req updateTaskStatusRequest
	if !bindJSON(c, &req) {
		return
	}
	task, err := h.tasks.UpdateStatus(c.Request.Context(), userID(c), id, req.Status, req.Position)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, task)
}

func (h *TasksHandler) delete(c *gin.Context) {
	id, ok := parseParamUUID(c, "id")
	if !ok {
		return
	}
	if err := h.tasks.Delete(c.Request.Context(), userID(c), id); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"id": id})
}
