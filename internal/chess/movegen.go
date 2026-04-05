package chess

import "fmt"

func (p *Position) GenerateMoves() []Move {
	// TODO: 50 move rule, draw by material, draw by repetition

	legalMoves := make([]Move, 0, maxMoveCount)

	for sq := range p.ours.Ones() {
		piece := p.at(sq)
		p.generateMovesForPiece(piece, sq, &legalMoves)
	}

	return legalMoves
}

func (p *Position) generateMovesForPiece(piece Piece, at bitBoard, legalMoves *[]Move) {
	switch piece.Type() {
	case Pawn:
		p.generatePawnMoves(at, legalMoves)
	case Bishop:
		p.generateBishopMoves(at, legalMoves)
	case King:
		p.generateKingMoves(at, legalMoves)
	case Knight:
		p.generateKnightMoves(at, legalMoves)
	case Queen:
		p.generateQueenMoves(at, legalMoves)
	case Rook:
		p.generateRookMoves(at, legalMoves)
	default:
		panic(fmt.Sprintf("unexpected chess.Piece to generate moves for: %#v", piece))
	}
}

const doublePushRank bitBoard = 0xff00
const promoteRank bitBoard = 0xff000000000000

func (p *Position) generatePawnMoves(from bitBoard, legalMoves *[]Move) {
	all := p.all()

	pushed := from.Up()
	doublePushed := from.Up().Up()

	if from&doublePushRank != 0 && all&pushed == 0 && all&doublePushed == 0 {
		// can double push

		m := Move{
			From:    p.toWhitePerspective(from),
			To:      p.toWhitePerspective(doublePushed),
			Special: DoublePawnPush,
		}

		if p.pawnCheckSquares&bitBoard(m.To) != 0 {
			m.Special |= Check
		}

		if p.isLegalMove(m) {
			*legalMoves = append(*legalMoves, m)
		}
	}

	// moves tracks the possible moves temporarily to defer
	// promotion decision
	moves := make([]Move, 0, 3)

	if all&pushed == 0 {
		// can push

		m := Move{
			From:    p.toWhitePerspective(from),
			To:      p.toWhitePerspective(pushed),
			Special: NoSpecial,
		}

		if p.pawnCheckSquares&bitBoard(m.To) != 0 {
			m.Special |= Check
		}

		if p.isLegalMove(m) {
			moves = append(moves, m)
		}
	}

	takes := []bitBoard{
		pushed.Left(),
		pushed.Right(),
	}

	enPassantTarget := p.toCurrentPerspective(p.enPassantTarget)

	for _, t := range takes {
		if t == 0 {
			// edge of the board
			continue
		}

		if p.theirs&t == 0 && enPassantTarget != t {
			continue
		}

		taken := Empty
		special := Captures

		if enPassantTarget == t {
			taken = Pawn
			special |= EnPassant
		} else {
			taken = p.at(t)
		}

		m := Move{
			From:    p.toWhitePerspective(from),
			To:      p.toWhitePerspective(t),
			Special: special,
			Takes:   taken,
		}

		if p.pawnCheckSquares&t != 0 {
			m.Special |= Check
		}

		if p.isLegalMove(m) {
			moves = append(moves, m)
		}
	}

	if from&promoteRank == 0 {
		// no promotion
		*legalMoves = append(*legalMoves, moves...)
		return
	}

	promotions := []MoveSpecial{
		PromoteQueen,
		PromoteRook,
		PromoteBishop,
		PromoteKnight,
	}

	for _, prom := range promotions {
		for _, m := range moves {
			m.Special = m.Special | prom
			*legalMoves = append(*legalMoves, m)
		}
	}
}

// bitboards to check when trying to castle, all from the perspective of the current player
const (
	whiteCastleKingEmpty      = bitBoard(0b00000110) // f1,g1
	whiteCastleKingNoCheck    = whiteCastleKingEmpty
	whiteCastleKingKingTarget = bitBoard(0b00000010)

	whiteCastleQueenEmpty      = bitBoard(0b01110000) // b1, c1, d1
	whiteCastleQueenNoCheck    = bitBoard(0b00110000) // c1, d1
	whiteCastleQueenKingTarget = bitBoard(0b00100000)

	blackCastleQueenEmpty      = bitBoard(0b00001110) // d8, c8, b8
	blackCastleQueenNoCheck    = bitBoard(0b00001100) // d8, c8
	blackCastleQueenKingTarget = bitBoard(0b00000100)

	blackCastleKingEmpty      = bitBoard(0b01100000)
	blackCastleKingNoCheck    = blackCastleKingEmpty
	blackCastleKingKingTarget = bitBoard(0b01000000)
)

// notAttacked returns true if none of the squares set in the given bitboard are attacked
func (p *Position) notAttacked(bb bitBoard) bool {
	for sq := range bb.Ones() {
		if p.attacksTo.get(sq) != 0 {
			return false
		}
	}
	return true
}

func (p *Position) canCastleWhiteKing() bool {
	return p.castling.Has(CastleWhiteKing) && p.all()&whiteCastleKingEmpty == 0 && p.notAttacked(whiteCastleKingNoCheck)
}

func (p *Position) canCastleWhiteQueen() bool {
	return p.castling.Has(CastleWhiteQueen) && p.all()&whiteCastleQueenEmpty == 0 && p.notAttacked(whiteCastleQueenNoCheck)
}
func (p *Position) canCastleBlackKing() bool {
	return p.castling.Has(CastleBlackKing) && p.all()&blackCastleKingEmpty == 0 && p.notAttacked(blackCastleKingNoCheck)
}

