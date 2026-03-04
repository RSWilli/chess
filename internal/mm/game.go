// package mm contains all logic for matchmaking
package mm

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/rswilli/chess/internal/chess"
	"github.com/rswilli/chess/internal/uci"
)

type player string

const (
	playerWhite player = "white"
	playerBlack player = "black"
)

type Game struct {
	positionLock sync.Mutex
	position     *chess.Position

	current player

	white *uci.Client
	black *uci.Client
}

func (g *Game) currentClient() *uci.Client {
	switch g.current {
	case playerBlack:
		return g.black
	case playerWhite:
		return g.white
	default:
		panic(fmt.Sprintf("unexpected mm.player: %#v", g.current))
	}
}

var ErrCheckMate = errors.New("checkmate")
var ErrStaleMate = errors.New("stalemate")

// Move lets the current side think about the move and performs it.
func (g *Game) Move() error {
	if g.position.IsCheckMate() {
		return ErrCheckMate
	}

	if g.position.IsDraw() {
		return ErrStaleMate
	}

	var err error
	current := g.currentClient()

	err = current.Position(g.position.FEN(), nil)

	if err != nil {
		return fmt.Errorf("could not set position for current player %s: %w", g.current, err)
	}

	bestmove, err := current.Go(uci.GoOptions{
		MoveTime: 1 * time.Second,
	})

	if err != nil {
		return fmt.Errorf("error while engine go command from player %s: %w", g.current, err)
	}

	move, err := g.position.ParseMove(bestmove.BestMove)

	if err != nil {
		return fmt.Errorf("received invalid move in current position from player %s, could not parse: %w", g.current, err)
	}

	g.positionLock.Lock()
	g.position.DoMove(move)
	g.positionLock.Unlock()

	if g.current == playerWhite {
		g.current = playerBlack
	} else {
		g.current = playerWhite
	}

	return nil
}

func (g *Game) Position() *chess.Position {
	g.positionLock.Lock()
	defer g.positionLock.Unlock()

	return g.position.Copy()
}

func NewGame(white, black *uci.Client) *Game {
	black.NewGame()
	white.NewGame()

	return &Game{
		position: chess.NewPosition(),
		current:  playerWhite,
		white:    white,
		black:    black,
	}
}
