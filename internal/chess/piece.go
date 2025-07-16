package chess

import "fmt"

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
