package uci_test

import (
	"fmt"
	"testing"

	"github.com/rswilli/chess/internal/uci"
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

	err = sf.Position(uci.StartPositionFEN, nil)

	if err != nil {
		t.Fatalf("could not set position: %v", err)
	}

	res, err := sf.Perft(5)

	if err != nil {
		t.Fatalf("failed to run perft: %v", err)
	}

	fmt.Printf("%#v\n", res)

	gores, err := sf.Go(uci.GoOptions{
		Depth: 5,
	})

	if err != nil {
		t.Fatalf("failed to run go: %v", err)
	}

	fmt.Printf("%#v\n", gores)
}
