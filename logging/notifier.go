package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// GotifyService is a struct that implements the Notifier interface, allowing it to send notifications to a Gotify server.
type GotifyService struct {
	URL      string
	Token    string
	Title    string
	Priority *int
}

// Validate checks the configuration of the GotifyService and ensures that it has the necessary information to send notifications.
func (g *GotifyService) Validate() (bool, error) {
	if g == nil {
		return false, fmt.Errorf("Notifier cannot be nil!")
	}

	defaultPriority := 5

	if g.Title == "" {
		g.Title = "Logger Wrapper Notification"
	}

	if g.URL == "" && g.Token != "" {
		return false, fmt.Errorf("No Gotify URL is set for sending notifications with provided token")
	}

	if g.Priority == nil || *g.Priority < 1 || *g.Priority > 10 {
		g.Priority = &defaultPriority
	}

	if g.Token == "" && g.URL != "" {
		return false, fmt.Errorf("Gotify token is required for sending notifications to %s", g.URL)
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
