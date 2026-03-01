package chess

import (
	"testing"
)

// ============================================================
// Bitboard Utility Tests
// ============================================================

func TestBitboardSetClearGet(t *testing.T) {
	var bb Bitboard

	// Set A1
	bb = SetBit(bb, A1)
	if !GetBit(bb, A1) {
		t.Error("A1 should be set")
	}
	if GetBit(bb, A2) {
		t.Error("A2 should not be set")
	}

	// Set H8
	bb = SetBit(bb, H8)
	if !GetBit(bb, H8) {
		t.Error("H8 should be set")
	}

	// Clear A1
	bb = ClearBit(bb, A1)
	if GetBit(bb, A1) {
		t.Error("A1 should be cleared")
	}
	if !GetBit(bb, H8) {
		t.Error("H8 should still be set")
	}
}

func TestPopCount(t *testing.T) {
	tests := []struct {
		bb   Bitboard
		want int
	}{
		{0, 0},
		{1, 1},
		{0xFF, 8},
		{Rank1, 8},
		{Rank1 | Rank8, 16},
		{^Bitboard(0), 64},
	}
	for _, tc := range tests {
		if got := PopCount(tc.bb); got != tc.want {
			t.Errorf("PopCount(0x%X) = %d, want %d", uint64(tc.bb), got, tc.want)
		}
	}
}

func TestLSBAndPopLSB(t *testing.T) {
	// LSB of empty
	if LSB(0) != 64 {
		t.Error("LSB of 0 should be 64")
	}

	bb := SetBit(0, D4)
	bb = SetBit(bb, G7)

	sq := LSB(bb)
	if sq != D4 {
		t.Errorf("LSB should be D4 (%d), got %d", D4, sq)
	}

	sq = PopLSB(&bb)
	if sq != D4 {
		t.Errorf("PopLSB should return D4 (%d), got %d", D4, sq)
	}
	if GetBit(bb, D4) {
		t.Error("D4 should be cleared after PopLSB")
	}
	if !GetBit(bb, G7) {
		t.Error("G7 should still be set after PopLSB")
	}
}

func TestFileRank(t *testing.T) {
	if FileOf(E4) != 4 {
		t.Errorf("FileOf(E4) = %d, want 4", FileOf(E4))
	}
	if RankOf(E4) != 3 {
		t.Errorf("RankOf(E4) = %d, want 3", RankOf(E4))
	}
	if FileOf(A1) != 0 {
		t.Errorf("FileOf(A1) = %d, want 0", FileOf(A1))
	}
	if RankOf(H8) != 7 {
		t.Errorf("RankOf(H8) = %d, want 7", RankOf(H8))
	}
}

// ============================================================
// Attack Table Tests
// ============================================================

func TestPawnAttacks(t *testing.T) {
	// White pawn on E4 attacks D5 and F5
	attacks := PawnAttacks[White][E4]
	if !GetBit(attacks, D5) {
		t.Error("White pawn on E4 should attack D5")
	}
	if !GetBit(attacks, F5) {
		t.Error("White pawn on E4 should attack F5")
	}
	if PopCount(attacks) != 2 {
		t.Errorf("White pawn on E4 should have 2 attacks, got %d", PopCount(attacks))
	}

	// White pawn on A2 attacks only B3
	attacks = PawnAttacks[White][A2]
	if !GetBit(attacks, B3) {
		t.Error("White pawn on A2 should attack B3")
	}
	if PopCount(attacks) != 1 {
		t.Errorf("White pawn on A2 should have 1 attack, got %d", PopCount(attacks))
	}

	// Black pawn on E5 attacks D4 and F4
	attacks = PawnAttacks[Black][E5]
	if !GetBit(attacks, D4) {
		t.Error("Black pawn on E5 should attack D4")
	}
	if !GetBit(attacks, F4) {
		t.Error("Black pawn on E5 should attack F4")
	}
}

