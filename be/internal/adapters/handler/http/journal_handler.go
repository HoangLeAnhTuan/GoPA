package http

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopa/internal/core/domain"
	"gopa/internal/services"
	"gopa/pkg/utils/response"
)

type JournalHandler struct {
	journals *services.JournalService
}

func NewJournalHandler(journals *services.JournalService) *JournalHandler {
	return &JournalHandler{journals: journals}
}

func (h *JournalHandler) RegisterRoutes(api *gin.RouterGroup, authenticate gin.HandlerFunc) {
	j := api.Group("/journals", authenticate)
	j.GET("", h.list)
	j.POST("", h.create)
	j.GET("/search", h.search)
	j.GET("/stats", h.stats)
	j.GET("/:id", h.get)
	j.PATCH("/:id", h.update)
	j.DELETE("/:id", h.delete)
	j.POST("/:id/link", h.link)
	j.DELETE("/:id/link/:linked_id", h.unlink)
}

func (h *JournalHandler) list(c *gin.Context) {
	tag := c.Query("tag")
	if tag == "" {
		tag = c.Query("tags")
	}
	filter := domain.JournalFilter{
		Term:  c.Query("q"),
		Tag:   tag,
		Limit: queryLimit(c, 20),
	}
	if value := c.Query("mood"); value != "" {
		mood := domain.JournalMood(value)
		filter.Mood = &mood
	}
	if value := c.Query("from"); value != "" {
		if parsed, err := time.Parse("2006-01-02", value); err == nil {
			filter.From = &parsed
		} else {
			response.Fail(c, domain.ErrValidation)
			return
		}
	}
	if value := c.Query("to"); value != "" {
		if parsed, err := time.Parse("2006-01-02", value); err == nil {
			filter.To = &parsed
		} else {
			response.Fail(c, domain.ErrValidation)
			return
		}
	}
	items, err := h.journals.List(c.Request.Context(), userID(c), filter)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, items)
}

func (h *JournalHandler) search(c *gin.Context) {
	h.list(c)
}

func (h *JournalHandler) stats(c *gin.Context) {
	st, err := h.journals.Stats(c.Request.Context(), userID(c))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, st)
}

func (h *JournalHandler) get(c *gin.Context) {
	id, ok := parseParamUUID(c, "id")
	if !ok {
		return
	}
	item, err := h.journals.Get(c.Request.Context(), userID(c), id)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, item)
}

type journalRequest struct {
	Title         string              `json:"title" binding:"required,max=300"`
	Content       string              `json:"content"`
	Tags          []string            `json:"tags"`
	Mood          *domain.JournalMood `json:"mood"`
	EnergyLevel   *int16              `json:"energy_level"`
	Pinned        *bool               `json:"pinned"`
	PublishedDate *FlexibleDate       `json:"published_date"`
}

func (h *JournalHandler) create(c *gin.Context) {
	var req journalRequest
	if !bindJSON(c, &req) {
		return
	}
	val, err := h.journals.Create(c.Request.Context(), userID(c), services.JournalInput{
		Title:         req.Title,
		Content:       req.Content,
		Tags:          req.Tags,
		Mood:          req.Mood,
		EnergyLevel:   req.EnergyLevel,
		Pinned:        req.Pinned,
		PublishedDate: req.PublishedDate.Time(),
	})
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, val)
}

func (h *JournalHandler) update(c *gin.Context) {
	id, ok := parseParamUUID(c, "id")
	if !ok {
		return
	}
	var req journalRequest
	if !bindJSON(c, &req) {
		return
	}
	val, err := h.journals.Update(c.Request.Context(), userID(c), id, services.JournalInput{
		Title:         req.Title,
		Content:       req.Content,
		Tags:          req.Tags,
		Mood:          req.Mood,
		EnergyLevel:   req.EnergyLevel,
		Pinned:        req.Pinned,
		PublishedDate: req.PublishedDate.Time(),
	})
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, val)
}

func (h *JournalHandler) delete(c *gin.Context) {
	id, ok := parseParamUUID(c, "id")
	if !ok {
		return
	}
	if err := h.journals.Delete(c.Request.Context(), userID(c), id); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"id": id})
}

type linkRequest struct {
	LinkedID uuid.UUID `json:"linked_id" binding:"required"`
}

func (h *JournalHandler) link(c *gin.Context) {
	id, ok := parseParamUUID(c, "id")
	if !ok {
		return
	}
	var req linkRequest
	if !bindJSON(c, &req) {
		return
	}
	if err := h.journals.Link(c.Request.Context(), userID(c), id, req.LinkedID); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"linked": true})
}

func (h *JournalHandler) unlink(c *gin.Context) {
	id, ok := parseParamUUID(c, "id")
	if !ok {
		return
	}
	linkedID, ok := parseParamUUID(c, "linked_id")
	if !ok {
		return
	}
	if err := h.journals.Unlink(c.Request.Context(), userID(c), id, linkedID); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"unlinked": true})
}
