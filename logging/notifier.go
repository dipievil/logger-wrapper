package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// GotifyService is a struct that implements the Notifier interface, allowing it to send notifications to a Gotify server.
type GotifyService struct {
	URL      string
	Token    string
	Title    string
	Priority *int
}

// Validate checks the configuration of the GotifyService.
// It ensures that the URL and token are properly set and
// that the priority is within the valid range. It also
// checks if the Gotify server is reachable.
func (g *GotifyService) Validate() (bool, error) {
	if g == nil {
		return false, fmt.Errorf("Notifier cannot be nil!")
	}

	if g.URL == "" && g.Token == "" {
		return false, nil
	}

	if g.Title == "" {
		g.Title = "Logger Wrapper Notification"
	}

	if g.URL == "" && g.Token != "" {
		return false, fmt.Errorf("No Gotify URL is set for sending notifications with provided token")
	}

	if g.Token == "" && g.URL != "" {
		return false, fmt.Errorf("Gotify token is required for sending notifications to %s", g.URL)
	}

	g.URL = strings.TrimRight(g.URL, "/")

	err := g.checkUrl()
	if err != nil {
		return false, err
	}

	defaultPriority := 5
	if g.Priority == nil || *g.Priority < 1 || *g.Priority > 10 {
		g.Priority = &defaultPriority
	}
	return true, nil
}

func (g *GotifyService) GetHost() string {
	return g.URL
}

func (g *GotifyService) Notify(message string) error {

	client := &http.Client{Timeout: 5 * time.Second}

	payload := map[string]any{
		"message":  message,
		"priority": g.Priority,
		"title":    g.Title,
	}

	body, _ := json.Marshal(payload)

	resp, err := client.Post(g.URL+"/message?token="+g.Token, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (g *GotifyService) checkUrl() error {
	if g.URL != "" {
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Get(g.URL + "/version")
		if err != nil {
			return fmt.Errorf("Failed to reach Gotify server at %s: %v", g.URL, err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Gotify server at %s returned non-OK status: %s", g.URL, resp.Status)
		}
	}
	return nil
}
