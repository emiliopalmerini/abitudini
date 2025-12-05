package record

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/epalmerini/abitudini/internal/habit"
)

type mockRecordHandlerService struct {
	contributions []ContributionDay
	habit         *habit.Habit
	err           error
}

func (m *mockRecordHandlerService) MarkDoneToday(habitID int) error {
	return m.err
}

func (m *mockRecordHandlerService) GetContributionData(habitID int, from, to time.Time) ([]ContributionDay, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.contributions, nil
}

func (m *mockRecordHandlerService) GetHabit(habitID int) (*habit.Habit, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.habit, nil
}

func TestRecordMarkDoneToday_Success(t *testing.T) {
	h := &habit.Habit{ID: 1, Description: "Test"}
	service := &mockRecordHandlerService{habit: h}
	handler := NewHandler(service)

	req := httptest.NewRequest("POST", "/api/habits/1/done-today", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.MarkDoneToday(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Errorf("expected HTML content type, got %s", ct)
	}
	if !strings.Contains(w.Body.String(), "success") {
		t.Error("expected success message in response")
	}
}

func TestRecordMarkDoneToday_WrongMethod(t *testing.T) {
	handler := NewHandler(&mockRecordHandlerService{})
	req := httptest.NewRequest("GET", "/api/habits/1/done-today", nil)
	w := httptest.NewRecorder()

	handler.MarkDoneToday(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestRecordMarkDoneToday_InvalidID(t *testing.T) {
	handler := NewHandler(&mockRecordHandlerService{})
	req := httptest.NewRequest("POST", "/api/habits/abc/done-today", nil)
	req.SetPathValue("id", "abc")
	w := httptest.NewRecorder()

	handler.MarkDoneToday(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestRecordMarkDoneToday_ServiceError(t *testing.T) {
	service := &mockRecordHandlerService{err: errors.New("record failed")}
	handler := NewHandler(service)

	req := httptest.NewRequest("POST", "/api/habits/1/done-today", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.MarkDoneToday(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}

func TestRecordMarkDoneToday_GetHabitError(t *testing.T) {
	service := &mockRecordHandlerService{err: errors.New("fetch failed")}
	handler := NewHandler(service)

	req := httptest.NewRequest("POST", "/api/habits/1/done-today", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.MarkDoneToday(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}

func TestGetContribution_Success(t *testing.T) {
	now := time.Now()
	contributions := []ContributionDay{
		{Date: now, Completed: true},
		{Date: now.AddDate(0, 0, -1), Completed: false},
	}
	service := &mockRecordHandlerService{contributions: contributions}
	handler := NewHandler(service)

	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	
	req := httptest.NewRequest("GET", "/api/habits/1/contribution?from="+yesterday+"&to="+today, nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.GetContribution(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "contribution") {
		t.Error("expected contribution HTML in response")
	}
}

func TestRecordGetContribution_WrongMethod(t *testing.T) {
	handler := NewHandler(&mockRecordHandlerService{})
	req := httptest.NewRequest("POST", "/api/habits/1/contribution", nil)
	w := httptest.NewRecorder()

	handler.GetContribution(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestRecordGetContribution_InvalidID(t *testing.T) {
	handler := NewHandler(&mockRecordHandlerService{})
	req := httptest.NewRequest("GET", "/api/habits/xyz/contribution", nil)
	req.SetPathValue("id", "xyz")
	w := httptest.NewRecorder()

	handler.GetContribution(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestRecordGetContribution_NoDateParams(t *testing.T) {
	now := time.Now()
	contributions := []ContributionDay{
		{Date: now, Completed: true},
	}
	service := &mockRecordHandlerService{contributions: contributions}
	handler := NewHandler(service)

	// No from/to params - should use defaults
	req := httptest.NewRequest("GET", "/api/habits/1/contribution", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.GetContribution(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestRecordGetContribution_Empty(t *testing.T) {
	service := &mockRecordHandlerService{contributions: []ContributionDay{}}
	handler := NewHandler(service)

	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	
	req := httptest.NewRequest("GET", "/api/habits/1/contribution?from="+yesterday+"&to="+today, nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.GetContribution(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestRecordGetContribution_ServiceError(t *testing.T) {
	service := &mockRecordHandlerService{err: errors.New("fetch failed")}
	handler := NewHandler(service)

	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	
	req := httptest.NewRequest("GET", "/api/habits/1/contribution?from="+yesterday+"&to="+today, nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.GetContribution(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}

func TestRecordGetContribution_InvalidDateFormat(t *testing.T) {
	service := &mockRecordHandlerService{contributions: []ContributionDay{}}
	handler := NewHandler(service)

	// Invalid date format should use defaults
	req := httptest.NewRequest("GET", "/api/habits/1/contribution?from=invalid&to=invalid", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.GetContribution(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}


