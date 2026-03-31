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

// Square constants, as untyped ints so they can be used as bitboards or squares:
const (
	a8 = 1 << iota
	b8
	c8
	d8
	e8
	f8
	g8
	h8

	a7
	b7
	c7
	d7
	e7
	f7
	g7
	h7

	a6
	b6
	c6
	d6
	e6
	f6
	g6
	h6

	a5
	b5
	c5
	d5
	e5
	f5
	g5
	h5

	a4
	b4
	c4
	d4
	e4
	f4
	g4
	h4

	a3
	b3
	c3
	d3
	e3
	f3
	g3
	h3

	a2
	b2
	c2
	d2
	e2
	f2
	g2
	h2

	a1
	b1
	c1
	d1
	e1
	f1
	g1
	h1
)
