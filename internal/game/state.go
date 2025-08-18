package game

import (
	"slices"

	"github.com/rswilli/chess/internal/chess"
	"github.com/rswilli/chess/internal/pubsub"
	"github.com/rswilli/chess/internal/www"
)

// State tracks the state of the game
//
// it allows others to subscribe to changes on it
type State struct {
	lock pubsub.RWLock

	currentBoard *chess.Position

	// currentSquare is the square that was selected for the next move
	currentSquare chess.Square
}

func (s *State) Render() ([]byte, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var moves []chess.Square
	var promotion bool

	for _, m := range s.currentBoard.GenerateMoves() {
		if m.From != s.currentSquare {
			continue
		}

		promotion = m.Special.Has(chess.PromoteBishop | chess.PromoteKnight | chess.PromoteQueen | chess.PromoteRook)
		moves = append(moves, m.To)
	}

	return www.RenderBoard(www.Data{
		Board:    s.currentBoard,
		Selected: s.currentSquare,

		MoveTargets: moves,
		Promotion:   promotion,
	})
}

func New() *State {
	s := &State{
		currentBoard:  chess.New(),
		currentSquare: chess.InvalidSquare,
	}

	s.currentBoard.GenerateMoves()

	return s
}

// DoSquare either selects the field for making a move or performs a move from the
// last selected square to the current one
func (s *State) DoSquare(square chess.Square, special chess.MoveSpecial) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.currentSquare == chess.InvalidSquare {

		hasMove := slices.ContainsFunc(s.currentBoard.GenerateMoves(), func(m chess.Move) bool {
			return m.From == square
		})

		if !hasMove {
			return
		}

		s.currentSquare = square
		return
	}

	if square == s.currentSquare {
		// selected again -> deselect
		s.currentSquare = chess.InvalidSquare
		return
	}

	var moveI int

	if special == chess.NoSpecial {
		moveI = slices.IndexFunc(s.currentBoard.GenerateMoves(), func(m chess.Move) bool {
			return m.From == s.currentSquare && m.To == square
		})
	} else {
		moveI = slices.IndexFunc(s.currentBoard.GenerateMoves(), func(m chess.Move) bool {
			// search for a move that also contains the special move (aka promotion)
			// take in mind that captures can also promote
			return m.From == s.currentSquare && m.To == square && m.Special.Has(special)
		})
	}

	if moveI == -1 {
		s.currentSquare = chess.InvalidSquare
		return
	}

	s.currentBoard.DoMove(s.currentBoard.GenerateMoves()[moveI])
	s.currentBoard.GenerateMoves()

	s.currentSquare = chess.InvalidSquare
}

func (s *State) Subscribe() chan struct{} {
	return s.lock.Subscribe()
}

func (s *State) Unsubscribe(ch chan struct{}) {
	s.lock.Unsubscribe(ch)
}
