package streak

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/epalmerini/abitudini/internal/habit"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetRecordsByHabit(habitID int) ([]time.Time, error) {
	rows, err := s.db.Query(
		`SELECT record_date FROM records 
		 WHERE habit_id = ? 
		 ORDER BY record_date DESC 
		 LIMIT 100`,
		habitID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recordDates []time.Time
	for rows.Next() {
		var dateStr string
		if err := rows.Scan(&dateStr); err != nil {
			continue
		}
		if date, err := time.Parse("2006-01-02", dateStr); err == nil {
			recordDates = append(recordDates, date)
		}
	}

	return recordDates, rows.Err()
}

func (s *Store) GetHabitByID(habitID int) (*habit.Habit, error) {
	h := &habit.Habit{}
	var startDate string
	var createdAt string

	err := s.db.QueryRow(
		`SELECT id, description, start_date, color, created_at 
		 FROM habits WHERE id = ?`,
		habitID,
	).Scan(&h.ID, &h.Description, &startDate, &h.Color, &createdAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("habit not found")
		}
		return nil, fmt.Errorf("failed to get habit: %w", err)
	}

	h.StartDate, _ = time.Parse("2006-01-02", startDate)
	h.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	return h, nil
}
