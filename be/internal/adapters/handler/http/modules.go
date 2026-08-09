package http

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopa/internal/core/domain"
	"gopa/internal/services"
)

type ModuleHandler struct {
	tasks        *services.TaskService
	vocabularies *services.VocabularyService
	assets       *services.AssetService
	finance      *services.FinanceService
	journals     *services.JournalService
	pomodoro     *services.PomodoroService
}

func NewModuleHandler(tasks *services.TaskService, vocabularies *services.VocabularyService, assets *services.AssetService, finance *services.FinanceService, journals *services.JournalService, pomodoro *services.PomodoroService) *ModuleHandler {
	return &ModuleHandler{tasks: tasks, vocabularies: vocabularies, assets: assets, finance: finance, journals: journals, pomodoro: pomodoro}
}
func (h *ModuleHandler) RegisterRoutes(api *gin.RouterGroup, authenticate gin.HandlerFunc) {
	protected := api.Group("", authenticate)
	protected.GET("/tasks", h.listTasks)
	protected.POST("/tasks", h.createTask)
	protected.PATCH("/tasks/:id", h.updateTask)
	protected.DELETE("/tasks/:id", h.deleteTask)
	protected.GET("/vocabularies", h.listVocabularies)
	protected.POST("/vocabularies", h.createVocabulary)
	protected.GET("/vocabularies/review-queue", h.reviewQueue)
	protected.GET("/vocabularies/:id", h.getVocabulary)
	protected.PATCH("/vocabularies/:id", h.updateVocabulary)
	protected.DELETE("/vocabularies/:id", h.deleteVocabulary)
	protected.POST("/vocabularies/:id/reviews", h.reviewVocabulary)
	protected.GET("/vehicles", h.listVehicles)
	protected.POST("/vehicles", h.createVehicle)
	protected.GET("/vehicles/:id", h.getVehicle)
	protected.PATCH("/vehicles/:id", h.updateVehicle)
	protected.GET("/vehicles/:id/logs", h.listVehicleLogs)
	protected.POST("/vehicles/:id/logs", h.createVehicleLog)
	protected.GET("/vehicles/:id/maintenance-status", h.maintenanceStatus)
	protected.GET("/network-nodes", h.listNetworkNodes)
	protected.POST("/network-nodes", h.createNetworkNode)
	protected.PATCH("/network-nodes/:id", h.updateNetworkNode)
	protected.DELETE("/network-nodes/:id", h.deleteNetworkNode)
	h.registerFinanceRoutes(protected)
	protected.GET("/pomodoro", h.getPomodoro)
	protected.PUT("/pomodoro", h.setPomodoro)
	protected.POST("/pomodoro/stop", h.stopPomodoro)
	protected.GET("/journals", h.listJournals)
	protected.GET("/journals/search", h.listJournals)
	protected.GET("/journals/stats", h.journalStats)
	protected.POST("/journals", h.createJournal)
	protected.POST("/journals/:id/link", h.linkJournal)
	protected.DELETE("/journals/:id/link/:linked_id", h.unlinkJournal)
	protected.GET("/journals/:id", h.getJournal)
	protected.PATCH("/journals/:id", h.updateJournal)
	protected.DELETE("/journals/:id", h.deleteJournal)
}

func (h *ModuleHandler) listTasks(c *gin.Context) {
	tasks, err := h.tasks.List(c, userID(c))
	respond(c, tasks, err, http.StatusOK)
}
func (h *ModuleHandler) createTask(c *gin.Context) {
	var request taskRequest
	if !bind(c, &request) {
		return
	}
	task, err := h.tasks.Create(c, userID(c), request.input())
	respond(c, task, err, http.StatusCreated)
}
func (h *ModuleHandler) updateTask(c *gin.Context) {
	id, ok := parameterUUID(c)
	if !ok {
		return
	}
	var request taskRequest
	if !bind(c, &request) {
		return
	}
	task, err := h.tasks.Update(c, userID(c), id, request.updateInput())
	respond(c, task, err, http.StatusOK)
}
func (h *ModuleHandler) deleteTask(c *gin.Context) {
	id, ok := parameterUUID(c)
	if !ok {
		return
	}
	err := h.tasks.Delete(c, userID(c), id)
	respond(c, gin.H{}, err, http.StatusOK)
}

