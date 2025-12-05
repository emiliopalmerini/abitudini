package habit

import (
	"errors"
	"testing"
)

type mockHabitStore struct {
	habits []Habit
	habit  *Habit
	id     int
	err    error
}

func (m *mockHabitStore) Create(h *Habit, schedule *Schedule) (int, error) {
	if m.err != nil {
		return 0, m.err
	}
	return m.id, nil
}

func (m *mockHabitStore) Update(h *Habit, schedule *Schedule) error {
	return m.err
}

func (m *mockHabitStore) GetByID(habitID int) (*Habit, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.habit, nil
}

func (m *mockHabitStore) GetAll() ([]Habit, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.habits, nil
}

func (m *mockHabitStore) Delete(habitID int) error {
	return m.err
}

type mockRecordService struct {
	completed bool
	err       error
}

func (m *mockRecordService) IsCompletedToday(habitID int) (bool, error) {
	if m.err != nil {
		return false, m.err
	}
	return m.completed, nil
}

func TestHabitCreate_Success(t *testing.T) {
	store := &mockHabitStore{id: 42}
	s := NewService(store)

	habit := &Habit{Description: "Test"}
	schedule := &Schedule{}
	
	id, err := s.Create(habit, schedule)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if id != 42 {
		t.Errorf("expected ID 42, got %d", id)
	}
}

func TestHabitCreate_Error(t *testing.T) {
	store := &mockHabitStore{err: errors.New("create failed")}
	s := NewService(store)

	_, err := s.Create(&Habit{}, &Schedule{})
	if err == nil {
		t.Error("expected error when create fails")
	}
}

func TestHabitUpdate_Success(t *testing.T) {
	store := &mockHabitStore{}
	s := NewService(store)

	habit := &Habit{ID: 1}
	schedule := &Schedule{}
	
	err := s.Update(habit, schedule)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestHabitUpdate_Error(t *testing.T) {
	store := &mockHabitStore{err: errors.New("update failed")}
	s := NewService(store)

	err := s.Update(&Habit{}, &Schedule{})
	if err == nil {
		t.Error("expected error when update fails")
	}
}

func TestHabitGetByID_Success(t *testing.T) {
	expected := &Habit{ID: 1, Description: "Test"}
	store := &mockHabitStore{habit: expected}
	s := NewService(store)

	habit, err := s.GetByID(1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if habit.ID != 1 {
		t.Errorf("expected ID 1, got %d", habit.ID)
	}
	if habit.Description != "Test" {
		t.Errorf("expected description 'Test', got '%s'", habit.Description)
	}
}

func TestHabitGetByID_Error(t *testing.T) {
	store := &mockHabitStore{err: errors.New("not found")}
	s := NewService(store)

	_, err := s.GetByID(1)
	if err == nil {
		t.Error("expected error when habit not found")
	}
}

func TestHabitGetAll_Success(t *testing.T) {
	habits := []Habit{
		{ID: 1, Description: "Habit 1"},
		{ID: 2, Description: "Habit 2"},
	}
	store := &mockHabitStore{habits: habits}
	s := NewService(store)

	result, err := s.GetAll()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 habits, got %d", len(result))
	}
}

func TestHabitGetAll_WithRecordService_AllCompleted(t *testing.T) {
	habits := []Habit{
		{ID: 1, Description: "Habit 1"},
		{ID: 2, Description: "Habit 2"},
	}
	store := &mockHabitStore{habits: habits}
	recordService := &mockRecordService{completed: true}
	s := NewService(store, recordService)

	result, err := s.GetAll()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 habits, got %d", len(result))
	}

	// All should be marked completed
	for _, h := range result {
		if !h.CompletedToday {
			t.Error("expected CompletedToday to be true")
		}
	}
}

func TestHabitGetAll_WithRecordService_NoneCompleted(t *testing.T) {
	habits := []Habit{
		{ID: 1, Description: "Habit 1"},
		{ID: 2, Description: "Habit 2"},
	}
	store := &mockHabitStore{habits: habits}
	recordService := &mockRecordService{completed: false}
	s := NewService(store, recordService)

	result, err := s.GetAll()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// All should be marked incomplete
	for _, h := range result {
		if h.CompletedToday {
			t.Error("expected CompletedToday to be false")
		}
	}
}

func TestHabitGetAll_WithRecordService_Error(t *testing.T) {
	habits := []Habit{
		{ID: 1, Description: "Habit 1"},
	}
	store := &mockHabitStore{habits: habits}
	recordService := &mockRecordService{err: errors.New("check failed")}
	s := NewService(store, recordService)

	result, err := s.GetAll()
	if err != nil {
		t.Errorf("expected no error (should skip errors), got %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 habit, got %d", len(result))
	}
	// Should skip the error and leave CompletedToday as false
	if result[0].CompletedToday {
		t.Error("expected CompletedToday to be false after error")
	}
}

func TestHabitGetAll_Error(t *testing.T) {
	store := &mockHabitStore{err: errors.New("fetch failed")}
	s := NewService(store)

	_, err := s.GetAll()
	if err == nil {
		t.Error("expected error when fetch fails")
	}
}

func TestHabitGetAll_WithoutRecordService(t *testing.T) {
	habits := []Habit{
		{ID: 1, Description: "Habit 1", CompletedToday: true},
		{ID: 2, Description: "Habit 2"},
	}
	store := &mockHabitStore{habits: habits}
	s := NewService(store) // No record service

	result, err := s.GetAll()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// First habit should keep its initial state
	if !result[0].CompletedToday {
		t.Error("expected first habit CompletedToday to be true")
	}
	// Second should be false (default)
	if result[1].CompletedToday {
		t.Error("expected second habit CompletedToday to be false")
	}
}

func TestHabitDelete_Success(t *testing.T) {
	store := &mockHabitStore{}
	s := NewService(store)

	err := s.Delete(1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestHabitDelete_Error(t *testing.T) {
	store := &mockHabitStore{err: errors.New("delete failed")}
	s := NewService(store)

	err := s.Delete(1)
	if err == nil {
		t.Error("expected error when delete fails")
	}
}

func TestNewService_WithoutRecordService(t *testing.T) {
	store := &mockHabitStore{}
	s := NewService(store)

	if s.recordService != nil {
		t.Error("expected recordService to be nil")
	}
}

func TestNewService_WithRecordService(t *testing.T) {
	store := &mockHabitStore{}
	recordService := &mockRecordService{}
	s := NewService(store, recordService)

	if s.recordService == nil {
		t.Error("expected recordService to be set")
	}
}


