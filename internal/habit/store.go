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

func (s *Store) Create(h *Habit, schedule *Schedule) (int, error) {
	result, err := s.db.Exec(
		`INSERT INTO habits (description, frequency, start_date, color)
		 VALUES (?, ?, ?, ?)`,
		h.Description,
		string(h.Frequency),
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

	if schedule != nil {
		if err := s.insertSchedule(int(habitID), schedule); err != nil {
			return 0, err
		}
	}

	return int(habitID), nil
}

func (s *Store) Update(h *Habit, schedule *Schedule) error {
	_, err := s.db.Exec(
		`UPDATE habits 
		 SET description = ?, frequency = ?, start_date = ?, color = ?
		 WHERE id = ?`,
		h.Description,
		string(h.Frequency),
		h.StartDate.Format("2006-01-02"),
		h.Color,
		h.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update habit: %w", err)
	}

	// Delete old schedule
	if _, err := s.db.Exec(`DELETE FROM habit_schedule WHERE habit_id = ?`, h.ID); err != nil {
		return fmt.Errorf("failed to delete old schedule: %w", err)
	}

	// Insert new schedule
	if schedule != nil {
		if err := s.insertSchedule(h.ID, schedule); err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) insertSchedule(habitID int, schedule *Schedule) error {
	if len(schedule.DaysOfWeek) > 0 {
		for _, day := range schedule.DaysOfWeek {
			if _, err := s.db.Exec(
				`INSERT INTO habit_schedule (habit_id, day_of_week) VALUES (?, ?)`,
				habitID, day,
			); err != nil {
				return fmt.Errorf("failed to insert weekly schedule: %w", err)
			}
		}
	}

	if len(schedule.DaysOfMonth) > 0 {
		for _, day := range schedule.DaysOfMonth {
			if _, err := s.db.Exec(
				`INSERT INTO habit_schedule (habit_id, day_of_month) VALUES (?, ?)`,
				habitID, day,
			); err != nil {
				return fmt.Errorf("failed to insert monthly schedule: %w", err)
			}
		}
	}

	return nil
}

func (s *Store) GetByID(habitID int) (*Habit, error) {
	h := &Habit{}
	var startDate string
	var createdAt string
	var frequencyStr string

	err := s.db.QueryRow(
		`SELECT id, description, frequency, start_date, color, created_at 
		 FROM habits WHERE id = ?`,
		habitID,
	).Scan(&h.ID, &h.Description, &frequencyStr, &startDate, &h.Color, &createdAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("habit not found")
		}
		return nil, fmt.Errorf("failed to get habit: %w", err)
	}

	h.Frequency = Frequency(frequencyStr)
	h.StartDate, _ = time.Parse("2006-01-02", startDate)
	h.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)

	if schedule, err := s.getSchedule(habitID); err == nil {
		h.Schedule = *schedule
	}

	return h, nil
}

func (s *Store) GetAll() ([]Habit, error) {
	rows, err := s.db.Query(
		`SELECT id, description, frequency, start_date, color, created_at 
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
		var frequencyStr string

		if err := rows.Scan(&h.ID, &h.Description, &frequencyStr, &startDate, &h.Color, &createdAt); err != nil {
			return nil, fmt.Errorf("failed to scan habit: %w", err)
		}

		h.Frequency = Frequency(frequencyStr)
		h.StartDate, _ = time.Parse("2006-01-02", startDate)
		h.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)

		if schedule, err := s.getSchedule(h.ID); err == nil {
			h.Schedule = *schedule
		}

		habits = append(habits, h)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating habits: %w", err)
	}

	return habits, nil
}

func (s *Store) getSchedule(habitID int) (*Schedule, error) {
	schedule := &Schedule{
		DaysOfWeek:  []int{},
		DaysOfMonth: []int{},
	}

	rows, err := s.db.Query(
		`SELECT day_of_week, day_of_month FROM habit_schedule WHERE habit_id = ?`,
		habitID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dayOfWeek sql.NullInt64
		var dayOfMonth sql.NullInt64

		if err := rows.Scan(&dayOfWeek, &dayOfMonth); err != nil {
			return nil, err
		}

		if dayOfWeek.Valid {
			schedule.DaysOfWeek = append(schedule.DaysOfWeek, int(dayOfWeek.Int64))
		}
		if dayOfMonth.Valid {
			schedule.DaysOfMonth = append(schedule.DaysOfMonth, int(dayOfMonth.Int64))
		}
	}

	return schedule, rows.Err()
}

func (s *Store) IsValidForDate(h *Habit, date time.Time) bool {
	switch h.Frequency {
	case "daily":
		return true
	case "weekly":
		if len(h.Schedule.DaysOfWeek) == 0 {
			return true
		}
		dayOfWeek := int(date.Weekday())
		for _, d := range h.Schedule.DaysOfWeek {
			if d == dayOfWeek {
				return true
			}
		}
		return false
	case "monthly":
		if len(h.Schedule.DaysOfMonth) == 0 {
			return true
		}
		dayOfMonth := date.Day()
		for _, d := range h.Schedule.DaysOfMonth {
			if d == dayOfMonth {
				return true
			}
		}
		return false
	}
	return false
}

func (s *Store) Delete(habitID int) error {
	// Delete habit records (cascades due to FK constraint)
	if _, err := s.db.Exec(`DELETE FROM records WHERE habit_id = ?`, habitID); err != nil {
		return fmt.Errorf("failed to delete records: %w", err)
	}

	// Delete habit schedule entries
	if _, err := s.db.Exec(`DELETE FROM habit_schedule WHERE habit_id = ?`, habitID); err != nil {
		return fmt.Errorf("failed to delete schedule: %w", err)
	}

	// Delete habit
	if _, err := s.db.Exec(`DELETE FROM habits WHERE id = ?`, habitID); err != nil {
		return fmt.Errorf("failed to delete habit: %w", err)
	}

	return nil
}
