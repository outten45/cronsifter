package collector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Event contains the information to be sent to the collector.
type Event struct {
	Service string
	Host    string
	State   string
	Time    int
	Message string
	Tags    []string
	Token   string
}

// NewEvent creates a collector.Event that contains some default
// values set like the Host and Time.
func NewEvent(service, state, message, token string, tags []string) *Event {
	e := &Event{
		State:   state,
		Message: message,
		Tags:    tags,
	}
	if h, err := os.Hostname(); err != nil {
		e.Host = fmt.Sprintf("Unknown(%v)", err)
	} else {
		e.Host = h
	}
	now := time.Now()
	e.Time = int(now.Unix())
	return e
}

// JSONReader returns a io.Reader with the json in it.
func (e *Event) JSONReader() (io.Reader, error) {
	eventB, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(eventB)
	return reader, nil
}

// ParseEvent takes the json string of an event and returns the Event.
func ParseEvent(str string) (*Event, error) {
	e := &Event{}
	if err := json.Unmarshal([]byte(str), &e); err != nil {
		return e, err
	}
	return e, nil
}
