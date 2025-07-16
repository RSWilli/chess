package chess

import (
	"strings"
)

type Board struct {
	// max 50, see https://www.chessprogramming.org/Fifty-move_Rule
	HalfmoveClock uint8
	Moves         int

	PlayerInTurn    Piece
	Castling        CastlingAbility
	EnPassantTarget Square

	// pos is an array of Pieces from rank 8 through rank 1, file a to h
	pos [64]Piece
}

func NewBoard() *Board {
	b, err := NewBoardFromFEN(DefaultFen)

	if err != nil {
		panic(err)
	}

	return b
}

func (b Board) String() string {
	sb := strings.Builder{}

	for rank := range 8 {
		for file := range 8 {
			sb.WriteRune(b.pos[8*rank+file].Rune())
		}
		sb.WriteByte('\n')
	}
	switch b.PlayerInTurn {
	case White:
		sb.WriteString("White\n")
	case Black:
		sb.WriteString("Black\n")
	default:
		panic("unexpected player turn")
	}

	sb.WriteString(b.Castling.String())
	if b.EnPassantTarget != InvalidSquare {
		sb.WriteString("\nEn passant to: ")
		sb.WriteString(b.EnPassantTarget.String())
	}

	return sb.String()
}

func (b *Board) Square(sq Square) Piece {
	return b.pos[sq.Index()]
}

func (b *Board) Copy() *Board {
	c := *b

	return &c
}
