package logging

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNotifierValidate(t *testing.T) {
	tests := []struct {
		name     string
		notifier *GotifyService
		wantErr  bool
	}{
		{
			name:     "Nil notifier",
			notifier: nil,
			wantErr:  true,
		},
		{
			name: "Valid notifier",
			notifier: &GotifyService{
				URL:   "http://localhost:8080",
				Token: "testtoken",
				Title: "Test Notification",
			},
			wantErr: false,
		},
		{
			name: "Missing URL with token",
			notifier: &GotifyService{
				Token: "testtoken",
				Title: "Test Notification",
			},
			wantErr: true,
		},
		{
			name: "Missing token with URL",
			notifier: &GotifyService{
				URL:   "http://localhost:8080",
				Title: "Test Notification",
			},
			wantErr: true,
		},
		{
			name: "Invalid priority",
			notifier: &GotifyService{
				URL:      "http://localhost:8080",
				Token:    "testtoken",
				Title:    "Test Notification",
				Priority: func() *int { p := 11; return &p }(),
			},
			wantErr: false, // Should default to 5
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.notifier.Validate()
			if got && tt.wantErr {
				t.Errorf("Expected error but got success")
			}
			if !got && !tt.wantErr {
				t.Errorf("Expected success but got error: %v", err)
			}
		})
	}
}

func TestNotifierNotifyInvalidHost(t *testing.T) {
	notifier := &GotifyService{
		URL:   "http://localhost:8080",
		Token: "testtoken",
		Title: "Test Notification",
	}

	err := notifier.Notify("This is a test notification")
	if err == nil {
		t.Errorf("Expected error but got success")
	} else {
		if strings.Contains(err.Error(), "connect: connection refused") {
			t.Log("Received expected connection refused error")
		}
		t.Log("Notification failed as expected with error:", err)
	}
}

func TestGetHost(t *testing.T) {

	host := "http://localhost:8080"

	notifier := &GotifyService{
		URL: host,
	}

	gotHost := notifier.GetHost()
	expected := host
	if gotHost != expected {
		t.Errorf("Expected GetHost to return '%s', but got '%s'", expected, gotHost)
	} else {
		t.Log("GetHost returned the expected value")
	}
}

func TestCheckUrl(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{"OK response", http.StatusOK, false},
		{"Non-OK response", http.StatusInternalServerError, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			g := &GotifyService{URL: server.URL}
			err := g.checkUrl()

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
