package http

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopa/internal/core/domain"
	"gopa/internal/services"
	"gopa/pkg/utils"
	"gopa/pkg/utils/response"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// Stubs for test router
type testTaskRepo struct {
	tasks map[uuid.UUID]domain.Task
}

func newTestTaskRepo() *testTaskRepo {
	return &testTaskRepo{tasks: make(map[uuid.UUID]domain.Task)}
}
func (r *testTaskRepo) List(_ context.Context, userID uuid.UUID) ([]domain.Task, error) {
	var res []domain.Task
	for _, t := range r.tasks {
		if t.UserID == userID {
			res = append(res, t)
		}
	}
	return res, nil
}
func (r *testTaskRepo) Get(_ context.Context, userID, id uuid.UUID) (domain.Task, error) {
	t, ok := r.tasks[id]
	if !ok || t.UserID != userID {
		return domain.Task{}, domain.ErrNotFound
	}
	return t, nil
}
func (r *testTaskRepo) Create(_ context.Context, t domain.Task) (domain.Task, error) {
	r.tasks[t.ID] = t
	return t, nil
}
func (r *testTaskRepo) Update(_ context.Context, t domain.Task) (domain.Task, error) {
	if _, ok := r.tasks[t.ID]; !ok {
		return domain.Task{}, domain.ErrNotFound
	}
	r.tasks[t.ID] = t
	return t, nil
}
func (r *testTaskRepo) UpdateStatus(_ context.Context, userID, id uuid.UUID, status domain.TaskStatus, targetPos *int64) (domain.Task, error) {
	t, ok := r.tasks[id]
	if !ok || t.UserID != userID {
		return domain.Task{}, domain.ErrNotFound
	}
	t.Status = status
	if targetPos != nil {
		t.Position = *targetPos
	}
	r.tasks[id] = t
	return t, nil
}
func (r *testTaskRepo) Delete(_ context.Context, userID, id uuid.UUID) error {
	t, ok := r.tasks[id]
	if !ok || t.UserID != userID {
		return domain.ErrNotFound
	}
	delete(r.tasks, id)
	return nil
}
func (r *testTaskRepo) NextPosition(_ context.Context, _ uuid.UUID, _ domain.TaskStatus) (int64, error) {
	return 1000, nil
}

type testVocabRepo struct {
	vocabs map[uuid.UUID]domain.Vocabulary
}

func newTestVocabRepo() *testVocabRepo {
	return &testVocabRepo{vocabs: make(map[uuid.UUID]domain.Vocabulary)}
}
func (r *testVocabRepo) List(_ context.Context, userID uuid.UUID, _ int) ([]domain.Vocabulary, error) {
	var res []domain.Vocabulary
	for _, v := range r.vocabs {
		if v.UserID == userID {
			res = append(res, v)
		}
	}
	return res, nil
}
func (r *testVocabRepo) ListDue(_ context.Context, _ uuid.UUID, _ int, _ time.Time) ([]domain.Vocabulary, error) {
	return nil, nil
}
func (r *testVocabRepo) Get(_ context.Context, userID, id uuid.UUID) (domain.Vocabulary, error) {
	v, ok := r.vocabs[id]
	if !ok || v.UserID != userID {
		return domain.Vocabulary{}, domain.ErrNotFound
	}
	return v, nil
}
func (r *testVocabRepo) Create(_ context.Context, v domain.Vocabulary) (domain.Vocabulary, error) {
	r.vocabs[v.ID] = v
	return v, nil
}
func (r *testVocabRepo) BulkCreate(_ context.Context, userID uuid.UUID, items []domain.Vocabulary) ([]domain.Vocabulary, error) {
	for _, item := range items {
		item.UserID = userID
		r.vocabs[item.ID] = item
	}
	return items, nil
}
func (r *testVocabRepo) Update(_ context.Context, v domain.Vocabulary) (domain.Vocabulary, error) {
	r.vocabs[v.ID] = v
	return v, nil
}
func (r *testVocabRepo) Delete(_ context.Context, userID, id uuid.UUID) error {
	delete(r.vocabs, id)
	return nil
}
func (r *testVocabRepo) RecordReview(_ context.Context, _ domain.VocabularyReview, updated domain.Vocabulary) error {
	r.vocabs[updated.ID] = updated
	return nil
}
func (r *testVocabRepo) CreateLearningSession(_ context.Context, s domain.LearningSession) (domain.LearningSession, error) {
	return s, nil
}
func (r *testVocabRepo) UpdateLearningSession(_ context.Context, s domain.LearningSession) (domain.LearningSession, error) {
	return s, nil
}
func (r *testVocabRepo) GetStats(_ context.Context, _ uuid.UUID) (domain.VocabularyStats, error) {
	return domain.VocabularyStats{TotalWords: 5}, nil
}

type testPomodoroStore struct {
	state *domain.PomodoroState
}

