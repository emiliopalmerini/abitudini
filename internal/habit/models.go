package habit

import "time"

type Habit struct {
	ID               int       `json:"id"`
	Description      string    `json:"description"`
	StartDate        time.Time `json:"start_date"`
	Color            string    `json:"color"`
	CreatedAt        time.Time `json:"created_at"`
	CompletedToday   bool      `json:"completed_today"`
}
