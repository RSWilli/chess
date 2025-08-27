package chess_test

import (
	"fmt"
	"testing"

	"github.com/rswilli/chess/internal/chess"
)

func TestBoard(t *testing.T) {
	b := chess.NewGame()

	fmt.Println(b)
}

func TestBoardFENs(t *testing.T) {
	fens := []string{
		"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
		"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2",
		"rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2",
	}

	for _, fen := range fens {
		b, err := chess.NewGameFromFEN(fen)

		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(b)
	}
}
