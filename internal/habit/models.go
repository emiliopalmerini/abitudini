package habit

import "time"

type Frequency string

const (
	FrequencyDaily   Frequency = "daily"
	FrequencyWeekly  Frequency = "weekly"
	FrequencyMonthly Frequency = "monthly"
)

type Habit struct {
	ID               int       `json:"id"`
	Description      string    `json:"description"`
	Frequency        Frequency `json:"frequency"`
	StartDate        time.Time `json:"start_date"`
	Color            string    `json:"color"`
	CreatedAt        time.Time `json:"created_at"`
	Schedule         Schedule  `json:"schedule,omitempty"`
	CompletedToday   bool      `json:"completed_today"`
}

// Schedule holds the days/dates when habit should be done
// For weekly: array of day numbers (0-6, Sunday=0)
// For monthly: array of day numbers (1-31)
// For daily: empty (applies every day)
type Schedule struct {
	DaysOfWeek  []int `json:"days_of_week,omitempty"`
	DaysOfMonth []int `json:"days_of_month,omitempty"`
}
