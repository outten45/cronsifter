package cronsifter

import (
	"log"
	"net/http"

	"github.com/outten45/cronsifter/collector"
)

// Server contains the server to send the notification to.
type Server struct {
	URL string
}

// NewEvent creates a collector.Event that contains some default
// values set like the Service, Host and Time.
func NewEvent(state, description, token string, tags []string) *collector.Event {
	e := &collector.Event{
		State:       state,
		Description: description,
		Tags:        tags,
	}
	return e
}

// Notify post the Event to the given server.
func Notify(s *Server, e *collector.Event) {
	resp, err := http.Post(s.URL, "application/json", e.JSONReader())
	if err != nil {
		log.Printf("Error sending to server: %#v\n", err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Status code returned was not ok. It was %v. [%s]\n", resp.StatusCode, resp.Status)
	}
}
