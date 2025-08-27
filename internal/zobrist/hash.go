// package zobrist implements zobrist hashing https://www.chessprogramming.org/Zobrist_Hashing for hashing the contents of a chessboard
package zobrist

import (
	"math/rand/v2"
)

var gen = rand.NewChaCha8([32]byte{1, 4, 0, 8, 1, 9, 9, 7})

var randoms []uint64

func init() {
	for range lastOffset {
		randoms = append(randoms, gen.Uint64())
	}

	Default = NewBoard(
		BlackRook, BlackKnight, BlackBishop, BlackQueen, BlackKing, BlackBishop, BlackKnight, BlackRook,
		BlackPawn, BlackPawn, BlackPawn, BlackPawn, BlackPawn, BlackPawn, BlackPawn, BlackPawn,
		-1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1,
		WhitePawn, WhitePawn, WhitePawn, WhitePawn, WhitePawn, WhitePawn, WhitePawn, WhitePawn,
		WhiteRook, WhiteKnight, WhiteBishop, WhiteQueen, WhiteKing, WhiteBishop, WhiteKnight, WhiteRook,
	) ^ WhiteCastleKing ^ BlackCastleKing ^ WhiteCastleQueen ^ BlackCastleQueen
}

// The offsets for the random hash values. Piece types are spaced out in a way that there is enough space for the square offsets between them.
const (
	WhiteKing = iota * 64
	BlackKing
	WhiteQueen
	BlackQueen
	WhiteBishop
	BlackBishop
	WhiteKnight
	BlackKnight
	WhiteRook
	BlackRook
	WhitePawn
	BlackPawn

	// piecesOffsetEnd is the first index that can be used for the other numbers, as until this index is used by [BlackPawn]
	piecesOffsetEnd
)

const (
	BlackToMove = piecesOffsetEnd + iota

	WhiteCastleKing
	BlackCastleKing
	WhiteCastleQueen
	BlackCastleQueen

	EnPassantAFile
	EnPassantBFile
	EnPassantCFile
	EnPassantDFile
	EnPassantEFile
	EnPassantFFile
	EnPassantGFile
	EnPassantHFile

	lastOffset
)

// SwitchSide can be applied when the player changes, useful for incremental hash updates
const SwitchSide = BlackToMove

// Hash is a 64-bit zobrist hash implementing zobrist hashing https://www.chessprogramming.org/Zobrist_Hashing
type Hash uint64

// Update updates the hash and returns the new one. This is used for incremental updates, so the given offsets will toggle on and also off
func (h Hash) Update(offsets ...int) Hash {
	for _, off := range offsets {
		h ^= Hash(randoms[off])
	}

	return h
}

// NewBoard creates a new hash from the given pieces. pieces is expected to be a length 64 board containing -1 for empty pieces and
// the offsets into the hash table from this package. The square offsets are calculated automatically.
func NewBoard(pieces ...int) Hash {
	var h Hash

	if len(pieces) != 64 {
		panic("expected 64 pieces exactly")
	}

	for sq, piece := range pieces {
		if piece < 0 {
			continue
		}
		h ^= Hash(randoms[piece+sq])
	}

	return h
}

// Default is the hash for the starting position
var Default Hash
