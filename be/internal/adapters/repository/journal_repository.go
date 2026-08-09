package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gopa/internal/core/domain"
	"gopa/pkg/database"
)

type JournalRepository struct{ db *sqlx.DB }
type journalRow struct {
	ID            uuid.UUID           `db:"id"`
	UserID        uuid.UUID           `db:"user_id"`
	Title         string              `db:"title"`
	Content       string              `db:"content"`
	Tags          []byte              `db:"tags"`
	Mood          *domain.JournalMood `db:"mood"`
	PublishedDate time.Time           `db:"published_date"`
	CreatedAt     time.Time           `db:"created_at"`
	UpdatedAt     time.Time           `db:"updated_at"`
}

const journalColumns = "id, user_id, title, content, tags, mood, published_date, created_at, updated_at"

func NewJournalRepository(db *sqlx.DB) *JournalRepository { return &JournalRepository{db: db} }

func (r *JournalRepository) List(ctx context.Context, userID uuid.UUID, filter domain.JournalFilter) ([]domain.Journal, error) {
	where := []string{"user_id = $1"}
	args := []any{userID}
	if filter.Term != "" {
		args = append(args, "%"+filter.Term+"%")
		where = append(where, fmt.Sprintf("(title ILIKE $%[1]d OR content ILIKE $%[1]d OR tags::text ILIKE $%[1]d)", len(args)))
	}
	if filter.Tag != "" {
		args = append(args, filter.Tag)
		where = append(where, fmt.Sprintf("tags @> jsonb_build_array($%d::text)", len(args)))
	}
	if filter.Mood != nil {
		args = append(args, *filter.Mood)
		where = append(where, fmt.Sprintf("mood = $%d", len(args)))
	}
	if filter.From != nil {
		args = append(args, filter.From.UTC())
		where = append(where, fmt.Sprintf("published_date >= $%d::date", len(args)))
	}
	if filter.To != nil {
		args = append(args, filter.To.UTC())
		where = append(where, fmt.Sprintf("published_date <= $%d::date", len(args)))
	}
	args = append(args, filter.Limit)
	query := fmt.Sprintf(`SELECT %s FROM journals WHERE %s ORDER BY published_date DESC, id DESC LIMIT $%d`, journalColumns, strings.Join(where, " AND "), len(args))
	var rows []journalRow
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("list journals: %w", err)
	}
	return journalsFromRows(rows)
}

