package record

import "time"

type Record struct {
	ID          int       `json:"id"`
	HabitID     int       `json:"habit_id"`
	RecordDate  time.Time `json:"record_date"`
	CompletedAt time.Time `json:"completed_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type ContributionDay struct {
	Date      time.Time `json:"date"`
	Completed bool      `json:"completed"`
}
