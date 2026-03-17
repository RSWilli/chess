package chess

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const DefaultFen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// maxMovesFen is the position with the most number of legal moves for white
// const maxMovesFen = "R6R/3Q4/1Q4Q1/4Q3/2Q4Q/Q4Q2/pp1Q4/kBNN1KB1 w - - 0 1"
// const maxMoveCount = 218

// maxMoveCount adds a bit of padding to the maximum moves from the FEN
//
// This is used to initialize the array backing the moves slice
const maxMoveCount = 256

var fenPieceTranslation = map[rune]Piece{
	'r': BlackRook,
	'n': BlackKnight,
	'b': BlackBishop,
	'q': BlackQueen,
	'k': BlackKing,
	'p': BlackPawn,

	'R': WhiteRook,
	'N': WhiteKnight,
	'B': WhiteBishop,
	'Q': WhiteQueen,
	'K': WhiteKing,
	'P': WhitePawn,
}

var fenEmptyTranslation = map[rune]int{
	'1': 1,
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
}

var fenColorTranslation = map[string]Piece{
	"b": Black,
	"w": White,
}

const fenRankSeparator = "/"

var fenCastlingAbilityTranslation = map[rune]CastlingAbility{
	'K': CastleWhiteKing,
	'Q': CastleWhiteQueen,
	'k': CastleBlackKing,
	'q': CastleBlackQueen,
}

var ErrMalformedFEN = errors.New("given FEN is malformed")

// NewPositionFromFEN parses the given FEN string as defined in https://www.chessprogramming.org/Forsyth-Edwards_Notation
//
// Additionally a list of moves can be passed, which are performed on the position before being returned
func NewPositionFromFEN(fen string, moves []string) (*Position, error) {
	parts := strings.Split(fen, " ")

	if len(parts) != 6 {
		return nil, fmt.Errorf("%w: expected FEN with 6 parts", ErrMalformedFEN)
	}

	p := &Position{}

	ranks := strings.Split(parts[0], fenRankSeparator)

	if len(ranks) != 8 {
		return nil, fmt.Errorf("%w: expected 8 ranks", ErrMalformedFEN)
	}

	i := 0
	for _, rank := range ranks {
		rowLength := 8
		for _, c := range rank {
			piece, ok := fenPieceTranslation[c]

			if ok {
				p.set(NewSquareFromIndex(i), piece)
				i++
				rowLength--

				continue
			}

			emptySquares, ok := fenEmptyTranslation[c]

			if ok && emptySquares <= rowLength {
				i += emptySquares
				continue
			}

			if ok {
				return nil, fmt.Errorf("%w: can't skip %d squares, only %d left in row", ErrMalformedFEN, emptySquares, rowLength)
			}

			return nil, fmt.Errorf("%w: unexpected FEN char %c", ErrMalformedFEN, c)
		}
	}

	player, ok := fenColorTranslation[parts[1]]

	if !ok {
		return nil, fmt.Errorf("%w: expected a color, got %s", ErrMalformedFEN, parts[1])
	}

	p.PlayerInTurn = player

	// castling
	if parts[2] == "-" {
		p.castling = NoCastling
	} else if len(parts[2]) <= 4 {
		for _, c := range parts[2] {
			castling, ok := fenCastlingAbilityTranslation[c]

			if !ok {
				return nil, fmt.Errorf("%w: invalid castling ability %c", ErrMalformedFEN, c)
			}

			p.castling |= castling
		}
	} else {
		return nil, fmt.Errorf("%w: expected max 4 chars for castling ability, or '-'", ErrMalformedFEN)
	}

	// en passant
	if parts[3] == "-" {
		p.enPassantTarget = InvalidSquare
	} else {
		tile, err := ParseSquare(parts[3])

		if err != nil {
			return nil, fmt.Errorf("%w: failed to parse en passant target square %w", ErrMalformedFEN, err)
		}

		if tile.Rank()+1 != 3 && tile.Rank()+1 != 6 {
			return nil, fmt.Errorf("%w: en passant target square must be either rank 3 or 6", ErrMalformedFEN)
		}

		p.enPassantTarget = tile
	}

	// halfmove
	halfMoves, err := strconv.ParseUint(parts[4], 10, 8)

	if err != nil {
		return nil, fmt.Errorf("%w: half move counter invp.Dalid: %s", ErrMalformedFEN, parts[4])
	}

	p.HalfmoveClock = uint8(halfMoves)

	// fullmove
	fullMoves, err := strconv.ParseUint(parts[5], 10, 32)

	if err != nil {
		return nil, fmt.Errorf("%w: full move counter invalid: %s", ErrMalformedFEN, parts[5])
	}

	if fullMoves == 0 {
		return nil, fmt.Errorf("%w: full move counter must start at 1", ErrMalformedFEN)
	}

	p.Moves = int(fullMoves)

	p.hashFull()

	p.computeAll()

	for _, m := range moves {
		move, err := p.ParseMove(m)

		if err != nil {
			return nil, fmt.Errorf("could not parse move: %v", err)
		}

		p.DoMove(move)
	}

	return p, nil
}

// FEN returns a string version of the board according to https://www.chessprogramming.org/Forsyth-Edwards_Notation
func (p *Position) FEN() string {
	var s strings.Builder
	empties := 0

	// board:	if b.
	for row := range 8 {
		for col := range 8 {
			p := p.index(row*8 + col)

			if p == Empty {
				empties++
				continue
			} else if empties > 0 {
				fmt.Fprintf(&s, "%d", empties)
				empties = 0
			}

			uncoloredPiece := p &^ (Black | White)
			color := p & (Black | White)

			var piece rune

			switch uncoloredPiece {
			case Bishop:
				piece = 'b'
			case King:
				piece = 'k'
			case Knight:
				piece = 'n'
			case Pawn:
				piece = 'p'
			case Queen:
				piece = 'q'
			case Rook:
				piece = 'r'
			default:
				panic(fmt.Sprintf("unexpected chess.Piece: %#v", p))
			}

			if color == White {
				// map to upper case
				piece += 'A' - 'a'
			}

			s.WriteRune(piece)
		}

		if empties > 0 {
			fmt.Fprintf(&s, "%d", empties)
		}

		if row != 7 {
			s.WriteRune('/')
		}
		empties = 0
	}

	s.WriteRune(' ')

	switch p.PlayerInTurn {
	case Black:
		s.WriteRune('b')
	case White:
		s.WriteRune('w')
	default:
		panic(fmt.Sprintf("unexpected player in turn: %#v", p.PlayerInTurn))
	}

	s.WriteRune(' ')

	if p.castling.Has(CastleWhiteKing) {
		s.WriteRune('K')
	}
	if p.castling.Has(CastleWhiteQueen) {
		s.WriteRune('Q')
	}
	if p.castling.Has(CastleBlackKing) {
		s.WriteRune('k')
	}
	if p.castling.Has(CastleBlackQueen) {
		s.WriteRune('q')
	}
	if p.castling == NoCastling {
		s.WriteRune('-')
	}

	s.WriteRune(' ')

	if p.enPassantTarget != InvalidSquare {
		s.WriteString(p.enPassantTarget.String())
	} else {
		s.WriteRune('-')
	}

	s.WriteRune(' ')

	fmt.Fprintf(&s, "%d", p.HalfmoveClock)

	s.WriteRune(' ')

	fmt.Fprintf(&s, "%d", p.Moves)

	return s.String()
}
