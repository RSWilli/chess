package chess

import (
	"fmt"

	"github.com/rswilli/chess/internal/zobrist"
)

type CastlingAbility uint8

const (
	NoCastling CastlingAbility = 0

	CastleWhiteQueen CastlingAbility = 1 << (iota - 1)
	CastleWhiteKing

	CastleBlackQueen
	CastleBlackKing
)

func (c CastlingAbility) Has(other CastlingAbility) bool {
	return c&other != 0
}
func (c CastlingAbility) String() string {
	if c == NoCastling {
		return "-"
	}

	s := ""

	if c.Has(CastleWhiteKing) {
		s += "K"
	}
	if c.Has(CastleWhiteQueen) {
		s += "Q"
	}
	if c.Has(CastleBlackKing) {
		s += "k"
	}
	if c.Has(CastleBlackQueen) {
		s += "q"
	}

	return s
}

func (c CastlingAbility) zobrist() int {
	switch c {
	case CastleBlackKing:
		return zobrist.BlackCastleKing
	case CastleBlackQueen:
		return zobrist.BlackCastleQueen
	case CastleWhiteKing:
		return zobrist.WhiteCastleKing
	case CastleWhiteQueen:
		return zobrist.WhiteCastleQueen
	default:
		panic(fmt.Sprintf("unexpected chess.CastlingAbility: %#v", c))
	}
}
