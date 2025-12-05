package habit

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockHandlerService struct {
	createID int
	habits   []Habit
	habit    *Habit
	err      error
}

func (m *mockHandlerService) Create(h *Habit) (int, error) {
	if m.err != nil {
		return 0, m.err
	}
	return m.createID, nil
}

func (m *mockHandlerService) Update(h *Habit) error {
	return m.err
}

func (m *mockHandlerService) GetByID(habitID int) (*Habit, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.habit, nil
}

func (m *mockHandlerService) GetAll() ([]Habit, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.habits, nil
}

func (m *mockHandlerService) Delete(habitID int) error {
	return m.err
}

func TestCreate_Success(t *testing.T) {
	habit := &Habit{ID: 1, Description: "Test", Color: "blue"}
	service := &mockHandlerService{createID: 1, habit: habit}
	handler := NewHandler(service)

	req := httptest.NewRequest("POST", "/api/habits", strings.NewReader("description=Test&start_date=2025-01-01&color=blue"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Errorf("expected HTML content type, got %s", ct)
	}
}

func TestCreate_WrongMethod(t *testing.T) {
	handler := NewHandler(&mockHandlerService{})
	req := httptest.NewRequest("GET", "/api/habits", nil)
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}



func TestCreate_ServiceError(t *testing.T) {
	service := &mockHandlerService{err: errors.New("service failed")}
	handler := NewHandler(service)

	req := httptest.NewRequest("POST", "/api/habits", strings.NewReader("description=Test&start_date=2025-01-01&color=blue"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}

func TestGetAll_Success(t *testing.T) {
	habits := []Habit{
		{ID: 1, Description: "Habit 1"},
		{ID: 2, Description: "Habit 2"},
	}
	service := &mockHandlerService{habits: habits}
	handler := NewHandler(service)

	req := httptest.NewRequest("GET", "/api/habits", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestGetAll_WrongMethod(t *testing.T) {
	handler := NewHandler(&mockHandlerService{})
	req := httptest.NewRequest("POST", "/api/habits", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestGetAll_Error(t *testing.T) {
	service := &mockHandlerService{err: errors.New("service failed")}
	handler := NewHandler(service)

	req := httptest.NewRequest("GET", "/api/habits", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}

func TestGetAll_Empty(t *testing.T) {
	service := &mockHandlerService{habits: []Habit{}}
	handler := NewHandler(service)

	req := httptest.NewRequest("GET", "/api/habits", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "empty-state") {
		t.Error("expected empty state HTML in response")
	}
}

func TestGetByID_Success(t *testing.T) {
	habit := &Habit{ID: 1, Description: "Test"}
	service := &mockHandlerService{habit: habit}
	handler := NewHandler(service)

	req := httptest.NewRequest("GET", "/api/habits/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestGetByID_WrongMethod(t *testing.T) {
	handler := NewHandler(&mockHandlerService{})
	req := httptest.NewRequest("POST", "/api/habits/1", nil)
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestGetByID_InvalidID(t *testing.T) {
	handler := NewHandler(&mockHandlerService{})
	req := httptest.NewRequest("GET", "/api/habits/abc", nil)
	req.SetPathValue("id", "abc")
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	service := &mockHandlerService{err: errors.New("not found")}
	handler := NewHandler(service)

	req := httptest.NewRequest("GET", "/api/habits/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestUpdate_Success(t *testing.T) {
	habit := &Habit{ID: 1, Description: "Updated"}
	service := &mockHandlerService{habit: habit}
	handler := NewHandler(service)

	req := httptest.NewRequest("PUT", "/api/habits/1", strings.NewReader("description=Updated&start_date=2025-01-01&color=blue"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestUpdate_WrongMethod(t *testing.T) {
	handler := NewHandler(&mockHandlerService{})
	req := httptest.NewRequest("GET", "/api/habits/1", nil)
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestUpdate_InvalidID(t *testing.T) {
	handler := NewHandler(&mockHandlerService{})
	req := httptest.NewRequest("PUT", "/api/habits/xyz", nil)
	req.SetPathValue("id", "xyz")
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}



func TestUpdate_ServiceError(t *testing.T) {
	service := &mockHandlerService{err: errors.New("service failed")}
	handler := NewHandler(service)

	req := httptest.NewRequest("PUT", "/api/habits/1", strings.NewReader("description=Test&start_date=2025-01-01&color=blue"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}

func TestDelete_Success(t *testing.T) {
	service := &mockHandlerService{}
	handler := NewHandler(service)

	req := httptest.NewRequest("DELETE", "/api/habits/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestDelete_WrongMethod(t *testing.T) {
	handler := NewHandler(&mockHandlerService{})
	req := httptest.NewRequest("GET", "/api/habits/1", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestDelete_InvalidID(t *testing.T) {
	handler := NewHandler(&mockHandlerService{})
	req := httptest.NewRequest("DELETE", "/api/habits/bad", nil)
	req.SetPathValue("id", "bad")
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestDelete_ServiceError(t *testing.T) {
	service := &mockHandlerService{err: errors.New("service failed")}
	handler := NewHandler(service)

	req := httptest.NewRequest("DELETE", "/api/habits/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}


