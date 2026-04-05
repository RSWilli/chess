// package mm contains all logic for matchmaking
package mm

import (
	"fmt"
	"time"

	"github.com/rswilli/chess/internal/chess"
	"github.com/rswilli/chess/internal/uci"
	"github.com/rswilli/chess/internal/uci/search"
)

type player string

const (
	playerWhite player = "white"
	playerBlack player = "black"
)

type GameState int

const (
	Running GameState = iota
	CheckMate
	StaleMate
	Error
)

type Game struct {
	Position *chess.Position

	state GameState

	startFEN string
	history  []string

	current player

	white uci.Engine
	black uci.Engine
}

func (g *Game) Stop() {
	g.white.Stop()
	g.black.Stop()
}

func (g *Game) Current() uci.Engine {
	switch g.current {
	case playerBlack:
		return g.black
	case playerWhite:
		return g.white
	default:
		panic(fmt.Sprintf("unexpected mm.player: %#v", g.current))
	}
}

// Move lets the current side think about the move and performs it.
func (g *Game) Move() error {
	if g.state != Running {
		return fmt.Errorf("game over")
	}

	var err error
	current := g.Current()

	err = current.Position(g.startFEN, g.history)

	if err != nil {
		g.state = Error

		return fmt.Errorf("could not set position for current player %s: %w", g.current, err)
	}

	bestmove, _ := current.Go(search.Options{
		MoveTime: 1 * time.Second,
	})

	move, err := g.Position.ParseMove(bestmove)

	if err != nil {
		g.state = Error
		return fmt.Errorf("received invalid move %s in current position from player %s, could not parse: %w", bestmove, g.current, err)
	}

	g.Position.DoMove(move)

	if g.current == playerWhite {
		g.current = playerBlack
	} else {
		g.current = playerWhite
	}

	g.history = append(g.history, bestmove)

	if g.Position.IsCheckMate() {
		g.state = CheckMate
	}

	if g.Position.IsDraw() {
		g.state = StaleMate
	}

	return nil
}

func (g *Game) State() GameState {
	return g.state
}

func NewGame(white, black uci.Engine) *Game {
	g, err := NewGameWithFEN(chess.DefaultFen, white, black)

	if err != nil {
		panic("could not parse default FEN for new game")
	}

	return g
}

func NewGameWithFEN(fen string, white, black uci.Engine) (*Game, error) {
	black.NewGame()
	white.NewGame()

	pos, err := chess.NewPositionFromFEN(fen, nil)

	if err != nil {
		return nil, err
	}

	err = white.Position(fen, nil)

	if err != nil {
		return nil, fmt.Errorf("could not set white's position: %v", err)
	}

	err = black.Position(fen, nil)

	if err != nil {
		return nil, fmt.Errorf("could not set black's position: %v", err)
	}

	return &Game{
		Position: pos,
		startFEN: fen,
		current:  playerWhite,
		white:    white,
		black:    black,
	}, nil
}
