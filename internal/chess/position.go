package chess

import (
	"fmt"

	"github.com/rswilli/chess/internal/zobrist"
)

// Position represents a state of a chess game. This can be used by engines to traverse the game tree.
type Position struct {
	_ [0]func() // equal guard

	board

	Moves int

	// history contains the previous positions in a stack, useful for [Position.UndoMove]
	history []board

	// Calculated properties, will be reset after a move is done:

	// attacksTo maps each square to a [BitBoard] where all the attacking pieces are 1
	attacksTo squareLookup[BitBoard]
	// attacksFrom maps each square to a [BitBoard] where all the attacked squares are 1
	attacksFrom squareLookup[BitBoard]

	// xRayKingAttacks contains all lines of attacks that could create a check when a piece moves.
	// this is needed to detect pinned pieces
	xRayKingAttacks [8]attackRay
}

func NewPosition() *Position {
	p, err := NewPositionFromFEN(DefaultFen)
	// b, err := NewGameFromFEN("r3kbnr/1b3ppp/pqn5/1pp1P3/3p4/1BN2N2/PP2QPPP/R1BR2K1 w kq - 0 1")

	if err != nil {
		panic(err)
	}

	return p
}

func (p *Position) DoMove(m Move) {
	p.history = append(p.history, p.board)

	pawnMove := false

	if m.Special.Has(CastleLong | CastleShort) {
		// Castling move:
		if p.PlayerInTurn == White && m.Special == CastleShort {
			p.unset(e1)
			p.unset(h1)

			p.set(whiteCastleKingKingTarget, WhiteKing)
			p.set(whiteCastleKingRookTarget, WhiteRook)

			p.removeCastling(CastleWhiteKing)
			p.removeCastling(CastleWhiteQueen)
		}
		if p.PlayerInTurn == White && m.Special == CastleLong {
			p.unset(e1)
			p.unset(a1)

			p.set(whiteCastleQueenKingTarget, WhiteKing)
			p.set(whiteCastleQueenRookTarget, WhiteRook)

			p.removeCastling(CastleWhiteKing)
			p.removeCastling(CastleWhiteQueen)
		}
		if p.PlayerInTurn == Black && m.Special == CastleShort {
			p.unset(e8)
			p.unset(h8)

			p.set(blackCastleKingKingTarget, BlackKing)
			p.set(blackCastleKingRookTarget, BlackRook)

			p.removeCastling(CastleBlackKing)
			p.removeCastling(CastleBlackQueen)
		}
		if p.PlayerInTurn == Black && m.Special == CastleLong {
			p.unset(e8)
			p.unset(a8)

			p.set(blackCastleQueenKingTarget, BlackKing)
			p.set(blackCastleQueenRookTarget, BlackRook)

			p.removeCastling(CastleBlackKing)
			p.removeCastling(CastleBlackQueen)
		}

		p.clearEnpassant()
	} else {
		piece := p.Square(m.From)

		if piece == Pawn {
			pawnMove = true
		}

		// clear the old square
		p.unset(m.From)
		p.unset(m.To)

		switch {
		case m.Special.Has(PromoteQueen):
			piece = Queen | p.PlayerInTurn
		case m.Special.Has(PromoteRook):
			piece = Rook | p.PlayerInTurn
		case m.Special.Has(PromoteKnight):
			piece = Knight | p.PlayerInTurn
		case m.Special.Has(PromoteBishop):
			piece = Bishop | p.PlayerInTurn
		}

		// remove the en-passant captured pawn, no need to check for the piece type since the en passant
		// square is always empty, so no other move can capture on it
		if m.Special.Has(Captures) && m.To == p.enPassantTarget && p.PlayerInTurn == White {
			p.unset(Square(BitBoard(m.To).Down()))
		} else if m.Special.Has(Captures) && m.To == p.enPassantTarget && p.PlayerInTurn == Black {
			p.unset(Square(BitBoard(m.To).Up()))
		}

		// save the en passant square for the move generation of the en passant moves
		if m.Special.Has(DoublePawnPush) && p.PlayerInTurn == White {
			p.setEnpassant(Square(BitBoard(m.From).Up()))
		} else if m.Special.Has(DoublePawnPush) && p.PlayerInTurn == Black {
			p.setEnpassant(Square(BitBoard(m.From).Down()))
		} else {
			p.clearEnpassant()
		}

		// prevent castling moves, but only if set because the hash update depends on it:
		if p.castling.Has(CastleWhiteQueen) && (m.From == a1 || m.From == e1 || m.To == a1) {
			p.removeCastling(CastleWhiteQueen)
		}

		if p.castling.Has(CastleWhiteKing) && (m.From == h1 || m.From == e1 || m.To == h1) {
			p.removeCastling(CastleWhiteKing)
		}

		if p.castling.Has(CastleBlackQueen) && (m.From == a8 || m.From == e8 || m.To == a8) {
			p.removeCastling(CastleBlackQueen)
		}

		if p.castling.Has(CastleBlackKing) && (m.From == h8 || m.From == e8 || m.To == h8) {
			p.removeCastling(CastleBlackKing)
		}

		p.set(m.To, piece)
	}

	// halfmove clock for 50 move rule, see https://www.chessprogramming.org/Halfmove_Clock
	previousCastling := p.history[len(p.history)-1].castling
	nowCastling := p.castling

	lostCastling := previousCastling != nowCastling

	if !pawnMove || lostCastling {
		p.HalfmoveClock++
	} else {
		p.HalfmoveClock = 0
	}

	if p.PlayerInTurn == Black {
		// Fullmove counter only increments after black played
		p.Moves++
	}

	// Move done, reset state and recalculate:
	if p.PlayerInTurn == White {
		p.PlayerInTurn = Black
	} else {
		p.PlayerInTurn = White
	}

	p.HashKey = p.HashKey.Update(zobrist.SwitchSide)

	// TODO: do this more cleverly, e.g. incrementally update
	p.computeAll()
}

