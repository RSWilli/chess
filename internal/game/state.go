package game

import (
	"time"

	"github.com/rswilli/chess/internal/chess"
	"github.com/rswilli/chess/internal/pubsub"
)

// State is a struct that holds a reference to the current board
//
// it allows others to subscribe to changes on it
type State struct {
	lock pubsub.RWLock

	currentBoard *chess.Board
}

func New() *State {
	s := &State{
		currentBoard: chess.NewBoard(),
	}

	go func() {
		// test go func to publish
		fens := []string{
			"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
			"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2",
			"rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2",
		}

		i := 0
		for {
			s.lock.Lock()
			s.currentBoard, _ = chess.NewBoardFromFEN(fens[i])
			s.lock.Unlock()

			i = (i + 1) % len(fens)

			time.Sleep(3 * time.Second)
		}
	}()

	return s
}

// Board returns a copy of the current board
func (s *State) Board() *chess.Board {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.currentBoard.Copy()
}

func (s *State) Subscribe() chan struct{} {
	return s.lock.Subscribe()
}

func (s *State) Unsubscribe(ch chan struct{}) {
	s.lock.Unsubscribe(ch)
}
