package habit

import (
	"testing"
	"time"

	"github.com/epalmerini/abitudini/internal/testhelpers"
)

func TestStore_CreateAndGetByID(t *testing.T) {
	db := testhelpers.NewTestDB(t)
	store := NewStore(db)

	habit := &Habit{
		Description: "Test Habit",
		Frequency:   FrequencyDaily,
		StartDate:   time.Now(),
		Color:       "blue",
	}

	id, err := store.Create(habit, nil)
	if err != nil {
		t.Fatalf("failed to create habit: %v", err)
	}
	if id == 0 {
		t.Error("expected non-zero habit ID")
	}

	retrieved, err := store.GetByID(id)
	if err != nil {
		t.Fatalf("failed to get habit: %v", err)
	}
	if retrieved.ID != id {
		t.Errorf("expected ID %d, got %d", id, retrieved.ID)
	}
	if retrieved.Description != "Test Habit" {
		t.Errorf("expected description 'Test Habit', got '%s'", retrieved.Description)
	}
}

func TestStore_GetAll(t *testing.T) {
	db := testhelpers.NewTestDB(t)
	store := NewStore(db)

	// Create multiple habits
	for i := 0; i < 3; i++ {
		_, err := store.Create(&Habit{
			Description: "Habit " + string(rune(i)),
			Frequency:   FrequencyDaily,
			StartDate:   time.Now(),
			Color:       "blue",
		}, nil)
		if err != nil {
			t.Fatalf("failed to create habit: %v", err)
		}
	}

	habits, err := store.GetAll()
	if err != nil {
		t.Fatalf("failed to get all habits: %v", err)
	}
	if len(habits) != 3 {
		t.Errorf("expected 3 habits, got %d", len(habits))
	}
}

func TestStore_Update(t *testing.T) {
	db := testhelpers.NewTestDB(t)
	store := NewStore(db)

	habit := &Habit{
		Description: "Original",
		Frequency:   FrequencyDaily,
		StartDate:   time.Now(),
		Color:       "blue",
	}

	id, err := store.Create(habit, nil)
	if err != nil {
		t.Fatalf("failed to create habit: %v", err)
	}

	habit.ID = id
	habit.Description = "Updated"

	err = store.Update(habit, nil)
	if err != nil {
		t.Fatalf("failed to update habit: %v", err)
	}

	retrieved, err := store.GetByID(id)
	if err != nil {
		t.Fatalf("failed to get habit: %v", err)
	}
	if retrieved.Description != "Updated" {
		t.Errorf("expected description 'Updated', got '%s'", retrieved.Description)
	}
}

func TestStore_Delete(t *testing.T) {
	db := testhelpers.NewTestDB(t)
	store := NewStore(db)

	habit := &Habit{
		Description: "Test",
		Frequency:   FrequencyDaily,
		StartDate:   time.Now(),
		Color:       "blue",
	}

	id, err := store.Create(habit, nil)
	if err != nil {
		t.Fatalf("failed to create habit: %v", err)
	}

	err = store.Delete(id)
	if err != nil {
		t.Fatalf("failed to delete habit: %v", err)
	}

	_, err = store.GetByID(id)
	if err == nil {
		t.Error("expected error when getting deleted habit")
	}
}

func TestStore_GetByID_NotFound(t *testing.T) {
	db := testhelpers.NewTestDB(t)
	store := NewStore(db)

	_, err := store.GetByID(999999)
	if err == nil {
		t.Error("expected error for non-existent habit")
	}
}

func TestStore_IsValidForDate_Daily(t *testing.T) {
	store := NewStore(nil)
	habit := &Habit{Frequency: FrequencyDaily}

	// Daily habits should be valid every day
	if !store.IsValidForDate(habit, time.Now()) {
		t.Error("daily habit should be valid")
	}
}

func TestStore_IsValidForDate_Weekly(t *testing.T) {
	store := NewStore(nil)
	habit := &Habit{
		Frequency: FrequencyWeekly,
		Schedule: Schedule{
			DaysOfWeek: []int{int(time.Monday)},
		},
	}

	// Find a Monday
	monday := time.Now()
	for monday.Weekday() != time.Monday {
		monday = monday.AddDate(0, 0, 1)
	}

	if !store.IsValidForDate(habit, monday) {
		t.Error("habit should be valid on Monday")
	}

	// Check non-matching day
	tuesday := monday.AddDate(0, 0, 1)
	if store.IsValidForDate(habit, tuesday) {
		t.Error("habit should not be valid on Tuesday")
	}
}

func TestStore_IsValidForDate_Weekly_NoSchedule(t *testing.T) {
	store := NewStore(nil)
	habit := &Habit{
		Frequency: FrequencyWeekly,
		Schedule:  Schedule{DaysOfWeek: []int{}},
	}

	// No schedule means valid every day
	if !store.IsValidForDate(habit, time.Now()) {
		t.Error("weekly habit with no schedule should be valid")
	}
}

func TestStore_IsValidForDate_Monthly(t *testing.T) {
	store := NewStore(nil)
	habit := &Habit{
		Frequency: FrequencyMonthly,
		Schedule: Schedule{
			DaysOfMonth: []int{15},
		},
	}

	// Day 15
	day15 := time.Date(2025, 12, 15, 0, 0, 0, 0, time.UTC)
	if !store.IsValidForDate(habit, day15) {
		t.Error("habit should be valid on day 15")
	}

	// Day 14
	day14 := time.Date(2025, 12, 14, 0, 0, 0, 0, time.UTC)
	if store.IsValidForDate(habit, day14) {
		t.Error("habit should not be valid on day 14")
	}
}

func TestStore_IsValidForDate_Monthly_NoSchedule(t *testing.T) {
	store := NewStore(nil)
	habit := &Habit{
		Frequency: FrequencyMonthly,
		Schedule:  Schedule{DaysOfMonth: []int{}},
	}

	// No schedule means valid every day
	if !store.IsValidForDate(habit, time.Now()) {
		t.Error("monthly habit with no schedule should be valid")
	}
}

func TestStore_IsValidForDate_Unknown(t *testing.T) {
	store := NewStore(nil)
	habit := &Habit{Frequency: "unknown"}

	// Unknown frequency should be invalid
	if store.IsValidForDate(habit, time.Now()) {
		t.Error("unknown frequency should be invalid")
	}
}
