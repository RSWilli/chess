package chess

import (
	"fmt"
	"iter"
	"strings"

	"github.com/rswilli/chess/internal/zobrist"
)

// board contains the board representation and the state of the board. This is tracked in the history
// of [Position] and contains all informations needed for evaluating the position.
type board struct {
	PlayerInTurn    Piece
	castling        CastlingAbility
	enPassantTarget Square

	// our pieces are the pieces of the color of the player in turn
	ours       bitBoard
	ourPawns   bitBoard
	ourKnights bitBoard
	ourBishops bitBoard
	ourRooks   bitBoard
	ourQueens  bitBoard
	ourKing    bitBoard

	// their pieces are the pieces of the color of the opponent
	theirs       bitBoard
	theirPawns   bitBoard
	theirKnights bitBoard
	theirBishops bitBoard
	theirRooks   bitBoard
	theirQueens  bitBoard
	theirKing    bitBoard

	HashKey zobrist.Hash
}

func (b *board) isValid() bool {
	if b.ourKing.Count() != 1 || b.theirKing.Count() != 1 {
		return false
	}

	return true
}

func (b *board) DoMove(m Move) {
	if m.Special.Has(CastleQueen | CastleKing) {
		// Castling move:
		if b.PlayerInTurn == White && m.Special == CastleKing {
			b.unset(m.From)
			b.unset(h1)

			b.set(g1, WhiteKing)
			b.set(f1, WhiteRook)

			b.removeCastling(CastleWhiteKing)
			b.removeCastling(CastleWhiteQueen)
		}
		if b.PlayerInTurn == White && m.Special == CastleQueen {
			b.unset(m.From)
			b.unset(a1)

			b.set(c1, WhiteKing)
			b.set(d1, WhiteRook)

			b.removeCastling(CastleWhiteKing)
			b.removeCastling(CastleWhiteQueen)
		}
		if b.PlayerInTurn == Black && m.Special == CastleKing {
			b.unset(m.From)
			b.unset(h8)

			b.set(g8, BlackKing)
			b.set(f8, BlackRook)

			b.removeCastling(CastleBlackKing)
			b.removeCastling(CastleBlackQueen)
		}
		if b.PlayerInTurn == Black && m.Special == CastleQueen {
			b.unset(m.From)
			b.unset(a8)

			b.set(c8, BlackKing)
			b.set(d8, BlackRook)

			b.removeCastling(CastleBlackKing)
			b.removeCastling(CastleBlackQueen)
		}

		b.clearEnpassant()
	} else {
		piece := b.Square(m.From)

		// clear the old square
		b.unset(m.From)
		b.unset(m.To)

		switch {
		case m.Special.Has(PromoteQueen):
			piece = Queen | b.PlayerInTurn
		case m.Special.Has(PromoteRook):
			piece = Rook | b.PlayerInTurn
		case m.Special.Has(PromoteKnight):
			piece = Knight | b.PlayerInTurn
		case m.Special.Has(PromoteBishop):
			piece = Bishop | b.PlayerInTurn
		}

		// remove the en-passant captured pawn, no need to check for the piece type since the en passant
		// square is always empty, so no other move can capture on it
		if m.Special.Has(Captures) && m.To == b.enPassantTarget && b.PlayerInTurn == White {
			b.unset(Square(bitBoard(m.To).Down()))
		} else if m.Special.Has(Captures) && m.To == b.enPassantTarget && b.PlayerInTurn == Black {
			b.unset(Square(bitBoard(m.To).Up()))
		}

		// save the en passant square for the move generation of the en passant moves
		if m.Special.Has(DoublePawnPush) && b.PlayerInTurn == White {
			b.setEnpassant(Square(bitBoard(m.From).Up()))
		} else if m.Special.Has(DoublePawnPush) && b.PlayerInTurn == Black {
			b.setEnpassant(Square(bitBoard(m.From).Down()))
		} else {
			b.clearEnpassant()
		}

		// prevent castling moves, but only if set because the hash update depends on it:
		if b.castling.Has(CastleWhiteQueen) && (m.From == a1 || m.From == e1 || m.To == a1) {
			b.removeCastling(CastleWhiteQueen)
		}

		if b.castling.Has(CastleWhiteKing) && (m.From == h1 || m.From == e1 || m.To == h1) {
			b.removeCastling(CastleWhiteKing)
		}

		if b.castling.Has(CastleBlackQueen) && (m.From == a8 || m.From == e8 || m.To == a8) {
			b.removeCastling(CastleBlackQueen)
		}

		if b.castling.Has(CastleBlackKing) && (m.From == h8 || m.From == e8 || m.To == h8) {
			b.removeCastling(CastleBlackKing)
		}

		b.set(m.To, piece)
	}

	// FIXME: the wiki says that losing castling increments the halfmove clock
	// // halfmove clock for 50 move rule, see https://www.chessprogramming.org/Halfmove_Clock
	// previousCastling := p.history[len(p.history)-1].castling
	// nowCastling := p.castling

	// lostCastling := previousCastling != nowCastling

	b.SwitchSide()

	b.HashKey = b.HashKey.Update(zobrist.SwitchSide)

	if !b.isValid() {
		panic(fmt.Sprintf("DoMove resulted in invalid board after %s", m.String()))
	}
}

