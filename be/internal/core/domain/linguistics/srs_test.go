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

	t.Run("Table-driven SM-2 quality score progression and mastery delta", func(t *testing.T) {
		tests := []struct {
			name            string
			initial         VocabularyState
			quality         int16
			expectedBox     int16
			expectedReps    int
			expectedInt     int
			expectedDelta   int16
			expectedMastery int16
		}{
			{
				name:            "Quality 0 failure",
				initial:         VocabularyState{CurrentBox: 3, EaseFactor: 2.5, IntervalDays: 10, RepetitionCount: 2, MasteryScore: 40},
				quality:         0,
				expectedBox:     1,
				expectedReps:    0,
				expectedInt:     1,
				expectedDelta:   -10,
				expectedMastery: 30,
			},
			{
				name:            "Quality 1 failure",
				initial:         VocabularyState{CurrentBox: 2, EaseFactor: 2.5, IntervalDays: 6, RepetitionCount: 1, MasteryScore: 5},
				quality:         1,
				expectedBox:     1,
				expectedReps:    0,
				expectedInt:     1,
				expectedDelta:   -10,
				expectedMastery: 0, // Clamped to 0
			},
			{
				name:            "Quality 2 failure boundary",
				initial:         VocabularyState{CurrentBox: 4, EaseFactor: 2.5, IntervalDays: 15, RepetitionCount: 3, MasteryScore: 60},
				quality:         2,
				expectedBox:     1,
				expectedReps:    0,
				expectedInt:     1,
				expectedDelta:   -10,
				expectedMastery: 50,
			},
			{
				name:            "Quality 3 pass boundary",
				initial:         VocabularyState{CurrentBox: 1, EaseFactor: 2.5, IntervalDays: 0, RepetitionCount: 0, MasteryScore: 0},
				quality:         3,
				expectedBox:     2,
				expectedReps:    1,
				expectedInt:     1,
				expectedDelta:   5, // (3-2)*5
				expectedMastery: 5,
			},
			{
				name:            "Quality 4 strong pass",
				initial:         VocabularyState{CurrentBox: 2, EaseFactor: 2.5, IntervalDays: 1, RepetitionCount: 1, MasteryScore: 20},
				quality:         4,
				expectedBox:     3,
				expectedReps:    2,
				expectedInt:     6,
				expectedDelta:   10, // (4-2)*5
				expectedMastery: 30,
			},
			{
				name:            "Quality 5 perfect score at max box",
				initial:         VocabularyState{CurrentBox: 5, EaseFactor: 2.5, IntervalDays: 30, RepetitionCount: 5, MasteryScore: 95},
				quality:         5,
				expectedBox:     5, // Max box is 5
				expectedReps:    6,
				expectedInt:     75,  // round(30 * 2.5) = 75 (before EF recalculation applied)
				expectedDelta:   15,  // (5-2)*5
				expectedMastery: 100, // Clamped to 100
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				res := CalculateNextReview(tc.initial, tc.quality, fixedNow)
				assert.Equal(t, tc.expectedBox, res.NextBox, "NextBox")
				assert.Equal(t, tc.expectedReps, res.RepetitionCount, "RepetitionCount")
				assert.Equal(t, tc.expectedInt, res.IntervalDays, "IntervalDays")
				assert.Equal(t, tc.expectedDelta, res.MasteryDelta, "MasteryDelta")
				assert.Equal(t, tc.expectedMastery, res.MasteryScore, "MasteryScore")
				assert.Equal(t, fixedNow.AddDate(0, 0, tc.expectedInt), res.NextReviewAt, "NextReviewAt")
			})
		}
	})

	t.Run("Input normalization and boundary clamping", func(t *testing.T) {
		// Quality < 0 normalized to 0
		resUnder := CalculateNextReview(VocabularyState{CurrentBox: 2, EaseFactor: 2.5}, -5, fixedNow)
		assert.Equal(t, int16(1), resUnder.NextBox)
		assert.Equal(t, int16(-10), resUnder.MasteryDelta)

		// Quality > 5 normalized to 5
		resOver := CalculateNextReview(VocabularyState{CurrentBox: 1, EaseFactor: 2.5}, 10, fixedNow)
		assert.Equal(t, int16(2), resOver.NextBox)
		assert.Equal(t, int16(15), resOver.MasteryDelta)

		// EaseFactor < 1.3 resets to DefaultEaseFactor (2.5)
		resInvalidEF := CalculateNextReview(VocabularyState{CurrentBox: 1, EaseFactor: 0.8}, 4, fixedNow)
		assert.Equal(t, 2.5, resInvalidEF.EaseFactor)

		// CurrentBox < 1 normalized to 1, CurrentBox > 5 normalized to 5
		resLowBox := CalculateNextReview(VocabularyState{CurrentBox: 0, EaseFactor: 2.5}, 4, fixedNow)
		assert.Equal(t, int16(2), resLowBox.NextBox)

		resHighBox := CalculateNextReview(VocabularyState{CurrentBox: 8, EaseFactor: 2.5}, 4, fixedNow)
		assert.Equal(t, int16(5), resHighBox.NextBox)

		// Repetition >= 2 with interval < 1 uses fallback interval = 1
		resSmallInterval := CalculateNextReview(VocabularyState{CurrentBox: 2, EaseFactor: 2.5, IntervalDays: 0, RepetitionCount: 2}, 4, fixedNow)
		assert.Equal(t, 3, resSmallInterval.IntervalDays) // round(1 * 2.5) = 3 (due to round(2.5) in IEEE 754)
	})

	t.Run("Continuous failures stay clamped at MinEaseFactor", func(t *testing.T) {
		state := VocabularyState{
			CurrentBox:      3,
			EaseFactor:      1.4,
			IntervalDays:    10,
			RepetitionCount: 2,
			MasteryScore:    50,
		}

		for i := 0; i < 5; i++ {
			res := CalculateNextReview(state, 0, fixedNow)
			assert.GreaterOrEqual(t, res.EaseFactor, MinEaseFactor)
			state.EaseFactor = res.EaseFactor
			state.CurrentBox = res.NextBox
			state.IntervalDays = res.IntervalDays
			state.RepetitionCount = res.RepetitionCount
			state.MasteryScore = res.MasteryScore
		}

		assert.Equal(t, MinEaseFactor, state.EaseFactor)
		assert.Equal(t, int16(1), state.CurrentBox)
		assert.Equal(t, 0, state.RepetitionCount)
		assert.Equal(t, int16(0), state.MasteryScore)
	})
}
