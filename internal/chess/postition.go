package chess

import (
	"fmt"
	"strings"
)

// position contains the board representation and the state of the position. This is tracked in the history
// of [Position] and contains all informations needed for move generation (exept draw by repetition).
type position struct {
	// max 50, see https://www.chessprogramming.org/Fifty-move_Rule
	HalfmoveClock uint8

	playerInTurn    Piece
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
}

type Position struct {
	_ [0]func() // equal guard, we have local cached values

	position

	Moves int

	// history contains the previous positions in a stack, useful for [Position.UndoMove]
	history []position

	// Calculated properties, will be reset after a move is done:

	// attacksTo maps each square to a [BitBoard] where all the attacking pieces are 1
	attacksTo squareLookup[BitBoard]
	// attacksFrom maps each square to a [BitBoard] where all the attacked squares are 1
	attacksFrom squareLookup[BitBoard]

	// xRayKingAttacks contains all lines of attacks that could create a check when a piece moves.
	// this is needed to detect pinned pieces
	xRayKingAttacks []attackRay
}

func NewPosition() *Position {
	// b, err := NewPositionFromFEN(DefaultFen)
	b, err := NewPositionFromFEN("2k5/8/8/1BR2q2/8/8/1PKP1N2/8 w - - 0 1")

	if err != nil {
		panic(err)
	}

	return b
}

