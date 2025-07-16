package chess

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const DefaultFen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

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

const fenRankSeparator = '/'

var fenCastlingAbilityTranslation = map[rune]CastlingAbility{
	'K': CastleWhiteKing,
	'Q': CastleWhiteQueen,
	'k': CastleBlackKing,
	'q': CastleBlackQueen,
}

var ErrMalformedFEN = errors.New("given FEN is malformed")

// NewBoardFromFEN parses the given FEN string as defined in https://www.chessprogramming.org/Forsyth-Edwards_Notation
func NewBoardFromFEN(in string) (*Board, error) {
	parts := strings.Split(in, " ")

	if len(parts) != 6 {
		return nil, fmt.Errorf("%w: expected FEN with 6 parts", ErrMalformedFEN)
	}

	b := &Board{}

	ranks := strings.Split(parts[0], "/")

	if len(ranks) != 8 {
		return nil, fmt.Errorf("%w: expected 8 ranks", ErrMalformedFEN)
	}

	i := 0
	for _, rank := range ranks {
		rowLength := 8
		for _, c := range rank {
			piece, ok := fenPieceTranslation[c]

			if ok {
				b.pos[i] = piece
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

	b.PlayerInTurn = player

	// castling
	if parts[2] == "-" {
		b.Castling = NoCastling
	} else if len(parts[2]) <= 4 {
		for _, c := range parts[2] {
			castling, ok := fenCastlingAbilityTranslation[c]

			if !ok {
				return nil, fmt.Errorf("%w: invalid castling ability %c", ErrMalformedFEN, c)
			}

			b.Castling |= castling
		}
	} else {
		return nil, fmt.Errorf("%w: expected max 4 chars for castling ability, or '-'", ErrMalformedFEN)
	}

	// en passant
	if parts[3] == "-" {
		b.EnPassantTarget = InvalidSquare
	} else {
		tile, err := ParseTile(parts[3])

		if err != nil {
			return nil, fmt.Errorf("%w: failed to parse en passant target square %w", ErrMalformedFEN, err)
		}

		if tile.Rank != Rank3 && tile.Rank != Rank6 {
			return nil, fmt.Errorf("%w: en passant target square must be either rank 3 or 6", ErrMalformedFEN)
		}

		b.EnPassantTarget = tile
	}

	// halfmove
	halfMoves, err := strconv.ParseUint(parts[4], 10, 8)

	if err != nil {
		return nil, fmt.Errorf("%w: half move counter invalid: %s", ErrMalformedFEN, parts[4])
	}

	b.HalfmoveClock = uint8(halfMoves)

	// fullmove
	fullMoves, err := strconv.ParseUint(parts[5], 10, 32)

	if err != nil {
		return nil, fmt.Errorf("%w: full move counter invalid: %s", ErrMalformedFEN, parts[5])
	}

	if fullMoves == 0 {
		return nil, fmt.Errorf("%w: full move counter must start at 1", ErrMalformedFEN)
	}

	b.Moves = int(fullMoves)

	return b, nil
}
