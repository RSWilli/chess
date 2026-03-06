package engine

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/rswilli/chess/internal/chess"
	"github.com/rswilli/chess/internal/uci"
)

type Engine struct {
	pos *chess.Position

	lock sync.Mutex
	// searching is used to check if the engine is searching but also to know when to abort the search in
	// the search loop
	searching         bool
	currentSearchOpts *uci.GoOptions
}

func NewEngine() *Engine {
	return &Engine{
		pos: chess.NewPosition(),
	}
}

// NewGame implements uci.Engine.
func (e *Engine) NewGame() error {
	e.pos = chess.NewPosition()

	// clear all caches here

	return nil
}

// Perft implements uci.Engine.
func (e *Engine) Perft(depth int) (uci.PerftResult, error) {
	e.lock.Lock()
	defer e.lock.Unlock()

	if e.searching {
		return uci.PerftResult{}, fmt.Errorf("search already running, wont do perft")
	}

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

// internalPerft is called from [Engine.Perft] and instead of listing the moves it only counts the positions until depth
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
	e.lock.Lock()
	defer e.lock.Unlock()

	if e.searching {
		return fmt.Errorf("search running cannot change position")
	}

	// TODO: the GUI passes an increasing list of moves on the starting position,
	// optimize this so full parsing is only done when needed
	pos, err := chess.NewPositionFromFEN(fen, moves)

	if err != nil {
		return fmt.Errorf("could not parse FEN: %w", err)
	}

	e.pos = pos

	return nil
}

// Ready implements uci.Engine.
func (e *Engine) Ready() error {
	// block here if not ready
	return nil
}

// Go implements uci.Engine.
func (e *Engine) Go(options uci.GoOptions) uci.Bestmove {
	e.lock.Lock()

	if e.searching {
		panic("double search")
	}

	e.searching = true
	e.currentSearchOpts = &options
	e.lock.Unlock()

	bm := e.search()

	e.lock.Lock()
	e.searching = true
	e.lock.Unlock()

	return bm
}

func (e *Engine) search() uci.Bestmove {
	moves := e.pos.GenerateMoves()

	if len(moves) == 0 {
		return uci.Bestmove{}
	}

	i := rand.IntN(len(moves))

	time.Sleep(250 * time.Millisecond)

	e.searching = false

	return uci.Bestmove{
		BestMove: moves[i].String(),
		Ponder:   "a1a2", // tmp
	}
}

// Stop implements uci.Engine.
func (e *Engine) Stop() {
	e.lock.Lock()
	defer e.lock.Unlock()

	e.searching = false
	e.currentSearchOpts = nil
}

var _ uci.Engine = &Engine{}
