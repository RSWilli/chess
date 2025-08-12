package chess

func (board *Board) GenerateMoves() []Move {
	if board.possibleMoves != nil {
		return board.possibleMoves
	}
	// TODO: 50 move rule, draw by material
	board.possibleMoves = make([]Move, 0, maxMoveCount)

	if board.PlayerInTurn == White {

		board.generateKingMoves(board.whiteKing)
		board.whitePawns.Each(board.generateWhitePawnMoves)
		board.whiteKnights.Each(board.generateKnightMoves)

		board.whiteBishops.Each(board.generateBishopMoves)
		board.whiteRooks.Each(board.generateRookMoves)
		board.whiteQueens.Each(board.generateQueenMoves)
	} else {
		board.generateKingMoves(board.blackKing)
		board.blackPawns.Each(board.generateBlackPawnMoves)
		board.blackKnights.Each(board.generateKnightMoves)

		board.blackBishops.Each(board.generateBishopMoves)
		board.blackRooks.Each(board.generateRookMoves)
		board.blackQueens.Each(board.generateQueenMoves)
	}

	return board.possibleMoves
}

func (board *Board) generatePawnMoves(from, pushed, doublePushed, doublePushRank, promoteRank, opposing BitBoard) {
	all := board.allPieces()

	if from&doublePushRank != 0 && all&pushed == 0 && all&doublePushed == 0 {
		// can double push

		board.possibleMoves = append(board.possibleMoves, Move{
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
		if opposing&t == 0 && board.EnPassantTarget != Square(t) {
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
		board.possibleMoves = append(board.possibleMoves, moves...)
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
			board.possibleMoves = append(board.possibleMoves, m)
		}
	}
}

const rank7BitBoard BitBoard = 0xff00
const rank2BitBoard BitBoard = 0xff000000000000

func (board *Board) generateWhitePawnMoves(bb BitBoard) {
	board.generatePawnMoves(bb, bb.Up(), bb.Up().Up(), rank2BitBoard, rank7BitBoard, board.blackPieces())
}

func (board *Board) generateBlackPawnMoves(bb BitBoard) {
	board.generatePawnMoves(bb, bb.Down(), bb.Down().Down(), rank7BitBoard, rank2BitBoard, board.whitePieces())
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

// notAttacked returns true if none of the squares set in the given bitboard are attacked
func (board *Board) notAttacked(bb BitBoard) bool {
	for sq := range bb.Ones() {
		if board.attacksTo.get(sq) != 0 {
			return false
		}
	}
	return true
}

func (board *Board) canCastleWhiteKing() bool {
	return board.Castling.Has(CastleWhiteKing) && board.all&whiteCastleKing == 0 && board.notAttacked(whiteCastleKing)
}

func (board *Board) canCastleWhiteQueen() bool {
	return board.Castling.Has(CastleWhiteQueen) && board.all&whiteCastleQueen == 0 && board.notAttacked(whiteCastleQueen)
}
func (board *Board) canCastleBlackKing() bool {
	return board.Castling.Has(CastleBlackKing) && board.all&blackCastleKing == 0 && board.notAttacked(blackCastleKing)
}

func (board *Board) canCastleBlackQueen() bool {
	return board.Castling.Has(CastleBlackQueen) && board.all&blackCastleQueen == 0 && board.notAttacked(blackCastleQueen)
}

func (board *Board) generateKingMoves(bb BitBoard) {
	if bb.Count() != 1 {
		panic("expected 1 king")
	}

	for t := range kingMoves(bb).Ones() {
		if t == 0 {
			// wrapped around
			continue
		}

		if board.attacksTo.get(t) != 0 {
			continue
		}

		if t&board.ours != 0 {
			// occupied
			continue
		}

		s := NoSpecial

		if t&board.theirs != 0 {
			s |= Captures
		}

		board.possibleMoves = append(board.possibleMoves, Move{
			From:    Square(bb),
			To:      Square(t),
			Special: s,
		})
	}

	if board.Castling == NoCastling {
		return
	}

	if board.attacksTo.get(bb) != 0 {
		// king in check, no castling allowed
		return
	}

	// white O-O
	if board.PlayerInTurn == White && board.canCastleWhiteKing() {
		board.possibleMoves = append(board.possibleMoves, Move{
			From:    Square(bb),
			To:      whiteCastleKingKingTarget,
			Special: CastleShort,
		})
	}

	// white O-O-O
	if board.PlayerInTurn == White && board.canCastleWhiteQueen() {
		board.possibleMoves = append(board.possibleMoves, Move{
			From:    Square(bb),
			To:      whiteCastleQueenKingTarget,
			Special: CastleLong,
		})
	}

	// black O-O
	if board.PlayerInTurn == Black && board.canCastleBlackKing() {
		board.possibleMoves = append(board.possibleMoves, Move{
			From:    Square(bb),
			To:      blackCastleKingKingTarget,
			Special: CastleShort,
		})
	}

	// black O-O-O
	if board.PlayerInTurn == Black && board.canCastleBlackQueen() {
		board.possibleMoves = append(board.possibleMoves, Move{
			From:    Square(bb),
			To:      blackCastleQueenKingTarget,
			Special: CastleLong,
		})
	}
}

func (board *Board) generateKnightMoves(bb BitBoard) {
	for t := range knightMoves(bb).Ones() {
		if t == 0 {
			// wrapped around
			continue
		}

		if t&board.ours != 0 {
			// occupied
			continue
		}

		s := NoSpecial

		if t&board.theirs != 0 {
			s |= Captures
		}

		// TODO: filter out moves with check and blocked moves
		board.possibleMoves = append(board.possibleMoves, Move{
			From:    Square(bb),
			To:      Square(t),
			Special: s,
		})
	}
}

func (board *Board) generateRookMoves(rooks BitBoard) {
	for rook := range rooks.Ones() {
		targets := rookMoves(rook, board.ours, board.theirs)

		for t := range targets.Ones() {
			s := NoSpecial

			if t&board.theirs != 0 {
				s |= Captures
			}

			board.possibleMoves = append(board.possibleMoves, Move{
				From:    Square(rook),
				To:      Square(t),
				Special: s,
			})
		}
	}
}

func (board *Board) generateQueenMoves(queens BitBoard) {
	for queen := range queens.Ones() {
		targets := queenMoves(queen, board.ours, board.theirs)

		for t := range targets.Ones() {
			s := NoSpecial

			if t&board.theirs != 0 {
				s |= Captures
			}

			board.possibleMoves = append(board.possibleMoves, Move{
				From:    Square(queen),
				To:      Square(t),
				Special: s,
			})
		}
	}
}

func (board *Board) generateBishopMoves(bishop BitBoard) {
	for bishop := range bishop.Ones() {
		targets := bishopMoves(bishop, board.ours, board.theirs)

		for t := range targets.Ones() {
			s := NoSpecial

			if t&board.theirs != 0 {
				s |= Captures
			}

			board.possibleMoves = append(board.possibleMoves, Move{
				From:    Square(bishop),
				To:      Square(t),
				Special: s,
			})
		}
	}
}
