package domain

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TaskStatus string
type TaskPriority string
type TaskCategory string

const (
	TaskStatusTodo       TaskStatus   = "TODO"
	TaskStatusInProgress TaskStatus   = "IN_PROGRESS"
	TaskStatusDone       TaskStatus   = "DONE"
	TaskPriorityLow      TaskPriority = "LOW"
	TaskPriorityMedium   TaskPriority = "MEDIUM"
	TaskPriorityHigh     TaskPriority = "HIGH"
	TaskCategoryWork     TaskCategory = "WORK"
	TaskCategoryStudy    TaskCategory = "STUDY"
	TaskCategoryLife     TaskCategory = "LIFE"
)

type Task struct {
	ID          uuid.UUID    `json:"id" db:"id"`
	UserID      uuid.UUID    `json:"user_id" db:"user_id"`
	Title       string       `json:"title" db:"title"`
	Description string       `json:"description" db:"description"`
	Status      TaskStatus   `json:"status" db:"status"`
	Priority    TaskPriority `json:"priority" db:"priority"`
	Category    TaskCategory `json:"category" db:"category"`
	DueDate     *time.Time   `json:"due_date" db:"due_date"`
	Position    int64        `json:"position" db:"position"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
}

func (t Task) Validate() error {
	if len(strings.TrimSpace(t.Title)) == 0 || len(t.Title) > 200 || !validTaskStatus(t.Status) || !validPriority(t.Priority) || !validCategory(t.Category) {
		return ErrValidation
	}
	return nil
}

func validTaskStatus(status TaskStatus) bool {
	return status == TaskStatusTodo || status == TaskStatusInProgress || status == TaskStatusDone
}
func validPriority(priority TaskPriority) bool {
	return priority == TaskPriorityLow || priority == TaskPriorityMedium || priority == TaskPriorityHigh
}
func validCategory(category TaskCategory) bool {
	return category == TaskCategoryWork || category == TaskCategoryStudy || category == TaskCategoryLife
}

type VocabularyLanguage string

const (
	VocabularyLanguageJapanese VocabularyLanguage = "JP"
	VocabularyLanguageEnglish  VocabularyLanguage = "EN"
)

type Vocabulary struct {
	ID              uuid.UUID          `json:"id" db:"id"`
	UserID          uuid.UUID          `json:"user_id" db:"user_id"`
	Language        VocabularyLanguage `json:"language" db:"language"`
	Word            string             `json:"word" db:"word"`
	Reading         *string            `json:"reading" db:"reading"`
	Meaning         string             `json:"meaning" db:"meaning"`
	ExampleSentence *string            `json:"example_sentence" db:"example_sentence"`
	CurrentBox      int16              `json:"current_box" db:"current_box"`
	NextReviewAt    time.Time          `json:"next_review_at" db:"next_review_at"`
	CreatedAt       time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" db:"updated_at"`
}

func (v Vocabulary) Validate() error {
	if (v.Language != VocabularyLanguageJapanese && v.Language != VocabularyLanguageEnglish) || len(strings.TrimSpace(v.Word)) == 0 || len(v.Word) > 300 || len(strings.TrimSpace(v.Meaning)) == 0 || len(v.Meaning) > 1000 || v.CurrentBox < 1 || v.CurrentBox > 5 {
		return ErrValidation
	}
	return nil
}

type ReviewSchedule struct {
	Box          int16
	NextReviewAt time.Time
}

func ScheduleReview(currentBox, quality int16, now time.Time) (ReviewSchedule, error) {
	if currentBox < 1 || currentBox > 5 || quality < 0 || quality > 5 {
		return ReviewSchedule{}, ErrValidation
	}
	if quality < 3 {
		return ReviewSchedule{Box: 1, NextReviewAt: now.Add(10 * time.Minute)}, nil
	}
	nextBox := currentBox + 1
	if nextBox > 5 {
		nextBox = 5
	}
	intervals := map[int16]time.Duration{1: 24 * time.Hour, 2: 3 * 24 * time.Hour, 3: 7 * 24 * time.Hour, 4: 14 * 24 * time.Hour, 5: 30 * 24 * time.Hour}
	return ReviewSchedule{Box: nextBox, NextReviewAt: now.Add(intervals[nextBox])}, nil
}

type Vehicle struct {
	ID             uuid.UUID `json:"id" db:"id"`
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
	Make           string    `json:"make" db:"make"`
	Model          string    `json:"model" db:"model"`
	Year           int16     `json:"year" db:"year"`
	Nickname       *string   `json:"nickname" db:"nickname"`
	CurrentMileage int       `json:"current_mileage" db:"current_mileage"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

func (v Vehicle) Validate() error {
	if len(strings.TrimSpace(v.Make)) == 0 || len(strings.TrimSpace(v.Model)) == 0 || v.Year < 1886 || v.Year > 2100 || v.CurrentMileage < 0 {
		return ErrValidation
	}
	return nil
}

type VehicleLogType string

const (
	VehicleLogMaintenance VehicleLogType = "MAINTENANCE"
	VehicleLogFuel        VehicleLogType = "FUEL"
	VehicleLogUpgrade     VehicleLogType = "UPGRADE"
)

type VehicleLog struct {
	ID          uuid.UUID      `json:"id" db:"id"`
	VehicleID   uuid.UUID      `json:"vehicle_id" db:"vehicle_id"`
	UserID      uuid.UUID      `json:"user_id" db:"user_id"`
	LogType     VehicleLogType `json:"log_type" db:"log_type"`
	Mileage     int            `json:"mileage" db:"mileage"`
	Description string         `json:"description" db:"description"`
	Cost        string         `json:"cost" db:"cost"`
	LogDate     time.Time      `json:"log_date" db:"log_date"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
}

