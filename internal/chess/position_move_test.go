package chess_test

import (
	"maps"
	"slices"
	"testing"

	"github.com/rswilli/chess/internal/chess"
	"github.com/rswilli/chess/internal/chesstest"
	"github.com/rswilli/chess/internal/uci"
)

func TestDoUndoMoves(t *testing.T) {
	chesstest.RunAll(t, func(t *testing.T, fen string) {
		board, err := chess.NewPositionFromFEN(fen)

		if err != nil {
			t.Fatal("could not init FEN: " + err.Error())
		}

		moves := board.GenerateMoves()

		copy := board.Copy()

		for _, m := range moves {
			copy.DoMove(m)
			copy.UndoMove()

			if !copy.Equals(board) {
				t.Logf("expected:\n%s", board.String())
				t.Logf("got:\n%s", copy.String())

				t.Fatalf("move %s was not undone correctly", m)
			}
		}
	})
}

func TestLegalMoveGen(t *testing.T) {
	stockfish, err := uci.NewStockfish()

	if err != nil {
		t.Fatal("could not init stockfish")
	}

	chesstest.RunAll(t, func(t *testing.T, fen string) {
		must(stockfish.NewGame())
		must(stockfish.Position(fen, nil))
		board, err := chess.NewPositionFromFEN(fen)

		if err != nil {
			t.Fatal("could not init FEN: " + err.Error())
		}

		expectedMoves, err := stockfish.Perft(1)
		expectedMoveList := slices.Sorted(maps.Keys(expectedMoves.Moves))

		if err != nil {
			t.Fatalf("Stockfish did not return moves: %v", err)
		}

		moves := board.GenerateMoves()

		generatedMoves := make(map[string]chess.Move)

		for _, m := range moves {
			generatedMoves[m.String()] = m
		}

		generatedMoveList := slices.Sorted(maps.Keys(generatedMoves))

		for m, move := range generatedMoves {
			_, ok := expectedMoves.Moves[m]

			if !ok {
				t.Logf("stockfish: %s", expectedMoveList)
				t.Logf("ours     : %s", generatedMoveList)
				t.Fatalf("generated a move that was not expected:\n%s\n%s\n%s\n%#v", fen, board.String(), m, move)
			}
		}

		for m := range expectedMoves.Moves {
			_, ok := generatedMoves[m]

			if !ok {
				t.Logf("stockfish: %s", expectedMoveList)
				t.Logf("ours     : %s", generatedMoveList)
				t.Fatalf("missing an expected move:\n%s\n%s\n%s", fen, board.String(), m)
			}
		}
	})
}

func must(err error) {
	if err != nil {
		panic("error: " + err.Error())
	}
}
