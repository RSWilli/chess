package chess

func (p *Position) GenerateMoves() []Move {
	if p.possibleMoves != nil {
		return p.possibleMoves
	}
	// TODO: 50 move rule, draw by material
	p.possibleMoves = make([]Move, 0, maxMoveCount)

	if p.PlayerInTurn == White {

		p.generateKingMoves(p.whiteKing)
		p.whitePawns.Each(p.generateWhitePawnMoves)
		p.whiteKnights.Each(p.generateKnightMoves)

		p.whiteBishops.Each(p.generateBishopMoves)
		p.whiteRooks.Each(p.generateRookMoves)
		p.whiteQueens.Each(p.generateQueenMoves)
	} else {
		p.generateKingMoves(p.blackKing)
		p.blackPawns.Each(p.generateBlackPawnMoves)
		p.blackKnights.Each(p.generateKnightMoves)

		p.blackBishops.Each(p.generateBishopMoves)
		p.blackRooks.Each(p.generateRookMoves)
		p.blackQueens.Each(p.generateQueenMoves)
	}

	return p.possibleMoves
}

func (board *Position) generatePawnMoves(from, pushed, doublePushed, doublePushRank, promoteRank, opposing BitBoard) {
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

func (p *Position) generateWhitePawnMoves(bb BitBoard) {
	p.generatePawnMoves(bb, bb.Up(), bb.Up().Up(), rank2BitBoard, rank7BitBoard, p.blackPieces())
}

func (p *Position) generateBlackPawnMoves(bb BitBoard) {
	p.generatePawnMoves(bb, bb.Down(), bb.Down().Down(), rank7BitBoard, rank2BitBoard, p.whitePieces())
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
func (p *Position) notAttacked(bb BitBoard) bool {
	for sq := range bb.Ones() {
		if p.attacksTo.get(sq) != 0 {
			return false
		}
	}
	return true
}

func (p *Position) canCastleWhiteKing() bool {
	return p.Castling.Has(CastleWhiteKing) && p.all&whiteCastleKing == 0 && p.notAttacked(whiteCastleKing)
}

func (p *Position) canCastleWhiteQueen() bool {
	return p.Castling.Has(CastleWhiteQueen) && p.all&whiteCastleQueen == 0 && p.notAttacked(whiteCastleQueen)
}
func (p *Position) canCastleBlackKing() bool {
	return p.Castling.Has(CastleBlackKing) && p.all&blackCastleKing == 0 && p.notAttacked(blackCastleKing)
}

func (p *Position) canCastleBlackQueen() bool {
	return p.Castling.Has(CastleBlackQueen) && p.all&blackCastleQueen == 0 && p.notAttacked(blackCastleQueen)
}

func (p *Position) generateKingMoves(bb BitBoard) {
	if bb.Count() != 1 {
		panic("expected 1 king")
	}

	for t := range kingMoves(bb).Ones() {
		if t == 0 {
			// wrapped around
			continue
		}

		if p.attacksTo.get(t) != 0 {
			continue
		}

		if t&p.ours != 0 {
			// occupied
			continue
		}

		s := NoSpecial

		if t&p.theirs != 0 {
			s |= Captures
		}

		p.possibleMoves = append(p.possibleMoves, Move{
			From:    Square(bb),
			To:      Square(t),
			Special: s,
		})
	}

	if p.Castling == NoCastling {
		return
	}

	if p.attacksTo.get(bb) != 0 {
		// king in check, no castling allowed
		return
	}

	// white O-O
	if p.PlayerInTurn == White && p.canCastleWhiteKing() {
		p.possibleMoves = append(p.possibleMoves, Move{
			From:    Square(bb),
			To:      whiteCastleKingKingTarget,
			Special: CastleShort,
		})
	}

	// white O-O-O
	if p.PlayerInTurn == White && p.canCastleWhiteQueen() {
		p.possibleMoves = append(p.possibleMoves, Move{
			From:    Square(bb),
			To:      whiteCastleQueenKingTarget,
			Special: CastleLong,
		})
	}

	// black O-O
	if p.PlayerInTurn == Black && p.canCastleBlackKing() {
		p.possibleMoves = append(p.possibleMoves, Move{
			From:    Square(bb),
			To:      blackCastleKingKingTarget,
			Special: CastleShort,
		})
	}

	// black O-O-O
	if p.PlayerInTurn == Black && p.canCastleBlackQueen() {
		p.possibleMoves = append(p.possibleMoves, Move{
			From:    Square(bb),
			To:      blackCastleQueenKingTarget,
			Special: CastleLong,
		})
	}
}

func (p *Position) generateKnightMoves(bb BitBoard) {
	for t := range knightMoves(bb).Ones() {
		if t == 0 {
			// wrapped around
			continue
		}

		if t&p.ours != 0 {
			// occupied
			continue
		}

		s := NoSpecial

		if t&p.theirs != 0 {
			s |= Captures
		}

		// TODO: filter out moves with check and blocked moves
		p.possibleMoves = append(p.possibleMoves, Move{
			From:    Square(bb),
			To:      Square(t),
			Special: s,
		})
	}
}

func (p *Position) generateRookMoves(rooks BitBoard) {
	for rook := range rooks.Ones() {
		targets := rookMoves(rook, p.ours, p.theirs)

		for t := range targets.Ones() {
			s := NoSpecial

			if t&p.theirs != 0 {
				s |= Captures
			}

			p.possibleMoves = append(p.possibleMoves, Move{
				From:    Square(rook),
				To:      Square(t),
				Special: s,
			})
		}
	}
}

func (p *Position) generateQueenMoves(queens BitBoard) {
	for queen := range queens.Ones() {
		targets := queenMoves(queen, p.ours, p.theirs)

		for t := range targets.Ones() {
			s := NoSpecial

			if t&p.theirs != 0 {
				s |= Captures
			}

			p.possibleMoves = append(p.possibleMoves, Move{
				From:    Square(queen),
				To:      Square(t),
				Special: s,
			})
		}
	}
}

func (p *Position) generateBishopMoves(bishop BitBoard) {
	for bishop := range bishop.Ones() {
		targets := bishopMoves(bishop, p.ours, p.theirs)

		for t := range targets.Ones() {
			s := NoSpecial

			if t&p.theirs != 0 {
				s |= Captures
			}

			p.possibleMoves = append(p.possibleMoves, Move{
				From:    Square(bishop),
				To:      Square(t),
				Special: s,
			})
		}
	}
}
