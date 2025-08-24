package main

import (
	"fmt"
	"os"

	"github.com/rswilli/chess/internal/chess"
	"github.com/rswilli/chess/internal/uci"
)

func main() {
	uciServer := uci.NewServer(os.Stdin, os.Stdout, chess.NewEngine())

	err := uciServer.Run()

	fmt.Fprintf(os.Stderr, "error from server: %v", err)
}
