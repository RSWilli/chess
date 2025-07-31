package chess

func (b *Board) GenerateMoves() []Move {
	b.PossibleMoves = make([]Move, 0, maxMoveCount)

	if b.PlayerInTurn == White {
		same := b.whitePieces()
		opposing := b.blackPieces()

		b.generateKingMoves(b.whiteKing)
		b.whitePawns.Each(b.generateWhitePawnMoves)
		b.whiteKnights.Each(b.generateKnightMoves)
		b.whiteBishops.Each(func(bb BitBoard) {
			b.generateBishopMoves(bb, same, opposing)
		})
		b.whiteRooks.Each(func(bb BitBoard) {
			b.generateRookMoves(bb, same, opposing)
		})
		b.whiteQueens.Each(func(bb BitBoard) {
			b.generateQueenMoves(bb, same, opposing)
		})
	} else {
		same := b.blackPieces()
		opposing := b.whitePieces()

		b.generateKingMoves(b.blackKing)
		b.blackPawns.Each(b.generateBlackPawnMoves)
		b.blackKnights.Each(b.generateKnightMoves)
		b.blackBishops.Each(func(bb BitBoard) {
			b.generateBishopMoves(bb, same, opposing)
		})
		b.blackRooks.Each(func(bb BitBoard) {
			b.generateRookMoves(bb, same, opposing)
		})
		b.blackQueens.Each(func(bb BitBoard) {
			b.generateQueenMoves(bb, same, opposing)
		})
	}

	return b.PossibleMoves
}

func (b *Board) generatePawnMoves(from, pushed, doublePushed, doublePushRank, promoteRank, opposing BitBoard) {
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

	takes := []BitBoard{
		pushed.Left(),
		pushed.Right(),
	}

	for _, t := range takes {
		if opposing&t == 0 && b.EnPassantTarget != Square(t) {
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

const rank7BitBoard BitBoard = 0xff00
const rank2BitBoard BitBoard = 0xff000000000000

func (b *Board) generateWhitePawnMoves(bb BitBoard) {
	b.generatePawnMoves(bb, bb.Up(), bb.Up().Up(), rank2BitBoard, rank7BitBoard, b.blackPieces())
}

func (b *Board) generateBlackPawnMoves(bb BitBoard) {
	b.generatePawnMoves(bb, bb.Down(), bb.Down().Down(), rank7BitBoard, rank2BitBoard, b.whitePieces())
}

func (b *Board) generateKingMoves(bb BitBoard) {
	// TODO: castling
	if bb.Count() != 1 {
		panic("expected 1 king")
	}

	targets := []BitBoard{
		bb.Up(),
		bb.Down(),
		bb.Up().Left(),
		bb.Down().Left(),
		bb.Up().Right(),
		bb.Down().Right(),
		bb.Left(),
		bb.Right(),
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

func (b *Board) generateKnightMoves(bb BitBoard) {
	targets := []BitBoard{
		bb.Up().Up().Left(),
		bb.Up().Up().Right(),

		bb.Down().Down().Left(),
		bb.Down().Down().Right(),

		bb.Up().Left().Left(),
		bb.Up().Right().Right(),

		bb.Down().Left().Left(),
		bb.Down().Right().Right(),
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

func (b *Board) generateRookMoves(rooks, same, opposing BitBoard) {
	for rook := range rooks.Ones() {
		targets := rookMoves(rook, same, opposing)

		for t := range targets.Ones() {
			// TODO: filter out moves with check and blocked moves
			b.PossibleMoves = append(b.PossibleMoves, Move{
				From:    Square(rook),
				To:      Square(t),
				Special: NoSpecial, // TODO: special takes moves
			})
		}
	}
}

func (b *Board) generateQueenMoves(queens, same, opposing BitBoard) {
	for queen := range queens.Ones() {
		targets := queenMoves(queen, same, opposing)

		for t := range targets.Ones() {
			// TODO: filter out moves with check and blocked moves
			b.PossibleMoves = append(b.PossibleMoves, Move{
				From:    Square(queen),
				To:      Square(t),
				Special: NoSpecial, // TODO: special takes moves
			})
		}
	}
}

func (b *Board) generateBishopMoves(bishop, same, opposing BitBoard) {
	for bishop := range bishop.Ones() {
		targets := bishopMoves(bishop, same, opposing)

		for t := range targets.Ones() {
			// TODO: filter out moves with check and blocked moves
			b.PossibleMoves = append(b.PossibleMoves, Move{
				From:    Square(bishop),
				To:      Square(t),
				Special: NoSpecial, // TODO: special takes moves
			})
		}
	}
}
