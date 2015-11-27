package cronsifter

import (
	"log"
	"net/http"

	"github.com/outten45/cronsifter/collector"
)

// Context contains the server to send the notification to.
type Context struct {
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
func Notify(c *Context, e *collector.Event) error {
	json, err := e.JSONReader()
	if err != nil {
		log.Printf("Unable to create JSON from: %#v", e)
		return err
	}
	resp, err := http.Post(c.URL, "application/json", json)
	if err != nil {
		log.Printf("Error sending to server: %#v\n", err)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Status code returned was not ok. It was %v. [%s]\n", resp.StatusCode, resp.Status)
	}
	return nil
}
