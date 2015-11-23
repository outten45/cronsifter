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
func (e *Event) JSONReader() io.Reader {
	eventB, _ := json.Marshal(e)
	reader := bytes.NewReader(eventB)
	return reader
}