func (s *testPomodoroStore) Get(_ context.Context, _ uuid.UUID) (domain.PomodoroState, error) {
	if s.state == nil {
		return domain.PomodoroState{}, domain.ErrNotFound
	}
	return *s.state, nil
}
func (s *testPomodoroStore) Set(_ context.Context, _ uuid.UUID, st domain.PomodoroState) error {
	s.state = &st
	return nil
}
func (s *testPomodoroStore) Delete(_ context.Context, _ uuid.UUID) error {
	s.state = nil
	return nil
}

type testPomodoroHistoryRepo struct{}

func (testPomodoroHistoryRepo) Create(_ context.Context, _ domain.PomodoroHistory) error {
	return nil
}
func (testPomodoroHistoryRepo) List(_ context.Context, _ uuid.UUID, _ int, _ string) ([]domain.PomodoroHistory, error) {
	return []domain.PomodoroHistory{
		{ID: uuid.New(), DurationSeconds: 1500, StartedAt: time.Now().UTC()},
	}, nil
}

type testAuthApp struct {
	user domain.User
}

func (a *testAuthApp) Register(_ context.Context, _ services.RegisterInput) (domain.User, error) {
	return a.user, nil
}
func (a *testAuthApp) Login(_ context.Context, _ services.LoginInput) (services.AuthResult, error) {
	return services.AuthResult{User: a.user, AccessToken: "test_token"}, nil
}
func (a *testAuthApp) Refresh(_ context.Context, _ string) (services.AuthResult, error) {
	return services.AuthResult{User: a.user, AccessToken: "test_token"}, nil
}
func (a *testAuthApp) Logout(_ context.Context, _ string) error { return nil }
func (a *testAuthApp) Me(_ context.Context, _ uuid.UUID) (domain.User, error) {
	return a.user, nil
}
func (a *testAuthApp) UpdateProfile(_ context.Context, _ uuid.UUID, input services.UpdateProfileInput) (domain.User, error) {
	if input.DisplayName != nil {
		a.user.DisplayName = *input.DisplayName
	}
	if input.AvatarURL != nil {
		a.user.AvatarURL = input.AvatarURL
	}
	return a.user, nil
}

func setupTestServer() (*gin.Engine, *utils.TokenManager, domain.User) {
	testUser := domain.User{
		ID:          uuid.New(),
		Email:       "test@example.com",
		DisplayName: "Test User",
		Role:        domain.RoleUser,
		CreatedAt:   time.Now().UTC(),
	}
	tokenMgr := utils.NewTokenManager("gopa", "a-very-secret-test-key-32-chars-long!", time.Hour)
	taskRepo := newTestTaskRepo()
	vocabRepo := newTestVocabRepo()
	pomoStore := &testPomodoroStore{}
	pomoHist := testPomodoroHistoryRepo{}

	authApp := &testAuthApp{user: testUser}

	tasksHandler := NewTasksHandler(services.NewTaskService(taskRepo))
	lingHandler := NewLinguisticsHandler(services.NewVocabularyService(vocabRepo))
	pomoHandler := NewPomodoroHandler(services.NewPomodoroService(pomoStore, pomoHist))
	authHandler := NewAuthHandler(authApp, time.Hour, false)

	router := NewRouter(
		slog.Default(),
		"http://localhost:5173",
		NewHealthHandler(readinessStub{}),
		authHandler,
		tasksHandler,
		lingHandler,
		nil,
		nil,
		pomoHandler,
		nil,
		tokenMgr,
		rateLimiterStub{},
	)
	return router, tokenMgr, testUser
}

func authHeader(t *testing.T, tokenMgr *utils.TokenManager, user domain.User) string {
	tok, err := tokenMgr.IssueAccessToken(user, uuid.New(), time.Now().UTC())
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	return "Bearer " + tok
}

