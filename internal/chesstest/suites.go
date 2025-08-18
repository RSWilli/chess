package chesstest

import (
	_ "embed"
	"strings"
)

//go:embed testpositions.txt
var positions string

var Suites []string

func init() {
	lines := strings.Split(positions, "\n")

	for _, l := range lines {
		if len(l) == 0 || l[0] == '#' {
			continue
		}

		Suites = append(Suites, l+" 0 1")
	}
}
