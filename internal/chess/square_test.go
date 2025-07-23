package chess_test

import (
	"fmt"
	"testing"

	"github.com/rswilli/chess/internal/chess"
)

func TestParseTile(t *testing.T) {
	squares := []string{
		"e4",
		"a8",
	}

	for _, square := range squares {
		tile, err := chess.ParseSquare(square)

		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(tile.Debug())
		fmt.Println(tile)

		if tile.String() != square {
			t.Fatalf("expected %s, got %s", square, tile.String())
		}
	}
}
