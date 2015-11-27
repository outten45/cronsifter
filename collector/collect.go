package collector

import (
	"bytes"
	"encoding/json"
	"io"
)

// Event contains the information to be sent to the collector.
type Event struct {
	Service     string
	Host        string
	State       string
	Time        int
	Description string
	Tags        []string
	Token       string
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
