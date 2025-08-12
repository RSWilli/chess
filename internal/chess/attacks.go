package chess

func calculateAttackMaps(all, king, queens, rooks, bishops, knights, pawns BitBoard) (attacksTo squareLookup[BitBoard], attacksFrom squareLookup[BitBoard]) {
	// we want to know all attacked pieces, so we treat all as enemy pieces for the
	// magic bitboard lookup
	same := BitBoard(0)
	opposing := all

	pawns.Each(func(b BitBoard) {
		attacksFrom.set(b, whitePawnAttacks(b))
	})
	knights.Each(func(b BitBoard) {
		attacksFrom.set(b, knightMoves(b))
	})
	bishops.Each(func(b BitBoard) {
		attacksFrom.set(b, bishopMoves(b, same, opposing))
	})
	rooks.Each(func(b BitBoard) {
		attacksFrom.set(b, rookMoves(b, same, opposing))
	})
	queens.Each(func(b BitBoard) {
		attacksFrom.set(b, queenMoves(b, same, opposing))
	})

	attacksFrom.set(king, kingMoves(king))

	return transpose(attacksFrom), attacksTo
}
