package chess

import (
	"iter"
	"strings"

	"github.com/rswilli/chess/internal/zobrist"
)

// board contains the board representation and the state of the board. This is tracked in the history
// of [Position] and contains all informations needed for move generation (exept draw by repetition).
type board struct {
	// max 50, see https://www.chessprogramming.org/Fifty-move_Rule
	HalfmoveClock uint8

	PlayerInTurn    Piece
	castling        CastlingAbility
	enPassantTarget Square

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

	HashKey zobrist.Hash
}

func (b board) String() string {
	sb := strings.Builder{}

	for rank := range 8 {
		sb.WriteString("+---+---+---+---+---+---+---+---+\n|")
		for file := range 8 {
			piece := b.Square(NewSquare(rank, file))
			sb.WriteRune(' ')
			sb.WriteRune(piece.Rune())
			sb.WriteRune(' ')
			sb.WriteRune('|')
		}
		sb.WriteByte('\n')
	}
	sb.WriteString("+---+---+---+---+---+---+---+---+\n")

	switch b.PlayerInTurn {
	case White:
		sb.WriteString("White\n")
	case Black:
		sb.WriteString("Black\n")
	default:
		panic("unexpected player turn")
	}

	sb.WriteString(b.castling.String())
	if b.enPassantTarget != InvalidSquare {
		sb.WriteString("\nEn passant to: ")
		sb.WriteString(b.enPassantTarget.String())
	}

	return sb.String()
}

func (b *board) Square(sq Square) Piece {
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
func (b *board) set(sq Square, piece Piece) {
	switch piece {
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

	// incrementally update the hashkey
	b.HashKey = b.HashKey.Update(piece.zobrist() + sq.index())
}

// unset removes the piece from the board at the given square
func (b *board) unset(sq Square) {
	piece := b.Square(sq)

	if piece != Empty {
		// incrementally update the hashkey
		b.HashKey = b.HashKey.Update(piece.zobrist() + sq.index())
	}

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

func (b *board) removeCastling(ab CastlingAbility) {
	if b.castling.Has(ab) {
		// the hash contains the castling ability, so update it
		b.HashKey = b.HashKey.Update(ab.zobrist())
		b.castling &^= ab
	}
}

func (b *board) clearEnpassant() {
	if b.enPassantTarget != InvalidSquare {
		// the hash contains the en passant square, so clear it
		b.HashKey = b.HashKey.Update(zobrist.EnPassantAFile + b.enPassantTarget.file())
		b.enPassantTarget = InvalidSquare
	}
}

func (b *board) setEnpassant(sq Square) {
	b.clearEnpassant()

	b.enPassantTarget = sq
	b.HashKey = b.HashKey.Update(zobrist.EnPassantAFile + b.enPassantTarget.file())
}

func (b *board) whitePieces() BitBoard {
	return b.whitePawns |
		b.whiteKnights |
		b.whiteBishops |
		b.whiteRooks |
		b.whiteQueens |
		b.whiteKing
}

func (b *board) blackPieces() BitBoard {
	return b.blackPawns |
		b.blackKnights |
		b.blackBishops |
		b.blackRooks |
		b.blackQueens |
		b.blackKing
}

var a1 = MustParseSquare("a1")
var e1 = MustParseSquare("e1")
var a8 = MustParseSquare("a8")
var h1 = MustParseSquare("h1")
var e8 = MustParseSquare("e8")
var h8 = MustParseSquare("h8")

// ours returns the pieces of the current player
func (b *board) ours() BitBoard {
	if b.PlayerInTurn == White {
		return b.whitePieces()
	} else {
		return b.blackPieces()
	}
}

// theirs returns the pieces of the current opponent player
func (b *board) theirs() BitBoard {
	if b.PlayerInTurn == Black {
		return b.whitePieces()
	} else {
		return b.blackPieces()
	}
}

// all returns a BitBoard containing all pieces
func (b *board) all() BitBoard {
	return b.whitePieces() | b.blackPieces()
}

func (b *board) index(i int) Piece {
	sq := Square(1 << i)
	return b.Square(sq)
}

// pieces returns an iterator over the position, row-wise starting from top left
func (b *board) pieces() iter.Seq[Piece] {
	return func(yield func(Piece) bool) {
		for i := range 64 {
			piece := b.index(i)
			if !yield(piece) {
				return
			}
		}
	}
}

// hashFull computes the [zobrist.Hash] for the [position]. It is way slower than incrementally updating the hash,
// which is what [Position.DoMove] does
func (b *board) hashFull() {

	pieces := make([]int, 0, 64)

	for p := range b.pieces() {
		pieces = append(pieces, p.zobrist())
	}

	b.HashKey = zobrist.NewBoard(pieces...)

	if b.castling.Has(CastleBlackKing) {
		b.HashKey = b.HashKey.Update(zobrist.BlackCastleKing)
	}
	if b.castling.Has(CastleBlackQueen) {
		b.HashKey = b.HashKey.Update(zobrist.BlackCastleQueen)
	}
	if b.castling.Has(CastleWhiteKing) {
		b.HashKey = b.HashKey.Update(zobrist.WhiteCastleKing)
	}
	if b.castling.Has(CastleWhiteQueen) {
		b.HashKey = b.HashKey.Update(zobrist.WhiteCastleQueen)
	}

	if b.enPassantTarget != 0 {
		b.HashKey = b.HashKey.Update(zobrist.EnPassantAFile + b.enPassantTarget.file())
	}

	if b.PlayerInTurn == Black {
		b.HashKey = b.HashKey.Update(zobrist.BlackToMove)
	}
}
