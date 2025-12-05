package habit

type StoreAdapter interface {
	Create(h *Habit, schedule *Schedule) (int, error)
	Update(h *Habit, schedule *Schedule) error
	GetByID(habitID int) (*Habit, error)
	GetAll() ([]Habit, error)
	Delete(habitID int) error
}

type RecordServiceAdapter interface {
	IsCompletedToday(habitID int) (bool, error)
}

type Service struct {
	store          StoreAdapter
	recordService  RecordServiceAdapter
}

func NewService(store StoreAdapter, recordService ...RecordServiceAdapter) *Service {
	s := &Service{store: store}
	if len(recordService) > 0 {
		s.recordService = recordService[0]
	}
	return s
}

func (s *Service) Create(h *Habit, schedule *Schedule) (int, error) {
	return s.store.Create(h, schedule)
}

func (s *Service) Update(h *Habit, schedule *Schedule) error {
	return s.store.Update(h, schedule)
}

func (s *Service) GetByID(habitID int) (*Habit, error) {
	return s.store.GetByID(habitID)
}

func (s *Service) GetAll() ([]Habit, error) {
	habits, err := s.store.GetAll()
	if err != nil {
		return nil, err
	}
	
	// Populate CompletedToday if recordService is available
	if s.recordService != nil {
		for i := range habits {
			completed, err := s.recordService.IsCompletedToday(habits[i].ID)
			if err == nil {
				habits[i].CompletedToday = completed
			}
		}
	}
	
	return habits, nil
}

func (s *Service) Delete(habitID int) error {
	return s.store.Delete(habitID)
}
