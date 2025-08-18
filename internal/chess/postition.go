package chess

import (
	"strings"
)

type Position struct {
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

	// Calculated properties, will be reset after a move is done:

	// attacksTo maps each square to a [BitBoard] where all the attacking pieces are 1
	attacksTo squareLookup[BitBoard]
	// attacksFrom maps each square to a [BitBoard] where all the attacked squares are 1
	attacksFrom squareLookup[BitBoard]

	// piece unions:
	// ours is a BitBoard where all pieces of the current player are unioned
	ours BitBoard
	// ours is a BitBoard where all pieces of the current opponent are unioned
	theirs BitBoard
	// all is a BitBoard where all pieces are unioned
	all BitBoard

	possibleMoves []Move
}

func New() *Position {
	b, err := NewFromFEN(DefaultFen)
	// b, err := NewBoardFromFEN(testFen)

	if err != nil {
		panic(err)
	}

	return b
}

func (p Position) String() string {
	sb := strings.Builder{}

	for rank := range 8 {
		for file := range 8 {
			piece := p.Square(NewSquare(rank, file))

			sb.WriteRune(piece.Rune())
		}
		sb.WriteByte('\n')
	}
	switch p.PlayerInTurn {
	case White:
		sb.WriteString("White\n")
	case Black:
		sb.WriteString("Black\n")
	default:
		panic("unexpected player turn")
	}

	sb.WriteString(p.Castling.String())
	if p.EnPassantTarget != InvalidSquare {
		sb.WriteString("\nEn passant to: ")
		sb.WriteString(p.EnPassantTarget.String())
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

func (p *Position) allPieces() BitBoard {
	return p.whitePieces() | p.blackPieces()
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

func (board *Position) DoMove(m Move) {
	if m.Special.Has(CastleLong | CastleShort) {
		// Castling move:
		if board.PlayerInTurn == White && m.Special == CastleShort {
			board.unset(e1)
			board.unset(h1)

			board.set(whiteCastleKingKingTarget, WhiteKing)
			board.set(whiteCastleKingRookTarget, WhiteRook)

			board.Castling &^= CastleWhiteKing | CastleWhiteQueen
		}
		if board.PlayerInTurn == White && m.Special == CastleLong {
			board.unset(e1)
			board.unset(a1)

			board.set(whiteCastleQueenKingTarget, WhiteKing)
			board.set(whiteCastleQueenRookTarget, WhiteRook)

			board.Castling &^= CastleWhiteKing | CastleWhiteQueen
		}
		if board.PlayerInTurn == Black && m.Special == CastleShort {
			board.unset(e8)
			board.unset(h8)

			board.set(blackCastleKingKingTarget, BlackKing)
			board.set(blackCastleKingRookTarget, BlackRook)

			board.Castling &^= CastleBlackKing | CastleBlackQueen
		}
		if board.PlayerInTurn == Black && m.Special == CastleLong {
			board.unset(e8)
			board.unset(a8)

			board.set(blackCastleQueenKingTarget, BlackKing)
			board.set(blackCastleQueenRookTarget, BlackRook)

			board.Castling &^= CastleBlackKing | CastleBlackQueen
		}
	} else {
		p := board.Square(m.From)

		// clear the old square
		board.unset(m.From)
		// must clear all other boards before setting the new one
		board.unset(m.To)

		switch {
		case m.Special.Has(PromoteQueen):
			p = Queen | board.PlayerInTurn
		case m.Special.Has(PromoteRook):
			p = Rook | board.PlayerInTurn
		case m.Special.Has(PromoteKnight):
			p = Knight | board.PlayerInTurn
		case m.Special.Has(PromoteBishop):
			p = Bishop | board.PlayerInTurn
		}

		// remove the en-passant captured pawn, no need to check for the piece type since the en passant
		// square is always empty, so no other move can capture on it
		if m.Special.Has(Captures) && m.To == board.EnPassantTarget && board.PlayerInTurn == White {
			board.unset(Square(BitBoard(m.To).Down()))
		} else if m.Special.Has(Captures) && m.To == board.EnPassantTarget && board.PlayerInTurn == Black {
			board.unset(Square(BitBoard(m.To).Up()))
		}

		// save the en passant square for the move generation of the en passant moves
		if m.Special.Has(DoublePawnPush) && board.PlayerInTurn == White {
			board.EnPassantTarget = Square(BitBoard(m.From).Up())
		} else if m.Special.Has(DoublePawnPush) && board.PlayerInTurn == Black {
			board.EnPassantTarget = Square(BitBoard(m.From).Down())
		} else {
			board.EnPassantTarget = InvalidSquare
		}

		// prevent castling moves:
		switch m.From {
		case a1:
			board.Castling &^= CastleWhiteQueen
		case e1:
			board.Castling &^= CastleWhiteQueen | CastleWhiteKing
		case h1:
			board.Castling &^= CastleWhiteKing
		case a8:
			board.Castling &^= CastleBlackQueen
		case e8:
			board.Castling &^= CastleBlackQueen | CastleBlackKing
		case h8:
			board.Castling &^= CastleBlackKing
		}

		switch m.To {
		case a1:
			board.Castling &^= CastleWhiteQueen
		case h1:
			board.Castling &^= CastleWhiteKing
		case a8:
			board.Castling &^= CastleBlackQueen
		case h8:
			board.Castling &^= CastleBlackKing
		}

		board.set(m.To, p)
	}

	// Move done, reset state and recalculate:
	board.possibleMoves = nil

	if board.PlayerInTurn == White {
		board.PlayerInTurn = Black
		board.ours = board.blackPieces()
		board.theirs = board.whitePieces()
		board.all = board.ours | board.theirs

		board.attacksTo, board.attacksFrom = calculateAttackMaps(
			board.all,
			board.whiteKing,
			board.whiteQueens,
			board.whiteRooks,
			board.whiteBishops,
			board.whiteKnights,
			board.whitePawns,
		)
	} else {
		board.PlayerInTurn = White
		board.ours = board.whitePieces()
		board.theirs = board.blackPieces()
		board.all = board.ours | board.theirs

		board.attacksTo, board.attacksFrom = calculateAttackMaps(
			board.all,
			board.blackKing,
			board.blackQueens,
			board.blackRooks,
			board.blackBishops,
			board.blackKnights,
			board.blackPawns,
		)
	}
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
