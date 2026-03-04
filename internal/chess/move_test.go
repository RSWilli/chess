package chess_test

import (
	"fmt"
	"testing"

	"github.com/rswilli/chess/internal/chess"
)

func TestParseMove(t *testing.T) {
	moves := []string{
		"e4e5",
		"a4e5",
		"g4h8",
		"a7a8q",
	}

	for _, m := range moves {
		move, err := chess.ParseMove(m)

		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(move)
	}
}
