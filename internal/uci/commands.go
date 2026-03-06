package uci

import (
	"bytes"
	"encoding"
	"fmt"
	"maps"
	"slices"
)

type info struct {
	key   string
	value string
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (i *info) UnmarshalText(text []byte) error {
	parts := bytes.SplitN(text, []byte(" "), 3)

	if string(parts[0]) != "info" {
		return fmt.Errorf("not an info message")
	}

	if len(parts) != 3 {
		return fmt.Errorf("info message wrong format")
	}

	*i = info{
		key:   string(parts[1]),
		value: string(parts[2]),
	}

	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (i *info) MarshalText() (text []byte, err error) {
	return fmt.Appendf(nil, "info %s %s", i.key, i.value), nil
}

func (i *info) String() string {
	return fmt.Sprintf("info %s %s", i.key, i.value)
}

var _ encoding.TextUnmarshaler = &info{}
var _ encoding.TextMarshaler = &info{}

type perftResult struct {
	total int
	moves map[string]int
}

// MarshalText implements encoding.TextMarshaler.
func (p *perftResult) String() string {
	var text []byte

	for _, move := range slices.Sorted(maps.Keys(p.moves)) {
		count := p.moves[move]

		text = fmt.Appendf(text, "%s: %d\n", move, count)
	}

	text = append(text, '\n')
	text = fmt.Appendf(text, "Nodes searched: %d\n\n", p.total)

	return string(text)
}

type bestmove struct {
	bestMove string
	ponder   string
}

func (bm *bestmove) String() string {
	if bm.ponder == "" {
		return fmt.Sprintf("bestmove %s", bm.bestMove)
	}
	return fmt.Sprintf("bestmove %s ponder %s", bm.bestMove, bm.ponder)
}
