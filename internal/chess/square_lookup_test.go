package chess

import (
	"math"
	"testing"
)

func BenchmarkStdMapSquareLookup(b *testing.B) {
	stdmap := make(map[bitBoard]bitBoard, 64)
	all := bitBoard(math.MaxUint64)

	var squares []bitBoard

	for sq := range all.Ones() {
		squares = append(squares, sq)

		stdmap[sq] = sq
	}

	for b.Loop() {
		for _, sq := range squares {
			_ = stdmap[sq]
		}
	}
}

func BenchmarkSpecialSquareLookup(b *testing.B) {
	var squareLookup squareLookup[bitBoard]

	all := bitBoard(math.MaxUint64)

	var squares []bitBoard

	for sq := range all.Ones() {
		squares = append(squares, sq)

		squareLookup.set(sq, sq)
	}

	for b.Loop() {
		for _, sq := range squares {
			_ = squareLookup.get(sq)
		}
	}
}
