package streak

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockStreakHandlerService struct {
	streak *Streak
	err    error
}

func (m *mockStreakHandlerService) GetByHabitID(habitID int) (*Streak, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.streak, nil
}

func TestStreakGetByHabitID_Success(t *testing.T) {
	streak := &Streak{HabitID: 1, CurrentCount: 5}
	service := &mockStreakHandlerService{streak: streak}
	handler := NewHandler(service)

	req := httptest.NewRequest("GET", "/api/habits/1/streak", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.GetByHabitID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Errorf("expected HTML content type, got %s", ct)
	}
	body := w.Body.String()
	if !strings.Contains(body, "5") {
		t.Error("expected streak count 5 in response")
	}
	if !strings.Contains(body, "days streak") {
		t.Error("expected 'days streak' label in response")
	}
}

func TestStreakGetByHabitID_SingleDay(t *testing.T) {
	streak := &Streak{HabitID: 1, CurrentCount: 1}
	service := &mockStreakHandlerService{streak: streak}
	handler := NewHandler(service)

	req := httptest.NewRequest("GET", "/api/habits/1/streak", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.GetByHabitID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "1") {
		t.Error("expected streak count 1 in response")
	}
	if !strings.Contains(body, "day streak") {
		t.Error("expected 'day streak' label in response")
	}
}

func TestStreakGetByHabitID_NoStreak(t *testing.T) {
	streak := &Streak{HabitID: 1, CurrentCount: 0}
	service := &mockStreakHandlerService{streak: streak}
	handler := NewHandler(service)

	req := httptest.NewRequest("GET", "/api/habits/1/streak", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.GetByHabitID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "0") {
		t.Error("expected streak count 0 in response")
	}
}

func TestGetByHabitID_WrongMethod(t *testing.T) {
	handler := NewHandler(&mockStreakHandlerService{})
	req := httptest.NewRequest("POST", "/api/habits/1/streak", nil)
	w := httptest.NewRecorder()

	handler.GetByHabitID(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestGetByHabitID_InvalidID(t *testing.T) {
	handler := NewHandler(&mockStreakHandlerService{})
	req := httptest.NewRequest("GET", "/api/habits/abc/streak", nil)
	req.SetPathValue("id", "abc")
	w := httptest.NewRecorder()

	handler.GetByHabitID(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestStreakGetByHabitID_ServiceError(t *testing.T) {
	service := &mockStreakHandlerService{err: errors.New("service failed")}
	handler := NewHandler(service)

	req := httptest.NewRequest("GET", "/api/habits/1/streak", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.GetByHabitID(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}

func TestStreakGetByHabitID_LargeStreak(t *testing.T) {
	streak := &Streak{HabitID: 1, CurrentCount: 365}
	service := &mockStreakHandlerService{streak: streak}
	handler := NewHandler(service)

	req := httptest.NewRequest("GET", "/api/habits/1/streak", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.GetByHabitID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "365") {
		t.Error("expected streak count 365 in response")
	}
}
