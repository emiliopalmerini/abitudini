package record

import (
	"errors"
	"testing"
	"time"

	"github.com/epalmerini/abitudini/internal/habit"
)

type mockRecordStore struct {
	records []Record
	err     error
}

func (m *mockRecordStore) Record(habitID int, date time.Time) error {
	return m.err
}

func (m *mockRecordStore) GetByHabitAndDateRange(habitID int, from, to time.Time) ([]Record, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.records, nil
}

type mockHabitAdapter struct {
	habit *habit.Habit
	err   error
}

func (m *mockHabitAdapter) GetByID(habitID int) (*habit.Habit, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.habit, nil
}

func TestMarkDoneToday_Success(t *testing.T) {
	store := &mockRecordStore{}
	habitAdapter := &mockHabitAdapter{}
	s := NewService(store, habitAdapter)

	err := s.MarkDoneToday(1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestMarkDoneToday_Error(t *testing.T) {
	store := &mockRecordStore{err: errors.New("record failed")}
	habitAdapter := &mockHabitAdapter{}
	s := NewService(store, habitAdapter)

	err := s.MarkDoneToday(1)
	if err == nil {
		t.Error("expected error when recording fails")
	}
}

func TestGetRecords_Success(t *testing.T) {
	now := time.Now()
	records := []Record{
		{ID: 1, HabitID: 1, RecordDate: now},
		{ID: 2, HabitID: 1, RecordDate: now.AddDate(0, 0, -1)},
	}
	store := &mockRecordStore{records: records}
	s := NewService(store, nil)

	result, err := s.GetRecords(1, now.AddDate(0, 0, -7), now)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 records, got %d", len(result))
	}
}

func TestGetRecords_Error(t *testing.T) {
	store := &mockRecordStore{err: errors.New("fetch failed")}
	s := NewService(store, nil)

	_, err := s.GetRecords(1, time.Now().AddDate(0, 0, -7), time.Now())
	if err == nil {
		t.Error("expected error when fetch fails")
	}
}

func TestGetContributionData_Success(t *testing.T) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterday := today.AddDate(0, 0, -1)
	
	records := []Record{
		{ID: 1, HabitID: 1, RecordDate: today},
		{ID: 2, HabitID: 1, RecordDate: yesterday},
	}
	store := &mockRecordStore{records: records}
	s := NewService(store, nil)

	from := today.AddDate(0, 0, -3)
	contributions, err := s.GetContributionData(1, from, today)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	
	if len(contributions) != 4 { // 3 days back + today
		t.Errorf("expected 4 contribution days, got %d", len(contributions))
	}

	// Check that completed days are marked correctly
	completedCount := 0
	for _, c := range contributions {
		if c.Completed {
			completedCount++
		}
	}
	if completedCount != 2 {
		t.Errorf("expected 2 completed days, got %d", completedCount)
	}
}

func TestGetContributionData_NoRecords(t *testing.T) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	
	store := &mockRecordStore{records: []Record{}}
	s := NewService(store, nil)

	from := today.AddDate(0, 0, -3)
	contributions, err := s.GetContributionData(1, from, today)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	
	if len(contributions) != 4 {
		t.Errorf("expected 4 contribution days, got %d", len(contributions))
	}

	// All should be incomplete
	for _, c := range contributions {
		if c.Completed {
			t.Error("expected all days to be incomplete")
			break
		}
	}
}

func TestGetContributionData_Error(t *testing.T) {
	store := &mockRecordStore{err: errors.New("fetch failed")}
	s := NewService(store, nil)

	_, err := s.GetContributionData(1, time.Now().AddDate(0, 0, -7), time.Now())
	if err == nil {
		t.Error("expected error when fetch fails")
	}
}

func TestGetHabit_Success(t *testing.T) {
	expectedHabit := &habit.Habit{ID: 1, Description: "Test"}
	adapter := &mockHabitAdapter{habit: expectedHabit}
	store := &mockRecordStore{records: []Record{}}
	s := NewService(store, adapter)

	h, err := s.GetHabit(1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if h.ID != 1 {
		t.Errorf("expected habit ID 1, got %d", h.ID)
	}
}

func TestGetHabit_Error(t *testing.T) {
	adapter := &mockHabitAdapter{err: errors.New("habit not found")}
	s := NewService(nil, adapter)

	_, err := s.GetHabit(1)
	if err == nil {
		t.Error("expected error when habit not found")
	}
}

func TestIsCompletedToday_Completed(t *testing.T) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	
	records := []Record{
		{ID: 1, HabitID: 1, RecordDate: today},
	}
	store := &mockRecordStore{records: records}
	s := NewService(store, nil)

	completed, err := s.IsCompletedToday(1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !completed {
		t.Error("expected completed to be true")
	}
}

func TestIsCompletedToday_NotCompleted(t *testing.T) {
	store := &mockRecordStore{records: []Record{}}
	s := NewService(store, nil)

	completed, err := s.IsCompletedToday(1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if completed {
		t.Error("expected completed to be false")
	}
}

func TestIsCompletedToday_Error(t *testing.T) {
	store := &mockRecordStore{err: errors.New("fetch failed")}
	s := NewService(store, nil)

	_, err := s.IsCompletedToday(1)
	if err == nil {
		t.Error("expected error when fetch fails")
	}
}
