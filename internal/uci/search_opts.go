package uci

import (
	"time"
)

// GoOptions is a struct containing the go parameters, see https://official-stockfish.github.io/docs/stockfish-wiki/UCI-&-Commands.html#go
type GoOptions struct {
	// SearchMoves restricts search to these moves only.
	SearchMoves []string
	// Ponder tells the engine to start pondering mode
	Ponder bool
	// WhiteTime is the total time white has left
	WhiteTime time.Duration
	// BlackTime is the total time black has left
	BlackTime time.Duration
	// WhiteIncrement is the increment white gets after a move
	WhiteIncrement time.Duration
	// BlackIncrement is the increment black gets after a move
	BlackIncrement time.Duration
	// MovesToGo is the amount of moves until the next time control
	MovesToGo int
	// Depth tells the engine to only search until this depth
	Depth int
	// Nodes tells the engine to stop the search after a certain amount of nodes is approx. reached
	Nodes int
	// Mate tells the engine to search for mate in X
	Mate int
	// MoveTime tells the engine to stop the search after this time
	MoveTime time.Duration
	// Infinite stop the search only after stop is received
	Infinite bool
}
