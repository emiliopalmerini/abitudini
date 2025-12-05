package habit

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

func (s *Store) Create(h *Habit) (int, error) {
	result, err := s.db.Exec(
		`INSERT INTO habits (description, start_date, color)
		 VALUES (?, ?, ?)`,
		h.Description,
		h.StartDate.Format("2006-01-02"),
		h.Color,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create habit: %w", err)
	}

	habitID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get habit id: %w", err)
	}

	return int(habitID), nil
}

func (s *Store) Update(h *Habit) error {
	_, err := s.db.Exec(
		`UPDATE habits 
		 SET description = ?, start_date = ?, color = ?
		 WHERE id = ?`,
		h.Description,
		h.StartDate.Format("2006-01-02"),
		h.Color,
		h.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update habit: %w", err)
	}

	return nil
}



func (s *Store) GetByID(habitID int) (*Habit, error) {
	h := &Habit{}
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

func (s *Store) GetAll() ([]Habit, error) {
	rows, err := s.db.Query(
		`SELECT id, description, start_date, color, created_at 
		 FROM habits ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get habits: %w", err)
	}
	defer rows.Close()

	var habits []Habit
	for rows.Next() {
		h := Habit{}
		var startDate string
		var createdAt string

		if err := rows.Scan(&h.ID, &h.Description, &startDate, &h.Color, &createdAt); err != nil {
			return nil, fmt.Errorf("failed to scan habit: %w", err)
		}

		h.StartDate, _ = time.Parse("2006-01-02", startDate)
		h.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)

		habits = append(habits, h)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating habits: %w", err)
	}

	return habits, nil
}

func (s *Store) IsValidForDate(h *Habit, date time.Time) bool {
	return true
}

func (s *Store) Delete(habitID int) error {
	// Delete habit records (cascades due to FK constraint)
	if _, err := s.db.Exec(`DELETE FROM records WHERE habit_id = ?`, habitID); err != nil {
		return fmt.Errorf("failed to delete records: %w", err)
	}

	// Delete habit
	if _, err := s.db.Exec(`DELETE FROM habits WHERE id = ?`, habitID); err != nil {
		return fmt.Errorf("failed to delete habit: %w", err)
	}

	return nil
}
