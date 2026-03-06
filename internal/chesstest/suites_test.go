package chesstest_test

import (
	"testing"

	"github.com/rswilli/chess/internal/chess"
	"github.com/rswilli/chess/internal/chesstest"
)

func Test(t *testing.T) {
	for _, fen := range chesstest.Suites {
		board, err := chess.NewPositionFromFEN(fen, nil)

		if err != nil {
			t.Fatalf("invalid FEN %s: %v", fen, err)
		}

		_ = board

		t.Log(board.String())
	}
}
