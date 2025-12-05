package streak

import (
	"errors"
	"testing"
	"time"

	"github.com/epalmerini/abitudini/internal/habit"
)

type mockStreakStore struct {
	habit   *habit.Habit
	records []time.Time
	habitErr error
	recordsErr error
}

func (m *mockStreakStore) GetHabitByID(habitID int) (*habit.Habit, error) {
	if m.habitErr != nil {
		return nil, m.habitErr
	}
	return m.habit, nil
}

func (m *mockStreakStore) GetRecordsByHabit(habitID int) ([]time.Time, error) {
	if m.recordsErr != nil {
		return nil, m.recordsErr
	}
	return m.records, nil
}

func TestGetByHabitID_Success(t *testing.T) {
	s := NewService(&mockStreakStore{
		records: []time.Time{
			time.Now(),
			time.Now().AddDate(0, 0, -1),
		},
	})

	streak, err := s.GetByHabitID(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if streak.HabitID != 1 {
		t.Errorf("expected HabitID 1, got %d", streak.HabitID)
	}
	if streak.CurrentCount == 0 {
		t.Error("expected CurrentCount > 0")
	}
}

func TestGetByHabitID_HabitNotFound(t *testing.T) {
	s := NewService(&mockStreakStore{
		recordsErr: errors.New("habit not found"),
	})

	streak, err := s.GetByHabitID(1)
	// Should return zero streak on error, not an error
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if streak.CurrentCount != 0 {
		t.Errorf("expected CurrentCount 0, got %d", streak.CurrentCount)
	}
}

func TestGetByHabitID_NoRecords(t *testing.T) {
	s := NewService(&mockStreakStore{
		records: []time.Time{},
	})

	streak, err := s.GetByHabitID(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if streak.CurrentCount != 0 {
		t.Errorf("expected CurrentCount 0 for no records, got %d", streak.CurrentCount)
	}
}

func TestGetByHabitID_RecordsFetchError(t *testing.T) {
	s := NewService(&mockStreakStore{
		recordsErr: errors.New("fetch records failed"),
	})

	streak, err := s.GetByHabitID(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if streak.CurrentCount != 0 {
		t.Errorf("expected CurrentCount 0 on fetch error, got %d", streak.CurrentCount)
	}
}

func TestCalculateDailyStreak_Consecutive(t *testing.T) {
	s := NewService(nil)
	now := time.Now()
	records := []time.Time{
		now,
		now.AddDate(0, 0, -1),
		now.AddDate(0, 0, -2),
		now.AddDate(0, 0, -3),
	}

	count := s.calculateDailyStreak(now, records)
	if count != 4 {
		t.Errorf("expected streak of 4, got %d", count)
	}
}

func TestCalculateDailyStreak_Broken(t *testing.T) {
	s := NewService(nil)
	now := time.Now()
	records := []time.Time{
		now,
		now.AddDate(0, 0, -1),
		now.AddDate(0, 0, -3), // gap here
	}

	count := s.calculateDailyStreak(now, records)
	if count != 2 {
		t.Errorf("expected streak of 2 (broken), got %d", count)
	}
}

func TestCalculateDailyStreak_Empty(t *testing.T) {
	s := NewService(nil)
	count := s.calculateDailyStreak(time.Now(), []time.Time{})
	if count != 0 {
		t.Errorf("expected streak of 0 for empty records, got %d", count)
	}
}




