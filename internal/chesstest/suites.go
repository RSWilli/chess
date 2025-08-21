package chesstest

import (
	_ "embed"
	"strings"
	"testing"
)

//go:embed testpositions.txt
var positions string

var Suites []string

func init() {
	for l := range strings.SplitSeq(positions, "\n") {
		if len(l) == 0 || l[0] == '#' {
			continue
		}

		Suites = append(Suites, l+" 0 1")
	}
}

func RunAll(t *testing.T, f func(t *testing.T, fen string)) {
	for _, fen := range Suites {
		t.Run(fen, func(t *testing.T) {
			f(t, fen)
		})
	}
}
