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
	h, err := s.store.GetHabitByID(habitID)
	if err != nil {
		return nil, err
	}

	recordDates, err := s.store.GetRecordsByHabit(habitID)
	if err != nil {
		return &Streak{HabitID: habitID, CurrentCount: 0}, nil
	}

	count := s.calculateStreak(string(h.Frequency), time.Now(), recordDates)
	return &Streak{HabitID: habitID, CurrentCount: count}, nil
}

func (s *Service) calculateStreak(frequency string, today time.Time, recordDates []time.Time) int {
	if len(recordDates) == 0 {
		return 0
	}

	var count int

	switch frequency {
	case "daily":
		count = s.calculateDailyStreak(today, recordDates)
	case "weekly":
		count = s.calculateWeeklyStreak(today, recordDates)
	case "monthly":
		count = s.calculateMonthlyStreak(today, recordDates)
	}

	return count
}

func (s *Service) calculateDailyStreak(today time.Time, recordDates []time.Time) int {
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

func (s *Service) calculateWeeklyStreak(today time.Time, recordDates []time.Time) int {
	count := 0
	currentWeek := s.getWeekStart(today)
	recordMap := make(map[string]bool)

	for _, record := range recordDates {
		recordMap[record.Format("2006-01-02")] = true
	}

	// Check if this week has at least one record
	hasCurrentWeek := false
	for i := 0; i < 7; i++ {
		if recordMap[currentWeek.AddDate(0, 0, i).Format("2006-01-02")] {
			hasCurrentWeek = true
			break
		}
	}

	if !hasCurrentWeek {
		return 0
	}

	count = 1
	currentWeek = currentWeek.AddDate(0, 0, -7)

	for len(recordDates) > 0 {
		found := false

		for _, record := range recordDates {
			if record.After(currentWeek) && record.Before(currentWeek.AddDate(0, 0, 7)) {
				found = true
				break
			}
		}

		if found {
			count++
			currentWeek = currentWeek.AddDate(0, 0, -7)
		} else {
			break
		}
	}

	return count
}

func (s *Service) calculateMonthlyStreak(today time.Time, recordDates []time.Time) int {
	count := 0
	currentMonth := s.getMonthStart(today)
	recordMap := make(map[string]bool)

	for _, record := range recordDates {
		recordMap[record.Format("2006-01-02")] = true
	}

	// Check if this month has at least one record
	hasCurrentMonth := false
	for i := 0; i < s.daysInMonth(currentMonth); i++ {
		if recordMap[currentMonth.AddDate(0, 0, i).Format("2006-01-02")] {
			hasCurrentMonth = true
			break
		}
	}

	if !hasCurrentMonth {
		return 0
	}

	count = 1
	currentMonth = currentMonth.AddDate(0, -1, 0)

	for len(recordDates) > 0 {
		found := false
		monthEnd := currentMonth.AddDate(0, 1, 0)

		for _, record := range recordDates {
			if record.After(currentMonth) && record.Before(monthEnd) {
				found = true
				break
			}
		}

		if found {
			count++
			currentMonth = currentMonth.AddDate(0, -1, 0)
		} else {
			break
		}
	}

	return count
}

func (s *Service) getWeekStart(t time.Time) time.Time {
	dayOfWeek := t.Weekday()
	offset := -int(dayOfWeek)
	return t.AddDate(0, 0, offset).Truncate(24 * time.Hour)
}

func (s *Service) getMonthStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func (s *Service) daysInMonth(t time.Time) int {
	return s.getMonthStart(t).AddDate(0, 1, 0).AddDate(0, 0, -1).Day()
}
