package streak

import (
	"time"

	"github.com/epalmerini/abitudini/internal/habit"
)

// StoreAdapter defines the interface for data access
type StoreAdapter interface {
	GetHabitByID(habitID int) (*habit.Habit, error)
	GetRecordsByHabit(habitID int) ([]time.Time, error)
}

type Service struct {
	store StoreAdapter
}

func NewService(store StoreAdapter) *Service {
	return &Service{store: store}
}

func (s *Service) GetByHabitID(habitID int) (*Streak, error) {
	recordDates, err := s.store.GetRecordsByHabit(habitID)
	if err != nil {
		return &Streak{HabitID: habitID, CurrentCount: 0}, nil
	}

	count := s.calculateDailyStreak(time.Now(), recordDates)
	return &Streak{HabitID: habitID, CurrentCount: count}, nil
}

func (s *Service) calculateDailyStreak(today time.Time, recordDates []time.Time) int {
	if len(recordDates) == 0 {
		return 0
	}

	count := 0
	currentDate := today

	for _, record := range recordDates {
		if record.Format("2006-01-02") == currentDate.Format("2006-01-02") {
			count++
			currentDate = currentDate.AddDate(0, 0, -1)
		} else {
			break
		}
	}

	return count
}
