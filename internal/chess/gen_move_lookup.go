//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"math/rand/v2"
	"os"
	"text/template"

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

// rookMoveTargets outputs the possible rook move targets given the opposite color occupies the given squares
//
// This means that the moves will end on top of the first blocking piece, signifying a capture move
func rookMoveTargets(sq, occupied chess.BitBoard) chess.BitBoard {
	m := chess.BitBoard(0)

	for c := sq.Right(); c != 0; c = c.Right() {
		m |= c

		if c&occupied != 0 {
			break
		}
	}
	for c := sq.Left(); c != 0; c = c.Left() {
		m |= c

		if c&occupied != 0 {
			break
		}
	}
	for c := sq.Up(); c != 0; c = c.Up() {
		m |= c

		if c&occupied != 0 {
			break
		}
	}
	for c := sq.Down(); c != 0; c = c.Down() {
		m |= c

		if c&occupied != 0 {
			break
		}
	}

	return m
}

// bishopMoveTargets outputs the possible rook move targets given the opposite color occupies the given squares
//
// This means that the moves will end on top of the first blocking piece, signifying a capture move
func bishopMoveTargets(sq, occupied chess.BitBoard) chess.BitBoard {
	m := chess.BitBoard(0)

	for c := sq.Up().Right(); c != 0; c = c.Up().Right() {
		m |= c

		if c&occupied != 0 {
			break
		}
	}
	for c := sq.Up().Left(); c != 0; c = c.Up().Left() {
		m |= c

		if c&occupied != 0 {
			break
		}
	}
	for c := sq.Down().Right(); c != 0; c = c.Down().Right() {
		m |= c

		if c&occupied != 0 {
			break
		}
	}
	for c := sq.Down().Left(); c != 0; c = c.Down().Left() {
		m |= c

		if c&occupied != 0 {
			break
		}
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
func findMagicNumber(occupations []chess.BitBoard, bits chess.BitBoard) chess.BitBoard {
nextmagic:
	for {
		magic := newRandom()

		// track the taken indices, we don't want collisions
		taken := make([]bool, 1<<bits)

		for _, occupancy := range occupations {
			index := hash(occupancy, magic, bits)

			if taken[index] {
				continue nextmagic // hash collision
			}

			taken[index] = true
		}

		return magic
	}
}

func hash(b, magic, bits chess.BitBoard) int {
	return int((b * magic) >> (64 - bits))
}

// newRandom returns a random uint64 for the multiplication part of the hash
func newRandom() chess.BitBoard {
	// https://www.chessprogramming.org/Looking_for_Magics#Feeding_in_Randoms
	// a good approach is to have random numbers with very few nonzero bits

	a, b, c := rand.Uint64(), rand.Uint64(), rand.Uint64()

	return chess.BitBoard(a & b & c)
}

var srcTemplate = `// this file is generated by gen_move_lookup.go
package chess

var rookMoveMasks = squareLookup[BitBoard]{
	store: [64]BitBoard{
		{{ range $m := .Rook.MoveMasks }}
		{{- printf "%#x," $m }}
		{{ end -}}
	},
}

var rookMoveTargets = squareLookup[hashTable]{
	store: [64]hashTable{
		{{ range $m := .Rook.Magics -}}
		{
			magic: {{ printf "%#x" $m.Magic }},
			bits: {{ printf "%#x" $m.Bits }},
			data:  []BitBoard{
				{{ range $v := $m.Values }}
				{{- printf "%d: %#x," $v.Key $v.Value }}
				{{ end -}}
			},
		},
		{{ end }}
	},
}

var bishopMoveMasks = squareLookup[BitBoard]{
	store: [64]BitBoard{
		{{ range $m := .Bishop.MoveMasks }}
		{{- printf "%#x," $m }}
		{{ end -}}
	},
}

var bishopMoveTargets = squareLookup[hashTable]{
	store: [64]hashTable{
		{{ range $m := .Bishop.Magics -}}
		{
			magic: {{ printf "%#x" $m.Magic }},
			bits: {{ printf "%#x" $m.Bits }},
			data:  []BitBoard{
				{{ range $v := $m.Values }}
				{{- printf "%d: %#x," $v.Key $v.Value }}
				{{ end -}}
			},
		},
		{{ end }}
	},
}
`

type entry struct {
	Key   int
	Value uint64
}

type magic struct {
	Magic  uint64
	Bits   uint64
	Values []entry
}

type pieceData struct {
	MoveMasks []uint64
	Magics    []magic
}

type data struct {
	Rook   pieceData
	Bishop pieceData
}

func main() {
	var buf bytes.Buffer

	tpl := template.Must(template.New("tpl").Parse(srcTemplate))

	data := data{}

	for i := range 64 {
		fmt.Println("lookup for magic ", i)
		sq := chess.BitBoard(1 << i)

		rookMask := rookMoveMask(sq)
		bishopMask := bishopMoveMask(sq)

		rookBits := rookBits[i]
		bishopBits := bishopBits[i]

		rookOcc := allPossibleOccupationsForSquare(rookMask)
		bishopOcc := allPossibleOccupationsForSquare(bishopMask)

		rookMagic := findMagicNumber(rookOcc, rookBits)
		bishopMagic := findMagicNumber(bishopOcc, bishopBits)

		data.Rook.MoveMasks = append(data.Rook.MoveMasks, uint64(rookMask))
		data.Bishop.MoveMasks = append(data.Bishop.MoveMasks, uint64(bishopMask))

		rookMagicData := magic{
			Magic: uint64(rookMagic),
			Bits:  uint64(rookBits),
		}

		bishopMagicData := magic{
			Magic: uint64(bishopMagic),
			Bits:  uint64(bishopBits),
		}

		for _, occ := range rookOcc {
			rookMagicData.Values = append(rookMagicData.Values, entry{
				Key:   hash(occ, rookMagic, rookBits),
				Value: uint64(rookMoveTargets(sq, occ)),
			})
		}

		for _, occ := range bishopOcc {
			bishopMagicData.Values = append(bishopMagicData.Values, entry{
				Key:   hash(occ, bishopMagic, bishopBits),
				Value: uint64(bishopMoveTargets(sq, occ)),
			})
		}

		data.Rook.Magics = append(data.Rook.Magics, rookMagicData)
		data.Bishop.Magics = append(data.Bishop.Magics, bishopMagicData)
	}

	err := tpl.Execute(&buf, data)

	out, err := format.Source(buf.Bytes())

	if err != nil {
		fmt.Println("failed to format generated source:", err)
		os.Exit(1)
	}

	f, err := os.Create("moves_lookup.gen.go")

	if err != nil {
		fmt.Println("failed to open file:", err)
		os.Exit(1)
	}
	defer f.Close()

	_, err = f.Write(out)

	if err != nil {
		fmt.Println("failed to write file:", err)
		os.Exit(1)
	}
}
