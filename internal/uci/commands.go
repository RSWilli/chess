package uci

import (
	"bytes"
	"encoding"
	"fmt"
	"maps"
	"slices"
)

type Info struct {
	Key   string
	Value string
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (i *Info) UnmarshalText(text []byte) error {
	parts := bytes.SplitN(text, []byte(" "), 3)

	if string(parts[0]) != "info" {
		return fmt.Errorf("not an info message")
	}

	if len(parts) != 3 {
		return fmt.Errorf("info message wrong format")
	}

	*i = Info{
		Key:   string(parts[1]),
		Value: string(parts[2]),
	}

	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (i *Info) MarshalText() (text []byte, err error) {
	return fmt.Appendf(nil, "info %s %s", i.Key, i.Value), nil
}

var _ encoding.TextUnmarshaler = &Info{}
var _ encoding.TextMarshaler = &Info{}

type PerftResult struct {
	Total int
	Moves map[string]int
}

// MarshalText implements encoding.TextMarshaler.
func (p *PerftResult) String() string {
	var text []byte

	for _, move := range slices.Sorted(maps.Keys(p.Moves)) {
		count := p.Moves[move]

		text = fmt.Appendf(text, "%s: %d\n", move, count)
	}

	text = append(text, '\n')
	text = fmt.Appendf(text, "Nodes searched: %d\n\n", p.Total)

	return string(text)
}

type GoResponse struct {
	BestMove string
	Ponder   string
}
