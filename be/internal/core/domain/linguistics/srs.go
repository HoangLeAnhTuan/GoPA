package linguistics

import (
	"math"
	"time"
)

// DefaultEaseFactor is the standard starting ease factor in the SM-2 algorithm.
const DefaultEaseFactor = 2.5

// MinEaseFactor is the minimum lower bound for the ease factor in SM-2.
const MinEaseFactor = 1.3

// VocabularyState represents the persistent review parameters of a vocabulary item.
type VocabularyState struct {
	CurrentBox      int16
	EaseFactor      float64
	IntervalDays    int
	RepetitionCount int
	MasteryScore    int16
}

// SRSResult contains the calculated next state after evaluating a review event.
type SRSResult struct {
	NextBox         int16
	EaseFactor      float64
	IntervalDays    int
	RepetitionCount int
	NextReviewAt    time.Time
	MasteryScore    int16
	MasteryDelta    int16
}

// CalculateNextReview calculates the next SuperMemo-2 (SM-2) scheduling state.
// It is a pure domain function with no I/O dependencies.
func CalculateNextReview(current VocabularyState, quality int16, now time.Time) SRSResult {
	// Normalize quality score into valid SM-2 range [0, 5]
	if quality < 0 {
		quality = 0
	} else if quality > 5 {
		quality = 5
	}

	easeFactor := current.EaseFactor
	if easeFactor < MinEaseFactor {
		easeFactor = DefaultEaseFactor
	}

	currentBox := current.CurrentBox
	if currentBox < 1 {
		currentBox = 1
	} else if currentBox > 5 {
		currentBox = 5
	}

	var (
		nextBox         int16
		nextInterval    int
		nextRepetitions int
		masteryDelta    int16
	)

	if quality >= 3 {
		switch current.RepetitionCount {
		case 0:
			nextInterval = 1
		case 1:
			nextInterval = 6
		default:
			prevInterval := current.IntervalDays
			if prevInterval < 1 {
				prevInterval = 1
			}
			nextInterval = int(math.Round(float64(prevInterval) * easeFactor))
			if nextInterval < 1 {
				nextInterval = 1
			}
		}

		nextRepetitions = current.RepetitionCount + 1
		nextBox = currentBox + 1
		if nextBox > 5 {
			nextBox = 5
		}

		// Mastery score reward scaled by quality
		masteryDelta = (quality - 2) * 5
	} else {
		nextRepetitions = 0
		nextInterval = 1
		nextBox = 1
		masteryDelta = -10
	}

	// Update Ease Factor: EF' = EF + (0.1 - (5 - q) * (0.08 + (5 - q) * 0.02))
	diff := float64(5 - quality)
	easeFactor = easeFactor + (0.1 - diff*(0.08+diff*0.02))
	if easeFactor < MinEaseFactor {
		easeFactor = MinEaseFactor
	}
	// Round ease factor to 2 decimal places
	easeFactor = math.Round(easeFactor*100) / 100

	newMastery := int(current.MasteryScore) + int(masteryDelta)
	if newMastery < 0 {
		newMastery = 0
	} else if newMastery > 100 {
		newMastery = 100
	}

	nextReviewAt := now.UTC().AddDate(0, 0, nextInterval)

	return SRSResult{
		NextBox:         nextBox,
		EaseFactor:      easeFactor,
		IntervalDays:    nextInterval,
		RepetitionCount: nextRepetitions,
		NextReviewAt:    nextReviewAt,
		MasteryScore:    int16(newMastery),
		MasteryDelta:    masteryDelta,
	}
}