func (r *JournalRepository) Get(ctx context.Context, userID, id uuid.UUID) (domain.Journal, error) {
	var row journalRow
	if err := r.db.GetContext(ctx, &row, `SELECT `+journalColumns+` FROM journals WHERE id = $1 AND user_id = $2`, id, userID); err != nil {
		return domain.Journal{}, mapNotFound(err, "get journal")
	}
	return journalFromRow(row)
}
func (r *JournalRepository) Create(ctx context.Context, journal domain.Journal) (domain.Journal, error) {
	if err := r.write(ctx, journal, "journal.created", false); err != nil {
		return domain.Journal{}, err
	}
	return journal, nil
}
func (r *JournalRepository) Update(ctx context.Context, journal domain.Journal) (domain.Journal, error) {
	if err := r.write(ctx, journal, "journal.updated", true); err != nil {
		return domain.Journal{}, err
	}
	return journal, nil
}
func (r *JournalRepository) write(ctx context.Context, j domain.Journal, eventType string, update bool) error {
	tags, err := json.Marshal(j.Tags)
	if err != nil {
		return fmt.Errorf("marshal journal tags: %w", err)
	}
	return database.WithTx(ctx, r.db, func(tx *sqlx.Tx) error {
		var result sql.Result
		if update {
			result, err = tx.ExecContext(ctx, `UPDATE journals SET title=$1, content=$2, tags=$3, mood=$4, published_date=$5, updated_at=$6 WHERE id=$7 AND user_id=$8`, j.Title, j.Content, tags, j.Mood, j.PublishedDate, j.UpdatedAt, j.ID, j.UserID)
		} else {
			result, err = tx.ExecContext(ctx, `INSERT INTO journals (id,user_id,title,content,tags,mood,published_date,search_vector,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,to_tsvector('english',''),$8,$8)`, j.ID, j.UserID, j.Title, j.Content, tags, j.Mood, j.PublishedDate, j.CreatedAt)
		}
		if err != nil {
			return fmt.Errorf("write journal: %w", err)
		}
		if update {
			affected, rowsErr := result.RowsAffected()
			if rowsErr != nil {
				return fmt.Errorf("count updated journal: %w", rowsErr)
			}
			if affected == 0 {
				return domain.ErrNotFound
			}
		}
		payload, err := json.Marshal(map[string]uuid.UUID{"journal_id": j.ID, "user_id": j.UserID})
		if err != nil {
			return fmt.Errorf("marshal journal event: %w", err)
		}
		if _, err = tx.ExecContext(ctx, `INSERT INTO outbox_events (id,event_type,routing_key,payload,occurred_at,created_at) VALUES ($1,$2,$2,$3,$4,$4)`, uuid.New(), eventType, payload, j.UpdatedAt); err != nil {
			return fmt.Errorf("create journal outbox event: %w", err)
		}
		return nil
	})
}
func (r *JournalRepository) Delete(ctx context.Context, userID, id uuid.UUID) error {
	return database.WithTx(ctx, r.db, func(tx *sqlx.Tx) error {
		result, err := tx.ExecContext(ctx, `DELETE FROM journals WHERE id=$1 AND user_id=$2`, id, userID)
		if err != nil {
			return fmt.Errorf("delete journal: %w", err)
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("count deleted journal: %w", err)
		}
		if affected == 0 {
			return domain.ErrNotFound
		}
		payload, err := json.Marshal(map[string]uuid.UUID{"journal_id": id, "user_id": userID})
		if err != nil {
			return fmt.Errorf("marshal journal delete event: %w", err)
		}
		_, err = tx.ExecContext(ctx, `INSERT INTO outbox_events(id,event_type,routing_key,payload,occurred_at,created_at) VALUES($1,'journal.deleted','journal.deleted',$2,now(),now())`, uuid.New(), payload)
		if err != nil {
			return fmt.Errorf("create journal delete event: %w", err)
		}
		return nil
	})
}
func (r *JournalRepository) Link(ctx context.Context, userID, journalID, linkedID uuid.UUID) error {
	return r.modifyLink(ctx, userID, journalID, linkedID, true)
}
func (r *JournalRepository) Unlink(ctx context.Context, userID, journalID, linkedID uuid.UUID) error {
	return r.modifyLink(ctx, userID, journalID, linkedID, false)
}
func (r *JournalRepository) modifyLink(ctx context.Context, userID, journalID, linkedID uuid.UUID, create bool) error {
	return database.WithTx(ctx, r.db, func(tx *sqlx.Tx) error {
		var count int
		if err := tx.GetContext(ctx, &count, `SELECT count(*) FROM journals WHERE user_id=$1 AND id IN ($2,$3)`, userID, journalID, linkedID); err != nil {
			return fmt.Errorf("verify journal link ownership: %w", err)
		}
		if count != 2 {
			return domain.ErrNotFound
		}
		if create {
			_, err := tx.ExecContext(ctx, `INSERT INTO journal_links(journal_id,linked_id,user_id) VALUES($1,$2,$3),($2,$1,$3) ON CONFLICT DO NOTHING`, journalID, linkedID, userID)
			if err != nil {
				return fmt.Errorf("insert journal link: %w", err)
			}
			return nil
		}
		_, err := tx.ExecContext(ctx, `DELETE FROM journal_links WHERE user_id=$1 AND ((journal_id=$2 AND linked_id=$3) OR (journal_id=$3 AND linked_id=$2))`, userID, journalID, linkedID)
		if err != nil {
			return fmt.Errorf("delete journal link: %w", err)
		}
		return nil
	})
}
func (r *JournalRepository) Stats(ctx context.Context, userID uuid.UUID) (domain.JournalStats, error) {
	var rows []struct {
		PublishedDate time.Time           `db:"published_date"`
		Content       string              `db:"content"`
		Mood          *domain.JournalMood `db:"mood"`
	}
	if err := r.db.SelectContext(ctx, &rows, `SELECT published_date,content,mood FROM journals WHERE user_id=$1`, userID); err != nil {
		return domain.JournalStats{}, fmt.Errorf("list journal stats: %w", err)
	}
	stats := domain.JournalStats{Entries: len(rows), MoodCounts: map[domain.JournalMood]int{}}
	dates := make([]time.Time, 0, len(rows))
	for _, row := range rows {
		dates = append(dates, row.PublishedDate)
		stats.WordCount += len(strings.Fields(row.Content))
		if row.Mood != nil {
			stats.MoodCounts[*row.Mood]++
		}
	}
	stats.CurrentStreak, stats.LongestStreak = domain.CalculateJournalStreaks(dates, time.Now().UTC())
	return stats, nil
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
	return domain.Journal{ID: row.ID, UserID: row.UserID, Title: row.Title, Content: row.Content, Tags: tags, Mood: row.Mood, PublishedDate: row.PublishedDate, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt}, nil
}
