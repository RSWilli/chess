package sse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// Event is the struct capable of being marshaled into the correct format expected by a Content-Type: "text/event-stream"
type Event struct {
	ID    string
	Event string
	Data  any
	Retry time.Duration
}

func (e Event) Bytes() []byte {
	var s bytes.Buffer
	if e.ID != "" {
		fmt.Fprintf(&s, "id: %s\n", e.ID)
	}
	if e.Event != "" {
		fmt.Fprintf(&s, "event: %s\n", e.Event)
	}
	if e.Retry > 0 {
		fmt.Fprintf(&s, "retry: %d\n", e.Retry/time.Millisecond)
	}
	fmt.Fprintf(&s, "data: ")

	if err := json.NewEncoder(&s).Encode(e.Data); err != nil {
		panic("json encoding event payload failed")
	}

	fmt.Fprintf(&s, "\n\n")
	return s.Bytes()
}
