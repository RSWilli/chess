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

	PossibleMoves []Move
}

// const testFen = "8/8/8/3Nn3/8/8/8/3K1k2 w KQkq - 0 1"

func NewBoard() *Board {
	b, err := NewBoardFromFEN(DefaultFen)
	// b, err := NewBoardFromFEN(testFen)

	if err != nil {
		panic(err)
	}

	return b
}

func (b Board) String() string {
	sb := strings.Builder{}

	for rank := range 8 {
		for file := range 8 {
			piece := b.Square(NewSquare(rank, file))

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
	case b.whitePawns.has(sq):
		return WhitePawn
	case b.whiteKnights.has(sq):
		return WhiteKnight
	case b.whiteBishops.has(sq):
		return WhiteBishop
	case b.whiteRooks.has(sq):
		return WhiteRook
	case b.whiteQueens.has(sq):
		return WhiteQueen
	case b.whiteKing.has(sq):
		return WhiteKing
	case b.blackPawns.has(sq):
		return BlackPawn
	case b.blackKnights.has(sq):
		return BlackKnight
	case b.blackBishops.has(sq):
		return BlackBishop
	case b.blackRooks.has(sq):
		return BlackRook
	case b.blackQueens.has(sq):
		return BlackQueen
	case b.blackKing.has(sq):
		return BlackKing
	default:
		return Empty
	}
}

// set sets the given piece on the given square
func (b *Board) set(sq Square, p Piece) {
	switch p {
	case WhitePawn:
		b.whitePawns = b.whitePawns.set(sq)
	case WhiteKnight:
		b.whiteKnights = b.whiteKnights.set(sq)
	case WhiteBishop:
		b.whiteBishops = b.whiteBishops.set(sq)
	case WhiteRook:
		b.whiteRooks = b.whiteRooks.set(sq)
	case WhiteQueen:
		b.whiteQueens = b.whiteQueens.set(sq)
	case WhiteKing:
		b.whiteKing = b.whiteKing.set(sq)
	case BlackPawn:
		b.blackPawns = b.blackPawns.set(sq)
	case BlackKnight:
		b.blackKnights = b.blackKnights.set(sq)
	case BlackBishop:
		b.blackBishops = b.blackBishops.set(sq)
	case BlackRook:
		b.blackRooks = b.blackRooks.set(sq)
	case BlackQueen:
		b.blackQueens = b.blackQueens.set(sq)
	case BlackKing:
		b.blackKing = b.blackKing.set(sq)
	default:
		panic("unexpected piece received by set")
	}
}

// unset removes the piece from the board at the given square
func (b *Board) unset(sq Square) {
	b.whitePawns = b.whitePawns.unset(sq)
	b.whiteKnights = b.whiteKnights.unset(sq)
	b.whiteBishops = b.whiteBishops.unset(sq)
	b.whiteRooks = b.whiteRooks.unset(sq)
	b.whiteQueens = b.whiteQueens.unset(sq)
	b.whiteKing = b.whiteKing.unset(sq)

	b.blackPawns = b.blackPawns.unset(sq)
	b.blackKnights = b.blackKnights.unset(sq)
	b.blackBishops = b.blackBishops.unset(sq)
	b.blackRooks = b.blackRooks.unset(sq)
	b.blackQueens = b.blackQueens.unset(sq)
	b.blackKing = b.blackKing.unset(sq)
}

func (b *Board) allPieces() bitBoard {
	return b.whitePieces() | b.blackPieces()
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

func (b *Board) DoMove(m Move) {
	// TODO: castling
	p := b.Square(m.From)

	// clear the old square
	b.unset(m.From)
	// must clear all other boards before setting the new one
	b.unset(m.To)

	switch {
	case m.Special.Has(PromoteQueen):
		p = Queen | b.PlayerInTurn
	case m.Special.Has(PromoteRook):
		p = Rook | b.PlayerInTurn
	case m.Special.Has(PromoteKnight):
		p = Knight | b.PlayerInTurn
	case m.Special.Has(PromoteBishop):
		p = Bishop | b.PlayerInTurn
	}

	// remove the en-passant captured pawn, no need to check for the piece type since the en passant
	// square is always empty, so no other move can capture on it
	if m.Special.Has(Captures) && m.To == b.EnPassantTarget && b.PlayerInTurn == White {
		b.unset(Square(bitBoard(m.To).down()))
	} else if m.Special.Has(Captures) && m.To == b.EnPassantTarget && b.PlayerInTurn == Black {
		b.unset(Square(bitBoard(m.To).up()))
	}

	// save the en passant square for the move generation of the en passant moves
	if m.Special.Has(DoublePawnPush) && b.PlayerInTurn == White {
		b.EnPassantTarget = Square(bitBoard(m.From).up())
	} else if m.Special.Has(DoublePawnPush) && b.PlayerInTurn == Black {
		b.EnPassantTarget = Square(bitBoard(m.From).down())
	} else {
		b.EnPassantTarget = InvalidSquare
	}

	b.set(m.To, p)

	if b.PlayerInTurn == White {
		b.PlayerInTurn = Black
	} else {
		b.PlayerInTurn = White
	}
}