func TestKnightAttacks(t *testing.T) {
	// Knight on E4 should have 8 possible attacks
	attacks := KnightAttacks[E4]
	expected := []int{D6, F6, C5, G5, C3, G3, D2, F2}
	for _, sq := range expected {
		if !GetBit(attacks, sq) {
			t.Errorf("Knight on E4 should attack %s", SquareToString[sq])
		}
	}
	if PopCount(attacks) != 8 {
		t.Errorf("Knight on E4 should have 8 attacks, got %d", PopCount(attacks))
	}

	// Knight on A1 should have only 2 attacks
	attacks = KnightAttacks[A1]
	if PopCount(attacks) != 2 {
		t.Errorf("Knight on A1 should have 2 attacks, got %d", PopCount(attacks))
	}
}

func TestKingAttacks(t *testing.T) {
	// King on E4 should have 8 attacks
	attacks := KingAttacks[E4]
	if PopCount(attacks) != 8 {
		t.Errorf("King on E4 should have 8 attacks, got %d", PopCount(attacks))
	}

	// King on A1 should have 3 attacks
	attacks = KingAttacks[A1]
	if PopCount(attacks) != 3 {
		t.Errorf("King on A1 should have 3 attacks, got %d", PopCount(attacks))
	}
}

func TestSlidingAttacks(t *testing.T) {
	// Bishop on E4 with empty board
	attacks := GetBishopAttacks(E4, 0)
	// Should reach all 4 diagonals to the edge
	if !GetBit(attacks, D5) || !GetBit(attacks, A8) {
		t.Error("Bishop on E4 should reach D5 and A8 on empty board")
	}
	if !GetBit(attacks, H7) {
		t.Error("Bishop on E4 should reach H7 on empty board")
	}

	// Rook on E4 with empty board
	attacks = GetRookAttacks(E4, 0)
	if !GetBit(attacks, E1) || !GetBit(attacks, E8) {
		t.Error("Rook on E4 should reach E1 and E8 on empty board")
	}
	if !GetBit(attacks, A4) || !GetBit(attacks, H4) {
		t.Error("Rook on E4 should reach A4 and H4 on empty board")
	}

	// Rook on E4 with blocker on E6
	occ := SetBit(0, E6)
	attacks = GetRookAttacks(E4, Bitboard(occ))
	if !GetBit(attacks, E5) || !GetBit(attacks, E6) {
		t.Error("Rook on E4 should see E5 and E6 (blocker)")
	}
	if GetBit(attacks, E7) {
		t.Error("Rook on E4 should NOT see E7 (blocked by E6)")
	}
}

// ============================================================
// Position Tests
// ============================================================

func TestInitialPosition(t *testing.T) {
	pos := NewPosition()

	// White pawns on rank 2
	if pos.Pieces[White][Pawn] != Rank2 {
		t.Error("White pawns should be on rank 2")
	}
	// Black pawns on rank 7
	if pos.Pieces[Black][Pawn] != Rank7 {
		t.Error("Black pawns should be on rank 7")
	}
	// White king on E1
	if !GetBit(pos.Pieces[White][King], E1) {
		t.Error("White king should be on E1")
	}
	// Black king on E8
	if !GetBit(pos.Pieces[Black][King], E8) {
		t.Error("Black king should be on E8")
	}

	// 16 white pieces, 16 black pieces
	if PopCount(pos.Occupancy[White]) != 16 {
		t.Errorf("White should have 16 pieces, got %d", PopCount(pos.Occupancy[White]))
	}
	if PopCount(pos.Occupancy[Black]) != 16 {
		t.Errorf("Black should have 16 pieces, got %d", PopCount(pos.Occupancy[Black]))
	}
	if PopCount(pos.Occupancy[2]) != 32 {
		t.Errorf("Total occupancy should be 32, got %d", PopCount(pos.Occupancy[2]))
	}

	// Side to move
	if pos.SideToMove != White {
		t.Error("White should move first")
	}

	// Full castling rights
	if pos.CastlingRights != (WhiteKingSide | WhiteQueenSide | BlackKingSide | BlackQueenSide) {
		t.Error("Full castling rights should be set")
	}

	// King timers
	if pos.WhiteKingTime != DefaultKingTime {
		t.Errorf("White king time should be %d, got %d", DefaultKingTime, pos.WhiteKingTime)
	}
}

