package chesstest_test

import (
	"fmt"
	"testing"

	"github.com/rswilli/chess/internal/chess"
	"github.com/rswilli/chess/internal/chesstest"
)

func Test(t *testing.T) {
	for _, fen := range chesstest.Suites {
		board, err := chess.NewFromFEN(fen)

		if err != nil {
			t.Fatalf("invalid FEN %s: %v", fen, err)
		}

		_ = board

		fmt.Println(board.String())
	}
}
