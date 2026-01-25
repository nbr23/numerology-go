package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFilterDigits(t *testing.T) {
	tests := []struct {
		input string
		want  []int
	}{
		{"123", []int{1, 2, 3}},
		{"abc123def", []int{1, 2, 3}},
		{"1a2b3c", []int{1, 2, 3}},
		{"", []int(nil)},
		{"abc", []int(nil)},
		{"000", []int{0, 0, 0}},
		{"12-34", []int{1, 2, 3, 4}},
		{"(123)", []int{1, 2, 3}},
	}

	for _, tt := range tests {
		got := filterDigits(tt.input)
		if len(got) != len(tt.want) {
			t.Errorf("filterDigits(%q) = %v, want %v", tt.input, got, tt.want)
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("filterDigits(%q) = %v, want %v", tt.input, got, tt.want)
				break
			}
		}
	}
}

func TestHandler(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantBody   bool
	}{
		{"valid input finds solution", "/23/2023", http.StatusOK, true},
		{"valid input finds solution with letters", "/23/abc2def0ghi2jkl3", http.StatusOK, true},
		{"impossible input", "/23/19", http.StatusNotFound, false},
		{"empty path", "/", http.StatusNotFound, false},
		{"invalid target", "/abc/123", http.StatusBadRequest, false},
		{"direct 23", "/23/23", http.StatusOK, true},
		{"complex input", "/23/123456", http.StatusOK, true},
		{"different target", "/42/123456", http.StatusOK, true},
		{"target 10", "/10/19", http.StatusOK, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			handler(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("handler(%s) status = %d, want %d", tt.path, w.Code, tt.wantStatus)
			}
			if tt.wantBody && w.Body.Len() == 0 {
				t.Errorf("handler(%s) expected body but got empty", tt.path)
			}
		})
	}
}

func TestHandlerResponseContent(t *testing.T) {
	req := httptest.NewRequest("GET", "/23/23", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/plain" {
		t.Errorf("expected Content-Type text/plain, got %s", contentType)
	}

	body := w.Body.String()
	if body != "23" {
		t.Errorf("expected body '23', got '%s'", body)
	}
}

func TestHandlerExpressionValidity(t *testing.T) {
	paths := []string{"/23/123456", "/23/2023", "/23/110615", "/23/987654"}

	for _, path := range paths {
		req := httptest.NewRequest("GET", path, nil)
		w := httptest.NewRecorder()

		handler(w, req)

		if w.Code == http.StatusOK {
			expr := w.Body.String()
			if !isValidExpression(expr) {
				t.Errorf("handler(%s) returned invalid expression: %s", path, expr)
			}
		}
	}
}

func TestHandlerDefaultDate(t *testing.T) {
	req := httptest.NewRequest("GET", "/23", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	expectedDate := time.Now().Format("02012006")
	if w.Code == http.StatusNotFound {
		t.Logf("No solution found for today's date: %s", expectedDate)
	} else if w.Code != http.StatusOK {
		t.Errorf("handler(/23) unexpected status = %d", w.Code)
	}
}

func isValidExpression(expr string) bool {
	if len(expr) == 0 {
		return false
	}
	hasDigit := false
	for _, c := range expr {
		if c >= '0' && c <= '9' {
			hasDigit = true
		} else if c != '+' && c != '-' && c != '*' && c != '/' {
			return false
		}
	}
	return hasDigit
}
