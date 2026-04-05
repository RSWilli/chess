package chess_test

import (
	"fmt"
	"testing"

	"github.com/rswilli/chess/internal/chess"
	"github.com/rswilli/chess/internal/chesstest"
)

func TestBoard(t *testing.T) {
	b := chess.NewPosition()

	fmt.Println(b)
}

func TestBoardFENs(t *testing.T) {
	chesstest.RunAll(t, func(t *testing.T, fen string) {
		b, err := chess.NewPositionFromFEN(fen, nil)

		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(b.ASCIIArt())

		marshaledFEN := b.FEN()

		if fen != marshaledFEN {
			t.Fatalf("expected %q, got: %q", fen, marshaledFEN)
		}
	})
}
