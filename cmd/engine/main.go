package main

import (
	"fmt"
	"os"

	"github.com/rswilli/chess/internal/engine"
	"github.com/rswilli/chess/internal/uci"
)

func main() {
	uciServer := uci.NewServer(os.Stdin, os.Stdout, engine.NewEngine())

	err := uciServer.Run()

	fmt.Fprintf(os.Stderr, "error from server: %v", err)
}