func (l VehicleLog) Validate() error {
	if (l.LogType != VehicleLogMaintenance && l.LogType != VehicleLogFuel && l.LogType != VehicleLogUpgrade) || l.Mileage < 0 || len(strings.TrimSpace(l.Description)) == 0 || len(l.Description) > 2000 {
		return ErrValidation
	}
	return nil
}

type MaintenanceStatus string

const (
	MaintenanceUpcoming MaintenanceStatus = "UPCOMING"
	MaintenanceDue      MaintenanceStatus = "DUE"
	MaintenanceOverdue  MaintenanceStatus = "OVERDUE"
)

type MaintenanceResult struct {
	Status            MaintenanceStatus `json:"status"`
	RemainingDistance int               `json:"remaining_distance"`
	NextMileage       int               `json:"next_mileage"`
}

func CalculateMaintenance(currentMileage int, lastMaintenanceMileage *int) MaintenanceResult {
	base := 0
	if lastMaintenanceMileage != nil {
		base = *lastMaintenanceMileage
	}
	next := base + 5000
	remaining := next - currentMileage
	status := MaintenanceUpcoming
	if remaining == 0 {
		status = MaintenanceDue
	}
	if remaining < 0 {
		status = MaintenanceOverdue
	}
	return MaintenanceResult{Status: status, RemainingDistance: remaining, NextMileage: next}
}

type ConnectionType string

const (
	ConnectionLAN  ConnectionType = "LAN"
	ConnectionWiFi ConnectionType = "WIFI"
)

type NetworkNode struct {
	ID             uuid.UUID      `json:"id" db:"id"`
	UserID         uuid.UUID      `json:"user_id" db:"user_id"`
	DeviceName     string         `json:"device_name" db:"device_name"`
	MACAddress     *string        `json:"mac_address" db:"mac_address"`
	StaticIP       *string        `json:"static_ip" db:"static_ip"`
	VLAN           *int16         `json:"vlan_tag" db:"vlan_tag"`
	ConnectionType ConnectionType `json:"connection_type" db:"connection_type"`
	ParentNodeID   *uuid.UUID     `json:"parent_node_id" db:"parent_node_id"`
	Location       *string        `json:"location" db:"location"`
	Notes          *string        `json:"notes" db:"notes"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at" db:"updated_at"`
}

func (n *NetworkNode) Validate() error {
	if len(strings.TrimSpace(n.DeviceName)) == 0 || len(n.DeviceName) > 200 || (n.ConnectionType != ConnectionLAN && n.ConnectionType != ConnectionWiFi) || (n.VLAN != nil && (*n.VLAN < 1 || *n.VLAN > 4094)) || (n.ParentNodeID != nil && *n.ParentNodeID == n.ID) {
		return ErrValidation
	}
	if n.MACAddress != nil {
		parsed, err := net.ParseMAC(*n.MACAddress)
		if err != nil {
			return ErrValidation
		}
		normalized := strings.ToUpper(parsed.String())
		n.MACAddress = &normalized
	}
	if n.StaticIP != nil && net.ParseIP(*n.StaticIP) == nil {
		return ErrValidation
	}
	return nil
}

type Journal struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	Tags      []string  `json:"tags" db:"tags"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type PomodoroStatus string

const (
	PomodoroRunning PomodoroStatus = "running"
	PomodoroPaused  PomodoroStatus = "paused"
)

type PomodoroState struct {
	Status          PomodoroStatus `json:"status"`
	StartedAt       time.Time      `json:"started_at"`
	DurationSeconds int            `json:"duration_seconds"`
	TaskID          *uuid.UUID     `json:"task_id,omitempty"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func (s PomodoroState) Validate() error {
	if (s.Status != PomodoroRunning && s.Status != PomodoroPaused) || s.DurationSeconds < 60 || s.DurationSeconds > 14400 || s.StartedAt.IsZero() {
		return ErrValidation
	}
	return nil
}

type PomodoroHistory struct {
	ID              uuid.UUID  `db:"id"`
	EventID         uuid.UUID  `db:"event_id"`
	UserID          uuid.UUID  `db:"user_id"`
	TaskID          *uuid.UUID `db:"task_id"`
	StartedAt       time.Time  `db:"started_at"`
	EndedAt         time.Time  `db:"ended_at"`
	DurationSeconds int        `db:"duration_seconds"`
	CreatedAt       time.Time  `db:"created_at"`
}

func (j *Journal) Validate() error {
	if len(strings.TrimSpace(j.Title)) == 0 || len(j.Title) > 300 || len(j.Content) > 100000 || len(j.Tags) > 20 {
		return ErrValidation
	}
	seen := map[string]bool{}
	for index, tag := range j.Tags {
		tag = strings.TrimSpace(strings.ToLower(tag))
		if tag == "" || len(tag) > 50 || seen[tag] {
			return ErrValidation
		}
		seen[tag] = true
		j.Tags[index] = tag
	}
	return nil
}

func ParseCost(raw string) error {
	if strings.TrimSpace(raw) == "" {
		return ErrValidation
	}
	for _, character := range raw {
		if (character < '0' || character > '9') && character != '.' {
			return fmt.Errorf("%w: invalid cost", ErrValidation)
		}
	}
	return nil
}
