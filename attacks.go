package chess

// Precomputed attack tables for all piece types.
// Initialized at package load time via init().

var (
	// PawnAttacks[side][square] — attack bitboard for a pawn on the given square.
	PawnAttacks [2][64]Bitboard

	// KnightAttacks[square] — attack bitboard for a knight on the given square.
	KnightAttacks [64]Bitboard

	// KingAttacks[square] — attack bitboard for a king on the given square.
	KingAttacks [64]Bitboard
)

func init() {
	initPawnAttacks()
	initKnightAttacks()
	initKingAttacks()
}

// initPawnAttacks precomputes pawn attack masks for both sides.
func initPawnAttacks() {
	for sq := 0; sq < 64; sq++ {
		bb := Bitboard(1) << sq

		// White pawn attacks (north-east and north-west)
		PawnAttacks[White][sq] = ((bb & NotFileA) << 7) | ((bb & NotFileH) << 9)

		// Black pawn attacks (south-east and south-west)
		PawnAttacks[Black][sq] = ((bb & NotFileH) >> 7) | ((bb & NotFileA) >> 9)
	}
}

// initKnightAttacks precomputes knight attack masks.
func initKnightAttacks() {
	for sq := 0; sq < 64; sq++ {
		bb := Bitboard(1) << sq

		var attacks Bitboard

		// 2 up, 1 right
		attacks |= (bb & NotFileH) << 17
		// 2 up, 1 left
		attacks |= (bb & NotFileA) << 15
		// 1 up, 2 right
		attacks |= (bb & NotFileGH) << 10
		// 1 up, 2 left
		attacks |= (bb & NotFileAB) << 6

		// 2 down, 1 right
		attacks |= (bb & NotFileH) >> 15
		// 2 down, 1 left
		attacks |= (bb & NotFileA) >> 17
		// 1 down, 2 right
		attacks |= (bb & NotFileGH) >> 6
		// 1 down, 2 left
		attacks |= (bb & NotFileAB) >> 10

		KnightAttacks[sq] = attacks
	}
}

// initKingAttacks precomputes king attack masks.
func initKingAttacks() {
	for sq := 0; sq < 64; sq++ {
		bb := Bitboard(1) << sq

		var attacks Bitboard

		attacks |= bb << 8             // up
		attacks |= bb >> 8             // down
		attacks |= (bb & NotFileH) << 1 // right
		attacks |= (bb & NotFileA) >> 1 // left
		attacks |= (bb & NotFileH) << 9 // up-right
		attacks |= (bb & NotFileA) << 7 // up-left
		attacks |= (bb & NotFileH) >> 7 // down-right
		attacks |= (bb & NotFileA) >> 9 // down-left

		KingAttacks[sq] = attacks
	}
}

// GetBishopAttacks generates bishop attacks on-the-fly using ray tracing.
// Iterates along each diagonal until hitting an occupied square (blocker).
func GetBishopAttacks(sq int, occupancy Bitboard) Bitboard {
	var attacks Bitboard
	r, f := RankOf(sq), FileOf(sq)

	// North-East
	for tr, tf := r+1, f+1; tr <= 7 && tf <= 7; tr, tf = tr+1, tf+1 {
		s := tr*8 + tf
		attacks = SetBit(attacks, s)
		if GetBit(occupancy, s) {
			break
		}
	}
	// North-West
	for tr, tf := r+1, f-1; tr <= 7 && tf >= 0; tr, tf = tr+1, tf-1 {
		s := tr*8 + tf
		attacks = SetBit(attacks, s)
		if GetBit(occupancy, s) {
			break
		}
	}
	// South-East
	for tr, tf := r-1, f+1; tr >= 0 && tf <= 7; tr, tf = tr-1, tf+1 {
		s := tr*8 + tf
		attacks = SetBit(attacks, s)
		if GetBit(occupancy, s) {
			break
		}
	}
	// South-West
	for tr, tf := r-1, f-1; tr >= 0 && tf >= 0; tr, tf = tr-1, tf-1 {
		s := tr*8 + tf
		attacks = SetBit(attacks, s)
		if GetBit(occupancy, s) {
			break
		}
	}

	return attacks
}

// GetRookAttacks generates rook attacks on-the-fly using ray tracing.
func GetRookAttacks(sq int, occupancy Bitboard) Bitboard {
	var attacks Bitboard
	r, f := RankOf(sq), FileOf(sq)

	// North
	for tr := r + 1; tr <= 7; tr++ {
		s := tr*8 + f
		attacks = SetBit(attacks, s)
		if GetBit(occupancy, s) {
			break
		}
	}
	// South
	for tr := r - 1; tr >= 0; tr-- {
		s := tr*8 + f
		attacks = SetBit(attacks, s)
		if GetBit(occupancy, s) {
			break
		}
	}
	// East
	for tf := f + 1; tf <= 7; tf++ {
		s := r*8 + tf
		attacks = SetBit(attacks, s)
		if GetBit(occupancy, s) {
			break
		}
	}
	// West
	for tf := f - 1; tf >= 0; tf-- {
		s := r*8 + tf
		attacks = SetBit(attacks, s)
		if GetBit(occupancy, s) {
			break
		}
	}

	return attacks
}

// GetQueenAttacks generates queen attacks as the union of bishop and rook attacks.
func GetQueenAttacks(sq int, occupancy Bitboard) Bitboard {
	return GetBishopAttacks(sq, occupancy) | GetRookAttacks(sq, occupancy)
}

// IsSquareAttacked returns true if the given square is attacked by the given side.
// Uses reverse attack lookup: checks if a piece of each type on the target square
// would attack any actual piece of that type belonging to the attacker.
func IsSquareAttacked(sq int, attackerSide int, pos *Position) bool {
	// Pawn attacks: check if any enemy pawn attacks this square
	// We look at where a pawn of the DEFENDER color on `sq` would attack —
	// if any attacker pawn is there, then `sq` is attacked.
	defenderSide := attackerSide ^ 1
	pawns := pos.Pieces[attackerSide][Pawn]
	if PawnAttacks[defenderSide][sq]&pawns != 0 {
		return true
	}

	// Knight attacks
	knights := pos.Pieces[attackerSide][Knight]
	if KnightAttacks[sq]&knights != 0 {
		return true
	}

	// King attacks
	king := pos.Pieces[attackerSide][King]
	if KingAttacks[sq]&king != 0 {
		return true
	}

	// Bishop / Queen attacks (diagonal)
	bishopsQueens := pos.Pieces[attackerSide][Bishop] | pos.Pieces[attackerSide][Queen]
	if GetBishopAttacks(sq, pos.Occupancy[2])&bishopsQueens != 0 {
		return true
	}

	// Rook / Queen attacks (straight)
	rooksQueens := pos.Pieces[attackerSide][Rook] | pos.Pieces[attackerSide][Queen]
	if GetRookAttacks(sq, pos.Occupancy[2])&rooksQueens != 0 {
		return true
	}

	return false
}

// IsKingInCheck returns true if the given side's king is in check.
func IsKingInCheck(side int, pos *Position) bool {
	kingSq := LSB(pos.Pieces[side][King])
	if kingSq == 64 {
		return false // no king found (should not happen in valid game)
	}
	return IsSquareAttacked(kingSq, side^1, pos)
}
