package chess

import (
	"encoding"
	"fmt"
	"math/bits"
	"strings"
)

// Square is a bitboard with just a single one set.
//
// it is always normalized, meaning viewed from the perspective of the white player.
type Square bitBoard

// InvalidSquare represents an invalid square.
const InvalidSquare Square = 0

const files = "abcdefgh"
const ranks = "87654321"

// index returns the index of the square starting from the top left corner of the board
func (t Square) index() int {
	return 64 - bits.Len64(uint64(t))
}

// Rank returns the index of the Rank of the square from 1 to 8
func (t Square) Rank() int {
	i := t.index()
	return int(uint64(i) / 8)
}

// File returns the index of the File of the square from a (1) to h (8)
func (t Square) File() int {
	i := t.index()
	return int(uint64(i) % 8)
}

func (t Square) Debug() string {
	return bitBoard(t).String()
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

// NewSquareFromIndex will return a [Square] from the given index row wise. Index 0 is a8, and index 63 is h1.
func NewSquareFromIndex(i int) Square {
	return Square(1 << (63 - i))
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

// all squares pre parsed
var (
	a1 = MustParseSquare("a1")
	a2 = MustParseSquare("a2")
	a3 = MustParseSquare("a3")
	a4 = MustParseSquare("a4")
	a5 = MustParseSquare("a5")
	a6 = MustParseSquare("a6")
	a7 = MustParseSquare("a7")
	a8 = MustParseSquare("a8")

	b1 = MustParseSquare("b1")
	b2 = MustParseSquare("b2")
	b3 = MustParseSquare("b3")
	b4 = MustParseSquare("b4")
	b5 = MustParseSquare("b5")
	b6 = MustParseSquare("b6")
	b7 = MustParseSquare("b7")
	b8 = MustParseSquare("b8")

	c1 = MustParseSquare("c1")
	c2 = MustParseSquare("c2")
	c3 = MustParseSquare("c3")
	c4 = MustParseSquare("c4")
	c5 = MustParseSquare("c5")
	c6 = MustParseSquare("c6")
	c7 = MustParseSquare("c7")
	c8 = MustParseSquare("c8")

	d1 = MustParseSquare("d1")
	d2 = MustParseSquare("d2")
	d3 = MustParseSquare("d3")
	d4 = MustParseSquare("d4")
	d5 = MustParseSquare("d5")
	d6 = MustParseSquare("d6")
	d7 = MustParseSquare("d7")
	d8 = MustParseSquare("d8")

	e1 = MustParseSquare("e1")
	e2 = MustParseSquare("e2")
	e3 = MustParseSquare("e3")
	e4 = MustParseSquare("e4")
	e5 = MustParseSquare("e5")
	e6 = MustParseSquare("e6")
	e7 = MustParseSquare("e7")
	e8 = MustParseSquare("e8")

	f1 = MustParseSquare("f1")
	f2 = MustParseSquare("f2")
	f3 = MustParseSquare("f3")
	f4 = MustParseSquare("f4")
	f5 = MustParseSquare("f5")
	f6 = MustParseSquare("f6")
	f7 = MustParseSquare("f7")
	f8 = MustParseSquare("f8")

	g1 = MustParseSquare("g1")
	g2 = MustParseSquare("g2")
	g3 = MustParseSquare("g3")
	g4 = MustParseSquare("g4")
	g5 = MustParseSquare("g5")
	g6 = MustParseSquare("g6")
	g7 = MustParseSquare("g7")
	g8 = MustParseSquare("g8")

	h1 = MustParseSquare("h1")
	h2 = MustParseSquare("h2")
	h3 = MustParseSquare("h3")
	h4 = MustParseSquare("h4")
	h5 = MustParseSquare("h5")
	h6 = MustParseSquare("h6")
	h7 = MustParseSquare("h7")
	h8 = MustParseSquare("h8")
)

var (
	aFile = a1 | a2 | a3 | a4 | a5 | a6 | a7 | a8
	bFile = b1 | b2 | b3 | b4 | b5 | b6 | b7 | b8
	cFile = c1 | c2 | c3 | c4 | c5 | c6 | c7 | c8
	dFile = d1 | d2 | d3 | d4 | d5 | d6 | d7 | d8
	eFile = e1 | e2 | e3 | e4 | e5 | e6 | e7 | e8
	fFile = f1 | f2 | f3 | f4 | f5 | f6 | f7 | f8
	gFile = g1 | g2 | g3 | g4 | g5 | g6 | g7 | g8
	hFile = h1 | h2 | h3 | h4 | h5 | h6 | h7 | h8
)

var (
	rank1 = a1 | b1 | c1 | d1 | e1 | f1 | g1 | h1
	rank2 = a2 | b2 | c2 | d2 | e2 | f2 | g2 | h2
	rank3 = a3 | b3 | c3 | d3 | e3 | f3 | g3 | h3
	rank4 = a4 | b4 | c4 | d4 | e4 | f4 | g4 | h4
	rank5 = a5 | b5 | c5 | d5 | e5 | f5 | g5 | h5
	rank6 = a6 | b6 | c6 | d6 | e6 | f6 | g6 | h6
	rank7 = a7 | b7 | c7 | d7 | e7 | f7 | g7 | h7
	rank8 = a8 | b8 | c8 | d8 | e8 | f8 | g8 | h8
)
