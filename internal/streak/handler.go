package streak

import (
	"fmt"
	"net/http"

	"github.com/epalmerini/abitudini/internal/shared"
)

// HandlerService interface for dependency injection
type HandlerService interface {
	GetByHabitID(habitID int) (*Streak, error)
}

type Handler struct {
	shared.BaseHandler
	service HandlerService
}

func NewHandler(service HandlerService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetByHabitID(w http.ResponseWriter, r *http.Request) {
	if !h.ValidateMethod(w, r, http.MethodGet) {
		return
	}

	habitID, err := h.ExtractIntPathParam(r, "id")
	if err != nil {
		h.WriteError(w, "Invalid habit ID", http.StatusBadRequest)
		return
	}

	streak, err := h.service.GetByHabitID(habitID)
	if err != nil {
		h.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	streakLabel := "day streak"
	if streak.CurrentCount != 1 {
		streakLabel = "days streak"
	}

	html := fmt.Sprintf(`<div class="streak-display">
		<span class="streak-count">%d</span>
		<span class="streak-label">%s</span>
	</div>`, streak.CurrentCount, streakLabel)

	h.WriteHTML(w, html)
}
