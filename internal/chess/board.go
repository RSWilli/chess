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

// const testFen = "8/8/8/3Nn3/8/8/8/3K1k2 w KQkq - 0 1"

func NewBoard() *Board {
	b, err := NewBoardFromFEN(DefaultFen)
	// b, err := NewBoardFromFEN(testFen)

	if err != nil {
		panic(err)
	}

	return b
}

func (board Board) String() string {
	sb := strings.Builder{}

	for rank := range 8 {
		for file := range 8 {
			piece := board.Square(NewSquare(rank, file))

			sb.WriteRune(piece.Rune())
		}
		sb.WriteByte('\n')
	}
	switch board.PlayerInTurn {
	case White:
		sb.WriteString("White\n")
	case Black:
		sb.WriteString("Black\n")
	default:
		panic("unexpected player turn")
	}

	sb.WriteString(board.Castling.String())
	if board.EnPassantTarget != InvalidSquare {
		sb.WriteString("\nEn passant to: ")
		sb.WriteString(board.EnPassantTarget.String())
	}

	return sb.String()
}

func (board *Board) Square(sq Square) Piece {
	switch {
	case board.whitePawns.Has(sq):
		return WhitePawn
	case board.whiteKnights.Has(sq):
		return WhiteKnight
	case board.whiteBishops.Has(sq):
		return WhiteBishop
	case board.whiteRooks.Has(sq):
		return WhiteRook
	case board.whiteQueens.Has(sq):
		return WhiteQueen
	case board.whiteKing.Has(sq):
		return WhiteKing
	case board.blackPawns.Has(sq):
		return BlackPawn
	case board.blackKnights.Has(sq):
		return BlackKnight
	case board.blackBishops.Has(sq):
		return BlackBishop
	case board.blackRooks.Has(sq):
		return BlackRook
	case board.blackQueens.Has(sq):
		return BlackQueen
	case board.blackKing.Has(sq):
		return BlackKing
	default:
		return Empty
	}
}

// set sets the given piece on the given square
func (board *Board) set(sq Square, p Piece) {
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
func (board *Board) unset(sq Square) {
	board.whitePawns = board.whitePawns.Unset(sq)
	board.whiteKnights = board.whiteKnights.Unset(sq)
	board.whiteBishops = board.whiteBishops.Unset(sq)
	board.whiteRooks = board.whiteRooks.Unset(sq)
	board.whiteQueens = board.whiteQueens.Unset(sq)
	board.whiteKing = board.whiteKing.Unset(sq)

	board.blackPawns = board.blackPawns.Unset(sq)
	board.blackKnights = board.blackKnights.Unset(sq)
	board.blackBishops = board.blackBishops.Unset(sq)
	board.blackRooks = board.blackRooks.Unset(sq)
	board.blackQueens = board.blackQueens.Unset(sq)
	board.blackKing = board.blackKing.Unset(sq)
}

func (board *Board) allPieces() BitBoard {
	return board.whitePieces() | board.blackPieces()
}

func (board *Board) whitePieces() BitBoard {
	return board.whitePawns |
		board.whiteKnights |
		board.whiteBishops |
		board.whiteRooks |
		board.whiteQueens |
		board.whiteKing
}

func (board *Board) blackPieces() BitBoard {
	return board.blackPawns |
		board.blackKnights |
		board.blackBishops |
		board.blackRooks |
		board.blackQueens |
		board.blackKing
}

var a1 = MustParseSquare("a1")
var e1 = MustParseSquare("e1")
var a8 = MustParseSquare("a8")
var h1 = MustParseSquare("h1")
var e8 = MustParseSquare("e8")
var h8 = MustParseSquare("h8")

func (board *Board) DoMove(m Move) {
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
func (board *Board) IsCheck() bool {
	if board.PlayerInTurn == White {
		return board.attacksTo.get(board.whiteKing) != 0
	} else {
		return board.attacksTo.get(board.blackKing) != 0
	}
}

// IsCheckMate is meant to be called by the visualization and returns true if the current player is in checkmate
func (board *Board) IsCheckMate() bool {
	return board.IsCheck() && len(board.GenerateMoves()) == 0
}

// IsDraw is meant to be called by the visualization and returns true if the game is drewn
func (board *Board) IsDraw() bool {
	return !board.IsCheck() && len(board.GenerateMoves()) == 0
}
