package chess

import (
	"fmt"

	"github.com/rswilli/chess/internal/zobrist"
)

type Piece uint8

func (p Piece) Is(other Piece) bool {
	return p&other != 0
}

func (p Piece) Rune() rune {
	switch p {
	case Empty:
		return ' '

	case WhiteKing:
		return '♔'
	case WhiteQueen:
		return '♕'
	case WhiteRook:
		return '♖'
	case WhiteBishop:
		return '♗'
	case WhiteKnight:
		return '♘'
	case WhitePawn:
		return '♙'

	case BlackKing:
		return '♚'
	case BlackQueen:
		return '♛'
	case BlackRook:
		return '♜'
	case BlackBishop:
		return '♝'
	case BlackKnight:
		return '♞'
	case BlackPawn:
		return '♟'
	default:
		panic(fmt.Sprintf("unexpected chess.Piece: %#v", p))
	}
}

// zobrist returns the zobrist offset, see [zobrist.NewBoard]
func (p Piece) zobrist() int {
	switch p {
	case Empty:
		return -1

	case WhiteKing:
		return zobrist.WhiteKing
	case WhiteQueen:
		return zobrist.WhiteQueen
	case WhiteRook:
		return zobrist.WhiteRook
	case WhiteBishop:
		return zobrist.WhiteBishop
	case WhiteKnight:
		return zobrist.WhiteKnight
	case WhitePawn:
		return zobrist.WhitePawn

	case BlackKing:
		return zobrist.BlackKing
	case BlackQueen:
		return zobrist.BlackQueen
	case BlackRook:
		return zobrist.BlackRook
	case BlackBishop:
		return zobrist.BlackBishop
	case BlackKnight:
		return zobrist.BlackKnight
	case BlackPawn:
		return zobrist.BlackPawn
	default:
		panic(fmt.Sprintf("unexpected chess.Piece: %#v", p))
	}
}

const (
	Empty Piece = iota

	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
)

const (
	White Piece = 1 << (iota + 4)
	Black
)

const (
	WhitePawn   = White | Pawn
	WhiteKnight = White | Knight
	WhiteBishop = White | Bishop
	WhiteRook   = White | Rook
	WhiteQueen  = White | Queen
	WhiteKing   = White | King

	BlackPawn   = Black | Pawn
	BlackKnight = Black | Knight
	BlackBishop = Black | Bishop
	BlackRook   = Black | Rook
	BlackQueen  = Black | Queen
	BlackKing   = Black | King
)
