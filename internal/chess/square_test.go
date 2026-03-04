package chess_test

import (
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

		t.Log(tile.Debug())
		t.Log(tile)

		if tile.String() != square {
			t.Fatalf("expected %s, got %s", square, tile.String())
		}
	}
}