func (p Position) String() string {
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

	switch p.playerInTurn {
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

func (p *Position) Square(sq Square) Piece {
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
func (board *Position) set(sq Square, p Piece) {
	switch p {
	case WhitePawn:
		board.whitePawns = board.whitePawns.Set(sq)
	case WhiteKnight:
		board.whiteKnights = board.whiteKnights.Set(sq)
	case WhiteBishop:
		board.whiteBishops = board.whiteBishops.Set(sq)
	case WhiteRook:
		board.whiteRooks = board.whiteRooks.Set(sq)
	case WhiteQueen:
		board.whiteQueens = board.whiteQueens.Set(sq)
	case WhiteKing:
		board.whiteKing = board.whiteKing.Set(sq)
	case BlackPawn:
		board.blackPawns = board.blackPawns.Set(sq)
	case BlackKnight:
		board.blackKnights = board.blackKnights.Set(sq)
	case BlackBishop:
		board.blackBishops = board.blackBishops.Set(sq)
	case BlackRook:
		board.blackRooks = board.blackRooks.Set(sq)
	case BlackQueen:
		board.blackQueens = board.blackQueens.Set(sq)
	case BlackKing:
		board.blackKing = board.blackKing.Set(sq)
	default:
		panic("unexpected piece received by set")
	}
}

// unset removes the piece from the board at the given square
func (p *Position) unset(sq Square) {
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

func (p *Position) whitePieces() BitBoard {
	return p.whitePawns |
		p.whiteKnights |
		p.whiteBishops |
		p.whiteRooks |
		p.whiteQueens |
		p.whiteKing
}

func (p *Position) blackPieces() BitBoard {
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

func (p *Position) DoMove(m Move) {
	p.history = append(p.history, p.position)

	pawnMove := false

	if m.Special.Has(CastleLong | CastleShort) {
		// Castling move:
		if p.playerInTurn == White && m.Special == CastleShort {
			p.unset(e1)
			p.unset(h1)

			p.set(whiteCastleKingKingTarget, WhiteKing)
			p.set(whiteCastleKingRookTarget, WhiteRook)

			p.castling &^= CastleWhiteKing | CastleWhiteQueen
		}
		if p.playerInTurn == White && m.Special == CastleLong {
			p.unset(e1)
			p.unset(a1)

			p.set(whiteCastleQueenKingTarget, WhiteKing)
			p.set(whiteCastleQueenRookTarget, WhiteRook)

			p.castling &^= CastleWhiteKing | CastleWhiteQueen
		}
		if p.playerInTurn == Black && m.Special == CastleShort {
			p.unset(e8)
			p.unset(h8)

			p.set(blackCastleKingKingTarget, BlackKing)
			p.set(blackCastleKingRookTarget, BlackRook)

			p.castling &^= CastleBlackKing | CastleBlackQueen
		}
		if p.playerInTurn == Black && m.Special == CastleLong {
			p.unset(e8)
			p.unset(a8)

			p.set(blackCastleQueenKingTarget, BlackKing)
			p.set(blackCastleQueenRookTarget, BlackRook)

			p.castling &^= CastleBlackKing | CastleBlackQueen
		}
	} else {
		piece := p.Square(m.From)

		if piece == Pawn {
			pawnMove = true
		}

		// clear the old square
		p.unset(m.From)
		// must clear all other boards before setting the new one
		p.unset(m.To)

		switch {
		case m.Special.Has(PromoteQueen):
			piece = Queen | p.playerInTurn
		case m.Special.Has(PromoteRook):
			piece = Rook | p.playerInTurn
		case m.Special.Has(PromoteKnight):
			piece = Knight | p.playerInTurn
		case m.Special.Has(PromoteBishop):
			piece = Bishop | p.playerInTurn
		}

		// remove the en-passant captured pawn, no need to check for the piece type since the en passant
		// square is always empty, so no other move can capture on it
		if m.Special.Has(Captures) && m.To == p.enPassantTarget && p.playerInTurn == White {
			p.unset(Square(BitBoard(m.To).Down()))
		} else if m.Special.Has(Captures) && m.To == p.enPassantTarget && p.playerInTurn == Black {
			p.unset(Square(BitBoard(m.To).Up()))
		}

		// save the en passant square for the move generation of the en passant moves
		if m.Special.Has(DoublePawnPush) && p.playerInTurn == White {
			p.enPassantTarget = Square(BitBoard(m.From).Up())
		} else if m.Special.Has(DoublePawnPush) && p.playerInTurn == Black {
			p.enPassantTarget = Square(BitBoard(m.From).Down())
		} else {
			p.enPassantTarget = InvalidSquare
		}

		// prevent castling moves:
		switch m.From {
		case a1:
			p.castling &^= CastleWhiteQueen
		case e1:
			p.castling &^= CastleWhiteQueen | CastleWhiteKing
		case h1:
			p.castling &^= CastleWhiteKing
		case a8:
			p.castling &^= CastleBlackQueen
		case e8:
			p.castling &^= CastleBlackQueen | CastleBlackKing
		case h8:
			p.castling &^= CastleBlackKing
		}

		switch m.To {
		case a1:
			p.castling &^= CastleWhiteQueen
		case h1:
			p.castling &^= CastleWhiteKing
		case a8:
			p.castling &^= CastleBlackQueen
		case h8:
			p.castling &^= CastleBlackKing
		}

		p.set(m.To, piece)
	}

	// halfmove clock for 50 move rule, see https://www.chessprogramming.org/Halfmove_Clock
	previousCastling := p.history[len(p.history)-1].castling & (CastleWhiteQueen | CastleWhiteKing | CastleBlackQueen | CastleBlackKing)
	nowCastling := p.castling & (CastleWhiteQueen | CastleWhiteKing | CastleBlackQueen | CastleBlackKing)

	lostCastling := previousCastling != nowCastling

	if !pawnMove || lostCastling {
		p.HalfmoveClock++
	} else {
		p.HalfmoveClock = 0
	}

	// Move done, reset state and recalculate:
	if p.playerInTurn == White {
		p.playerInTurn = Black
	} else {
		p.playerInTurn = White
	}

	p.reset()
}

// ours returns the pieces of the current player
func (p *Position) ours() BitBoard {
	if p.playerInTurn == White {
		return p.whitePieces()
	} else {
		return p.blackPieces()
	}
}

// theirs returns the pieces of the current opponent player
func (p *Position) theirs() BitBoard {
	if p.playerInTurn == Black {
		return p.whitePieces()
	} else {
		return p.blackPieces()
	}
}

// all returns a BitBoard containing all pieces
func (p *Position) all() BitBoard {
	return p.whitePieces() | p.blackPieces()
}

func (p *Position) reset() {
	p.attacksFrom = squareLookup[BitBoard]{}
	p.attacksTo = squareLookup[BitBoard]{}
	p.xRayKingAttacks = nil
}

func (p *Position) computeAll() {
	p.reset()

	if p.playerInTurn == Black {
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

		p.xRayKingAttacks = calculateXRayAttacks(
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

		p.xRayKingAttacks = calculateXRayAttacks(
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

	p.position = lastPos

	p.reset()
}

func (p *Position) Equals(other *Position) bool {
	return p.position == other.position
}

func (p *Position) Key(m Move) uint64 {
	return uint64(p.all()) ^ uint64(p.castling) ^ uint64(p.enPassantTarget) ^ uint64(p.playerInTurn)
}

// IsCheck is meant to be called by the visualization and returns true if the current player is in check
func (p *Position) IsCheck() bool {
	if p.playerInTurn == White {
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

	historyCopy := make([]position, len(c.history))
	copy(historyCopy, c.history)

	c.reset()

	return &c
}

// maxMovesSinglePiece is the number of max possible moves of a single piece (queen at e4, 27 Moves) with added padding
// for cases I missed :)
const maxMovesSinglePiece = 30

// ParseMove parses the move in algebraic notation (see [ParseMove]) and takes into account the current
// position to account for special moves
func (p Position) ParseMove(in string) (Move, error) {
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
		if legalMove.From == m.From && legalMove.To == m.To {
			return legalMove, nil
		}
	}

	return Move{}, fmt.Errorf("given move is not a legal move")
}
