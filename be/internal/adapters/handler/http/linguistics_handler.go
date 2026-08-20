package http

import (
	"io"

	"github.com/gin-gonic/gin"
	"gopa/internal/core/domain"
	"gopa/internal/services"
	"gopa/pkg/utils/response"
)

type LinguisticsHandler struct {
	vocabularies *services.VocabularyService
}

func NewLinguisticsHandler(vocabularies *services.VocabularyService) *LinguisticsHandler {
	return &LinguisticsHandler{vocabularies: vocabularies}
}

func (h *LinguisticsHandler) RegisterRoutes(api *gin.RouterGroup, authenticate gin.HandlerFunc) {
	vocab := api.Group("/vocabularies", authenticate)
	vocab.GET("", h.list)
	vocab.POST("", h.create)
	vocab.GET("/review-queue", h.reviewQueue)
	vocab.GET("/stats", h.getStats)
	vocab.POST("/import", h.importVocabularies)
	vocab.GET("/:id", h.get)
	vocab.PATCH("/:id", h.update)
	vocab.DELETE("/:id", h.delete)
	vocab.POST("/:id/reviews", h.recordReview)
	vocab.POST("/:id/audio", h.generateAudio)

	sessions := api.Group("/learning-sessions", authenticate)
	sessions.POST("", h.startSession)
	sessions.PATCH("/:id", h.endSession)
}

func (h *LinguisticsHandler) list(c *gin.Context) {
	items, err := h.vocabularies.List(c.Request.Context(), userID(c), queryLimit(c, 20))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, items)
}

func (h *LinguisticsHandler) reviewQueue(c *gin.Context) {
	items, err := h.vocabularies.Due(c.Request.Context(), userID(c), queryLimit(c, 20))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, items)
}

func (h *LinguisticsHandler) getStats(c *gin.Context) {
	stats, err := h.vocabularies.GetStats(c.Request.Context(), userID(c))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, stats)
}

func (h *LinguisticsHandler) get(c *gin.Context) {
	id, ok := parseParamUUID(c, "id")
	if !ok {
		return
	}
	vocab, err := h.vocabularies.Get(c.Request.Context(), userID(c), id)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, vocab)
}

type vocabularyRequest struct {
	Language           domain.VocabularyLanguage `json:"language" binding:"required"`
	Word               string                    `json:"word" binding:"required,max=300"`
	Reading            *string                   `json:"reading"`
	Meaning            string                    `json:"meaning" binding:"required,max=1000"`
	ExampleSentence    *string                   `json:"example_sentence"`
	ExampleTranslation *string                   `json:"example_translation"`
	Tags               []string                  `json:"tags"`
	DifficultyLevel    string                    `json:"difficulty_level"`
}

func (h *LinguisticsHandler) create(c *gin.Context) {
	var req vocabularyRequest
	if !bindJSON(c, &req) {
		return
	}
	created, err := h.vocabularies.Create(c.Request.Context(), userID(c), services.VocabularyInput{
		Language:           req.Language,
		Word:               req.Word,
		Reading:            req.Reading,
		Meaning:            req.Meaning,
		ExampleSentence:    req.ExampleSentence,
		ExampleTranslation: req.ExampleTranslation,
		Tags:               req.Tags,
		DifficultyLevel:    req.DifficultyLevel,
	})
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, created)
}

func (h *LinguisticsHandler) update(c *gin.Context) {
	id, ok := parseParamUUID(c, "id")
	if !ok {
		return
	}
	var req vocabularyRequest
	if !bindJSON(c, &req) {
		return
	}
	updated, err := h.vocabularies.Update(c.Request.Context(), userID(c), id, services.VocabularyInput{
		Language:           req.Language,
		Word:               req.Word,
		Reading:            req.Reading,
		Meaning:            req.Meaning,
		ExampleSentence:    req.ExampleSentence,
		ExampleTranslation: req.ExampleTranslation,
		Tags:               req.Tags,
		DifficultyLevel:    req.DifficultyLevel,
	})
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, updated)
}

func (h *LinguisticsHandler) delete(c *gin.Context) {
	id, ok := parseParamUUID(c, "id")
	if !ok {
		return
	}
	if err := h.vocabularies.Delete(c.Request.Context(), userID(c), id); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"id": id})
}

type reviewRequest struct {
	Quality int16 `json:"quality" binding:"gte=0,lte=5"`
}

func (h *LinguisticsHandler) recordReview(c *gin.Context) {
	id, ok := parseParamUUID(c, "id")
	if !ok {
		return
	}
	var req reviewRequest
	if !bindJSON(c, &req) {
		return
	}
	if err := h.vocabularies.Review(c.Request.Context(), userID(c), id, req.Quality); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"reviewed": true})
}

func (h *LinguisticsHandler) generateAudio(c *gin.Context) {
	id, ok := parseParamUUID(c, "id")
	if !ok {
		return
	}
	vocab, err := h.vocabularies.Get(c.Request.Context(), userID(c), id)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"vocabulary_id": vocab.ID, "audio_url": vocab.AudioURL})
}

func (h *LinguisticsHandler) importVocabularies(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil || len(body) == 0 {
		response.Fail(c, domain.ErrValidation)
		return
	}
	format := c.DefaultQuery("format", "json")
	count, err := h.vocabularies.ImportVocabularies(c.Request.Context(), userID(c), body, format)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, gin.H{"imported_count": count})
}

type startSessionRequest struct {
	Language    domain.VocabularyLanguage `json:"language" binding:"required"`
	SessionType string                    `json:"session_type"`
}

func (h *LinguisticsHandler) startSession(c *gin.Context) {
	var req startSessionRequest
	if !bindJSON(c, &req) {
		return
	}
	session, err := h.vocabularies.StartSession(c.Request.Context(), userID(c), req.Language, req.SessionType)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, session)
}

type endSessionRequest struct {
	ItemsReviewed   int `json:"items_reviewed" binding:"gte=0"`
	ItemsCorrect    int `json:"items_correct" binding:"gte=0"`
	DurationSeconds int `json:"duration_seconds" binding:"gte=0"`
}

func (h *LinguisticsHandler) endSession(c *gin.Context) {
	id, ok := parseParamUUID(c, "id")
	if !ok {
		return
	}
	var req endSessionRequest
	if !bindJSON(c, &req) {
		return
	}
	session, err := h.vocabularies.EndSession(c.Request.Context(), userID(c), id, req.ItemsReviewed, req.ItemsCorrect, req.DurationSeconds)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, session)
}
