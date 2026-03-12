package cli

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/audstanley/david/app"
)

func TestWrapRecovery(t *testing.T) {
	cfg := &app.Config{
		Cors: app.Cors{
			Origin:      "http://example.com",
			Credentials: true,
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := wrapRecovery(handler, cfg)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	wrapped.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestWrapRecoveryPanic(t *testing.T) {
	cfg := &app.Config{
		Cors: app.Cors{
			Origin:      "http://example.com",
			Credentials: true,
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	wrapped := wrapRecovery(handler, cfg)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Should not crash
	wrapped.ServeHTTP(w, req)

	// Recovery should prevent panic from propagating
	if w.Code != http.StatusOK {
		t.Logf("Recovery handled panic, status code: %d", w.Code)
	}
}
