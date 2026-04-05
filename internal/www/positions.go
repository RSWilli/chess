package www

import (
	_ "embed"
	"strings"
)

//go:embed positions.txt
var positions string

type startPosition struct {
	Name string
	FEN  string
}

var startPositions []startPosition

func init() {
	for line := range strings.SplitSeq(positions, "\n") {
		if line == "" {
			continue
		}

		name, fen, _ := strings.Cut(line, ": ")
		startPositions = append(startPositions, startPosition{
			Name: name,
			FEN:  fen,
		})
	}
}