func TestPieceOnSquare(t *testing.T) {
	pos := NewPosition()

	side, piece := pos.PieceOnSquare(E1)
	if side != White || piece != King {
		t.Errorf("E1 should have White King, got side=%d piece=%d", side, piece)
	}

	side, piece = pos.PieceOnSquare(E8)
	if side != Black || piece != King {
		t.Errorf("E8 should have Black King, got side=%d piece=%d", side, piece)
	}

	side, piece = pos.PieceOnSquare(E4)
	if piece != NoPiece {
		t.Errorf("E4 should be empty, got side=%d piece=%d", side, piece)
	}
}

// ============================================================
// Move Encoding Tests
// ============================================================

func TestMoveEncoding(t *testing.T) {
	// Normal move: e2-e4, pawn, no promotion, double push
	m := NewMove(E2, E4, Pawn, 0, FlagDoublePush)
	if m.From() != E2 {
		t.Errorf("From should be E2, got %d", m.From())
	}
	if m.To() != E4 {
		t.Errorf("To should be E4, got %d", m.To())
	}
	if m.Piece() != Pawn {
		t.Errorf("Piece should be Pawn, got %d", m.Piece())
	}
	if !m.IsDoublePush() {
		t.Error("Should be double push")
	}
	if m.IsCapture() {
		t.Error("Should not be capture")
	}

	// Capture with promotion
	m2 := NewMove(A7, B8, Pawn, Queen, FlagCapture)
	if m2.From() != A7 {
		t.Errorf("From should be A7, got %d", m2.From())
	}
	if m2.To() != B8 {
		t.Errorf("To should be B8, got %d", m2.To())
	}
	if !m2.IsCapture() {
		t.Error("Should be capture")
	}
	if !m2.IsPromotion() {
		t.Error("Should be promotion")
	}
	if m2.Promotion() != Queen {
		t.Errorf("Promotion should be Queen, got %d", m2.Promotion())
	}

	// Castling
	m3 := NewMove(E1, G1, King, 0, FlagCastling)
	if !m3.IsCastle() {
		t.Error("Should be castling")
	}

	// String representation
	if m.String() != "e2e4" {
		t.Errorf("Expected 'e2e4', got '%s'", m.String())
	}
	if m2.String() != "a7b8q" {
		t.Errorf("Expected 'a7b8q', got '%s'", m2.String())
	}
}

// ============================================================
// Legal Move Generation Tests
// ============================================================

func TestLegalMovesInitialPosition(t *testing.T) {
	pos := NewPosition()
	moves := GenerateLegalMoves(pos)

	// Standard starting position has exactly 20 legal moves:
	// 16 pawn moves (8 single + 8 double) + 4 knight moves
	if len(moves) != 20 {
		t.Errorf("Initial position should have 20 legal moves, got %d", len(moves))
		for _, m := range moves {
			t.Logf("  %s", m.MoveString())
		}
	}
}

func TestIsSquareAttackedInitial(t *testing.T) {
	pos := NewPosition()

	// E2 is not attacked by black in starting position
	if IsSquareAttacked(E2, Black, pos) {
		t.Error("E2 should not be attacked by black initially")
	}

	// E3 is attacked by white pawns (D2 and F2)
	// Actually E3 is not attacked by pawns from d2/f2 — let me check:
	// Pawn on D2 attacks C3 and E3, Pawn on F2 attacks E3 and G3
	if !IsSquareAttacked(E3, White, pos) {
		t.Error("E3 should be attacked by white (pawns on d2 and f2)")
	}
}

