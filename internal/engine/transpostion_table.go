package engine

import (
	"github.com/rswilli/chess/internal/chess"
	"github.com/rswilli/chess/internal/zobrist"
)

type transpositionTable[T any] struct {
	m map[zobrist.Hash]T
}

func newTranspostionTable[T any]() *transpositionTable[T] {
	return &transpositionTable[T]{
		m: make(map[zobrist.Hash]T),
	}
}

func (tt *transpositionTable[T]) set(pos *chess.Position, v T) {
	tt.m[pos.HashKey] = v
}

func (tt *transpositionTable[T]) get(pos *chess.Position) (v T, ok bool) {
	v, ok = tt.m[pos.HashKey]

	return v, ok
}
