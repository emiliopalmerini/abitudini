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
		habit: &habit.Habit{ID: 1, Frequency: "daily"},
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
		habitErr: errors.New("habit not found"),
	})

	_, err := s.GetByHabitID(1)
	if err == nil {
		t.Error("expected error for non-existent habit")
	}
}

func TestGetByHabitID_NoRecords(t *testing.T) {
	s := NewService(&mockStreakStore{
		habit:   &habit.Habit{ID: 1, Frequency: "daily"},
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
		habit:      &habit.Habit{ID: 1, Frequency: "daily"},
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

func TestCalculateWeeklyStreak_WithRecords(t *testing.T) {
	s := NewService(nil)
	
	// Get a Monday reference
	now := time.Now()
	for now.Weekday() != time.Monday {
		now = now.AddDate(0, 0, -1)
	}
	
	// Add records for this week and previous weeks
	records := []time.Time{
		now.AddDate(0, 0, 2), // Wednesday of this week
		now.AddDate(0, 0, -5), // Friday of previous week
		now.AddDate(0, 0, -12), // Friday of 2 weeks ago
	}

	count := s.calculateWeeklyStreak(now, records)
	if count < 1 {
		t.Errorf("expected streak >= 1 for weekly, got %d", count)
	}
}

func TestCalculateWeeklyStreak_NoCurrentWeek(t *testing.T) {
	s := NewService(nil)
	
	now := time.Now()
	// Records only in past weeks, not this week
	records := []time.Time{
		now.AddDate(0, 0, -7),
		now.AddDate(0, 0, -14),
	}

	count := s.calculateWeeklyStreak(now, records)
	if count != 0 {
		t.Errorf("expected streak of 0 (no current week), got %d", count)
	}
}

func TestCalculateWeeklyStreak_Empty(t *testing.T) {
	s := NewService(nil)
	count := s.calculateWeeklyStreak(time.Now(), []time.Time{})
	if count != 0 {
		t.Errorf("expected streak of 0 for empty records, got %d", count)
	}
}

func TestCalculateMonthlyStreak_WithRecords(t *testing.T) {
	s := NewService(nil)
	
	now := time.Now()
	// Record this month, and previous months
	records := []time.Time{
		now,
		now.AddDate(0, -1, 0),
		now.AddDate(0, -2, 0),
	}

	count := s.calculateMonthlyStreak(now, records)
	if count < 1 {
		t.Errorf("expected streak >= 1 for monthly, got %d", count)
	}
}

func TestCalculateMonthlyStreak_NoCurrentMonth(t *testing.T) {
	s := NewService(nil)
	
	now := time.Now()
	// Records only in past months
	records := []time.Time{
		now.AddDate(0, -1, 0),
		now.AddDate(0, -2, 0),
	}

	count := s.calculateMonthlyStreak(now, records)
	if count != 0 {
		t.Errorf("expected streak of 0 (no current month), got %d", count)
	}
}

func TestCalculateMonthlyStreak_Empty(t *testing.T) {
	s := NewService(nil)
	count := s.calculateMonthlyStreak(time.Now(), []time.Time{})
	if count != 0 {
		t.Errorf("expected streak of 0 for empty records, got %d", count)
	}
}

func TestGetWeekStart(t *testing.T) {
	s := NewService(nil)
	
	// Use a known date: Friday, Dec 5, 2025
	date := time.Date(2025, 12, 5, 15, 30, 0, 0, time.UTC)
	weekStart := s.getWeekStart(date)
	
	// Should return Sunday, Dec 1, 2025 (since Sunday is weekday 0)
	expected := time.Date(2025, 11, 30, 0, 0, 0, 0, time.UTC)
	if weekStart != expected {
		t.Errorf("expected week start %v, got %v", expected, weekStart)
	}
}

func TestGetMonthStart(t *testing.T) {
	s := NewService(nil)
	
	date := time.Date(2025, 12, 15, 15, 30, 0, 0, time.UTC)
	monthStart := s.getMonthStart(date)
	
	expected := time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)
	if monthStart != expected {
		t.Errorf("expected month start %v, got %v", expected, monthStart)
	}
}

func TestDaysInMonth(t *testing.T) {
	s := NewService(nil)
	
	tests := []struct {
		date     time.Time
		expected int
	}{
		{time.Date(2025, 12, 15, 0, 0, 0, 0, time.UTC), 31}, // December
		{time.Date(2025, 2, 15, 0, 0, 0, 0, time.UTC), 28},  // February (non-leap)
		{time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC), 29},  // February (leap year)
	}

	for _, test := range tests {
		days := s.daysInMonth(test.date)
		if days != test.expected {
			t.Errorf("expected %d days in %s, got %d", test.expected, test.date.Month(), days)
		}
	}
}

func TestCalculateStreak_UnknownFrequency(t *testing.T) {
	s := NewService(nil)
	now := time.Now()
	records := []time.Time{now, now.AddDate(0, 0, -1)}

	// Unknown frequency should return 0
	count := s.calculateStreak("unknown", now, records)
	if count != 0 {
		t.Errorf("expected streak of 0 for unknown frequency, got %d", count)
	}
}