func TestApplyMoveBasic(t *testing.T) {
	pos := NewPosition()

	// 1. e4
	move := NewMove(E2, E4, Pawn, 0, FlagDoublePush)
	newPos := ApplyMove(pos, move)

	// Pawn should be on E4, not E2
	if GetBit(newPos.Pieces[White][Pawn], E2) {
		t.Error("E2 should be empty after e4")
	}
	if !GetBit(newPos.Pieces[White][Pawn], E4) {
		t.Error("E4 should have white pawn after e4")
	}

	// En passant should be E3
	if newPos.EnPassant != E3 {
		t.Errorf("En passant should be E3 (%d), got %d", E3, newPos.EnPassant)
	}

	// Side to move should be Black
	if newPos.SideToMove != Black {
		t.Error("Side to move should be Black after White's move")
	}

	// Original position should be unchanged
	if !GetBit(pos.Pieces[White][Pawn], E2) {
		t.Error("Original position should still have pawn on E2")
	}
}

func TestEnPassant(t *testing.T) {
	pos := NewPosition()

	// Set up en passant scenario:
	// 1. e4 d5 2. e5 f5 (now white can play exf6 en passant)
	pos = ApplyMove(pos, NewMove(E2, E4, Pawn, 0, FlagDoublePush))
	pos = ApplyMove(pos, NewMove(D7, D5, Pawn, 0, FlagDoublePush))
	pos = ApplyMove(pos, NewMove(E4, E5, Pawn, 0, 0))
	pos = ApplyMove(pos, NewMove(F7, F5, Pawn, 0, FlagDoublePush))

	// En passant square should be F6
	if pos.EnPassant != F6 {
		t.Errorf("En passant should be F6 (%d), got %d", F6, pos.EnPassant)
	}

	// White should have en passant capture available
	moves := GenerateLegalMoves(pos)
	foundEP := false
	for _, m := range moves {
		if m.IsEnPassant() && m.From() == E5 && m.To() == F6 {
			foundEP = true
			break
		}
	}
	if !foundEP {
		t.Error("En passant capture e5xf6 should be available")
	}

	// Apply the en passant capture
	epMove := NewMove(E5, F6, Pawn, 0, FlagCapture|FlagEnPassant)
	newPos := ApplyMove(pos, epMove)

	// White pawn should be on F6
	if !GetBit(newPos.Pieces[White][Pawn], F6) {
		t.Error("White pawn should be on F6 after en passant")
	}
	// Black pawn on F5 should be removed
	if GetBit(newPos.Pieces[Black][Pawn], F5) {
		t.Error("Black pawn on F5 should be captured by en passant")
	}
}

func TestCastling(t *testing.T) {
	// Set up a position where white can castle kingside
	pos := &Position{
		SideToMove:     White,
		CastlingRights: WhiteKingSide | WhiteQueenSide | BlackKingSide | BlackQueenSide,
		EnPassant:      NoSquare,
		FullMoveNumber: 1,
		WhiteKingTime:  DefaultKingTime,
		BlackKingTime:  DefaultKingTime,
	}
	pos.Pieces[White][King] = SetBit(0, E1)
	pos.Pieces[White][Rook] = SetBit(0, H1) | SetBit(0, A1)
	pos.Pieces[Black][King] = SetBit(0, E8)
	// Add some pawns to make it realistic
	pos.Pieces[White][Pawn] = Rank2
	pos.Pieces[Black][Pawn] = Rank7
	pos.UpdateOccupancies()

	moves := GenerateLegalMoves(pos)

	// Should find kingside castling
	foundKingSide := false
	foundQueenSide := false
	for _, m := range moves {
		if m.IsCastle() && m.To() == G1 {
			foundKingSide = true
		}
		if m.IsCastle() && m.To() == C1 {
			foundQueenSide = true
		}
	}
	if !foundKingSide {
		t.Error("White should be able to castle kingside")
	}
	if !foundQueenSide {
		t.Error("White should be able to castle queenside")
	}

	// Apply kingside castling
	castleMove := NewMove(E1, G1, King, 0, FlagCastling)
	newPos := ApplyMove(pos, castleMove)

	// King should be on G1
	if !GetBit(newPos.Pieces[White][King], G1) {
		t.Error("King should be on G1 after kingside castling")
	}
	// Rook should be on F1
	if !GetBit(newPos.Pieces[White][Rook], F1) {
		t.Error("Rook should be on F1 after kingside castling")
	}
	// Castling rights should be lost for white
	if newPos.CastlingRights&(WhiteKingSide|WhiteQueenSide) != 0 {
		t.Error("White should lose castling rights after castling")
	}
}

