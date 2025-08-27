package chess

import (
	"iter"
	"strings"

	"github.com/rswilli/chess/internal/zobrist"
)

// position contains the board representation and the state of the position. This is tracked in the history
// of [Game] and contains all informations needed for move generation (exept draw by repetition).
type position struct {
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

func (p position) String() string {
	sb := strings.Builder{}

	for rank := range 8 {
		sb.WriteString("+---+---+---+---+---+---+---+---+\n|")
		for file := range 8 {
			piece := p.Square(NewSquare(rank, file))
			sb.WriteRune(' ')
			sb.WriteRune(piece.Rune())
			sb.WriteRune(' ')
			sb.WriteRune('|')
		}
		sb.WriteByte('\n')
	}
	sb.WriteString("+---+---+---+---+---+---+---+---+\n")

	switch p.PlayerInTurn {
	case White:
		sb.WriteString("White\n")
	case Black:
		sb.WriteString("Black\n")
	default:
		panic("unexpected player turn")
	}

	sb.WriteString(p.castling.String())
	if p.enPassantTarget != InvalidSquare {
		sb.WriteString("\nEn passant to: ")
		sb.WriteString(p.enPassantTarget.String())
	}

	return sb.String()
}

func (p *position) Square(sq Square) Piece {
	switch {
	case p.whitePawns.Has(sq):
		return WhitePawn
	case p.whiteKnights.Has(sq):
		return WhiteKnight
	case p.whiteBishops.Has(sq):
		return WhiteBishop
	case p.whiteRooks.Has(sq):
		return WhiteRook
	case p.whiteQueens.Has(sq):
		return WhiteQueen
	case p.whiteKing.Has(sq):
		return WhiteKing
	case p.blackPawns.Has(sq):
		return BlackPawn
	case p.blackKnights.Has(sq):
		return BlackKnight
	case p.blackBishops.Has(sq):
		return BlackBishop
	case p.blackRooks.Has(sq):
		return BlackRook
	case p.blackQueens.Has(sq):
		return BlackQueen
	case p.blackKing.Has(sq):
		return BlackKing
	default:
		return Empty
	}
}

// set sets the given piece on the given square
func (p *position) set(sq Square, piece Piece) {
	switch piece {
	case WhitePawn:
		p.whitePawns = p.whitePawns.Set(sq)
	case WhiteKnight:
		p.whiteKnights = p.whiteKnights.Set(sq)
	case WhiteBishop:
		p.whiteBishops = p.whiteBishops.Set(sq)
	case WhiteRook:
		p.whiteRooks = p.whiteRooks.Set(sq)
	case WhiteQueen:
		p.whiteQueens = p.whiteQueens.Set(sq)
	case WhiteKing:
		p.whiteKing = p.whiteKing.Set(sq)
	case BlackPawn:
		p.blackPawns = p.blackPawns.Set(sq)
	case BlackKnight:
		p.blackKnights = p.blackKnights.Set(sq)
	case BlackBishop:
		p.blackBishops = p.blackBishops.Set(sq)
	case BlackRook:
		p.blackRooks = p.blackRooks.Set(sq)
	case BlackQueen:
		p.blackQueens = p.blackQueens.Set(sq)
	case BlackKing:
		p.blackKing = p.blackKing.Set(sq)
	default:
		panic("unexpected piece received by set")
	}

	// incrementally update the hashkey
	p.HashKey = p.HashKey.Update(piece.zobrist() + sq.index())
}

// unset removes the piece from the board at the given square
func (p *position) unset(sq Square) {
	piece := p.Square(sq)

	if piece != Empty {
		// incrementally update the hashkey
		p.HashKey = p.HashKey.Update(piece.zobrist() + sq.index())
	}

	p.whitePawns = p.whitePawns.Unset(sq)
	p.whiteKnights = p.whiteKnights.Unset(sq)
	p.whiteBishops = p.whiteBishops.Unset(sq)
	p.whiteRooks = p.whiteRooks.Unset(sq)
	p.whiteQueens = p.whiteQueens.Unset(sq)
	p.whiteKing = p.whiteKing.Unset(sq)

	p.blackPawns = p.blackPawns.Unset(sq)
	p.blackKnights = p.blackKnights.Unset(sq)
	p.blackBishops = p.blackBishops.Unset(sq)
	p.blackRooks = p.blackRooks.Unset(sq)
	p.blackQueens = p.blackQueens.Unset(sq)
	p.blackKing = p.blackKing.Unset(sq)
}

func (p *position) removeCastling(ab CastlingAbility) {
	if p.castling.Has(ab) {
		// the hash contains the castling ability, so update it
		p.HashKey = p.HashKey.Update(ab.zobrist())
		p.castling &^= ab
	}
}

func (p *position) whitePieces() BitBoard {
	return p.whitePawns |
		p.whiteKnights |
		p.whiteBishops |
		p.whiteRooks |
		p.whiteQueens |
		p.whiteKing
}

func (p *position) blackPieces() BitBoard {
	return p.blackPawns |
		p.blackKnights |
		p.blackBishops |
		p.blackRooks |
		p.blackQueens |
		p.blackKing
}

var a1 = MustParseSquare("a1")
var e1 = MustParseSquare("e1")
var a8 = MustParseSquare("a8")
var h1 = MustParseSquare("h1")
var e8 = MustParseSquare("e8")
var h8 = MustParseSquare("h8")

// ours returns the pieces of the current player
func (p *position) ours() BitBoard {
	if p.PlayerInTurn == White {
		return p.whitePieces()
	} else {
		return p.blackPieces()
	}
}

// theirs returns the pieces of the current opponent player
func (p *position) theirs() BitBoard {
	if p.PlayerInTurn == Black {
		return p.whitePieces()
	} else {
		return p.blackPieces()
	}
}

// all returns a BitBoard containing all pieces
func (p *position) all() BitBoard {
	return p.whitePieces() | p.blackPieces()
}

// pieces returns an iterator over the position, row-wise starting from top left
func (p *position) pieces() iter.Seq[Piece] {
	return func(yield func(Piece) bool) {
		for i := range 64 {
			sq := Square(1 << i)
			piece := p.Square(sq)
			if !yield(piece) {
				return
			}
		}
	}
}

// hashFull computes the [zobrist.Hash] for the [position]. It is way slower than incrementally updating the hash,
// which is what [Game.DoMove] does
func (p *position) hashFull() {

	pieces := make([]int, 0, 64)

	for p := range p.pieces() {
		pieces = append(pieces, p.zobrist())
	}

	p.HashKey = zobrist.NewBoard(pieces...)

	if p.castling.Has(CastleBlackKing) {
		p.HashKey = p.HashKey.Update(zobrist.BlackCastleKing)
	}
	if p.castling.Has(CastleBlackQueen) {
		p.HashKey = p.HashKey.Update(zobrist.BlackCastleQueen)
	}
	if p.castling.Has(CastleWhiteKing) {
		p.HashKey = p.HashKey.Update(zobrist.WhiteCastleKing)
	}
	if p.castling.Has(CastleWhiteQueen) {
		p.HashKey = p.HashKey.Update(zobrist.WhiteCastleQueen)
	}

	if p.enPassantTarget != 0 {
		p.HashKey = p.HashKey.Update(zobrist.EnPassantAFile + p.enPassantTarget.file())
	}

	if p.PlayerInTurn == Black {
		p.HashKey = p.HashKey.Update(zobrist.BlackToMove)
	}
}
