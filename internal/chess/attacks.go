package chess

func calculateAttackMaps(all, king, queens, rooks, bishops, knights, pawns BitBoard, opponentPlayer Piece) (attacksTo squareLookup[BitBoard], attacksFrom squareLookup[BitBoard]) {
	// we want to know all attacked pieces, so we treat all as enemy pieces for the
	// magic bitboard lookup
	same := BitBoard(0)
	opposing := all

	if opponentPlayer == White {
		pawns.Each(func(b BitBoard) {
			attacksFrom.set(b, whitePawnAttacks(b))
		})
	} else {
		pawns.Each(func(b BitBoard) {
			attacksFrom.set(b, blackPawnAttacks(b))
		})
	}
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

// calculateSlidingKingAttacks returns a list of Bitboards that contain 1s everywhere that a piece must stand
// to prevent a check from a sliding piece. This is needed to detect pinned pieces as well as force blocking
// a check.
func (p *Game) calculateSlidingKingAttacks(king, opponentQueens, opponentRooks, opponentBishops BitBoard) {
	currentRay := 0

	// we simulate the king as a queen, and check the rays individually for the correct pieces
	// for that we don't use any of our pieces and join all sliders. If we hit a slider and the slider can move in
	// that direction, then the ray is an (maybe xray) sliding attack ray
	ours := BitBoard(0)
	theirSliders := opponentQueens | opponentRooks | opponentBishops

	allRays := queenMoves(king, ours, theirSliders)

	orthogonalKingRays := [4]BitBoard{
		northRays.get(king) & allRays,
		eastRays.get(king) & allRays,
		southRays.get(king) & allRays,
		westRays.get(king) & allRays,
	}

	diagonalKingRays := [4]BitBoard{
		northEastRays.get(king) & allRays,
		southEastRays.get(king) & allRays,
		northWestRays.get(king) & allRays,
		southWestRays.get(king) & allRays,
	}

	checkRay := func(piece, ray BitBoard) {
		if piece&ray == 0 {
			return
		}

		// the ray can be inversed and treated as if coming from this piece

		pinRay := ray &^ piece

		p.xRayKingAttacks[currentRay] = attackRay{
			from: piece,
			ray:  pinRay,
		}

		currentRay++
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
