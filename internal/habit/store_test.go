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
		StartDate:   time.Now(),
		Color:       "blue",
	}

	id, err := store.Create(habit)
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
			StartDate:   time.Now(),
			Color:       "blue",
		})
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
		StartDate:   time.Now(),
		Color:       "blue",
	}

	id, err := store.Create(habit)
	if err != nil {
		t.Fatalf("failed to create habit: %v", err)
	}

	habit.ID = id
	habit.Description = "Updated"

	err = store.Update(habit)
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
		StartDate:   time.Now(),
		Color:       "blue",
	}

	id, err := store.Create(habit)
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

func TestStore_IsValidForDate(t *testing.T) {
	store := NewStore(nil)
	habit := &Habit{}

	// All habits are valid for all dates now
	if !store.IsValidForDate(habit, time.Now()) {
		t.Error("habit should always be valid")
	}
}
