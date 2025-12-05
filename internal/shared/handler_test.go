package shared

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestWriteHTML(t *testing.T) {
	h := &BaseHandler{}
	w := httptest.NewRecorder()
	html := "<div>test</div>"

	h.WriteHTML(w, html)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Errorf("expected Content-Type 'text/html; charset=utf-8', got '%s'", ct)
	}
	if w.Body.String() != html {
		t.Errorf("expected body '%s', got '%s'", html, w.Body.String())
	}
}

func TestWriteError(t *testing.T) {
	h := &BaseHandler{}
	w := httptest.NewRecorder()
	message := "test error"
	status := http.StatusBadRequest

	h.WriteError(w, message, status)

	if w.Code != status {
		t.Errorf("expected status %d, got %d", status, w.Code)
	}
	if !strings.Contains(w.Body.String(), message) {
		t.Errorf("expected body to contain '%s', got '%s'", message, w.Body.String())
	}
}

func TestWriteError_InternalServerError(t *testing.T) {
	h := &BaseHandler{}
	w := httptest.NewRecorder()
	message := "internal error"
	status := http.StatusInternalServerError

	h.WriteError(w, message, status)

	if w.Code != status {
		t.Errorf("expected status %d, got %d", status, w.Code)
	}
}

func TestExtractIntPathParam_Valid(t *testing.T) {
	h := &BaseHandler{}
	req := httptest.NewRequest("GET", "/test/42", nil)
	// Manually set path value
	req.SetPathValue("id", "42")

	val, err := h.ExtractIntPathParam(req, "id")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if val != 42 {
		t.Errorf("expected 42, got %d", val)
	}
}

func TestExtractIntPathParam_Invalid(t *testing.T) {
	h := &BaseHandler{}
	req := httptest.NewRequest("GET", "/test/abc", nil)
	req.SetPathValue("id", "abc")

	_, err := h.ExtractIntPathParam(req, "id")

	if err == nil {
		t.Error("expected error for invalid integer")
	}
	if _, ok := err.(*strconv.NumError); !ok {
		t.Errorf("expected NumError, got %T", err)
	}
}

func TestExtractIntPathParam_Zero(t *testing.T) {
	h := &BaseHandler{}
	req := httptest.NewRequest("GET", "/test/0", nil)
	req.SetPathValue("id", "0")

	val, err := h.ExtractIntPathParam(req, "id")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if val != 0 {
		t.Errorf("expected 0, got %d", val)
	}
}

func TestExtractIntPathParam_Negative(t *testing.T) {
	h := &BaseHandler{}
	req := httptest.NewRequest("GET", "/test/-5", nil)
	req.SetPathValue("id", "-5")

	val, err := h.ExtractIntPathParam(req, "id")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if val != -5 {
		t.Errorf("expected -5, got %d", val)
	}
}

func TestValidateMethod_Match(t *testing.T) {
	h := &BaseHandler{}
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	valid := h.ValidateMethod(w, req, "GET")

	if !valid {
		t.Error("expected ValidateMethod to return true for matching method")
	}
	if w.Code != http.StatusOK {
		t.Errorf("expected no error written, got status %d", w.Code)
	}
}

func TestValidateMethod_Mismatch(t *testing.T) {
	h := &BaseHandler{}
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	valid := h.ValidateMethod(w, req, "POST")

	if valid {
		t.Error("expected ValidateMethod to return false for mismatched method")
	}
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestValidateMethod_AllMethods(t *testing.T) {
	methods := []string{
		"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS",
	}

	for _, method := range methods {
		h := &BaseHandler{}
		req := httptest.NewRequest(method, "/test", nil)
		w := httptest.NewRecorder()

		valid := h.ValidateMethod(w, req, method)

		if !valid {
			t.Errorf("expected ValidateMethod to return true for %s", method)
		}
	}
}
