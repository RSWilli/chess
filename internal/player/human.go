package player

import (
	"fmt"
	"log/slog"
	"slices"
	"sync"

	"github.com/rswilli/chess/internal/chess"
	"github.com/rswilli/chess/internal/uci"
)

// Human is an uci.Engine implementation that exposes additional methods that wait for user input
type Human struct {
	lock sync.Mutex

	pos *chess.Position

	// searching is true when the channels below are not closed, thus a search is running
	searching  bool
	stopSearch chan struct{}
	userMove   chan chess.Move

	// currentSquare is the square that was selected for the next move
	currentSquare chess.Square
}

func NewHuman() *Human {
	return &Human{
		pos: chess.NewPosition(),
	}
}

func (h *Human) CurrentSquare() chess.Square {
	h.lock.Lock()
	defer h.lock.Unlock()

	return h.currentSquare
}

func (h *Human) DoMove(movestr string) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	// in any case, deselect the current square
	h.currentSquare = chess.InvalidSquare

	move, err := h.pos.ParseMove(movestr)

	if err != nil {
		return fmt.Errorf("invalid square given: %w", err)
	}

	h.userMove <- move

	return nil
}

// DoSquare selects the field for making a move
func (h *Human) DoSquare(square chess.Square) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if !h.searching {
		return
	}

	moves := h.pos.GenerateMoves()

	hasMove := slices.ContainsFunc(moves, func(m chess.Move) bool {
		return m.From == square
	})

	if !hasMove {
		return
	}

	slog.Info("human set selected square", "square", square.String())

	h.currentSquare = square
}

// Go implements uci.Engine.
func (h *Human) Go(uci.GoOptions) uci.Bestmove {
	h.lock.Lock()

	if h.searching {
		panic("double search")
	}

	stopSearch := make(chan struct{})
	userMove := make(chan chess.Move)

	h.searching = true
	h.stopSearch = stopSearch
	h.userMove = userMove

	h.lock.Unlock()

	defer h.Stop()

	select {
	case <-stopSearch:
		return uci.Bestmove{} // null move
	case m := <-userMove: // closed will return null move
		// no ponder move
		return uci.Bestmove{
			BestMove: m.String(),
		}
	}
}

// NewGame implements uci.Engine.
func (h *Human) NewGame() error {
	h.pos = chess.NewPosition()

	return nil
}

// Perft implements uci.Engine.
func (h *Human) Perft(depth int) (uci.PerftResult, error) {
	return uci.PerftResult{}, fmt.Errorf("not implemented")
}

// Position implements uci.Engine.
func (h *Human) Position(fen string, moves []string) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.stopUnlocked()

	pos, err := chess.NewPositionFromFEN(fen, moves)

	if err != nil {
		return fmt.Errorf("could not parse FEN: %w", err)
	}

	h.pos = pos

	return nil
}

// Ready implements uci.Engine.
func (h *Human) Ready() error {
	// take the lock to wait for other actions to complete, only after that we are ready
	h.lock.Lock()
	defer h.lock.Unlock()
	return nil
}

// Stop implements uci.Engine.
func (h *Human) Stop() {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.stopUnlocked()
}

// stopUnlocked is the unlocked implementation of Stop that needs to be called with the lock held
func (h *Human) stopUnlocked() {
	if !h.searching {
		return
	}

	h.searching = false

	close(h.stopSearch)
	close(h.userMove)
}

var _ uci.Engine = (*Human)(nil)