func (h *ModuleHandler) listVocabularies(c *gin.Context) {
	values, err := h.vocabularies.List(c, userID(c), limit(c))
	respond(c, values, err, http.StatusOK)
}
func (h *ModuleHandler) reviewQueue(c *gin.Context) {
	values, err := h.vocabularies.Due(c, userID(c), limit(c))
	respond(c, values, err, http.StatusOK)
}
func (h *ModuleHandler) getVocabulary(c *gin.Context) {
	id, ok := parameterUUID(c)
	if !ok {
		return
	}
	value, err := h.vocabularies.Get(c, userID(c), id)
	respond(c, value, err, http.StatusOK)
}
func (h *ModuleHandler) createVocabulary(c *gin.Context) {
	var request vocabularyRequest
	if !bind(c, &request) {
		return
	}
	value, err := h.vocabularies.Create(c, userID(c), request.input())
	respond(c, value, err, http.StatusCreated)
}
func (h *ModuleHandler) updateVocabulary(c *gin.Context) {
	id, ok := parameterUUID(c)
	if !ok {
		return
	}
	var request vocabularyRequest
	if !bind(c, &request) {
		return
	}
	value, err := h.vocabularies.Update(c, userID(c), id, request.input())
	respond(c, value, err, http.StatusOK)
}
func (h *ModuleHandler) deleteVocabulary(c *gin.Context) {
	id, ok := parameterUUID(c)
	if !ok {
		return
	}
	respond(c, gin.H{}, h.vocabularies.Delete(c, userID(c), id), http.StatusOK)
}
func (h *ModuleHandler) reviewVocabulary(c *gin.Context) {
	id, ok := parameterUUID(c)
	if !ok {
		return
	}
	var request reviewRequest
	if !bind(c, &request) {
		return
	}
	respond(c, gin.H{}, h.vocabularies.Review(c, userID(c), id, request.Quality), http.StatusOK)
}

func (h *ModuleHandler) listVehicles(c *gin.Context) {
	values, err := h.assets.ListVehicles(c, userID(c))
	respond(c, values, err, http.StatusOK)
}
func (h *ModuleHandler) getVehicle(c *gin.Context) {
	id, ok := parameterUUID(c)
	if !ok {
		return
	}
	value, err := h.assets.GetVehicle(c, userID(c), id)
	respond(c, value, err, http.StatusOK)
}
func (h *ModuleHandler) createVehicle(c *gin.Context) {
	var request vehicleRequest
	if !bind(c, &request) {
		return
	}
	value, err := h.assets.CreateVehicle(c, userID(c), request.input())
	respond(c, value, err, http.StatusCreated)
}
func (h *ModuleHandler) updateVehicle(c *gin.Context) {
	id, ok := parameterUUID(c)
	if !ok {
		return
	}
	var request vehicleRequest
	if !bind(c, &request) {
		return
	}
	value, err := h.assets.UpdateVehicle(c, userID(c), id, request.input())
	respond(c, value, err, http.StatusOK)
}
func (h *ModuleHandler) listVehicleLogs(c *gin.Context) {
	id, ok := parameterUUID(c)
	if !ok {
		return
	}
	values, err := h.assets.ListVehicleLogs(c, userID(c), id)
	respond(c, values, err, http.StatusOK)
}
func (h *ModuleHandler) createVehicleLog(c *gin.Context) {
	id, ok := parameterUUID(c)
	if !ok {
		return
	}
	var request vehicleLogRequest
	if !bind(c, &request) {
		return
	}
	value, err := h.assets.CreateVehicleLog(c, userID(c), id, request.input())
	respond(c, value, err, http.StatusCreated)
}
func (h *ModuleHandler) maintenanceStatus(c *gin.Context) {
	id, ok := parameterUUID(c)
	if !ok {
		return
	}
	value, err := h.assets.Maintenance(c, userID(c), id)
	respond(c, value, err, http.StatusOK)
}

func (h *ModuleHandler) listNetworkNodes(c *gin.Context) {
	values, err := h.assets.ListNodes(c, userID(c))
	respond(c, values, err, http.StatusOK)
}
func (h *ModuleHandler) createNetworkNode(c *gin.Context) {
	var request networkNodeRequest
	if !bind(c, &request) {
		return
	}
	value, err := h.assets.CreateNode(c, userID(c), request.node())
	respond(c, value, err, http.StatusCreated)
}
func (h *ModuleHandler) updateNetworkNode(c *gin.Context) {
	id, ok := parameterUUID(c)
	if !ok {
		return
	}
	var request networkNodeRequest
	if !bind(c, &request) {
		return
	}
	value, err := h.assets.UpdateNode(c, userID(c), id, request.node())
	respond(c, value, err, http.StatusOK)
}
func (h *ModuleHandler) deleteNetworkNode(c *gin.Context) {
	id, ok := parameterUUID(c)
	if !ok {
		return
	}
	respond(c, gin.H{}, h.assets.DeleteNode(c, userID(c), id), http.StatusOK)
}

