package logging

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

// Notifier define o comportamento de qualquer serviço de alerta
type Notifier interface {
	Notify(message string) error
}

// GotifyService implementa o Notifier para o Gotify
type GotifyService struct {
	URL   string
	Token string
}

func (g *GotifyService) Notify(message string) error {
	client := &http.Client{Timeout: 5 * time.Second}
	
	payload := map[string]interface{}{
		"message":  message,
		"priority": 5,
		"title":    "Log Notification",
	}
	
	body, _ := json.Marshal(payload)
	resp, err := client.Post(g.URL+"/message?token="+g.Token, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}