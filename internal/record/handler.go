package record

import (
	"fmt"
	"net/http"
	"time"

	"github.com/epalmerini/abitudini/internal/habit"
	"github.com/epalmerini/abitudini/internal/shared"
)

// HandlerService interface for dependency injection
type HandlerService interface {
	MarkDoneToday(habitID int) error
	GetContributionData(habitID int, from, to time.Time) ([]ContributionDay, error)
	GetHabit(habitID int) (*habit.Habit, error)
}

type Handler struct {
	shared.BaseHandler
	service HandlerService
}

func NewHandler(service HandlerService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) MarkDoneToday(w http.ResponseWriter, r *http.Request) {
	if !h.ValidateMethod(w, r, http.MethodPost) {
		return
	}

	habitID, err := h.ExtractIntPathParam(r, "id")
	if err != nil {
		h.WriteError(w, "Invalid habit ID", http.StatusBadRequest)
		return
	}

	if err := h.service.MarkDoneToday(habitID); err != nil {
		h.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get updated habit and return it
	habitData, err := h.service.GetHabit(habitID)
	if err != nil {
		h.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := `<div class="success">Marked as done today</div>` + habit.RenderHabit(habitData)
	h.WriteHTML(w, response)
}

func (h *Handler) GetContribution(w http.ResponseWriter, r *http.Request) {
	if !h.ValidateMethod(w, r, http.MethodGet) {
		return
	}

	habitID, err := h.ExtractIntPathParam(r, "id")
	if err != nil {
		h.WriteError(w, "Invalid habit ID", http.StatusBadRequest)
		return
	}

	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	from, _ := time.Parse("2006-01-02", fromStr)
	to, _ := time.Parse("2006-01-02", toStr)

	if from.IsZero() || to.IsZero() {
		from = time.Now().AddDate(-1, 0, 0)
		to = time.Now()
	}

	contributions, err := h.service.GetContributionData(habitID, from, to)
	if err != nil {
		h.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Group contributions by week (Sunday to Saturday)
	html := `<div class="contribution-container">`
	
	// Add month headers if there are contributions
	if len(contributions) > 0 {
		html += `<div class="contribution-months">`
		lastMonth := -1
		lastYear := -1
		monthNames := []string{"", "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
		
		for i, day := range contributions {
			month := int(day.Date.Month())
			year := day.Date.Year()
			
			// Check if this is the first day of a month
			if (month != lastMonth || year != lastYear) && day.Date.Day() <= 7 {
				html += fmt.Sprintf(`<span class="contribution-month" style="grid-column: %d;">%s</span>`, 
					(i/7)+1, monthNames[month])
				lastMonth = month
				lastYear = year
			}
		}
		html += `</div>`
	}
	
	html += fmt.Sprintf(`<div id="contribution-%d" class="contribution-grid">`, habitID)
	
	var currentWeek []ContributionDay
	lastWeekStart := -1
	
	for i, day := range contributions {
		weekStart := day.Date.AddDate(0, 0, -int(day.Date.Weekday())).Unix()
		currentWeekStart := int(weekStart)
		
		if lastWeekStart != currentWeekStart && lastWeekStart != -1 {
			// Write current week
			html += `<div class="contribution-week">`
			for _, d := range currentWeek {
				level := "level-0"
				if d.Completed {
					level = "level-4"
				}
				html += fmt.Sprintf(`<div class="day %s" title="%s"></div>`,
					level,
					d.Date.Format("2006-01-02"),
				)
			}
			html += `</div>`
			currentWeek = []ContributionDay{}
		}
		
		currentWeek = append(currentWeek, day)
		lastWeekStart = currentWeekStart
		
		// Write last week
		if i == len(contributions)-1 {
			html += `<div class="contribution-week">`
			for _, d := range currentWeek {
				level := "level-0"
				if d.Completed {
					level = "level-4"
				}
				html += fmt.Sprintf(`<div class="day %s" title="%s"></div>`,
					level,
					d.Date.Format("2006-01-02"),
				)
			}
			html += `</div>`
		}
	}
	
	html += `</div></div>`

	h.WriteHTML(w, html)
}
