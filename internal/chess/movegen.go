package chess

func bitBoardForEachGenMove(f func(bb, same, opposing BitBoard), b, same, opposing BitBoard) {
	for i := range BitBoard(64) {
		if b&(1<<i) == 0 {
			continue
		}

		f(1<<i, same, opposing)
	}
}

func (b *Board) GenerateMoves() []Move {
	// TODO: 50 move rule, draw by material
	b.PossibleMoves = make([]Move, 0, maxMoveCount)

	if b.PlayerInTurn == White {
		same := b.whitePieces()
		opposing := b.blackPieces()

		opponentAttacks := b.blackAttackedSquares()

		b.generateKingMoves(b.whiteKing, same, opposing, opponentAttacks, b.PlayerInTurn)
		b.whitePawns.Each(b.generateWhitePawnMoves)
		bitBoardForEachGenMove(b.generateKnightMoves, b.whiteKnights, same, opposing)

		bitBoardForEachGenMove(b.generateBishopMoves, b.whiteBishops, same, opposing)
		bitBoardForEachGenMove(b.generateRookMoves, b.whiteRooks, same, opposing)
		bitBoardForEachGenMove(b.generateQueenMoves, b.whiteQueens, same, opposing)
	} else {
		same := b.blackPieces()
		opposing := b.whitePieces()

		opponentAttacks := b.whiteAttackedSquares()

		b.generateKingMoves(b.blackKing, same, opposing, opponentAttacks, b.PlayerInTurn)
		b.blackPawns.Each(b.generateBlackPawnMoves)
		bitBoardForEachGenMove(b.generateKnightMoves, b.blackKnights, same, opposing)

		bitBoardForEachGenMove(b.generateBishopMoves, b.blackBishops, same, opposing)
		bitBoardForEachGenMove(b.generateRookMoves, b.blackRooks, same, opposing)
		bitBoardForEachGenMove(b.generateQueenMoves, b.blackQueens, same, opposing)
	}

	return b.PossibleMoves
}

func (b *Board) generatePawnMoves(from, pushed, doublePushed, doublePushRank, promoteRank, opposing BitBoard) {
	all := b.allPieces()

	if from&doublePushRank != 0 && all&pushed == 0 && all&doublePushed == 0 {
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

	if all&pushed == 0 {
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

// bitboard to check when trying to castle
var (
	whiteCastleKing            = BitBoard(MustParseSquare("f1") | MustParseSquare("g1"))
	whiteCastleKingKingTarget  = MustParseSquare("g1")
	whiteCastleKingRookTarget  = MustParseSquare("f1")
	whiteCastleQueen           = BitBoard(MustParseSquare("b1") | MustParseSquare("c1") | MustParseSquare("d1"))
	whiteCastleQueenKingTarget = MustParseSquare("c1")
	whiteCastleQueenRookTarget = MustParseSquare("d1")

	blackCastleQueen           = BitBoard(MustParseSquare("b8") | MustParseSquare("c8") | MustParseSquare("d8"))
	blackCastleKingKingTarget  = MustParseSquare("g8")
	blackCastleKingRookTarget  = MustParseSquare("f8")
	blackCastleKing            = BitBoard(MustParseSquare("f8") | MustParseSquare("g8"))
	blackCastleQueenKingTarget = MustParseSquare("c8")
	blackCastleQueenRookTarget = MustParseSquare("d8")
)

func (b *Board) generateKingMoves(bb, same, opposing, opponentAttacks BitBoard, playerInTurn Piece) {
	// TODO: castling
	if bb.Count() != 1 {
		panic("expected 1 king")
	}

	for t := range kingMoves(bb).Ones() {
		if t == 0 {
			// wrapped around
			continue
		}

		if opponentAttacks&t != 0 {
			continue
		}

		if t&same != 0 {
			// occupied
			continue
		}

		s := NoSpecial

		if t&opposing != 0 {
			s |= Captures
		}

		b.PossibleMoves = append(b.PossibleMoves, Move{
			From:    Square(bb),
			To:      Square(t),
			Special: s,
		})
	}

	if b.Castling == NoCastling {
		return
	}

	if opponentAttacks&bb != 0 {
		// king in check, no castling allowed
		return
	}

	all := same | opposing

	// white O-O
	if playerInTurn == White && b.Castling.Has(CastleWhiteKing) && all&whiteCastleKing == 0 && opponentAttacks&whiteCastleKing == 0 {
		b.PossibleMoves = append(b.PossibleMoves, Move{
			From:    Square(bb),
			To:      whiteCastleKingKingTarget,
			Special: CastleShort,
		})
	}

	// white O-O-O
	if playerInTurn == White && b.Castling.Has(CastleWhiteQueen) && all&whiteCastleQueen == 0 && opponentAttacks&whiteCastleQueen == 0 {
		b.PossibleMoves = append(b.PossibleMoves, Move{
			From:    Square(bb),
			To:      whiteCastleQueenKingTarget,
			Special: CastleLong,
		})
	}

	// black O-O
	if playerInTurn == Black && b.Castling.Has(CastleBlackKing) && all&blackCastleKing == 0 && opponentAttacks&blackCastleKing == 0 {
		b.PossibleMoves = append(b.PossibleMoves, Move{
			From:    Square(bb),
			To:      blackCastleKingKingTarget,
			Special: CastleShort,
		})
	}

	// black O-O-O
	if playerInTurn == Black && b.Castling.Has(CastleBlackQueen) && all&blackCastleQueen == 0 && opponentAttacks&blackCastleQueen == 0 {
		b.PossibleMoves = append(b.PossibleMoves, Move{
			From:    Square(bb),
			To:      blackCastleQueenKingTarget,
			Special: CastleLong,
		})
	}
}

func (b *Board) generateKnightMoves(bb, same, opposing BitBoard) {
	for t := range knightMoves(bb).Ones() {
		if t == 0 {
			// wrapped around
			continue
		}

		if t&same != 0 {
			// occupied
			continue
		}

		s := NoSpecial

		if t&opposing != 0 {
			s |= Captures
		}

		// TODO: filter out moves with check and blocked moves
		b.PossibleMoves = append(b.PossibleMoves, Move{
			From:    Square(bb),
			To:      Square(t),
			Special: s,
		})
	}
}

func (b *Board) generateRookMoves(rooks, same, opposing BitBoard) {
	for rook := range rooks.Ones() {
		targets := rookMoves(rook, same, opposing)

		for t := range targets.Ones() {
			s := NoSpecial

			if t&opposing != 0 {
				s |= Captures
			}

			b.PossibleMoves = append(b.PossibleMoves, Move{
				From:    Square(rook),
				To:      Square(t),
				Special: s,
			})
		}
	}
}

func (b *Board) generateQueenMoves(queens, same, opposing BitBoard) {
	for queen := range queens.Ones() {
		targets := queenMoves(queen, same, opposing)

		for t := range targets.Ones() {
			s := NoSpecial

			if t&opposing != 0 {
				s |= Captures
			}

			b.PossibleMoves = append(b.PossibleMoves, Move{
				From:    Square(queen),
				To:      Square(t),
				Special: s,
			})
		}
	}
}

func (b *Board) generateBishopMoves(bishop, same, opposing BitBoard) {
	for bishop := range bishop.Ones() {
		targets := bishopMoves(bishop, same, opposing)

		for t := range targets.Ones() {
			s := NoSpecial

			if t&opposing != 0 {
				s |= Captures
			}

			b.PossibleMoves = append(b.PossibleMoves, Move{
				From:    Square(bishop),
				To:      Square(t),
				Special: s,
			})
		}
	}
}
