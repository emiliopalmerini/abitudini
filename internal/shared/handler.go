package shared

import (
	"net/http"
	"strconv"
)

type BaseHandler struct{}

// WriteHTML writes HTML response with proper content-type header
func (h *BaseHandler) WriteHTML(w http.ResponseWriter, html string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

// WriteError writes an HTTP error response
func (h *BaseHandler) WriteError(w http.ResponseWriter, message string, status int) {
	http.Error(w, message, status)
}

// ExtractIntPathParam extracts and converts a path parameter to integer
func (h *BaseHandler) ExtractIntPathParam(r *http.Request, paramName string) (int, error) {
	return strconv.Atoi(r.PathValue(paramName))
}

// ValidateMethod checks if request method matches expected method
func (h *BaseHandler) ValidateMethod(w http.ResponseWriter, r *http.Request, expectedMethod string) bool {
	if r.Method != expectedMethod {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return false
	}
	return true
}