func (h *ModuleHandler) getPomodoro(c *gin.Context) {
	state, err := h.pomodoro.Get(c, userID(c))
	respond(c, state, err, http.StatusOK)
}
func (h *ModuleHandler) setPomodoro(c *gin.Context) {
	var state domain.PomodoroState
	if !bind(c, &state) {
		return
	}
	respond(c, gin.H{}, h.pomodoro.Set(c, userID(c), state), http.StatusOK)
}
func (h *ModuleHandler) stopPomodoro(c *gin.Context) {
	respond(c, gin.H{}, h.pomodoro.Stop(c, userID(c)), http.StatusOK)
}

func (h *ModuleHandler) listJournals(c *gin.Context) {
	tag := c.Query("tag")
	if tag == "" {
		tag = c.Query("tags")
	}
	filter := domain.JournalFilter{Term: c.Query("q"), Tag: tag, Limit: limit(c)}
	if value := c.Query("mood"); value != "" {
		mood := domain.JournalMood(value)
		filter.Mood = &mood
	}
	if value := c.Query("from"); value != "" {
		if parsed, err := time.Parse("2006-01-02", value); err == nil {
			filter.From = &parsed
		} else {
			respond(c, nil, domain.ErrValidation, http.StatusBadRequest)
			return
		}
	}
	if value := c.Query("to"); value != "" {
		if parsed, err := time.Parse("2006-01-02", value); err == nil {
			filter.To = &parsed
		} else {
			respond(c, nil, domain.ErrValidation, http.StatusBadRequest)
			return
		}
	}
	values, err := h.journals.List(c, userID(c), filter)
	respond(c, values, err, http.StatusOK)
}
func (h *ModuleHandler) journalStats(c *gin.Context) {
	value, err := h.journals.Stats(c, userID(c))
	respond(c, value, err, http.StatusOK)
}
func (h *ModuleHandler) linkJournal(c *gin.Context) {
	journalID, ok := parameterUUID(c)
	if !ok {
		return
	}
	var request struct {
		LinkedID uuid.UUID `json:"linked_id" binding:"required"`
	}
	if !bind(c, &request) {
		return
	}
	respond(c, gin.H{}, h.journals.Link(c, userID(c), journalID, request.LinkedID), http.StatusOK)
}
func (h *ModuleHandler) unlinkJournal(c *gin.Context) {
	journalID, ok := parameterUUID(c)
	if !ok {
		return
	}
	linkedID, err := uuid.Parse(c.Param("linked_id"))
	if err != nil {
		writeError(c, domain.ErrValidation)
		return
	}
	respond(c, gin.H{}, h.journals.Unlink(c, userID(c), journalID, linkedID), http.StatusOK)
}
func (h *ModuleHandler) getJournal(c *gin.Context) {
	id, ok := parameterUUID(c)
	if !ok {
		return
	}
	value, err := h.journals.Get(c, userID(c), id)
	respond(c, value, err, http.StatusOK)
}
func (h *ModuleHandler) createJournal(c *gin.Context) {
	var request journalRequest
	if !bind(c, &request) {
		return
	}
	value, err := h.journals.Create(c, userID(c), request.input())
	respond(c, value, err, http.StatusCreated)
}
func (h *ModuleHandler) updateJournal(c *gin.Context) {
	id, ok := parameterUUID(c)
	if !ok {
		return
	}
	var request journalRequest
	if !bind(c, &request) {
		return
	}
	value, err := h.journals.Update(c, userID(c), id, request.input())
	respond(c, value, err, http.StatusOK)
}
func (h *ModuleHandler) deleteJournal(c *gin.Context) {
	id, ok := parameterUUID(c)
	if !ok {
		return
	}
	respond(c, gin.H{}, h.journals.Delete(c, userID(c), id), http.StatusOK)
}

type taskRequest struct {
	Title       string              `json:"title" binding:"required,max=200"`
	Description string              `json:"description"`
	Status      domain.TaskStatus   `json:"status" binding:"required"`
	Priority    domain.TaskPriority `json:"priority" binding:"required"`
	Category    domain.TaskCategory `json:"category" binding:"required"`
	DueDate     *time.Time          `json:"due_date"`
	Position    *int64              `json:"position"`
}

