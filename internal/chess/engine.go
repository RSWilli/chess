package chess

import (
	"fmt"
	"io"
	"slices"
	"sync"

	"github.com/rswilli/chess/internal/uci/search"
)

type Engine struct {
	logSink io.Writer

	pos *Position

	lock sync.Mutex
	// searching is used to check if the engine is searching but also to know when to abort the search in
	// the search loop
	searching         bool
	currentSearchOpts *search.Options

	// currentLine contains the current series of moves that are being evaluated
	currentLine []string
}

func NewEngine(logSink io.Writer) *Engine {
	if logSink == nil {
		logSink = io.Discard
	}

	return &Engine{
		pos:         NewPosition(),
		logSink:     logSink,
		currentLine: make([]string, 0, 20),
	}
}

// NewGame implements uci.Engine.
func (e *Engine) NewGame() error {
	e.pos = NewPosition()

	// clear all caches here

	e.currentLine = e.currentLine[0:0]

	e.logf("new game started")

	return nil
}

func (e *Engine) logf(f string, args ...any) {
	if e.logSink != nil {
		fmt.Fprintf(e.logSink, f, args...)
	}
}

// Perft implements uci.Engine.
func (e *Engine) Perft(depth int) (total int, moves map[string]int, err error) {
	e.lock.Lock()
	defer e.lock.Unlock()

	defer func() {
		if r := recover(); r != nil {
			e.logf("caught panic in line: %v\n", e.currentLine)
			err = fmt.Errorf("panicked in perft in engine: %v", r)
		}
	}()

	if e.searching {
		return 0, nil, fmt.Errorf("search already running, wont do perft")
	}

	moves = make(map[string]int)
	total = 0

	for _, m := range e.pos.GenerateMoves() {
		e.doMove(m)
		visited := 1
		if depth > 1 {
			visited = e.internalPerft(depth - 1)
		}
		e.undoMove()

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
		e.doMove(m)

		total += e.internalPerft(depth - 1)

		e.undoMove()
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

	e.currentLine = moves

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

	slices.SortFunc(moves, e.compareMoves)

	return moves[0].String()
}

// Stop implements uci.Engine.
func (e *Engine) Stop() {
	e.lock.Lock()
	defer e.lock.Unlock()

	e.searching = false
	e.currentSearchOpts = nil
}

// compareMoves is used as a sort function on the list of generated moves
//
// better moves are sorted before worse ones so the search will look at them first
func (e *Engine) compareMoves(a, b Move) int {
	const (
		better = -1
		equal  = 0
		worse  = 1
	)

	switch {
	case a.Special.Has(Check) && !b.Special.Has(Check):
		return better
	case !a.Special.Has(Check) && b.Special.Has(Check):
		return worse
	case a.Special.Has(Captures) && !b.Special.Has(Captures):
		return better
	case !a.Special.Has(Captures) && b.Special.Has(Captures):
		return better
	case a.Special.Has(Captures) && b.Special.Has(Captures):
		if a.Takes > b.Takes {
			return better
		} else if a.Takes < b.Takes {
			return worse
		}
	}

	return equal
}

func (e *Engine) doMove(m Move) {
	e.currentLine = append(e.currentLine, m.String())
	e.pos.DoMove(m)
}

func (e *Engine) undoMove() {
	e.pos.UndoMove()

	e.currentLine = e.currentLine[0 : len(e.currentLine)-1]
}