func (p *Position) canCastleBlackQueen() bool {
	return p.castling.Has(CastleBlackQueen) && p.all()&blackCastleQueenEmpty == 0 && p.notAttacked(blackCastleQueenNoCheck)
}

func (p *Position) generateKingMoves(bb bitBoard, legalMoves *[]Move) {
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
		taken := Empty

		if t&p.theirs != 0 {
			s |= Captures
			taken = p.at(t)
		}

		*legalMoves = append(*legalMoves, Move{
			From:    p.toWhitePerspective(bb),
			To:      p.toWhitePerspective(t),
			Special: s,
			Takes:   taken,
		})
	}

	if p.castling == NoCastling {
		return
	}

	if p.attacksTo.get(bb) != 0 {
		// king in check, no castling allowed
		return
	}

	// white O-O, aka e1g1
	if p.PlayerInTurn == White && p.canCastleWhiteKing() {
		*legalMoves = append(*legalMoves, Move{
			From:    p.toWhitePerspective(p.ourKing),
			To:      p.toWhitePerspective(whiteCastleKingKingTarget),
			Special: CastleKing,
		})
	}

	// white O-O-O, aka e1c1
	if p.PlayerInTurn == White && p.canCastleWhiteQueen() {
		*legalMoves = append(*legalMoves, Move{
			From:    p.toWhitePerspective(p.ourKing),
			To:      p.toWhitePerspective(whiteCastleQueenKingTarget),
			Special: CastleQueen,
		})
	}

	// black O-O, aka e8g8
	if p.PlayerInTurn == Black && p.canCastleBlackKing() {
		*legalMoves = append(*legalMoves, Move{
			From:    p.toWhitePerspective(p.ourKing),
			To:      p.toWhitePerspective(blackCastleKingKingTarget),
			Special: CastleKing,
		})
	}

	// black O-O-O, aka e8c8
	if p.PlayerInTurn == Black && p.canCastleBlackQueen() {
		*legalMoves = append(*legalMoves, Move{
			From:    p.toWhitePerspective(p.ourKing),
			To:      p.toWhitePerspective(blackCastleQueenKingTarget),
			Special: CastleQueen,
		})
	}
}

func (p *Position) generateKnightMoves(knight bitBoard, legalMoves *[]Move) {
	p.generatePieceMoves(knightMoves, knight, p.knightCheckSquares, legalMoves)
}

func (p *Position) generateRookMoves(rook bitBoard, legalMoves *[]Move) {
	p.generatePieceMoves(rookMoves, rook, p.rookCheckSquares, legalMoves)
}

func (p *Position) generateQueenMoves(queen bitBoard, legalMoves *[]Move) {
	p.generatePieceMoves(queenMoves, queen, p.rookCheckSquares|p.bishopCheckSquares, legalMoves)
}

func (p *Position) generateBishopMoves(bishop bitBoard, legalMoves *[]Move) {
	p.generatePieceMoves(bishopMoves, bishop, p.bishopCheckSquares, legalMoves)
}

// generatePieceMoves returns the moves for the given moveMaker function. This is used for all pieces except king and pawns
func (p *Position) generatePieceMoves(moveMaker func(sq bitBoard, same bitBoard, opposing bitBoard) bitBoard, piece bitBoard, checkSquares bitBoard, legalMoves *[]Move) {
	targets := moveMaker(piece, p.ours, p.theirs)

	capturing := targets & p.theirs
	noncapturing := targets &^ p.theirs

	for to := range noncapturing.Ones() {
		m := Move{
			From: p.toWhitePerspective(piece),
			To:   p.toWhitePerspective(to),
		}

		if !p.isLegalMove(m) {
			continue
		}

		if checkSquares&to != 0 {
			m.Special |= Check
		}

		*legalMoves = append(*legalMoves, m)
	}

	for to := range capturing.Ones() {
		m := Move{
			From:    p.toWhitePerspective(piece),
			To:      p.toWhitePerspective(to),
			Special: Captures,
			Takes:   p.at(to),
		}

		if !p.isLegalMove(m) {
			continue
		}

		if checkSquares&to != 0 {
			m.Special |= Check
		}

		*legalMoves = append(*legalMoves, m)
	}
}

// isLegalMove checks if the move would leave the king in check. This does not correctly
// check king moves or en passant captures.
func (p *Position) isLegalMove(m Move) bool {
	from := p.toCurrentPerspective(m.From)
	to := p.toCurrentPerspective(m.To)
	// enpassantTaken is the taken pawn while doing en passant
	var enpassantTaken bitBoard
	if m.Special.Has(EnPassant) {
		enpassantTaken = p.toCurrentPerspective(p.enPassantTarget).Down()
	}

	// move the piece, needed for pin detection:
	moved := (p.all() | to) &^ (from | enpassantTaken)

	var blockedSlidingAttackers bitBoard

	for _, attack := range p.xRayKingAttacks {
		if attack.from != 0 && to != attack.from && moved&attack.ray == 0 {
			// no piece is blocking the sliding piece attack or taking it, piece could have moved
			// out of a pin, or did not block the check -> not legal
			return false
		}

		blockedSlidingAttackers |= attack.from
	}

	// omit the attacks we blocked earlier
	kingAttacks := p.attacksTo.get(p.ourKing) &^ blockedSlidingAttackers

	if kingAttacks == 0 {
		return true
	}

	// move is only legal if we capture the attacker, either directly or via en passant
	return m.Special.Has(Captures) && (kingAttacks&^to) == 0 ||
		m.Special.Has(EnPassant) && (kingAttacks&^enpassantTaken) == 0
}
