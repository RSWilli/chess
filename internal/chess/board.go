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

	whitePawns   BitBoard
	whiteKnights BitBoard
	whiteBishops BitBoard
	whiteRooks   BitBoard
	whiteQueens  BitBoard
	whiteKing    BitBoard

	blackPawns   BitBoard
	blackKnights BitBoard
	blackBishops BitBoard
	blackRooks   BitBoard
	blackQueens  BitBoard
	blackKing    BitBoard

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
	case b.whitePawns.Has(sq):
		return WhitePawn
	case b.whiteKnights.Has(sq):
		return WhiteKnight
	case b.whiteBishops.Has(sq):
		return WhiteBishop
	case b.whiteRooks.Has(sq):
		return WhiteRook
	case b.whiteQueens.Has(sq):
		return WhiteQueen
	case b.whiteKing.Has(sq):
		return WhiteKing
	case b.blackPawns.Has(sq):
		return BlackPawn
	case b.blackKnights.Has(sq):
		return BlackKnight
	case b.blackBishops.Has(sq):
		return BlackBishop
	case b.blackRooks.Has(sq):
		return BlackRook
	case b.blackQueens.Has(sq):
		return BlackQueen
	case b.blackKing.Has(sq):
		return BlackKing
	default:
		return Empty
	}
}

// set sets the given piece on the given square
func (b *Board) set(sq Square, p Piece) {
	switch p {
	case WhitePawn:
		b.whitePawns = b.whitePawns.Set(sq)
	case WhiteKnight:
		b.whiteKnights = b.whiteKnights.Set(sq)
	case WhiteBishop:
		b.whiteBishops = b.whiteBishops.Set(sq)
	case WhiteRook:
		b.whiteRooks = b.whiteRooks.Set(sq)
	case WhiteQueen:
		b.whiteQueens = b.whiteQueens.Set(sq)
	case WhiteKing:
		b.whiteKing = b.whiteKing.Set(sq)
	case BlackPawn:
		b.blackPawns = b.blackPawns.Set(sq)
	case BlackKnight:
		b.blackKnights = b.blackKnights.Set(sq)
	case BlackBishop:
		b.blackBishops = b.blackBishops.Set(sq)
	case BlackRook:
		b.blackRooks = b.blackRooks.Set(sq)
	case BlackQueen:
		b.blackQueens = b.blackQueens.Set(sq)
	case BlackKing:
		b.blackKing = b.blackKing.Set(sq)
	default:
		panic("unexpected piece received by set")
	}
}

// unset removes the piece from the board at the given square
func (b *Board) unset(sq Square) {
	b.whitePawns = b.whitePawns.Unset(sq)
	b.whiteKnights = b.whiteKnights.Unset(sq)
	b.whiteBishops = b.whiteBishops.Unset(sq)
	b.whiteRooks = b.whiteRooks.Unset(sq)
	b.whiteQueens = b.whiteQueens.Unset(sq)
	b.whiteKing = b.whiteKing.Unset(sq)

	b.blackPawns = b.blackPawns.Unset(sq)
	b.blackKnights = b.blackKnights.Unset(sq)
	b.blackBishops = b.blackBishops.Unset(sq)
	b.blackRooks = b.blackRooks.Unset(sq)
	b.blackQueens = b.blackQueens.Unset(sq)
	b.blackKing = b.blackKing.Unset(sq)
}

func (b *Board) allPieces() BitBoard {
	return b.whitePieces() | b.blackPieces()
}

func (b *Board) whitePieces() BitBoard {
	return b.whitePawns |
		b.whiteKnights |
		b.whiteBishops |
		b.whiteRooks |
		b.whiteQueens |
		b.whiteKing
}

func (b *Board) blackPieces() BitBoard {
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
		b.unset(Square(BitBoard(m.To).Down()))
	} else if m.Special.Has(Captures) && m.To == b.EnPassantTarget && b.PlayerInTurn == Black {
		b.unset(Square(BitBoard(m.To).Up()))
	}

	// save the en passant square for the move generation of the en passant moves
	if m.Special.Has(DoublePawnPush) && b.PlayerInTurn == White {
		b.EnPassantTarget = Square(BitBoard(m.From).Up())
	} else if m.Special.Has(DoublePawnPush) && b.PlayerInTurn == Black {
		b.EnPassantTarget = Square(BitBoard(m.From).Down())
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
