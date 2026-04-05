package chess

// calculateAttackMaps returns [squareLookup]s of all attacked squares
func calculateAttackMaps(all, king, queens, rooks, bishops, knights, pawns bitBoard) (attacksTo squareLookup[bitBoard], attacksFrom squareLookup[bitBoard]) {
	// we want to know all attacked pieces, so we treat all as enemy pieces for the
	// magic bitboard lookup
	same := bitBoard(0)
	opposing := all

	// storeAttacks stores all attacked targets in the attackFrom and inverts them for the attacksTo
	storeAttacks := func(from, targets bitBoard) {
		attacksFrom.set(king, targets)

		for to := range targets.Ones() {
			attacksTo.set(to, attacksTo.get(to)|from)
		}
	}

	pawns.Each(func(b bitBoard) {
		storeAttacks(b, opponentPawnAttacks(b))
	})
	knights.Each(func(b bitBoard) {
		storeAttacks(b, knightMoves(b, same, opposing))
	})
	bishops.Each(func(b bitBoard) {
		storeAttacks(b, bishopMoves(b, same, opposing))
	})
	rooks.Each(func(b bitBoard) {
		storeAttacks(b, rookMoves(b, same, opposing))
	})
	queens.Each(func(b bitBoard) {
		storeAttacks(b, queenMoves(b, same, opposing))
	})

	storeAttacks(king, kingMoves(king))

	return attacksTo, attacksFrom
}

// calculateSlidingKingAttacks returns a list of Bitboards that contain 1s everywhere that a piece must stand
// to prevent a check from a sliding piece. This is needed to detect pinned pieces as well as force blocking
// a check.
func (p *Position) calculateSlidingKingAttacks(king, opponentQueens, opponentRooks, opponentBishops bitBoard) {

	// we simulate the king as a queen, and check the rays individually for the correct pieces
	// for that we don't use any of our pieces and join all sliders. If we hit a slider and the slider can move in
	// that direction, then the ray is an (maybe xray) sliding attack ray
	ours := bitBoard(0)
	theirSliders := opponentQueens | opponentRooks | opponentBishops

	allRays := queenMoves(king, ours, theirSliders)

	orthogonalKingRays := [4]bitBoard{
		northRays.get(king) & allRays,
		eastRays.get(king) & allRays,
		southRays.get(king) & allRays,
		westRays.get(king) & allRays,
	}

	diagonalKingRays := [4]bitBoard{
		northEastRays.get(king) & allRays,
		southEastRays.get(king) & allRays,
		northWestRays.get(king) & allRays,
		southWestRays.get(king) & allRays,
	}

	checkRay := func(piece, ray bitBoard) {
		if piece&ray == 0 {
			return
		}

		// the ray can be inversed and treated as if coming from this piece

		pinRay := ray &^ piece

		p.xRayKingAttacks = append(p.xRayKingAttacks, attackRay{
			from: piece,
			ray:  pinRay,
		})
	}

	opponentOrthoSliders := opponentRooks | opponentQueens
	opponentDiagSliders := opponentBishops | opponentQueens

	for rook := range opponentOrthoSliders.Ones() {
		for _, ray := range orthogonalKingRays {
			checkRay(rook, ray)
		}
	}

	for bishop := range opponentDiagSliders.Ones() {
		for _, ray := range diagonalKingRays {
			checkRay(bishop, ray)
		}
	}
}

// calculateCheckSquares computes the squares that would check the opponent. Useful for marking a move
// as a checking move, which is useful for move ordering in the engine.
func (p *Position) calculateCheckSquares(opponentKing, ours, theirs bitBoard) {
	// king as opponent pawn computes our pawn checking squares:
	p.pawnCheckSquares = opponentPawnAttacks(opponentKing)

	p.knightCheckSquares = knightMoves(opponentKing, ours, theirs)

	p.bishopCheckSquares = bishopMoves(opponentKing, ours, theirs)
	p.rookCheckSquares = rookMoves(opponentKing, ours, theirs)
}
