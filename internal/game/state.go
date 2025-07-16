package game

import (
	"log/slog"
	"sync"
	"time"

	"github.com/rswilli/chess/internal/chess"
)

// State is a struct that holds a reference to the current board
//
// it allows others to subscribe to changes on it
type State struct {
	lock sync.Mutex
	subs map[chan struct{}]struct{}

	currentBoard *chess.Board
}

func (s *State) Unsubscribe(ch chan struct{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.subs, ch)
}

func New() *State {
	s := &State{
		subs:         make(map[chan struct{}]struct{}),
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
			s.publish()
			s.lock.Unlock()

			i = (i + 1) % len(fens)

			time.Sleep(3 * time.Second)
		}
	}()

	return s
}

// Board returns a copy of the current board
func (s *State) Board() *chess.Board {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.currentBoard.Copy()
}

func (s *State) Subscribe() chan struct{} {
	s.lock.Lock()
	defer s.lock.Unlock()

	ch := make(chan struct{})

	s.subs[ch] = struct{}{}

	return ch
}

// publish must be called with the lock held
func (s *State) publish() {
	slog.Info("publishing update on game")
	for ch := range s.subs {
		ch <- struct{}{}
	}
}
