package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/rswilli/chess/internal/chess"
	"github.com/rswilli/chess/internal/uci"
)

var verbose = flag.Bool("v", false, "enable verbose logging")

func main() {
	flag.Parse()

	var enginelogSink io.Writer

	if *verbose {
		enginelogSink = os.Stderr
	}

	uciServer := uci.NewServer(os.Stdin, os.Stdout, chess.NewEngine(enginelogSink))

	err := uciServer.Run()

	fmt.Fprintf(os.Stderr, "error from server: %v", err)
}
