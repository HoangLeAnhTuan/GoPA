package domain

import (
	"testing"
	"time"
)

func TestScheduleReview_ResetsToBoxOneWhenQualityIsBelowThree(t *testing.T) {
	now := time.Date(2026, 8, 9, 8, 0, 0, 0, time.UTC)
	schedule, err := ScheduleReview(4, 2, now)
	if err != nil {
		t.Fatalf("schedule review: %v", err)
	}
	if schedule.Box != 1 || !schedule.NextReviewAt.Equal(now.Add(10*time.Minute)) {
		t.Fatalf("unexpected failed-review schedule: %#v", schedule)
	}
}

func TestCalculateMaintenance_ReturnsOverdueWhenCurrentMileageExceedsThreshold(t *testing.T) {
	lastMaintenance := 10_000
	status := CalculateMaintenance(15_200, &lastMaintenance)
	if status.Status != MaintenanceOverdue || status.RemainingDistance != -200 {
		t.Fatalf("unexpected maintenance status: %#v", status)
	}
}
