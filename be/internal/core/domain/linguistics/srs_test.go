package linguistics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalculateNextReview(t *testing.T) {
	fixedNow := time.Date(2026, 8, 20, 12, 0, 0, 0, time.UTC)

	t.Run("First successful review (Repetition=0, Quality=4)", func(t *testing.T) {
		initial := VocabularyState{
			CurrentBox:      1,
			EaseFactor:      DefaultEaseFactor,
			IntervalDays:    0,
			RepetitionCount: 0,
			MasteryScore:    0,
		}

		res := CalculateNextReview(initial, 4, fixedNow)

		assert.Equal(t, int16(2), res.NextBox)
		assert.Equal(t, 1, res.IntervalDays)
		assert.Equal(t, 1, res.RepetitionCount)
		assert.Equal(t, fixedNow.AddDate(0, 0, 1), res.NextReviewAt)
		assert.Equal(t, int16(10), res.MasteryScore) // (4-2)*5 = 10
		assert.Equal(t, int16(10), res.MasteryDelta)
		assert.Equal(t, 2.5, res.EaseFactor) // 2.5 + (0.1 - 1*(0.08+0.02)) = 2.5
	})

	t.Run("Second successful review (Repetition=1, Quality=5)", func(t *testing.T) {
		initial := VocabularyState{
			CurrentBox:      2,
			EaseFactor:      2.5,
			IntervalDays:    1,
			RepetitionCount: 1,
			MasteryScore:    10,
		}

		res := CalculateNextReview(initial, 5, fixedNow)

		assert.Equal(t, int16(3), res.NextBox)
		assert.Equal(t, 6, res.IntervalDays)
		assert.Equal(t, 2, res.RepetitionCount)
		assert.Equal(t, fixedNow.AddDate(0, 0, 6), res.NextReviewAt)
		assert.Equal(t, int16(25), res.MasteryScore) // 10 + 15 = 25
		assert.Equal(t, 2.6, res.EaseFactor)         // 2.5 + 0.1 = 2.6
	})

	t.Run("Third successful review (Repetition=2, Interval=6, Quality=4)", func(t *testing.T) {
		initial := VocabularyState{
			CurrentBox:      3,
			EaseFactor:      2.6,
			IntervalDays:    6,
			RepetitionCount: 2,
			MasteryScore:    25,
		}

		res := CalculateNextReview(initial, 4, fixedNow)

		// Interval = round(6 * 2.6) = round(15.6) = 16
		assert.Equal(t, 16, res.IntervalDays)
		assert.Equal(t, 3, res.RepetitionCount)
		assert.Equal(t, int16(4), res.NextBox)
	})

	t.Run("Failed review resets repetitions and box (Quality=1)", func(t *testing.T) {
		initial := VocabularyState{
			CurrentBox:      4,
			EaseFactor:      2.4,
			IntervalDays:    15,
			RepetitionCount: 3,
			MasteryScore:    50,
		}

		res := CalculateNextReview(initial, 1, fixedNow)

		assert.Equal(t, int16(1), res.NextBox)
		assert.Equal(t, 1, res.IntervalDays)
		assert.Equal(t, 0, res.RepetitionCount)
		assert.Equal(t, int16(40), res.MasteryScore) // 50 - 10 = 40
		assert.Equal(t, int16(-10), res.MasteryDelta)
		// Ease factor decreases but never drops below MinEaseFactor (1.3)
		assert.True(t, res.EaseFactor < 2.4)
		assert.True(t, res.EaseFactor >= MinEaseFactor)
	})

	t.Run("EaseFactor never drops below MinEaseFactor (1.3)", func(t *testing.T) {
		initial := VocabularyState{
			CurrentBox:      1,
			EaseFactor:      1.3,
			IntervalDays:    1,
			RepetitionCount: 0,
			MasteryScore:    5,
		}

		// Quality 0 creates significant ease factor penalty
		res := CalculateNextReview(initial, 0, fixedNow)

		assert.Equal(t, MinEaseFactor, res.EaseFactor)
		assert.Equal(t, int16(0), res.MasteryScore) // clamped at 0
	})
}
