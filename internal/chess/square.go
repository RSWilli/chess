package chess

import (
	"fmt"
	"strings"
)

type Rank uint8

const (
	Rank8 Rank = iota
	Rank7
	Rank6
	Rank5
	Rank4
	Rank3
	Rank2
	Rank1

	RankInvalid
)

type File uint8

const (
	FileA File = iota
	FileB
	FileC
	FileD
	FileE
	FileF
	FileG
	FileH

	FileInvalid
)

type Square struct {
	File File
	Rank Rank
}

var InvalidSquare = Square{FileInvalid, RankInvalid}

const files = "abcdefgh"
const ranks = "87654321"

func (t Square) String() string {
	return files[t.File:t.File+1] + ranks[t.Rank:t.Rank+1]
}

// Index returns the index of the Boards Position to look up this Tile
func (t Square) Index() int {
	return int(t.Rank)*8 + int(t.File)
}

func ParseTile(square string) (Square, error) {
	if len(square) != 2 {
		return Square{}, fmt.Errorf("invalid tile %s", square)
	}

	file := strings.IndexByte(files, square[0])
	rank := strings.IndexByte(ranks, square[1])

	if file == -1 || rank == -1 {
		return Square{}, fmt.Errorf("invalid tile %s", square)
	}

	return Square{
		Rank: Rank(rank),
		File: File(file),
	}, nil
}
