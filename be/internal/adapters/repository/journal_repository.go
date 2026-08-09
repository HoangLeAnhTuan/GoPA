package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gopa/internal/core/domain"
)

type JournalRepository struct{ db *sqlx.DB }
type journalRow struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	Tags      []byte    `db:"tags"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewJournalRepository(db *sqlx.DB) *JournalRepository { return &JournalRepository{db: db} }
func (r *JournalRepository) List(ctx context.Context, userID uuid.UUID, term, tag string, limit int) ([]domain.Journal, error) {
	where := []string{"user_id = $1"}
	args := []any{userID}
	if strings.TrimSpace(term) != "" {
		args = append(args, term)
		where = append(where, fmt.Sprintf("search_vector @@ websearch_to_tsquery('simple', $%d)", len(args)))
	}
	if strings.TrimSpace(tag) != "" {
		args = append(args, tag)
		where = append(where, fmt.Sprintf("tags @> jsonb_build_array($%d::text)", len(args)))
	}
	args = append(args, limit)
	query := fmt.Sprintf(`SELECT id, user_id, title, content, tags, created_at, updated_at FROM journals WHERE %s ORDER BY updated_at DESC LIMIT $%d`, strings.Join(where, " AND "), len(args))
	var rows []journalRow
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("list journals: %w", err)
	}
	return journalsFromRows(rows)
}
func (r *JournalRepository) Get(ctx context.Context, userID, id uuid.UUID) (domain.Journal, error) {
	var row journalRow
	if err := r.db.GetContext(ctx, &row, `SELECT id, user_id, title, content, tags, created_at, updated_at FROM journals WHERE id = $1 AND user_id = $2`, id, userID); err != nil {
		return domain.Journal{}, mapNotFound(err, "get journal")
	}
	return journalFromRow(row)
}
func (r *JournalRepository) Create(ctx context.Context, journal domain.Journal) (domain.Journal, error) {
	tags, err := json.Marshal(journal.Tags)
	if err != nil {
		return domain.Journal{}, fmt.Errorf("marshal journal tags: %w", err)
	}
	if _, err := r.db.ExecContext(ctx, `INSERT INTO journals (id, user_id, title, content, tags, search_vector, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, to_tsvector('simple', $3 || ' ' || $4), $6, $6)`, journal.ID, journal.UserID, journal.Title, journal.Content, tags, journal.CreatedAt); err != nil {
		return domain.Journal{}, fmt.Errorf("create journal: %w", err)
	}
	return journal, nil
}
func (r *JournalRepository) Update(ctx context.Context, journal domain.Journal) (domain.Journal, error) {
	tags, err := json.Marshal(journal.Tags)
	if err != nil {
		return domain.Journal{}, fmt.Errorf("marshal journal tags: %w", err)
	}
	result, err := r.db.ExecContext(ctx, `UPDATE journals SET title = $1, content = $2, tags = $3, search_vector = to_tsvector('simple', $1 || ' ' || $2), updated_at = $4 WHERE id = $5 AND user_id = $6`, journal.Title, journal.Content, tags, journal.UpdatedAt, journal.ID, journal.UserID)
	if err != nil {
		return domain.Journal{}, fmt.Errorf("update journal: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return domain.Journal{}, fmt.Errorf("count updated journal: %w", err)
	}
	if rows == 0 {
		return domain.Journal{}, domain.ErrNotFound
	}
	return journal, nil
}
func (r *JournalRepository) Delete(ctx context.Context, userID, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM journals WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("delete journal: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("count deleted journal: %w", err)
	}
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}
func journalsFromRows(rows []journalRow) ([]domain.Journal, error) {
	journals := make([]domain.Journal, 0, len(rows))
	for _, row := range rows {
		journal, err := journalFromRow(row)
		if err != nil {
			return nil, err
		}
		journals = append(journals, journal)
	}
	return journals, nil
}
func journalFromRow(row journalRow) (domain.Journal, error) {
	var tags []string
	if err := json.Unmarshal(row.Tags, &tags); err != nil {
		return domain.Journal{}, fmt.Errorf("unmarshal journal tags: %w", err)
	}
	return domain.Journal{ID: row.ID, UserID: row.UserID, Title: row.Title, Content: row.Content, Tags: tags, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt}, nil
}
