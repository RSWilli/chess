package uci_test

import (
	"fmt"
	"testing"

	"github.com/rswilli/chess/internal/uci"
	"github.com/rswilli/chess/internal/uci/search"
)

func Test(t *testing.T) {
	sf, err := uci.NewStockfish()

	if err != nil {
		t.Fatalf("could not initialize: %v", err)
	}
	defer sf.Close()

	err = sf.NewGame()
	if err != nil {
		t.Fatal("could not start new game")
	}

	err = sf.Position(uci.StartPosition, nil)

	if err != nil {
		t.Fatalf("could not set position: %v", err)
	}

	total, moves, err := sf.Perft(5)

	if err != nil {
		t.Fatalf("failed to run perft: %v", err)
	}

	fmt.Printf("%d %#v\n", total, moves)

	bm, ponder := sf.Go(search.Options{
		Depth: 5,
	})

	if err != nil {
		t.Fatalf("failed to run go: %v", err)
	}

	fmt.Printf("%s %s\n", bm, ponder)
}
