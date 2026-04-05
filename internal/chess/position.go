package chess

import (
	"fmt"
)

// Position represents a state of a chess game. This can be used by engines to traverse the game tree.
// It contains the informations needed for efficiently generating moves
type Position struct {
	_ [0]func() // equal guard

	board

	// HalveMoveClock gets incremented after each move that doesn't capture or move a pawn.
	// max 100, see https://www.chessprogramming.org/Fifty-move_Rule
	HalfmoveClock uint8

	// Moves gets incremented every time black played a move
	Moves int

	// history contains the previous positions in a stack, useful for [Position.UndoMove]
	history []board

	// Calculated properties, will be reset after a move is done:

	// attacksTo maps each square to a [BitBoard] where all the attacking pieces are 1
	attacksTo squareLookup[bitBoard]
	// attacksFrom maps each square to a [BitBoard] where all the attacked squares are 1
	attacksFrom squareLookup[bitBoard]

	// xRayKingAttacks contains all lines of attacks that could create a check when a piece moves.
	// this is needed to detect pinned pieces
	xRayKingAttacks []attackRay

	// Bitboards for testing for checks in the move generation.
	bishopCheckSquares bitBoard
	rookCheckSquares   bitBoard
	knightCheckSquares bitBoard
	pawnCheckSquares   bitBoard
}

func NewPosition() *Position {
	// p, err := NewPositionFromFEN(DefaultFen, nil)
	p, err := NewPositionFromFEN("8/8/8/KpP3r1/8/8/8/6k1 w - b6 0 1", nil)

	if err != nil {
		panic(err)
	}

	return p
}

func (p *Position) DoMove(m Move) {
	p.history = append(p.history, p.board)

	// will be reset later if applicable
	p.HalfmoveClock++

	p.board.DoMove(m)

	if !m.Special.Has(CastleQueen|CastleKing) && (p.Square(m.From) == Pawn || m.Special.Has(Captures)) {
		// 50 Move rule
		p.HalfmoveClock = 0
	}

	if p.PlayerInTurn == White {
		// Fullmove counter only increments after black played
		p.Moves++
	}

	// TODO: do this more cleverly, e.g. incrementally update
	p.computeAll()
}

func (p *Position) reset() {
	p.attacksFrom = squareLookup[bitBoard]{}
	p.attacksTo = squareLookup[bitBoard]{}

	// reuse the same cap, but trim the attacks:
	p.xRayKingAttacks = p.xRayKingAttacks[0:0]
}

func (p *Position) computeAll() {
	p.reset()

	p.attacksTo, p.attacksFrom = calculateAttackMaps(
		// we need to exclude our king so that it wont count as blocking an
		// enemy's sliding piece attack
		p.all()&^p.ourKing,
		p.theirKing,
		p.theirQueens,
		p.theirRooks,
		p.theirBishops,
		p.theirKnights,
		p.theirPawns,
	)

	p.calculateSlidingKingAttacks(
		p.ourKing,
		p.theirQueens,
		p.theirRooks,
		p.theirBishops,
	)

	p.calculateCheckSquares(p.theirKing, p.ours, p.theirs)
}

func (p *Position) UndoMove() {
	if len(p.history) == 0 {
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

	from := p.toCurrentPerspective(m.From)
	to := p.toCurrentPerspective(m.To)

	// fast path rejection

	if p.ours&from == 0 {
		return Move{}, fmt.Errorf("move not valid for current position: from is not a piece of current player")
	}

	if p.ours&to != 0 {
		return Move{}, fmt.Errorf("move target cannot be one of our pieces")
	}

	// slow path rejection, generate the moves for the piece

	piece := p.at(from)

	legalMoves := make([]Move, 0, maxMovesSinglePiece)

	p.generateMovesForPiece(piece, from, &legalMoves)

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
	return p.attacksTo.get(p.ourKing) != 0
}

// IsCheckMate is meant to be called by the visualization and returns true if the current player is in checkmate
func (p *Position) IsCheckMate() bool {
	return p.IsCheck() && len(p.GenerateMoves()) == 0
}

// IsDraw is meant to be called by the visualization and returns true if the game is drewn
func (p *Position) IsDraw() bool {
	if p.HalfmoveClock > 99 {
		return true
	}

	return !p.IsCheck() && len(p.GenerateMoves()) == 0
}

func (p *Position) Copy() *Position {
	c := *p

	return &c
}
