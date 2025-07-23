package chess

func (b *Board) GenerateMoves() []Move {
	b.PossibleMoves = make([]Move, 0, maxMoveCount)

	if b.PlayerInTurn == White {
		b.generateKingMoves(b.whiteKing)
		b.whitePawns.each(b.generateWhitePawnMoves)
		b.whiteKnights.each(b.generateKnightMoves)
		b.whiteBishops.each(b.generateBishopMoves)
		b.whiteRooks.each(b.generateRookMoves)
		b.whiteQueens.each(b.generateQueenMoves)
	} else {
		b.generateKingMoves(b.blackKing)
		b.blackPawns.each(b.generateBlackPawnMoves)
		b.blackKnights.each(b.generateKnightMoves)
		b.blackBishops.each(b.generateBishopMoves)
		b.blackRooks.each(b.generateRookMoves)
		b.blackQueens.each(b.generateQueenMoves)
	}

	return b.PossibleMoves
}

func (b *Board) generatePawnMoves(from, pushed, doublePushed, doublePushRank, promoteRank, capturable bitBoard) {
	if from&doublePushRank != 0 && b.allPieces()&pushed == 0 && b.allPieces()&doublePushed == 0 {
		// can double push

		b.PossibleMoves = append(b.PossibleMoves, Move{
			From:    Square(from),
			To:      Square(doublePushed),
			Special: DoublePawnPush,
		})
	}

	// moves tracks the possible moves temporarily to defer
	// promotionm decision
	moves := make([]Move, 0, 3)

	if b.allPieces()&pushed == 0 {
		// can push

		moves = append(moves, Move{
			From:    Square(from),
			To:      Square(pushed),
			Special: NoSpecial,
		})
	}

	takes := []bitBoard{
		pushed.left(),
		pushed.right(),
	}

	for _, t := range takes {
		if capturable&t == 0 && b.EnPassantTarget != Square(t) {
			continue
		}

		moves = append(moves, Move{
			From:    Square(from),
			To:      Square(t),
			Special: Captures,
		})
	}

	if from&promoteRank == 0 {
		// no promotion
		b.PossibleMoves = append(b.PossibleMoves, moves...)
		return
	}

	promotions := []MoveSpecial{
		PromoteQueen,
		PromoteRook,
		PromoteBishop,
		PromoteKnight,
	}

	for _, p := range promotions {
		for _, m := range moves {
			m.Special = m.Special | p
			b.PossibleMoves = append(b.PossibleMoves, m)
		}
	}
}

func (b *Board) generateWhitePawnMoves(bb bitBoard) {
	b.generatePawnMoves(bb, bb.up(), bb.up().up(), rank2BitBoard, rank7BitBoard, b.blackPieces())
}

func (b *Board) generateBlackPawnMoves(bb bitBoard) {
	b.generatePawnMoves(bb, bb.down(), bb.down().down(), rank7BitBoard, rank2BitBoard, b.whitePieces())
}

func (b *Board) generateKingMoves(bb bitBoard) {
	// TODO: castling
	if bb.count() != 1 {
		panic("expected 1 king")
	}

	targets := []bitBoard{
		bb.up(),
		bb.down(),
		bb.up().left(),
		bb.down().left(),
		bb.up().right(),
		bb.down().right(),
		bb.left(),
		bb.right(),
	}

	for _, t := range targets {
		if t == 0 {
			// wrapped around
			continue
		}

		// TODO: filter out moves with check and blocked moves
		b.PossibleMoves = append(b.PossibleMoves, Move{
			From:    Square(bb),
			To:      Square(t),
			Special: NoSpecial, // TODO: special takes moves
		})
	}
}

func (b *Board) generateKnightMoves(bb bitBoard) {
	targets := []bitBoard{
		bb.up().up().left(),
		bb.up().up().right(),

		bb.down().down().left(),
		bb.down().down().right(),

		bb.up().left().left(),
		bb.up().right().right(),

		bb.down().left().left(),
		bb.down().right().right(),
	}

	for _, t := range targets {
		if t == 0 {
			// wrapped around
			continue
		}

		// TODO: filter out moves with check and blocked moves
		b.PossibleMoves = append(b.PossibleMoves, Move{
			From:    Square(bb),
			To:      Square(t),
			Special: NoSpecial, // TODO: special takes moves
		})
	}
}

func (b *Board) generateRookMoves(bb bitBoard) {
	targets := rookMoves(bb)

	for t := range targets.ones() {
		// TODO: filter out moves with check and blocked moves
		b.PossibleMoves = append(b.PossibleMoves, Move{
			From:    Square(bb),
			To:      Square(t),
			Special: NoSpecial, // TODO: special takes moves
		})
	}
}

func (b *Board) generateQueenMoves(bb bitBoard) {
	targets := queenMoves(bb)

	for t := range targets.ones() {
		// TODO: filter out moves with check and blocked moves
		b.PossibleMoves = append(b.PossibleMoves, Move{
			From:    Square(bb),
			To:      Square(t),
			Special: NoSpecial, // TODO: special takes moves
		})
	}
}

func (b *Board) generateBishopMoves(bb bitBoard) {
	targets := bishopMoves(bb)

	for t := range targets.ones() {
		// TODO: filter out moves with check and blocked moves
		b.PossibleMoves = append(b.PossibleMoves, Move{
			From:    Square(bb),
			To:      Square(t),
			Special: NoSpecial, // TODO: special takes moves
		})
	}
}