func (r taskRequest) input() services.CreateTaskInput {
	return services.CreateTaskInput{Title: r.Title, Description: r.Description, Status: r.Status, Priority: r.Priority, Category: r.Category, DueDate: r.DueDate}
}
func (r taskRequest) updateInput() services.UpdateTaskInput {
	return services.UpdateTaskInput{Title: r.Title, Description: r.Description, Status: r.Status, Priority: r.Priority, Category: r.Category, DueDate: r.DueDate, Position: r.Position}
}

type vocabularyRequest struct {
	Language        domain.VocabularyLanguage `json:"language" binding:"required"`
	Word            string                    `json:"word" binding:"required,max=300"`
	Reading         *string                   `json:"reading"`
	Meaning         string                    `json:"meaning" binding:"required,max=1000"`
	ExampleSentence *string                   `json:"example_sentence"`
}

func (r vocabularyRequest) input() services.VocabularyInput {
	return services.VocabularyInput{Language: r.Language, Word: r.Word, Reading: r.Reading, Meaning: r.Meaning, ExampleSentence: r.ExampleSentence}
}

type reviewRequest struct {
	Quality int16 `json:"quality" binding:"gte=0,lte=5"`
}
type vehicleRequest struct {
	Make           string  `json:"make" binding:"required,max=100"`
	Model          string  `json:"model" binding:"required,max=100"`
	Year           int16   `json:"year" binding:"required"`
	Nickname       *string `json:"nickname"`
	CurrentMileage int     `json:"current_mileage" binding:"gte=0"`
}

func (r vehicleRequest) input() services.VehicleInput {
	return services.VehicleInput{Make: r.Make, Model: r.Model, Year: r.Year, Nickname: r.Nickname, CurrentMileage: r.CurrentMileage}
}

type vehicleLogRequest struct {
	LogType     domain.VehicleLogType `json:"log_type" binding:"required"`
	Mileage     int                   `json:"mileage" binding:"gte=0"`
	Description string                `json:"description" binding:"required,max=2000"`
	Cost        string                `json:"cost" binding:"required"`
	LogDate     time.Time             `json:"log_date" binding:"required"`
}

func (r vehicleLogRequest) input() services.VehicleLogInput {
	return services.VehicleLogInput{LogType: r.LogType, Mileage: r.Mileage, Description: r.Description, Cost: r.Cost, LogDate: r.LogDate}
}

type networkNodeRequest struct {
	DeviceName     string                `json:"device_name" binding:"required,max=200"`
	MACAddress     *string               `json:"mac_address"`
	StaticIP       *string               `json:"static_ip"`
	VLAN           *int16                `json:"vlan_tag"`
	ConnectionType domain.ConnectionType `json:"connection_type" binding:"required"`
	ParentNodeID   *uuid.UUID            `json:"parent_node_id"`
	Location       *string               `json:"location"`
	Notes          *string               `json:"notes"`
}

func (r networkNodeRequest) node() domain.NetworkNode {
	return domain.NetworkNode{DeviceName: r.DeviceName, MACAddress: r.MACAddress, StaticIP: r.StaticIP, VLAN: r.VLAN, ConnectionType: r.ConnectionType, ParentNodeID: r.ParentNodeID, Location: r.Location, Notes: r.Notes}
}

type journalRequest struct {
	Title         string              `json:"title" binding:"required,max=300"`
	Content       string              `json:"content"`
	Tags          []string            `json:"tags"`
	Mood          *domain.JournalMood `json:"mood"`
	PublishedDate *time.Time          `json:"published_date"`
}

func (r journalRequest) input() services.JournalInput {
	return services.JournalInput{Title: r.Title, Content: r.Content, Tags: r.Tags, Mood: r.Mood, PublishedDate: r.PublishedDate}
}

func userID(c *gin.Context) uuid.UUID { identity, _ := identityFromContext(c); return identity.UserID }
func bind(c *gin.Context, target any) bool {
	if err := c.ShouldBindJSON(target); err != nil {
		writeError(c, domain.ErrValidation)
		return false
	}
	return true
}
func parameterUUID(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		writeError(c, domain.ErrValidation)
		return uuid.Nil, false
	}
	return id, true
}
func limit(c *gin.Context) int {
	value, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		return 20
	}
	return value
}
func respond(c *gin.Context, data any, err error, status int) {
	if err != nil {
		slog.Error("module request failed", "request_id", requestID(c), "error", err)
		writeError(c, err)
		return
	}
	c.JSON(status, success(c, data))
}
