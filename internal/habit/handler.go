package habit

import (
	"net/http"
	"time"

	"github.com/epalmerini/abitudini/internal/shared"
)

// HandlerService interface for dependency injection
type HandlerService interface {
	Create(h *Habit, schedule *Schedule) (int, error)
	Update(h *Habit, schedule *Schedule) error
	GetByID(habitID int) (*Habit, error)
	GetAll() ([]Habit, error)
	Delete(habitID int) error
}

type Handler struct {
	shared.BaseHandler
	service HandlerService
}

func NewHandler(service HandlerService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if !h.ValidateMethod(w, r, http.MethodPost) {
		return
	}

	// Parse form-encoded request
	if err := r.ParseForm(); err != nil {
		h.WriteError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	description := r.FormValue("description")
	frequency := r.FormValue("frequency")
	startDateStr := r.FormValue("start_date")
	color := r.FormValue("color")

	startDate, _ := time.Parse("2006-01-02", startDateStr)

	// Convert to domain types
	domainHabit := &Habit{
		Description: description,
		Frequency:   Frequency(frequency),
		StartDate:   startDate,
		Color:       color,
	}

	domainSchedule := &Schedule{
		DaysOfWeek:  []int{},
		DaysOfMonth: []int{},
	}

	habitID, err := h.service.Create(domainHabit, domainSchedule)
	if err != nil {
		h.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get created habit and return HTML
	habit, err := h.service.GetByID(habitID)
	if err != nil {
		h.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.WriteHTML(w, RenderHabit(habit))
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	if !h.ValidateMethod(w, r, http.MethodGet) {
		return
	}

	domainHabits, err := h.service.GetAll()
	if err != nil {
		h.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(domainHabits) == 0 {
		h.WriteHTML(w, `<div class="empty-state">
			<div class="empty-state-icon">📝</div>
			<h3>No habits yet</h3>
			<p>Create your first habit to get started tracking your progress</p>
		</div>`)
		return
	}

	h.WriteHTML(w, RenderHabitsList(domainHabits))
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	if !h.ValidateMethod(w, r, http.MethodGet) {
		return
	}

	habitID, err := h.ExtractIntPathParam(r, "id")
	if err != nil {
		h.WriteError(w, "Invalid habit ID", http.StatusBadRequest)
		return
	}

	domainHabit, err := h.service.GetByID(habitID)
	if err != nil {
		h.WriteError(w, err.Error(), http.StatusNotFound)
		return
	}

	h.WriteHTML(w, RenderHabit(domainHabit))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if !h.ValidateMethod(w, r, http.MethodPut) {
		return
	}

	habitID, err := h.ExtractIntPathParam(r, "id")
	if err != nil {
		h.WriteError(w, "Invalid habit ID", http.StatusBadRequest)
		return
	}

	// Parse form-encoded request
	if err := r.ParseForm(); err != nil {
		h.WriteError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	description := r.FormValue("description")
	frequency := r.FormValue("frequency")
	startDateStr := r.FormValue("start_date")
	color := r.FormValue("color")

	startDate, _ := time.Parse("2006-01-02", startDateStr)

	// Convert to domain types
	domainHabit := &Habit{
		ID:          habitID,
		Description: description,
		Frequency:   Frequency(frequency),
		StartDate:   startDate,
		Color:       color,
	}

	domainSchedule := &Schedule{
		DaysOfWeek:  []int{},
		DaysOfMonth: []int{},
	}

	if err := h.service.Update(domainHabit, domainSchedule); err != nil {
		h.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get updated habit and return HTML
	habit, err := h.service.GetByID(habitID)
	if err != nil {
		h.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.WriteHTML(w, RenderHabit(habit))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if !h.ValidateMethod(w, r, http.MethodDelete) {
		return
	}

	habitID, err := h.ExtractIntPathParam(r, "id")
	if err != nil {
		h.WriteError(w, "Invalid habit ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(habitID); err != nil {
		h.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return empty response - the habit will be removed from DOM by HTMX
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
