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

	whitePawns   bitBoard
	whiteKnights bitBoard
	whiteBishops bitBoard
	whiteRooks   bitBoard
	whiteQueens  bitBoard
	whiteKing    bitBoard

	blackPawns   bitBoard
	blackKnights bitBoard
	blackBishops bitBoard
	blackRooks   bitBoard
	blackQueens  bitBoard
	blackKing    bitBoard
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
			piece := b.Square(mkSquare(rank, file))

			sb.WriteRune(piece.Rune())
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
	switch {
	case b.whitePawns.hasSquare(sq):
		return WhitePawn
	case b.whiteKnights.hasSquare(sq):
		return WhiteKnight
	case b.whiteBishops.hasSquare(sq):
		return WhiteBishop
	case b.whiteRooks.hasSquare(sq):
		return WhiteRook
	case b.whiteQueens.hasSquare(sq):
		return WhiteQueen
	case b.whiteKing.hasSquare(sq):
		return WhiteKing
	case b.blackPawns.hasSquare(sq):
		return BlackPawn
	case b.blackKnights.hasSquare(sq):
		return BlackKnight
	case b.blackBishops.hasSquare(sq):
		return BlackBishop
	case b.blackRooks.hasSquare(sq):
		return BlackRook
	case b.blackQueens.hasSquare(sq):
		return BlackQueen
	case b.blackKing.hasSquare(sq):
		return BlackKing
	default:
		return Empty
	}
}

// set sets the given piece on the given square
func (b *Board) set(sq Square, p Piece) {
	switch p {
	case WhitePawn:
		b.whitePawns = b.whitePawns.setSquare(sq)
	case WhiteKnight:
		b.whiteKnights = b.whiteKnights.setSquare(sq)
	case WhiteBishop:
		b.whiteBishops = b.whiteBishops.setSquare(sq)
	case WhiteRook:
		b.whiteRooks = b.whiteRooks.setSquare(sq)
	case WhiteQueen:
		b.whiteQueens = b.whiteQueens.setSquare(sq)
	case WhiteKing:
		b.whiteKing = b.whiteKing.setSquare(sq)
	case BlackPawn:
		b.blackPawns = b.blackPawns.setSquare(sq)
	case BlackKnight:
		b.blackKnights = b.blackKnights.setSquare(sq)
	case BlackBishop:
		b.blackBishops = b.blackBishops.setSquare(sq)
	case BlackRook:
		b.blackRooks = b.blackRooks.setSquare(sq)
	case BlackQueen:
		b.blackQueens = b.blackQueens.setSquare(sq)
	case BlackKing:
		b.blackKing = b.blackKing.setSquare(sq)
	default:
		panic("unexpected piece received by set")
	}
}

// unset removes the piece from the board at the given square
func (b *Board) unset(sq Square) {
	b.whitePawns = b.whitePawns.unsetSquare(sq)
	b.whiteKnights = b.whiteKnights.unsetSquare(sq)
	b.whiteBishops = b.whiteBishops.unsetSquare(sq)
	b.whiteRooks = b.whiteRooks.unsetSquare(sq)
	b.whiteQueens = b.whiteQueens.unsetSquare(sq)
	b.whiteKing = b.whiteKing.unsetSquare(sq)

	b.blackPawns = b.blackPawns.unsetSquare(sq)
	b.blackKnights = b.blackKnights.unsetSquare(sq)
	b.blackBishops = b.blackBishops.unsetSquare(sq)
	b.blackRooks = b.blackRooks.unsetSquare(sq)
	b.blackQueens = b.blackQueens.unsetSquare(sq)
	b.blackKing = b.blackKing.unsetSquare(sq)
}

func (b *Board) DoMove(m Move) {
	// TODO: castling, en passant, promotion
	p := b.Square(m.From)

	// clear the old square
	b.unset(m.From)
	// must clear all other boards before setting the new one
	b.unset(m.To)
	b.set(m.To, p)
}

func (b *Board) whitePieces() bitBoard {
	return b.whitePawns |
		b.whiteKnights |
		b.whiteBishops |
		b.whiteRooks |
		b.whiteQueens |
		b.whiteKing
}

func (b *Board) blackPieces() bitBoard {
	return b.blackPawns |
		b.blackKnights |
		b.blackBishops |
		b.blackRooks |
		b.blackQueens |
		b.blackKing
}
