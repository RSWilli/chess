package chess

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/rswilli/chess/internal/search"
)

type Engine struct {
	pos *Position

	lock sync.Mutex
	// searching is used to check if the engine is searching but also to know when to abort the search in
	// the search loop
	searching         bool
	currentSearchOpts *search.Options
}

func NewEngine() *Engine {
	return &Engine{
		pos: NewPosition(),
	}
}

// NewGame implements uci.Engine.
func (e *Engine) NewGame() error {
	e.pos = NewPosition()

	// clear all caches here

	return nil
}

// Perft implements uci.Engine.
func (e *Engine) Perft(depth int) (total int, moves map[string]int, err error) {
	e.lock.Lock()
	defer e.lock.Unlock()

	if e.searching {
		return 0, nil, fmt.Errorf("search already running, wont do perft")
	}

	moves = make(map[string]int)
	total = 0

	for _, m := range e.pos.GenerateMoves() {
		e.pos.DoMove(m)
		visited := 1
		if depth > 1 {
			visited = e.internalPerft(depth - 1)
		}
		e.pos.UndoMove()

		moves[m.String()] = visited
		total += visited
	}

	return total, moves, nil
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
	pos, err := NewPositionFromFEN(fen, moves)

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
func (e *Engine) Go(options search.Options) (bestmove string, ponder string) {
	e.lock.Lock()

	if e.searching {
		panic("double search")
	}

	e.searching = true
	e.currentSearchOpts = &options
	e.lock.Unlock()

	bm := e.search()

	e.lock.Lock()
	e.searching = false
	e.lock.Unlock()

	return bm, ""
}

func (e *Engine) search() string {
	moves := e.pos.GenerateMoves()

	if len(moves) == 0 {
		return ""
	}

	i := rand.IntN(len(moves))

	time.Sleep(250 * time.Millisecond)

	e.searching = false

	return moves[i].String()
}

// Stop implements uci.Engine.
func (e *Engine) Stop() {
	e.lock.Lock()
	defer e.lock.Unlock()

	e.searching = false
	e.currentSearchOpts = nil
}