func TestCastlingBlocked(t *testing.T) {
	// Position where pieces block castling
	pos := &Position{
		SideToMove:     White,
		CastlingRights: WhiteKingSide,
		EnPassant:      NoSquare,
		FullMoveNumber: 1,
		WhiteKingTime:  DefaultKingTime,
		BlackKingTime:  DefaultKingTime,
	}
	pos.Pieces[White][King] = SetBit(0, E1)
	pos.Pieces[White][Rook] = SetBit(0, H1)
	pos.Pieces[White][Bishop] = SetBit(0, F1) // blocks kingside
	pos.Pieces[Black][King] = SetBit(0, E8)
	pos.UpdateOccupancies()

	moves := GenerateLegalMoves(pos)
	for _, m := range moves {
		if m.IsCastle() {
			t.Error("Castling should be blocked by bishop on F1")
		}
	}
}

// ============================================================
// Check / Checkmate / Stalemate Tests
// ============================================================

func TestCheckDetection(t *testing.T) {
	pos := NewPosition()

	if IsKingInCheck(White, pos) {
		t.Error("White should not be in check at start")
	}
	if IsKingInCheck(Black, pos) {
		t.Error("Black should not be in check at start")
	}
}

func TestFoolsMate(t *testing.T) {
	// Fool's Mate: 1. f3 e5 2. g4 Qh4#
	gs := InitializeGame()
	var err error

	gs, err = ApplyMoveBySquares(gs, F2, F3, 0, 1000)
	if err != nil {
		t.Fatalf("f3 should be legal: %v", err)
	}

	gs, err = ApplyMoveBySquares(gs, E7, E5, 0, 1000)
	if err != nil {
		t.Fatalf("e5 should be legal: %v", err)
	}

	gs, err = ApplyMoveBySquares(gs, G2, G4, 0, 1000)
	if err != nil {
		t.Fatalf("g4 should be legal: %v", err)
	}

	gs, err = ApplyMoveBySquares(gs, D8, H4, 0, 1000)
	if err != nil {
		t.Fatalf("Qh4 should be legal: %v", err)
	}

	if !gs.IsCheck {
		t.Error("White should be in check")
	}
	if !gs.IsCheckmate {
		t.Error("Should be checkmate (Fool's Mate)")
	}
	if gs.GameResult != ResultBlackWins {
		t.Errorf("Black should win, got %s", gs.GameResult)
	}
	if !gs.IsGameOver {
		t.Error("Game should be over")
	}
}

func TestScholarsMate(t *testing.T) {
	// Scholar's Mate: 1. e4 e5 2. Bc4 Nc6 3. Qh5 Nf6?? 4. Qxf7#
	gs := InitializeGame()
	var err error

	steps := []struct {
		from, to int
	}{
		{E2, E4}, {E7, E5},
		{F1, C4}, {B8, C6},
		{D1, H5}, {G8, F6},
		{H5, F7}, // Qxf7#
	}

	for i, s := range steps {
		gs, err = ApplyMoveBySquares(gs, s.from, s.to, 0, 1000)
		if err != nil {
			t.Fatalf("Move %d (%s%s) should be legal: %v", i+1,
				SquareToString[s.from], SquareToString[s.to], err)
		}
	}

	if !gs.IsCheckmate {
		t.Error("Should be checkmate (Scholar's Mate)")
	}
	if gs.GameResult != ResultWhiteWins {
		t.Errorf("White should win, got %s", gs.GameResult)
	}
}

