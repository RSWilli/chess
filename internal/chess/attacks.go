package chess

// whiteAttackedSquares returns a [BitBoard] where all ones are a square that is attacked by white
func (b *Board) whiteAttackedSquares() BitBoard {
	// we want to know all attacked pieces, so we treat all as enemy pieces for the
	// magic bitboard lookup
	same := BitBoard(0)
	opposing := b.allPieces()

	attacked := BitBoard(0)

	b.whitePawns.Each(func(b BitBoard) {
		attacked |= whitePawnAttacks(b)
	})
	b.whiteKnights.Each(func(b BitBoard) {
		attacked |= knightMoves(b)
	})
	b.whiteBishops.Each(func(b BitBoard) {
		attacked |= bishopMoves(b, same, opposing)
	})
	b.whiteRooks.Each(func(b BitBoard) {
		attacked |= rookMoves(b, same, opposing)
	})
	b.whiteQueens.Each(func(b BitBoard) {
		attacked |= queenMoves(b, same, opposing)
	})

	attacked |= kingMoves(b.whiteKing)

	return attacked
}

// blackAttackedSquares returns a [BitBoard] where all ones are a square that is attacked by black
func (b *Board) blackAttackedSquares() BitBoard {
	// we want to know all attacked pieces, so we treat all as enemy pieces for the
	// magic bitboard lookup
	same := BitBoard(0)
	opposing := b.allPieces()

	attacked := BitBoard(0)

	b.blackPawns.Each(func(b BitBoard) {
		attacked |= blackPawnAttacks(b)
	})
	b.blackKnights.Each(func(b BitBoard) {
		attacked |= knightMoves(b)
	})
	b.blackBishops.Each(func(b BitBoard) {
		attacked |= bishopMoves(b, same, opposing)
	})
	b.blackRooks.Each(func(b BitBoard) {
		attacked |= rookMoves(b, same, opposing)
	})
	b.blackQueens.Each(func(b BitBoard) {
		attacked |= queenMoves(b, same, opposing)
	})

	attacked |= kingMoves(b.blackKing)

	return attacked
}
