package game

import (
	"github.com/rswilli/chess/internal/chess"
	"github.com/rswilli/chess/internal/pubsub"
	"github.com/rswilli/chess/internal/www"
)

// State tracks the state of the game
//
// it allows others to subscribe to changes on it
type State struct {
	lock pubsub.RWLock

	currentBoard *chess.Board

	// currentSquare is the square that was selected for the next move
	currentSquare chess.Square
}

func (s *State) Render() ([]byte, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return www.RenderBoard(www.Data{
		Board:    s.currentBoard,
		Selected: s.currentSquare,
	})
}

func New() *State {
	s := &State{
		currentBoard:  chess.NewBoard(),
		currentSquare: chess.InvalidSquare,
	}

	return s
}

// DoSquare either selects the field for making a move or performs a move from the
// last selected square to the current one
func (s *State) DoSquare(square chess.Square) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.currentSquare == chess.InvalidSquare && s.currentBoard.Square(square) != chess.Empty {
		s.currentSquare = square
		return
	}

	if s.currentSquare == chess.InvalidSquare {
		return
	}

	s.currentBoard.DoMove(chess.Move{
		From: s.currentSquare,
		To:   square,
	})

	s.currentSquare = chess.InvalidSquare
}

func (s *State) Subscribe() chan struct{} {
	return s.lock.Subscribe()
}

func (s *State) Unsubscribe(ch chan struct{}) {
	s.lock.Unsubscribe(ch)
}