func TestStalemateDetection(t *testing.T) {
	// Set up a known stalemate position:
	// White: King on A1
	// Black: King on C2, Queen on B3
	// White to move — no legal moves but not in check
	pos := &Position{
		SideToMove:     White,
		CastlingRights: 0,
		EnPassant:      NoSquare,
		FullMoveNumber: 1,
		WhiteKingTime:  DefaultKingTime,
		BlackKingTime:  DefaultKingTime,
	}
	pos.Pieces[White][King] = SetBit(0, A1)
	pos.Pieces[Black][King] = SetBit(0, C2)
	pos.Pieces[Black][Queen] = SetBit(0, B3)
	pos.UpdateOccupancies()

	gs := GetGameState(pos)

	if gs.IsCheck {
		t.Error("White should NOT be in check")
	}
	if !gs.IsStalemate {
		t.Errorf("Should be stalemate, got gameOver=%v check=%v moves=%d",
			gs.IsGameOver, gs.IsCheck, len(gs.LegalMoves))
	}
	if gs.GameResult != ResultDraw {
		t.Errorf("Result should be draw, got %s", gs.GameResult)
	}
}

func TestNotStalemateKingCanMove(t *testing.T) {
	// White King on A1, Black King on C3 — white can move to B1, A2, B2
	pos := &Position{
		SideToMove:     White,
		CastlingRights: 0,
		EnPassant:      NoSquare,
		FullMoveNumber: 1,
		WhiteKingTime:  DefaultKingTime,
		BlackKingTime:  DefaultKingTime,
	}
	pos.Pieces[White][King] = SetBit(0, A1)
	pos.Pieces[Black][King] = SetBit(0, C3)
	pos.UpdateOccupancies()

	gs := GetGameState(pos)
	if gs.IsStalemate {
		t.Error("Should NOT be stalemate — king can move")
	}
	if len(gs.LegalMoves) == 0 {
		t.Error("Should have legal moves")
	}
}

// ============================================================
// Timer Tests
// ============================================================

func TestTimerDeduction(t *testing.T) {
	pos := NewPosition()

	// Apply a move with 5000ms thinking time
	move := NewMove(E2, E4, Pawn, 0, FlagDoublePush)
	newPos := ApplyMoveWithTimer(pos, move, 5000)

	if newPos.WhiteKingTime != DefaultKingTime-5000 {
		t.Errorf("White king time should be %d, got %d",
			DefaultKingTime-5000, newPos.WhiteKingTime)
	}
	if newPos.BlackKingTime != DefaultKingTime {
		t.Errorf("Black king time should be unchanged (%d), got %d",
			DefaultKingTime, newPos.BlackKingTime)
	}
}

