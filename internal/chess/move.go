package chess

import (
	"errors"
	"fmt"
)

type MoveSpecial uint8

const (
	NoSpecial MoveSpecial = iota

	PromoteQueen
	PromoteRook
	PromoteBishop
	PromoteKnight

	CastleLong
	CastleShort
)

type Move struct {
	From    Square
	To      Square
	Special MoveSpecial
}

const castleShort = "O-O"
const castleLong = "O-O-O"

var ErrInvalidMove = errors.New("could not parse move")

var InvalidMove = Move{From: InvalidSquare, To: InvalidSquare}

// ParseMove parses a move given in pure coordinate notation
// see https://www.chessprogramming.org/Algebraic_Chess_Notation#Pure_coordinate_notation
func ParseMove(in string) (Move, error) {
	if in == castleLong {
		return Move{
			From:    InvalidSquare,
			To:      InvalidSquare,
			Special: CastleLong,
		}, nil
	}
	if in == castleShort {
		return Move{
			From:    InvalidSquare,
			To:      InvalidSquare,
			Special: CastleShort,
		}, nil
	}

	if len(in) != 4 && len(in) != 5 {
		return InvalidMove, ErrInvalidMove
	}

	from, err := ParseSquare(in[0:2])

	if err != nil {
		return InvalidMove, fmt.Errorf("%w: %w", ErrInvalidMove, err)
	}

	to, err := ParseSquare(in[2:4])

	if err != nil {
		return InvalidMove, fmt.Errorf("%w: %w", ErrInvalidMove, err)
	}

	m := Move{
		From:    from,
		To:      to,
		Special: NoSpecial,
	}

	if len(in) == 5 {
		switch in[4] {
		case 'q':
			m.Special = PromoteQueen
		case 'r':
			m.Special = PromoteRook
		case 'b':
			m.Special = PromoteBishop
		case 'n':
			m.Special = PromoteKnight
		default:
			return InvalidMove, fmt.Errorf("%w: invalid promotion %c", ErrInvalidMove, in[4])
		}
	}

	return m, nil
}

func (m Move) String() string {
	if m.Special == CastleLong {
		return castleLong
	}
	if m.Special == CastleShort {
		return castleShort
	}

	s := ""

	s += m.From.String()
	s += m.To.String()

	switch m.Special {
	case PromoteBishop:
		s += "b"
	case PromoteKnight:
		s += "k"
	case PromoteQueen:
		s += "q"
	case PromoteRook:
		s += "r"
	}

	return s
}
