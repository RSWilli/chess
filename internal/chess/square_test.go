package chess

import (
	"testing"
)

func TestParseTile(t *testing.T) {

	for rank := range ranks {
		for file := range files {
			square := files[file:file+1] + ranks[rank:rank+1]

			tile, err := ParseSquare(square)

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
}
