package record

import (
	"time"

	"github.com/epalmerini/abitudini/internal/habit"
)

// StoreAdapter defines the interface for data access
type StoreAdapter interface {
	Record(habitID int, date time.Time) error
	GetByHabitAndDateRange(habitID int, from, to time.Time) ([]Record, error)
}

// HabitAdapter defines the interface for habit access
type HabitAdapter interface {
	GetByID(habitID int) (*habit.Habit, error)
}

type Service struct {
	store        StoreAdapter
	habitService HabitAdapter
}

func NewService(store StoreAdapter, habitService HabitAdapter) *Service {
	return &Service{store: store, habitService: habitService}
}

func (s *Service) MarkDoneToday(habitID int) error {
	return s.store.Record(habitID, time.Now())
}

func (s *Service) GetRecords(habitID int, from, to time.Time) ([]Record, error) {
	return s.store.GetByHabitAndDateRange(habitID, from, to)
}

func (s *Service) GetContributionData(habitID int, from, to time.Time) ([]ContributionDay, error) {
	records, err := s.GetRecords(habitID, from, to)
	if err != nil {
		return nil, err
	}

	// Create a map of completed dates
	completedMap := make(map[string]bool)
	for _, record := range records {
		completedMap[record.RecordDate.Format("2006-01-02")] = true
	}

	// Generate all days in range
	var contributions []ContributionDay
	for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
		contributions = append(contributions, ContributionDay{
			Date:      d,
			Completed: completedMap[d.Format("2006-01-02")],
		})
	}

	return contributions, nil
}

func (s *Service) GetHabit(habitID int) (*habit.Habit, error) {
	h, err := s.habitService.GetByID(habitID)
	if err != nil {
		return nil, err
	}
	
	// Set CompletedToday flag
	completed, err := s.IsCompletedToday(habitID)
	if err != nil {
		return nil, err
	}
	h.CompletedToday = completed
	
	return h, nil
}

func (s *Service) IsCompletedToday(habitID int) (bool, error) {
	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1).Add(-time.Nanosecond)
	
	records, err := s.GetRecords(habitID, startOfDay, endOfDay)
	if err != nil {
		return false, err
	}
	return len(records) > 0, nil
}