func TestTasksHandler_CRUD(t *testing.T) {
	router, tokenMgr, user := setupTestServer()
	token := authHeader(t, tokenMgr, user)

	// 1. Create Task
	body := []byte(`{
		"title": "Build HTTP Handlers",
		"description": "Session 4.1 implementation",
		"status": "TODO",
		"priority": "HIGH",
		"category": "WORK",
		"estimated_pomodoros": 3
	}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewReader(body))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}

	var env response.Envelope
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("parse envelope: %v", err)
	}
	dataMap := env.Data.(map[string]any)
	taskID := dataMap["id"].(string)

	// 2. Get Task
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/"+taskID, nil)
	getReq.Header.Set("Authorization", token)
	gw := httptest.NewRecorder()
	router.ServeHTTP(gw, getReq)

	if gw.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", gw.Code)
	}

	// 3. Update Status
	statusBody := []byte(`{"status": "IN_PROGRESS"}`)
	statusReq := httptest.NewRequest(http.MethodPatch, "/api/v1/tasks/"+taskID+"/status", bytes.NewReader(statusBody))
	statusReq.Header.Set("Authorization", token)
	statusReq.Header.Set("Content-Type", "application/json")
	sw := httptest.NewRecorder()
	router.ServeHTTP(sw, statusReq)

	if sw.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for update status, got %d: %s", sw.Code, sw.Body.String())
	}

	// 4. Delete Task
	delReq := httptest.NewRequest(http.MethodDelete, "/api/v1/tasks/"+taskID, nil)
	delReq.Header.Set("Authorization", token)
	dw := httptest.NewRecorder()
	router.ServeHTTP(dw, delReq)

	if dw.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for delete, got %d", dw.Code)
	}
}

func TestLinguisticsHandler_ImportAndStats(t *testing.T) {
	router, tokenMgr, user := setupTestServer()
	token := authHeader(t, tokenMgr, user)

	// 1. Import JSON
	jsonData := []byte(`[
		{"word": "桜", "reading": "さくら", "meaning": "Cherry blossom", "language": "JP"},
		{"word": "山", "reading": "やま", "meaning": "Mountain", "language": "JP"}
	]`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/vocabularies/import?format=json", bytes.NewReader(jsonData))
	req.Header.Set("Authorization", token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for import, got %d: %s", w.Code, w.Body.String())
	}

	// 2. Stats
	statsReq := httptest.NewRequest(http.MethodGet, "/api/v1/vocabularies/stats", nil)
	statsReq.Header.Set("Authorization", token)
	sw := httptest.NewRecorder()
	router.ServeHTTP(sw, statsReq)

	if sw.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for stats, got %d", sw.Code)
	}
}

func TestPomodoroHandler_History(t *testing.T) {
	router, tokenMgr, user := setupTestServer()
	token := authHeader(t, tokenMgr, user)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pomodoro/history", nil)
	req.Header.Set("Authorization", token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for pomodoro history, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAuthHandler_UpdateProfile(t *testing.T) {
	router, tokenMgr, user := setupTestServer()
	token := authHeader(t, tokenMgr, user)

	patchBody := []byte(`{"display_name": "Antigravity Pro"}`)
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/me", bytes.NewReader(patchBody))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for update profile, got %d: %s", w.Code, w.Body.String())
	}

	var env response.Envelope
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("parse envelope: %v", err)
	}
	dataMap := env.Data.(map[string]any)
	if dataMap["display_name"] != "Antigravity Pro" {
		t.Fatalf("expected display name 'Antigravity Pro', got %v", dataMap["display_name"])
	}
}

func TestAuthenticate_RejectsQueryTokenOnRESTEndpoint(t *testing.T) {
	router, tokenMgr, user := setupTestServer()
	token := authHeader(t, tokenMgr, user)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil)
	query := request.URL.Query()
	query.Set("access_token", token)
	request.URL.RawQuery = query.Encode()
	responseRecorder := httptest.NewRecorder()

	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected query token to be rejected with 401, got %d", responseRecorder.Code)
	}
}

func TestBindJSON_RejectsUnknownFields(t *testing.T) {
	router, tokenMgr, user := setupTestServer()
	token := authHeader(t, tokenMgr, user)
	body := []byte(`{"title":"Known field","unexpected":"must fail"}`)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewReader(body))
	request.Header.Set("Authorization", token)
	request.Header.Set("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()

	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Fatalf("expected unknown JSON field to return 400, got %d: %s", responseRecorder.Code, responseRecorder.Body.String())
	}
}

func TestBindJSON_RejectsMultipleDocuments(t *testing.T) {
	router, tokenMgr, user := setupTestServer()
	token := authHeader(t, tokenMgr, user)
	body := []byte(`{"title":"First"}{"title":"Second"}`)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewReader(body))
	request.Header.Set("Authorization", token)
	request.Header.Set("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()

	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Fatalf("expected multiple JSON documents to return 400, got %d: %s", responseRecorder.Code, responseRecorder.Body.String())
	}
}

func TestFlexibleDate_UnmarshalJSON(t *testing.T) {
	type testPayload struct {
		Date *FlexibleDate `json:"date"`
	}

	// 1. HTML standard date input format (YYYY-MM-DD)
	input1 := []byte(`{"date":"2026-09-13"}`)
	var p1 testPayload
	if err := json.Unmarshal(input1, &p1); err != nil {
		t.Fatalf("failed to unmarshal YYYY-MM-DD date: %v", err)
	}
	if p1.Date == nil || p1.Date.Time().Year() != 2026 || p1.Date.Time().Month() != 9 || p1.Date.Time().Day() != 13 {
		t.Fatalf("unexpected date parsed: %v", p1.Date)
	}

	// 2. ISO/RFC3339 format
	input2 := []byte(`{"date":"2026-09-13T12:00:00Z"}`)
	var p2 testPayload
	if err := json.Unmarshal(input2, &p2); err != nil {
		t.Fatalf("failed to unmarshal RFC3339 date: %v", err)
	}
	if p2.Date == nil || p2.Date.Time().Hour() != 12 {
		t.Fatalf("unexpected hour parsed: %v", p2.Date)
	}

	// 3. Null date
	input3 := []byte(`{"date":null}`)
	var p3 testPayload
	if err := json.Unmarshal(input3, &p3); err != nil {
		t.Fatalf("failed to unmarshal null date: %v", err)
	}
	if p3.Date.Time() != nil {
		t.Fatalf("expected nil for null date, got %v", p3.Date)
	}
}
