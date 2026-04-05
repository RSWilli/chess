package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rswilli/chess/internal/chess"
	"github.com/rswilli/chess/internal/mm"
	"github.com/rswilli/chess/internal/uci"
)

//go:embed positions.txt
var positionsList string
var positions = strings.Split(positionsList, "\n")

var rounds = flag.Int("n", 10, "the amount of rounds to play")
var opponentPath = flag.String("opponent", "", "the opponent to play against, aka path to the engine binary or 'stockfish'")

func main() {
	flag.Parse()

	if *opponentPath == "" {
		fmt.Fprintln(os.Stderr, "no opponent given")
		flag.Usage()
		os.Exit(2)
	}

	var opponent uci.Engine
	var err error

	if *opponentPath == "stockfish" {
		opponent, err = uci.NewStockfish()
	} else {
		opponent, err = uci.NewClient(*opponentPath)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "could not create opponent: %v\n", err)
		flag.Usage()
		os.Exit(2)
	}

	local := chess.NewEngine(nil)

	var stats stats

	start := time.Now()

	for round := range *rounds {
		fmt.Printf("starting round %d\n", round+1)
		for _, fen := range positions {
			fmt.Printf("running game %s\n", fen)
			stats.update(runGame(fen, local, opponent), whiteWin)
			fmt.Printf("running game rematch %s\n", fen)
			stats.update(runGame(fen, opponent, local), blackWin)
		}
	}

	fmt.Printf("against %s:\nTotal: %d\n won %d, lost %d, draws %d\n", *opponentPath, stats.total(), stats.wins, stats.losses, stats.draws)
	fmt.Printf("Took: %s", time.Since(start).String())
}

type result int

const (
	draw result = iota
	whiteWin
	blackWin
)

type stats struct {
	wins   int
	draws  int
	losses int
}

func (s stats) total() int {
	return s.wins + s.draws + s.losses
}

func (s *stats) update(result, expected result) {
	switch result {
	case expected:
		s.wins++
		fmt.Println("won")
	case draw:
		s.draws++
		fmt.Println("draw")
	default:
		s.losses++
		fmt.Println("lost")
	}
}

func runGame(fen string, white, black uci.Engine) result {
	game, _ := mm.NewGameWithFEN(fen, white, black)

	for {
		game.Move()

		if game.Position.IsDraw() {
			return draw
		}

		if game.Position.IsCheckMate() {
			switch game.Position.PlayerInTurn {
			case chess.Black:
				return whiteWin
			case chess.White:
				return blackWin
			}
		}
	}
}
