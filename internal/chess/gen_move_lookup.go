//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"math/rand/v2"

	"github.com/rswilli/chess/internal/chess"
)

const afile chess.BitBoard = 0x0101010101010101
const hfile chess.BitBoard = 0x8080808080808080

const rank8 chess.BitBoard = 0xff
const rank1 chess.BitBoard = 0xff00000000000000

// rookMoveMask returns a bitmask that has all bits set where an opposing color piece
// would restrict the possible moves of a rook
//
// this means that we skip the last square in each direction, because the rook can move there
// regardless whether there is an opposing piece or not.
func rookMoveMask(sq chess.BitBoard) chess.BitBoard {
	m := chess.BitBoard(0)

	for c := sq.Right(); c != 0; c = c.Right() {
		m |= c &^ hfile
	}
	for c := sq.Left(); c != 0; c = c.Left() {
		m |= c &^ afile
	}
	for c := sq.Up(); c != 0; c = c.Up() {
		m |= c &^ rank8
	}
	for c := sq.Down(); c != 0; c = c.Down() {
		m |= c &^ rank1
	}

	return m
}

// bishopMoveMask returns a bitmask that has all bits set where an opposing color piece
// would restrict the possible moves of a bishop
//
// this means that we skip the last square in each direction, because the bishop can move there
// regardless whether there is an opposing piece or not.
func bishopMoveMask(sq chess.BitBoard) chess.BitBoard {
	m := chess.BitBoard(0)

	for c := sq.Up().Right(); c != 0; c = c.Up().Right() {
		m |= c &^ (rank8 | hfile)
	}
	for c := sq.Up().Left(); c != 0; c = c.Up().Left() {
		m |= c &^ (rank8 | afile)
	}
	for c := sq.Down().Right(); c != 0; c = c.Down().Right() {
		m |= c &^ (rank1 | hfile)
	}
	for c := sq.Down().Left(); c != 0; c = c.Down().Left() {
		m |= c &^ (rank1 | afile)
	}

	return m
}

func createOccupancyByIndex(mask chess.BitBoard, index int) chess.BitBoard {
	var res chess.BitBoard

	for m := range mask.Ones() {
		if index&1 == 1 {
			res |= m
		}
		index = index >> 1
	}

	if index != 0 {
		panic("mask was not big enough")
	}

	return res
}

func possibleOccupationsCount(mask chess.BitBoard) int {
	bits := mask.Count()

	return 1 << bits
}

func allPossibleOccupationsForSquare(mask chess.BitBoard) []chess.BitBoard {
	count := possibleOccupationsCount(mask)

	res := make([]chess.BitBoard, count)

	for i := range count {
		res[i] = createOccupancyByIndex(mask, i)
	}

	return res
}

var header = []byte(`
package chess

`)

// rookBits is the amount of bits needed for the perfect hashing of the rook moves
// source https://www.chessprogramming.org/Looking_for_Magics
var rookBits = [64]chess.BitBoard{
	12, 11, 11, 11, 11, 11, 11, 12,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	12, 11, 11, 11, 11, 11, 11, 12,
}

// bishopBits is the amount of bits needed for the perfect hashing of the bishop moves
// source https://www.chessprogramming.org/Looking_for_Magics
var bishopBits = [64]chess.BitBoard{
	6, 5, 5, 5, 5, 5, 5, 6,
	5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5,
	6, 5, 5, 5, 5, 5, 5, 6,
}

// const maxShiftBits = chess.BitBoard(12)

// findMagicNumber tries to find the best magic numbers for each of the squares
func findMagicNumber(sq chess.BitBoard, maskFn func(sq chess.BitBoard) chess.BitBoard, bits chess.BitBoard) chess.BitBoard {
	shift := 64 - bits

	occupations := allPossibleOccupationsForSquare(maskFn(sq))

nextmagic:
	for {
		magic := newRandom()

		// track the taken indices, we don't want collisions
		taken := make([]bool, 1<<bits)

		for _, occupancy := range occupations {
			index := hash(occupancy, magic, shift)

			if taken[index] {
				continue nextmagic // hash collision
			}

			taken[index] = true
		}

		return magic
	}
}

func hash(b, magic, shift chess.BitBoard) int {
	return int((b * magic) >> shift)
}

// newRandom returns a random uint64 for the multiplication part of the hash
func newRandom() chess.BitBoard {
	// https://www.chessprogramming.org/Looking_for_Magics#Feeding_in_Randoms
	// a good approach is to have random numbers with very few nonzero bits

	//lint:ignore SA4000 this is a false positive
	return chess.BitBoard(rand.Uint64() & rand.Uint64() & rand.Uint64())
}

func main() {
	var buf bytes.Buffer

	buf.Write(header)

	rookMagics := make([]chess.BitBoard, 64)
	bishopMagics := make([]chess.BitBoard, 64)

	for i := range 64 {
		sq := chess.BitBoard(1 << i)

		rookMagics[i] = findMagicNumber(sq, rookMoveMask, rookBits[i])
		bishopMagics[i] = findMagicNumber(sq, bishopMoveMask, bishopBits[i])

		fmt.Printf("rook magic: %#b\n", rookMagics[i])
		fmt.Printf("bishopmagic magic: %#b\n", bishopMagics[i])
	}
}
