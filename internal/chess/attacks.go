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

// calculateXRayAttacks returns a list of Bitboards that contain 1s everywhere that a piece must stand
// to prevent a check from a sliding piece. This is needed to detect pinned pieces as well as force blocking
// a check.
func calculateXRayAttacks(king, opponentQueens, opponentRooks, opponentBishops BitBoard) []attackRay {
	var attackRays []attackRay

	for queen := range opponentQueens.Ones() {
		if r, ok := rays[ray{from: king, to: queen}]; ok {
			attackRays = append(attackRays, attackRay{
				from: queen,
				ray:  r,
			})
		}
	}

	for rook := range opponentRooks.Ones() {
		if !rook.isSameFile(king) && !rook.isSameRank(king) {
			// a ray would be a bishop ray
			continue
		}

		if r, ok := rays[ray{from: king, to: rook}]; ok {
			attackRays = append(attackRays, attackRay{
				from: rook,
				ray:  r,
			})
		}
	}

	for bishop := range opponentBishops.Ones() {
		if bishop.isSameFile(king) || bishop.isSameRank(king) {
			// a ray would be a rook ray
			continue
		}

		if r, ok := rays[ray{from: king, to: bishop}]; ok {
			attackRays = append(attackRays, attackRay{
				from: bishop,
				ray:  r,
			})
		}
	}

	return attackRays
}
