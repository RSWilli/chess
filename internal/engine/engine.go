package engine

import (
	"fmt"

	"github.com/rswilli/chess/internal/chess"
	"github.com/rswilli/chess/internal/uci"
)

type Engine struct {
	pos *chess.Game
}

func NewEngine() *Engine {
	return &Engine{
		pos: chess.NewGame(),
	}
}

// Go implements uci.Engine.
func (e *Engine) Go(options uci.GoOptions) (uci.GoResponse, error) {
	panic("unimplemented")
}

// NewGame implements uci.Engine.
func (e *Engine) NewGame() error {
	e.pos = chess.NewGame()

	// clear all caches here

	return nil
}

// Perft implements uci.Engine.
func (e *Engine) Perft(depth int) (uci.PerftResult, error) {
	res := uci.PerftResult{
		Moves: make(map[string]int),
	}

	for _, m := range e.pos.GenerateMoves() {
		e.pos.DoMove(m)
		visited := 1
		if depth > 1 {
			visited = e.internalPerft(depth - 1)
		}
		e.pos.UndoMove()

		res.Moves[m.String()] = visited
		res.Total += visited
	}

	return res, nil
}

func (e *Engine) internalPerft(depth int) int {
	if depth < 1 {
		panic("depth wrong")
	}

	moves := e.pos.GenerateMoves()

	if depth == 1 {
		return len(moves)
	}

	total := 0

	for _, m := range moves {
		e.pos.DoMove(m)

		total += e.internalPerft(depth - 1)

		e.pos.UndoMove()
	}

	return total
}

// Position implements uci.Engine.
func (e *Engine) Position(fen string, moves []string) error {
	if fen == uci.StartPosition {
		e.pos = chess.NewGame()
		return nil
	} else {
		pos, err := chess.NewGameFromFEN(fen)

		if err != nil {
			return fmt.Errorf("could not parse FEN: %v", err)
		}

		e.pos = pos
	}

	for _, m := range moves {
		move, err := e.pos.ParseMove(m)

		if err != nil {
			return fmt.Errorf("could not parse move: %v", err)
		}

		e.pos.DoMove(move)
	}

	return nil
}

// Ready implements uci.Engine.
func (e *Engine) Ready() error {
	// block here if not ready
	return nil
}

var _ uci.Engine = &Engine{}