func (p *Position) reset() {
	p.attacksFrom = squareLookup[BitBoard]{}
	p.attacksTo = squareLookup[BitBoard]{}
	p.xRayKingAttacks = [8]attackRay{}
}

func (p *Position) computeAll() {
	p.reset()

	if p.PlayerInTurn == Black {
		p.attacksTo, p.attacksFrom = calculateAttackMaps(
			// we need to exclude our king so that it wont count as blocking an
			// enemy's sliding piece attack
			p.all()&^p.blackKing,
			p.whiteKing,
			p.whiteQueens,
			p.whiteRooks,
			p.whiteBishops,
			p.whiteKnights,
			p.whitePawns,
			White,
		)

		p.calculateSlidingKingAttacks(
			p.blackKing,
			p.whiteQueens,
			p.whiteRooks,
			p.whiteBishops,
		)
	} else {
		p.attacksTo, p.attacksFrom = calculateAttackMaps(
			// we need to exclude our king so that it wont count as blocking an
			// enemy's sliding piece attack
			p.all()&^p.whiteKing,
			p.blackKing,
			p.blackQueens,
			p.blackRooks,
			p.blackBishops,
			p.blackKnights,
			p.blackPawns,
			Black,
		)

		p.calculateSlidingKingAttacks(
			p.whiteKing,
			p.blackQueens,
			p.blackRooks,
			p.blackBishops,
		)
	}
}

func (p *Position) UndoMove() {
	if p.Moves == 0 || len(p.history) == 0 {
		// this is a programming error, so we panic here
		panic("cannot undo move, none was done")
	}

	lastPos := p.history[len(p.history)-1]
	p.history = p.history[:len(p.history)-1]

	p.board = lastPos

	// TODO: do this more cleverly, e.g. incrementally update
	p.computeAll()
}

// maxMovesSinglePiece is the number of max possible moves of a single piece (queen at e4, 27 Moves) with added padding
// for cases I missed :)
const maxMovesSinglePiece = 30

// ParseMove parses the move in algebraic notation (see [ParseMove]) and takes into account the current
// position to account for special moves
func (p *Position) ParseMove(in string) (Move, error) {
	m, err := ParseMove(in)

	if err != nil {
		return Move{}, err
	}

	// fast path rejection

	if p.ours()&BitBoard(m.From) == 0 {
		return Move{}, fmt.Errorf("move not valid for current position: from is not a piece of current player")
	}

	if p.ours()&BitBoard(m.To) != 0 {
		return Move{}, fmt.Errorf("move target cannot be one of our pieces")
	}

	// slow path rejection, generate the moves for the piece

	piece := p.Square(m.From)

	legalMoves := make([]Move, 0, maxMovesSinglePiece)

	p.generateMovesForPiece(piece, BitBoard(m.From), &legalMoves)

	for _, legalMove := range legalMoves {
		if legalMove.From == m.From && legalMove.To == m.To && (m.Special == NoSpecial || legalMove.Special.Has(m.Special)) {
			return legalMove, nil
		}
	}

	return Move{}, fmt.Errorf("given move is not a legal move: %s", in)
}

func (p *Position) Equals(other *Position) bool {
	return p.board == other.board
}

// IsCheck is meant to be called by the visualization and returns true if the current player is in check
func (p *Position) IsCheck() bool {
	if p.PlayerInTurn == White {
		return p.attacksTo.get(p.whiteKing) != 0
	} else {
		return p.attacksTo.get(p.blackKing) != 0
	}
}

// IsCheckMate is meant to be called by the visualization and returns true if the current player is in checkmate
func (p *Position) IsCheckMate() bool {
	return p.IsCheck() && len(p.GenerateMoves()) == 0
}

// IsDraw is meant to be called by the visualization and returns true if the game is drewn
func (p *Position) IsDraw() bool {
	return !p.IsCheck() && len(p.GenerateMoves()) == 0
}

func (p *Position) Copy() *Position {
	c := *p

	historyCopy := make([]board, len(c.history))
	copy(historyCopy, c.history)

	c.reset()

	return &c
}
