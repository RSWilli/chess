package chess

import (
	"errors"
	"fmt"
)

type Promotion uint8

const (
	NoPromote Promotion = iota
	PromoteQueen
	PromoteRook
	PromoteBishop
	PromoteKnight
)

type Castling uint8

const (
	NoCastle Castling = iota
	CastleLong
	CastleShort
)

type Move struct {
	From      Square
	To        Square
	Promotion Promotion
	Castling  Castling
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
			From:     InvalidSquare,
			To:       InvalidSquare,
			Castling: CastleLong,
		}, nil
	}
	if in == castleShort {
		return Move{
			From:     InvalidSquare,
			To:       InvalidSquare,
			Castling: CastleShort,
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
		From:      from,
		To:        to,
		Castling:  NoCastle,
		Promotion: NoPromote,
	}

	if len(in) == 5 {
		switch in[4] {
		case 'q':
			m.Promotion = PromoteQueen
		case 'r':
			m.Promotion = PromoteRook
		case 'b':
			m.Promotion = PromoteBishop
		case 'n':
			m.Promotion = PromoteKnight
		default:
			return InvalidMove, fmt.Errorf("%w: invalid promotion %c", ErrInvalidMove, in[4])
		}
	}

	return m, nil
}

func (m Move) String() string {
	if m.Castling == CastleLong {
		return castleLong
	}
	if m.Castling == CastleShort {
		return castleShort
	}

	s := ""

	s += m.From.String()
	s += m.To.String()

	switch m.Promotion {
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
