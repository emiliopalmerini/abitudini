package record

import (
	"database/sql"
	"fmt"
	"time"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Record(habitID int, date time.Time) error {
	dateStr := date.Format("2006-01-02")
	_, err := s.db.Exec(
		`INSERT OR REPLACE INTO records (habit_id, record_date, completed_at)
		 VALUES (?, ?, CURRENT_TIMESTAMP)`,
		habitID,
		dateStr,
	)
	if err != nil {
		return fmt.Errorf("failed to record completion: %w", err)
	}
	return nil
}

func (s *Store) GetByHabitAndDateRange(habitID int, from, to time.Time) ([]Record, error) {
	rows, err := s.db.Query(
		`SELECT id, habit_id, record_date, completed_at, created_at 
		 FROM records 
		 WHERE habit_id = ? AND record_date BETWEEN ? AND ?
		 ORDER BY record_date DESC`,
		habitID,
		from.Format("2006-01-02"),
		to.Format("2006-01-02"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get records: %w", err)
	}
	defer rows.Close()

	var records []Record
	for rows.Next() {
		r := Record{}
		var recordDate, completedAt, createdAt string

		if err := rows.Scan(&r.ID, &r.HabitID, &recordDate, &completedAt, &createdAt); err != nil {
			return nil, fmt.Errorf("failed to scan record: %w", err)
		}

		r.RecordDate, _ = time.Parse("2006-01-02", recordDate)
		r.CompletedAt, _ = time.Parse("2006-01-02 15:04:05", completedAt)
		r.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)

		records = append(records, r)
	}

	return records, rows.Err()
}
