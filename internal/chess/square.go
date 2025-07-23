package chess

import (
	"fmt"
	"math/bits"
	"strings"
)

// Square is a bitboard with just a single one set.
type Square bitBoard

// InvalidSquare represents an invalid square.
const InvalidSquare Square = 0

const files = "abcdefgh"
const ranks = "87654321"

// rank returns the index of the rank of the square from 8 to 1
func (t Square) rank() int {
	i := bits.Len64(uint64(t)) - 1
	return int(uint64(i) / 8)
}

// file returns the index of the file of the square from a to h
func (t Square) file() int {
	i := bits.Len64(uint64(t)) - 1
	return int(uint64(i) % 8)
}

func (t Square) Debug() string {
	return bitBoard(t).String()
}

func (t Square) String() string {
	if t == 0 {
		return "invalid square"
	}

	file := t.file()
	rank := t.rank()
	return files[file:file+1] + ranks[rank:rank+1]
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
