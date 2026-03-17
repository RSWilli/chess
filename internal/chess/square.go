package chess

import (
	"encoding"
	"fmt"
	"math/bits"
	"strings"
)

// Square is a bitboard with just a single one set.
type Square BitBoard

// InvalidSquare represents an invalid square.
const InvalidSquare Square = 0

const files = "abcdefgh"
const ranks = "87654321"

// index returns the index of the square starting from the top left corner of the board
func (t Square) index() int {
	return bits.Len64(uint64(t)) - 1
}

// Rank returns the index of the Rank of the square from 8 to 1
func (t Square) Rank() int {
	i := t.index()
	return int(uint64(i) / 8)
}

// File returns the index of the File of the square from a to h
func (t Square) File() int {
	i := t.index()
	return int(uint64(i) % 8)
}

func (t Square) Debug() string {
	return BitBoard(t).String()
}

func (t Square) String() string {
	if t == 0 {
		return "invalid square"
	}

	file := t.File()
	rank := t.Rank()
	return files[file:file+1] + ranks[rank:rank+1]
}

func MustParseSquare(square string) Square {
	sq, err := ParseSquare(square)

	if err != nil {
		panic(err)
	}

	return sq
}

func ParseSquare(square string) (Square, error) {
	if len(square) != 2 {
		return InvalidSquare, fmt.Errorf("invalid tile %s", square)
	}

	file := strings.IndexByte(files, square[0])
	rank := strings.IndexByte(ranks, square[1])

	if file == -1 || rank == -1 {
		return InvalidSquare, fmt.Errorf("invalid tile %s", square)
	}

	return NewSquare(rank, file), nil
}

func NewSquare(rankIndex, fileIndex int) Square {
	return NewSquareFromIndex(rankIndex*8 + fileIndex)
}

func NewSquareFromIndex(i int) Square {
	return Square(1 << i)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (t *Square) UnmarshalText(text []byte) error {
	sq, err := ParseSquare(string(text))

	if err != nil {
		return err
	}

	*t = sq
	return nil
}

var _ encoding.TextUnmarshaler = new(Square)