// SwitchSide changes the perspective of the [board]
func (b *board) SwitchSide() {
	b.PlayerInTurn = b.PlayerInTurn.SwitchColor()

	b.ours, b.theirs = b.theirs.rotate180(), b.ours.rotate180()

	b.ourPawns, b.theirPawns = b.theirPawns.rotate180(), b.ourPawns.rotate180()
	b.ourKnights, b.theirKnights = b.theirKnights.rotate180(), b.ourKnights.rotate180()
	b.ourBishops, b.theirBishops = b.theirBishops.rotate180(), b.ourBishops.rotate180()
	b.ourRooks, b.theirRooks = b.theirRooks.rotate180(), b.ourRooks.rotate180()
	b.ourQueens, b.theirQueens = b.theirQueens.rotate180(), b.ourQueens.rotate180()
	b.ourKing, b.theirKing = b.theirKing.rotate180(), b.ourKing.rotate180()
}

func (b board) ASCIIArt() string {
	sb := strings.Builder{}

	sb.WriteString("    a   b   c   d   e   f   g   h  \n")
	for rank := range 8 {
		sb.WriteString("  +---+---+---+---+---+---+---+---+\n")
		fmt.Fprintf(&sb, "%d |", 8-rank)
		for file := range 8 {
			piece := b.Square(NewSquare(rank, file))
			sb.WriteRune(' ')
			sb.WriteRune(piece.Rune())
			sb.WriteRune(' ')
			sb.WriteRune('|')
		}
		fmt.Fprintf(&sb, " %d\n", 8-rank)
	}
	sb.WriteString("  +---+---+---+---+---+---+---+---+\n")
	sb.WriteString("    a   b   c   d   e   f   g   h  \n")

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

func (b *board) toWhitePerspective(sq bitBoard) Square {
	switch b.PlayerInTurn {
	case White:
		return Square(sq)
	case Black:
		return Square(sq.rotate180())
	default:
		panic(fmt.Sprintf("invalid player in turn: %s", b.PlayerInTurn))
	}
}

func (b *board) toCurrentPerspective(sq Square) bitBoard {
	switch b.PlayerInTurn {
	case White:
		return bitBoard(sq)
	case Black:
		return bitBoard(sq).rotate180()
	default:
		panic(fmt.Sprintf("invalid player in turn: %s", b.PlayerInTurn))
	}
}

// Square returns the piece at the given square.
func (b *board) Square(sq Square) Piece {
	relativeSquare := b.toCurrentPerspective(sq)
	return b.at(relativeSquare)
}

// at returns the piece at the given BitBoard (Square), but from the current players perspective.
func (b *board) at(sq bitBoard) Piece {
	switch {
	case b.ourPawns.Has(sq):
		return Pawn | b.PlayerInTurn
	case b.ourKnights.Has(sq):
		return Knight | b.PlayerInTurn
	case b.ourBishops.Has(sq):
		return Bishop | b.PlayerInTurn
	case b.ourRooks.Has(sq):
		return Rook | b.PlayerInTurn
	case b.ourQueens.Has(sq):
		return Queen | b.PlayerInTurn
	case b.ourKing.Has(sq):
		return King | b.PlayerInTurn
		// opponent:
	case b.theirPawns.Has(sq):
		return Pawn | b.PlayerInTurn.SwitchColor()
	case b.theirKnights.Has(sq):
		return Knight | b.PlayerInTurn.SwitchColor()
	case b.theirBishops.Has(sq):
		return Bishop | b.PlayerInTurn.SwitchColor()
	case b.theirRooks.Has(sq):
		return Rook | b.PlayerInTurn.SwitchColor()
	case b.theirQueens.Has(sq):
		return Queen | b.PlayerInTurn.SwitchColor()
	case b.theirKing.Has(sq):
		return King | b.PlayerInTurn.SwitchColor()
	default:
		return Empty
	}
}

// set sets the given piece on the given square
func (b *board) set(sq Square, piece Piece) {
	relativeSquare := b.toCurrentPerspective(sq)

	opponent := b.PlayerInTurn.SwitchColor()

	switch piece.Color() {
	case b.PlayerInTurn:
		b.ours = b.ours.Set(relativeSquare)
	case opponent:
		b.theirs = b.theirs.Set(relativeSquare)
	default:
		panic("unexpected piece color")
	}

	switch piece {
	case b.PlayerInTurn | Pawn:
		b.ourPawns = b.ourPawns.Set(relativeSquare)
	case b.PlayerInTurn | Knight:
		b.ourKnights = b.ourKnights.Set(relativeSquare)
	case b.PlayerInTurn | Bishop:
		b.ourBishops = b.ourBishops.Set(relativeSquare)
	case b.PlayerInTurn | Rook:
		b.ourRooks = b.ourRooks.Set(relativeSquare)
	case b.PlayerInTurn | Queen:
		b.ourQueens = b.ourQueens.Set(relativeSquare)
	case b.PlayerInTurn | King:
		b.ourKing = b.ourKing.Set(relativeSquare)

	case opponent | Pawn:
		b.theirPawns = b.theirPawns.Set(relativeSquare)
	case opponent | Knight:
		b.theirKnights = b.theirKnights.Set(relativeSquare)
	case opponent | Bishop:
		b.theirBishops = b.theirBishops.Set(relativeSquare)
	case opponent | Rook:
		b.theirRooks = b.theirRooks.Set(relativeSquare)
	case opponent | Queen:
		b.theirQueens = b.theirQueens.Set(relativeSquare)
	case opponent | King:
		b.theirKing = b.theirKing.Set(relativeSquare)
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

	relativeSquare := b.toCurrentPerspective(sq)

	b.ours = b.ours.Unset(relativeSquare)
	b.theirs = b.theirs.Unset(relativeSquare)

	b.ourPawns = b.ourPawns.Unset(relativeSquare)
	b.ourKnights = b.ourKnights.Unset(relativeSquare)
	b.ourBishops = b.ourBishops.Unset(relativeSquare)
	b.ourRooks = b.ourRooks.Unset(relativeSquare)
	b.ourQueens = b.ourQueens.Unset(relativeSquare)
	b.ourKing = b.ourKing.Unset(relativeSquare)

	b.theirPawns = b.theirPawns.Unset(relativeSquare)
	b.theirKnights = b.theirKnights.Unset(relativeSquare)
	b.theirBishops = b.theirBishops.Unset(relativeSquare)
	b.theirRooks = b.theirRooks.Unset(relativeSquare)
	b.theirQueens = b.theirQueens.Unset(relativeSquare)
	b.theirKing = b.theirKing.Unset(relativeSquare)
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
		b.HashKey = b.HashKey.Update(zobrist.EnPassantAFile + b.enPassantTarget.File())
		b.enPassantTarget = InvalidSquare
	}
}

func (b *board) setEnpassant(sq Square) {
	b.clearEnpassant()

	b.enPassantTarget = sq
	b.HashKey = b.HashKey.Update(zobrist.EnPassantAFile + b.enPassantTarget.File())
}

// all returns a BitBoard containing all pieces
func (b *board) all() bitBoard {
	return b.ours | b.theirs
}

func (b *board) index(i int) Piece {
	sq := NewSquareFromIndex(i)
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

// hashFull computes the [zobrist.Hash] for the [board]. It is way slower than incrementally updating the hash,
// which is what [board.DoMove] does
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
		b.HashKey = b.HashKey.Update(zobrist.EnPassantAFile + b.enPassantTarget.File())
	}

	if b.PlayerInTurn == Black {
		b.HashKey = b.HashKey.Update(zobrist.BlackToMove)
	}
}
