package chess_test

import (
	"iter"
	"maps"
	"testing"

	"github.com/rswilli/chess/internal/chess"
	"github.com/rswilli/chess/internal/chesstest"
	"github.com/rswilli/chess/internal/uci"
)

func BenchmarkPerft(t *testing.B) {
	local := chess.NewEngine()

	total := 0
	for t.Loop() {
		for _, fen := range chesstest.Suites {
			local.Position(fen, nil)

			visited, _, err := local.Perft(5)

			if err != nil {
				t.Fatalf("perft ran into an error: %v", err)
			}

			total += visited
		}
	}

	t.Logf("visited %d positions in total", total)
}

func TestPerft(t *testing.T) {
	stockfish, err := uci.NewStockfish()

	if err != nil {
		t.Fatal(err)
	}

	depth := 5

	local := chess.NewEngine()

	chesstest.RunAll(t, func(t *testing.T, fen string) {
		comparePerft(t, stockfish, local, fen, depth, nil)
	})
}

func comparePerft(t *testing.T, stockfish, local uci.Engine, fen string, depth int, moves []string) {
	var err error
	err = stockfish.NewGame()

	if err != nil {
		t.Fatal(err)
	}

	err = local.NewGame()

	if err != nil {
		t.Fatal(err)
	}

	err = local.Position(fen, moves)
	if err != nil {
		t.Fatal(err)
	}
	err = stockfish.Position(fen, moves)
	if err != nil {
		t.Fatal(err)
	}

	sfTotal, sfMoves, err := stockfish.Perft(depth)

	if err != nil {
		t.Fatal(err)
	}

	localTotal, localMoves, err := local.Perft(depth)

	if err != nil {
		t.Fatal(err)
	}

	if sfTotal == localTotal && maps.Equal(sfMoves, localMoves) && maps.Equal(localMoves, sfMoves) {
		return
	}

	t.Logf("pos: %s moves %v", fen, moves)

	actual := localMoves
	expected := sfMoves

	tooMany := diff(actual, expected)
	missing := diff(expected, actual)

	switch {
	case len(tooMany) == 0 && len(missing) == 0:
		// ok
	case len(tooMany) > 0 && len(missing) == 0:
		t.Fatalf("local produced illegal moves: %v", tooMany)
		return
	case len(tooMany) == 0 && len(missing) > 0:
		t.Fatalf("local missing legal moves: %v", missing)
		return
	case len(tooMany) > 0 && len(missing) > 0:
		t.Fatalf("local produced illegal moves: %v and missing moves: %v", tooMany, missing)
		return
	}

	different := diffValues(actual, expected)

	t.Logf("local engine produced a wrong sum for the moves: %v", different)

	move := first(maps.Keys(different))

	doneMoves := append([]string(nil), moves...)
	doneMoves = append(doneMoves, move)

	comparePerft(t, stockfish, local, fen, depth-1, doneMoves)
}

func first[V any](i iter.Seq[V]) V {
	for i := range i {
		return i
	}

	panic("empty iter")
}

// diff computes a-b
func diff(a, b map[string]int) map[string]int {
	c := make(map[string]int, len(a))

	for k, v := range a {
		if _, ok := b[k]; ok {
			continue
		}

		c[k] = v
	}

	return c
}

// diffValues returns a map containing all keys from expected where the value in actual differs
func diffValues(actual, expected map[string]int) map[string]struct{} {
	c := make(map[string]struct{}, len(actual))

	for k, v1 := range expected {
		if v2, ok := actual[k]; ok && v1 == v2 {
			continue
		}

		c[k] = struct{}{}
	}

	return c
}
