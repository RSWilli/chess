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
func (e *Engine) Position(pos *chess.Position) error {
	e.lock.Lock()
	defer e.lock.Unlock()

	if e.searching {
		return fmt.Errorf("search running cannot change position")
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
func (e *Engine) Go(options uci.GoOptions) (<-chan uci.Bestmove, <-chan uci.Info) {
	e.lock.Lock()
	defer e.lock.Unlock()

	e.searching = true
	bm := make(chan uci.Bestmove)
	info := make(chan uci.Info)

	e.currentSearchOpts = &options

	go e.search(bm, info)

	return bm, info
}

func (e *Engine) search(ret chan uci.Bestmove, info chan uci.Info) {
	moves := e.pos.GenerateMoves()
	close(info) // Do we actually need to log info?

	if len(moves) == 0 {
		e.searching = false

		close(ret)
		return
	}

	i := rand.IntN(len(moves))

	time.Sleep(250 * time.Millisecond)

	ret <- uci.Bestmove{
		BestMove: moves[i].String(),
		Ponder:   "a1a2", // tmp
	}

	e.searching = false
}

// Stop implements uci.Engine.
func (e *Engine) Stop() {
	e.lock.Lock()
	defer e.lock.Unlock()

	e.searching = false
	e.currentSearchOpts = nil
}

var _ uci.Engine = &Engine{}