func TestTimerCheckBonus(t *testing.T) {
	// Set up a position where white can give check
	// White: Ke1, Qd1 — Black: Ke8
	// White plays Qd1-e2+ (check via e-file... actually need a real check)
	// Let's use a simple position: White Ke1, Qh5 — Black Ke8
	pos := &Position{
		SideToMove:     White,
		CastlingRights: 0,
		EnPassant:      NoSquare,
		FullMoveNumber: 1,
		WhiteKingTime:  30000,
		BlackKingTime:  30000,
	}
	pos.Pieces[White][King] = SetBit(0, E1)
	pos.Pieces[White][Queen] = SetBit(0, H5)
	pos.Pieces[Black][King] = SetBit(0, E8)
	pos.UpdateOccupancies()

	// Qh5-e5+ gives check (queen attacks along file e toward king on e8? No, e5 to e8 is rook-like)
	// Actually Qh5-f7 would be closer... Let's try Qh5-e2 — no that doesn't check.
	// Simply: Qh5-h8+ (check along 8th rank? No, king on e8).
	// Qh5-e8+ — queen goes to e8, that's just capture? No, nothing to capture.
	// Wait: Qh5 to e8 — is that legal? H5 to E8... diagonal? h5-g6-f7-e8, yes diagonal.
	// That gives check... wait, queen will be ON e8, same square as king. That's capture.
	// Let's try Qh5-f7, that's diagonal from h5, and from f7 the queen attacks e8.

	// Actually let's just use Qh5-e5 and the queen attacks e8 along the file.
	move, found := FindLegalMove(pos, H5, E5, 0)
	if !found {
		t.Fatal("Qh5-e5 should be a legal move")
	}

	newPos := ApplyMoveWithTimer(pos, move, 2000)

	// Verify check
	if !IsKingInCheck(Black, newPos) {
		t.Fatal("Black should be in check after Qe5")
	}

	// White (attacker): 30000 - 2000 (time) + 10000 (bonus) = 38000
	expectedWhite := int64(30000 - 2000 + 10000)
	if newPos.WhiteKingTime != expectedWhite {
		t.Errorf("White king time should be %d, got %d", expectedWhite, newPos.WhiteKingTime)
	}

	// Black (defender): 30000 - 5000 (penalty) = 25000
	expectedBlack := int64(30000 - 5000)
	if newPos.BlackKingTime != expectedBlack {
		t.Errorf("Black king time should be %d, got %d", expectedBlack, newPos.BlackKingTime)
	}
}

func TestTimerTimeout(t *testing.T) {
	pos := &Position{
		SideToMove:     White,
		CastlingRights: 0,
		EnPassant:      NoSquare,
		FullMoveNumber: 1,
		WhiteKingTime:  3000, // only 3 seconds
		BlackKingTime:  30000,
	}
	pos.Pieces[White][King] = SetBit(0, E1)
	pos.Pieces[White][Pawn] = SetBit(0, E2)
	pos.Pieces[Black][King] = SetBit(0, E8)
	pos.UpdateOccupancies()

	// White spends 5 seconds thinking — more than available
	move := NewMove(E2, E4, Pawn, 0, FlagDoublePush)
	newPos := ApplyMoveWithTimer(pos, move, 5000)

	// Timer should be 0 (clamped)
	if newPos.WhiteKingTime != 0 {
		t.Errorf("White king time should be 0, got %d", newPos.WhiteKingTime)
	}

	// Check timeout
	timeout, side := IsTimeout(newPos)
	if !timeout {
		t.Error("Should be timeout")
	}
	if side != White {
		t.Errorf("White should have timed out, got side %d", side)
	}

	// Game state should reflect timeout
	gs := GetGameState(newPos)
	if !gs.IsGameOver {
		t.Error("Game should be over on timeout")
	}
	if gs.GameResult != ResultBlackWins {
		t.Errorf("Black should win on white timeout, got %s", gs.GameResult)
	}
}

func TestGetRemainingTime(t *testing.T) {
	pos := NewPosition()
	wt, bt := GetRemainingTime(pos)
	if wt != DefaultKingTime || bt != DefaultKingTime {
		t.Errorf("Expected (%d, %d), got (%d, %d)",
			DefaultKingTime, DefaultKingTime, wt, bt)
	}
}

// ============================================================
// Move Rejection Tests
// ============================================================

func TestIllegalMoveRejected(t *testing.T) {
	gs := InitializeGame()

	// Try to move black's pawn on white's turn
	_, err := ApplyMoveBySquares(gs, E7, E5, 0, 1000)
	if err == nil {
		t.Error("Moving black's pawn on white's turn should be rejected")
	}

	// Try moving to an occupied square (e2 to d2)
	_, err = ApplyMoveBySquares(gs, E2, D2, 0, 1000)
	if err == nil {
		t.Error("Moving pawn to occupied square should be rejected")
	}

	// Try moving a knight to an invalid square
	_, err = ApplyMoveBySquares(gs, B1, B3, 0, 1000)
	if err == nil {
		t.Error("Invalid knight move should be rejected")
	}
}

