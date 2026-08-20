package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"gopa/internal/core/domain"
	"gopa/internal/services"
)

type mockJournalRepo struct {
	journals map[uuid.UUID]domain.Journal
	links    map[string]bool
}

func newMockJournalRepo() *mockJournalRepo {
	return &mockJournalRepo{
		journals: make(map[uuid.UUID]domain.Journal),
		links:    make(map[string]bool),
	}
}

func (m *mockJournalRepo) List(ctx context.Context, userID uuid.UUID, filter domain.JournalFilter) ([]domain.Journal, error) {
	res := make([]domain.Journal, 0, len(m.journals))
	for _, j := range m.journals {
		if j.UserID == userID {
			res = append(res, j)
		}
	}
	return res, nil
}

func (m *mockJournalRepo) Get(ctx context.Context, userID, id uuid.UUID) (domain.Journal, error) {
	j, ok := m.journals[id]
	if !ok || j.UserID != userID {
		return domain.Journal{}, domain.ErrNotFound
	}
	return j, nil
}

func (m *mockJournalRepo) Create(ctx context.Context, j domain.Journal) (domain.Journal, error) {
	m.journals[j.ID] = j
	return j, nil
}

func (m *mockJournalRepo) Update(ctx context.Context, j domain.Journal) (domain.Journal, error) {
	if _, ok := m.journals[j.ID]; !ok {
		return domain.Journal{}, domain.ErrNotFound
	}
	m.journals[j.ID] = j
	return j, nil
}

func (m *mockJournalRepo) Delete(ctx context.Context, userID, id uuid.UUID) error {
	if j, ok := m.journals[id]; !ok || j.UserID != userID {
		return domain.ErrNotFound
	}
	delete(m.journals, id)
	return nil
}

func (m *mockJournalRepo) Link(ctx context.Context, userID, journalID, linkedID uuid.UUID) error {
	j1, ok1 := m.journals[journalID]
	j2, ok2 := m.journals[linkedID]
	if !ok1 || !ok2 || j1.UserID != userID || j2.UserID != userID {
		return domain.ErrNotFound
	}
	m.links[journalID.String()+":"+linkedID.String()] = true
	m.links[linkedID.String()+":"+journalID.String()] = true
	return nil
}

func (m *mockJournalRepo) Unlink(ctx context.Context, userID, journalID, linkedID uuid.UUID) error {
	delete(m.links, journalID.String()+":"+linkedID.String())
	delete(m.links, linkedID.String()+":"+journalID.String())
	return nil
}

func (m *mockJournalRepo) Stats(ctx context.Context, userID uuid.UUID) (domain.JournalStats, error) {
	dates := make([]time.Time, 0)
	for _, j := range m.journals {
		if j.UserID == userID {
			dates = append(dates, j.PublishedDate)
		}
	}
	curr, long := domain.CalculateJournalStreaks(dates, time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC))
	return domain.JournalStats{CurrentStreak: curr, LongestStreak: long, Entries: len(dates)}, nil
}
func (m *mockJournalRepo) Search(ctx context.Context, userID uuid.UUID, query string) ([]domain.Journal, error) {
	results := make([]domain.Journal, 0)
	for _, j := range m.journals {
		if j.UserID == userID {
			results = append(results, j)
		}
	}
	return results, nil
}

func TestJournalService_Create_Validation(t *testing.T) {
	repo := newMockJournalRepo()
	svc := services.NewJournalService(repo)
	userID := uuid.New()

	// Empty title
	_, err := svc.Create(context.Background(), userID, services.JournalInput{Title: ""})
	if err != domain.ErrValidation {
		t.Errorf("expected ErrValidation for empty title, got %v", err)
	}

	// Invalid mood
	invalidMood := domain.JournalMood("SUPER_HAPPY")
	_, err = svc.Create(context.Background(), userID, services.JournalInput{
		Title: "Good day",
		Mood:  &invalidMood,
	})
	if err != domain.ErrValidation {
		t.Errorf("expected ErrValidation for invalid mood, got %v", err)
	}

	// Valid creation
	validMood := domain.JournalMoodPositive
	j, err := svc.Create(context.Background(), userID, services.JournalInput{
		Title: "Good day",
		Mood:  &validMood,
		Tags:  []string{"Tech", "Go"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if j.Title != "Good day" || len(j.Tags) != 2 || j.Tags[0] != "tech" {
		t.Errorf("unexpected journal state: %+v", j)
	}
}

func TestJournalService_Link_SelfLinkRejection(t *testing.T) {
	repo := newMockJournalRepo()
	svc := services.NewJournalService(repo)
	userID := uuid.New()
	id := uuid.New()

	err := svc.Link(context.Background(), userID, id, id)
	if err != domain.ErrValidation {
		t.Errorf("expected ErrValidation when linking to self, got %v", err)
	}
}

func TestCalculateJournalStreaks(t *testing.T) {
	today := time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC)

	// Continuous streak for past 3 days (today, yesterday, day before)
	dates := []time.Time{
		today,
		today.AddDate(0, 0, -1),
		today.AddDate(0, 0, -2),
		today.AddDate(0, 0, -5),
		today.AddDate(0, 0, -6),
		today.AddDate(0, 0, -7),
		today.AddDate(0, 0, -8),
	}

	curr, long := domain.CalculateJournalStreaks(dates, today)
	if curr != 3 {
		t.Errorf("expected current streak 3, got %d", curr)
	}
	if long != 4 {
		t.Errorf("expected longest streak 4, got %d", long)
	}
}