func TestGameOverPreventsMoving(t *testing.T) {
	// Play Fool's Mate then try another move
	gs := InitializeGame()
	var err error
	gs, _ = ApplyMoveBySquares(gs, F2, F3, 0, 1000)
	gs, _ = ApplyMoveBySquares(gs, E7, E5, 0, 1000)
	gs, _ = ApplyMoveBySquares(gs, G2, G4, 0, 1000)
	gs, err = ApplyMoveBySquares(gs, D8, H4, 0, 1000) // Qh4#
	if err != nil {
		t.Fatalf("Qh4 should succeed: %v", err)
	}

	// Try to play after game over
	_, err = ApplyMoveBySquares(gs, E2, E4, 0, 1000)
	if err == nil {
		t.Error("Should not be able to move after game is over")
	}
}

// ============================================================
// Promotion Tests
// ============================================================

func TestPromotion(t *testing.T) {
	// White pawn on A7, push to A8 for promotion
	pos := &Position{
		SideToMove:     White,
		CastlingRights: 0,
		EnPassant:      NoSquare,
		FullMoveNumber: 1,
		WhiteKingTime:  DefaultKingTime,
		BlackKingTime:  DefaultKingTime,
	}
	pos.Pieces[White][King] = SetBit(0, E1)
	pos.Pieces[White][Pawn] = SetBit(0, A7)
	pos.Pieces[Black][King] = SetBit(0, E8)
	pos.UpdateOccupancies()

	moves := GenerateLegalMoves(pos)

	// Should have 4 promotion options for a7-a8
	promoCount := 0
	for _, m := range moves {
		if m.From() == A7 && m.To() == A8 && m.IsPromotion() {
			promoCount++
		}
	}
	if promoCount != 4 {
		t.Errorf("Should have 4 promotion options for a7-a8, got %d", promoCount)
	}

	// Apply queen promotion
	promoMove, found := FindLegalMove(pos, A7, A8, Queen)
	if !found {
		t.Fatal("Queen promotion a7a8 should be legal")
	}
	newPos := ApplyMove(pos, promoMove)

	// Should have a queen on A8, no pawn
	if !GetBit(newPos.Pieces[White][Queen], A8) {
		t.Error("White queen should be on A8 after promotion")
	}
	if GetBit(newPos.Pieces[White][Pawn], A8) {
		t.Error("No pawn should remain on A8 after promotion")
	}
	if GetBit(newPos.Pieces[White][Pawn], A7) {
		t.Error("Pawn should be removed from A7")
	}
}

// ============================================================
// Full Game API Tests
// ============================================================

func TestInitializeGame(t *testing.T) {
	gs := InitializeGame()

	if gs.IsGameOver {
		t.Error("New game should not be over")
	}
	if gs.IsCheck {
		t.Error("No check at start")
	}
	if gs.GameResult != ResultOngoing {
		t.Errorf("Result should be ongoing, got %s", gs.GameResult)
	}
	if len(gs.LegalMoves) != 20 {
		t.Errorf("Should have 20 legal moves, got %d", len(gs.LegalMoves))
	}
}

func TestCopyPosition(t *testing.T) {
	pos := NewPosition()
	copy := pos.Copy()

	// Modify copy
	copy.Pieces[White][Pawn] = 0
	copy.UpdateOccupancies()

	// Original should be unchanged
	if pos.Pieces[White][Pawn] != Rank2 {
		t.Error("Original position should be unchanged after modifying copy")
	}
}

// ============================================================
// Benchmark
// ============================================================

func BenchmarkGenerateLegalMoves(b *testing.B) {
	pos := NewPosition()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateLegalMoves(pos)
	}
}

func BenchmarkApplyMove(b *testing.B) {
	pos := NewPosition()
	move := NewMove(E2, E4, Pawn, 0, FlagDoublePush)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ApplyMove(pos, move)
	}
}

func BenchmarkIsSquareAttacked(b *testing.B) {
	pos := NewPosition()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsSquareAttacked(E4, Black, pos)
	}
}
