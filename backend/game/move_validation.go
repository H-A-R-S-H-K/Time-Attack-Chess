package game

// move_validation.go — legal move generation and move validation.
// Generates pseudo-legal moves, then filters out those that leave own king in check.

// GeneratePseudoLegalMoves generates all pseudo-legal moves for the side to move.
// These may include moves that leave the king in check.
func GeneratePseudoLegalMoves(pos *Position) []Move {
	moves := make([]Move, 0, 64)
	side := pos.SideToMove
	enemy := side ^ 1
	ownOccupancy := pos.Occupancy[side]
	enemyOccupancy := pos.Occupancy[enemy]

	// --- PAWN MOVES ---
	pawns := pos.Pieces[side][Pawn]
	for pawns != 0 {
		from := PopLSB(&pawns)
		if side == White {
			moves = generateWhitePawnMoves(pos, from, ownOccupancy, enemyOccupancy, moves)
		} else {
			moves = generateBlackPawnMoves(pos, from, ownOccupancy, enemyOccupancy, moves)
		}
	}

	// --- KNIGHT MOVES ---
	knights := pos.Pieces[side][Knight]
	for knights != 0 {
		from := PopLSB(&knights)
		targets := KnightAttacks[from] & ^ownOccupancy
		for targets != 0 {
			to := PopLSB(&targets)
			flags := 0
			if GetBit(enemyOccupancy, to) {
				flags = FlagCapture
			}
			moves = append(moves, NewMove(from, to, Knight, 0, flags))
		}
	}

	// --- BISHOP MOVES ---
	bishops := pos.Pieces[side][Bishop]
	for bishops != 0 {
		from := PopLSB(&bishops)
		targets := GetBishopAttacks(from, pos.Occupancy[2]) & ^ownOccupancy
		for targets != 0 {
			to := PopLSB(&targets)
			flags := 0
			if GetBit(enemyOccupancy, to) {
				flags = FlagCapture
			}
			moves = append(moves, NewMove(from, to, Bishop, 0, flags))
		}
	}

	// --- ROOK MOVES ---
	rooks := pos.Pieces[side][Rook]
	for rooks != 0 {
		from := PopLSB(&rooks)
		targets := GetRookAttacks(from, pos.Occupancy[2]) & ^ownOccupancy
		for targets != 0 {
			to := PopLSB(&targets)
			flags := 0
			if GetBit(enemyOccupancy, to) {
				flags = FlagCapture
			}
			moves = append(moves, NewMove(from, to, Rook, 0, flags))
		}
	}

	// --- QUEEN MOVES ---
	queens := pos.Pieces[side][Queen]
	for queens != 0 {
		from := PopLSB(&queens)
		targets := GetQueenAttacks(from, pos.Occupancy[2]) & ^ownOccupancy
		for targets != 0 {
			to := PopLSB(&targets)
			flags := 0
			if GetBit(enemyOccupancy, to) {
				flags = FlagCapture
			}
			moves = append(moves, NewMove(from, to, Queen, 0, flags))
		}
	}

	// --- KING MOVES ---
	kingBB := pos.Pieces[side][King]
	if kingBB != 0 {
		from := LSB(kingBB)
		targets := KingAttacks[from] & ^ownOccupancy
		for targets != 0 {
			to := PopLSB(&targets)
			flags := 0
			if GetBit(enemyOccupancy, to) {
				flags = FlagCapture
			}
			moves = append(moves, NewMove(from, to, King, 0, flags))
		}

		// --- CASTLING ---
		moves = generateCastlingMoves(pos, side, from, moves)
	}

	return moves
}

// generateWhitePawnMoves generates pseudo-legal pawn moves for white.
func generateWhitePawnMoves(pos *Position, from int, ownOcc, enemyOcc Bitboard, moves []Move) []Move {
	allOcc := pos.Occupancy[2]

	// Single push
	to := from + 8
	if to <= 63 && !GetBit(allOcc, to) {
		if RankOf(to) == 7 { // promotion rank
			moves = append(moves, NewMove(from, to, Pawn, Queen, 0))
			moves = append(moves, NewMove(from, to, Pawn, Rook, 0))
			moves = append(moves, NewMove(from, to, Pawn, Bishop, 0))
			moves = append(moves, NewMove(from, to, Pawn, Knight, 0))
		} else {
			moves = append(moves, NewMove(from, to, Pawn, 0, 0))

			// Double push from rank 2
			if RankOf(from) == 1 {
				to2 := from + 16
				if !GetBit(allOcc, to2) {
					moves = append(moves, NewMove(from, to2, Pawn, 0, FlagDoublePush))
				}
			}
		}
	}

	// Captures
	attacks := PawnAttacks[White][from]
	captures := attacks & enemyOcc
	for captures != 0 {
		to := PopLSB(&captures)
		if RankOf(to) == 7 { // capture with promotion
			moves = append(moves, NewMove(from, to, Pawn, Queen, FlagCapture))
			moves = append(moves, NewMove(from, to, Pawn, Rook, FlagCapture))
			moves = append(moves, NewMove(from, to, Pawn, Bishop, FlagCapture))
			moves = append(moves, NewMove(from, to, Pawn, Knight, FlagCapture))
		} else {
			moves = append(moves, NewMove(from, to, Pawn, 0, FlagCapture))
		}
	}

	// En passant
	if pos.EnPassant != NoSquare && GetBit(attacks, pos.EnPassant) {
		moves = append(moves, NewMove(from, pos.EnPassant, Pawn, 0, FlagCapture|FlagEnPassant))
	}

	return moves
}

// generateBlackPawnMoves generates pseudo-legal pawn moves for black.
func generateBlackPawnMoves(pos *Position, from int, ownOcc, enemyOcc Bitboard, moves []Move) []Move {
	allOcc := pos.Occupancy[2]

	// Single push
	to := from - 8
	if to >= 0 && !GetBit(allOcc, to) {
		if RankOf(to) == 0 { // promotion rank
			moves = append(moves, NewMove(from, to, Pawn, Queen, 0))
			moves = append(moves, NewMove(from, to, Pawn, Rook, 0))
			moves = append(moves, NewMove(from, to, Pawn, Bishop, 0))
			moves = append(moves, NewMove(from, to, Pawn, Knight, 0))
		} else {
			moves = append(moves, NewMove(from, to, Pawn, 0, 0))

			// Double push from rank 7
			if RankOf(from) == 6 {
				to2 := from - 16
				if !GetBit(allOcc, to2) {
					moves = append(moves, NewMove(from, to2, Pawn, 0, FlagDoublePush))
				}
			}
		}
	}

	// Captures
	attacks := PawnAttacks[Black][from]
	captures := attacks & enemyOcc
	for captures != 0 {
		to := PopLSB(&captures)
		if RankOf(to) == 0 { // capture with promotion
			moves = append(moves, NewMove(from, to, Pawn, Queen, FlagCapture))
			moves = append(moves, NewMove(from, to, Pawn, Rook, FlagCapture))
			moves = append(moves, NewMove(from, to, Pawn, Bishop, FlagCapture))
			moves = append(moves, NewMove(from, to, Pawn, Knight, FlagCapture))
		} else {
			moves = append(moves, NewMove(from, to, Pawn, 0, FlagCapture))
		}
	}

	// En passant
	if pos.EnPassant != NoSquare && GetBit(attacks, pos.EnPassant) {
		moves = append(moves, NewMove(from, pos.EnPassant, Pawn, 0, FlagCapture|FlagEnPassant))
	}

	return moves
}

// generateCastlingMoves generates castling moves if legal conditions are met.
func generateCastlingMoves(pos *Position, side, kingFrom int, moves []Move) []Move {
	enemy := side ^ 1
	allOcc := pos.Occupancy[2]

	if side == White {
		// White kingside: e1-g1, f1 and g1 must be empty, e1/f1/g1 not attacked
		if pos.CastlingRights&WhiteKingSide != 0 {
			if !GetBit(allOcc, F1) && !GetBit(allOcc, G1) {
				if !IsSquareAttacked(E1, enemy, pos) &&
					!IsSquareAttacked(F1, enemy, pos) &&
					!IsSquareAttacked(G1, enemy, pos) {
					moves = append(moves, NewMove(kingFrom, G1, King, 0, FlagCastling))
				}
			}
		}
		// White queenside: e1-c1, b1/c1/d1 must be empty, e1/d1/c1 not attacked
		if pos.CastlingRights&WhiteQueenSide != 0 {
			if !GetBit(allOcc, D1) && !GetBit(allOcc, C1) && !GetBit(allOcc, B1) {
				if !IsSquareAttacked(E1, enemy, pos) &&
					!IsSquareAttacked(D1, enemy, pos) &&
					!IsSquareAttacked(C1, enemy, pos) {
					moves = append(moves, NewMove(kingFrom, C1, King, 0, FlagCastling))
				}
			}
		}
	} else {
		// Black kingside: e8-g8
		if pos.CastlingRights&BlackKingSide != 0 {
			if !GetBit(allOcc, F8) && !GetBit(allOcc, G8) {
				if !IsSquareAttacked(E8, enemy, pos) &&
					!IsSquareAttacked(F8, enemy, pos) &&
					!IsSquareAttacked(G8, enemy, pos) {
					moves = append(moves, NewMove(kingFrom, G8, King, 0, FlagCastling))
				}
			}
		}
		// Black queenside: e8-c8
		if pos.CastlingRights&BlackQueenSide != 0 {
			if !GetBit(allOcc, D8) && !GetBit(allOcc, C8) && !GetBit(allOcc, B8) {
				if !IsSquareAttacked(E8, enemy, pos) &&
					!IsSquareAttacked(D8, enemy, pos) &&
					!IsSquareAttacked(C8, enemy, pos) {
					moves = append(moves, NewMove(kingFrom, C8, King, 0, FlagCastling))
				}
			}
		}
	}

	return moves
}

// GenerateLegalMoves generates all fully legal moves for the side to move.
// Pseudo-legal moves are filtered by checking that the king is not left in check.
func GenerateLegalMoves(pos *Position) []Move {
	pseudoLegal := GeneratePseudoLegalMoves(pos)
	legal := make([]Move, 0, len(pseudoLegal))

	for _, move := range pseudoLegal {
		newPos := ApplyMove(pos, move)
		// After applying the move, check if the moving side's king is in check.
		// If yes, the move is illegal.
		if !IsKingInCheck(pos.SideToMove, newPos) {
			legal = append(legal, move)
		}
	}

	return legal
}

// ValidateMove checks if a specific move is legal in the given position.
// It searches for the move in the legal move list.
func ValidateMove(pos *Position, move Move) bool {
	legalMoves := GenerateLegalMoves(pos)
	for _, lm := range legalMoves {
		if lm.From() == move.From() && lm.To() == move.To() &&
			lm.Promotion() == move.Promotion() {
			return true
		}
	}
	return false
}

// FindLegalMove finds a matching legal move from the legal move list.
// This is useful for filling in flags that the caller might not know.
// Returns the move and true if found, zero and false otherwise.
func FindLegalMove(pos *Position, from, to, promotion int) (Move, bool) {
	legalMoves := GenerateLegalMoves(pos)
	for _, lm := range legalMoves {
		if lm.From() == from && lm.To() == to &&
			lm.Promotion() == promotion {
			return lm, true
		}
	}
	return 0, false
}

// ApplyMove applies a move to the position and returns a new position.
// The original position is not modified.
func ApplyMove(pos *Position, move Move) *Position {
	newPos := pos.Copy()
	side := pos.SideToMove
	enemy := side ^ 1
	from := move.From()
	to := move.To()
	piece := move.Piece()

	// Remove piece from source square
	newPos.Pieces[side][piece] = ClearBit(newPos.Pieces[side][piece], from)

	// Handle captures: remove captured piece
	if move.IsCapture() {
		if move.IsEnPassant() {
			// En passant: captured pawn is on a different square
			var capturedSq int
			if side == White {
				capturedSq = to - 8
			} else {
				capturedSq = to + 8
			}
			newPos.Pieces[enemy][Pawn] = ClearBit(newPos.Pieces[enemy][Pawn], capturedSq)
		} else {
			// Normal capture: find and remove the captured piece
			for pt := Pawn; pt <= King; pt++ {
				if GetBit(newPos.Pieces[enemy][pt], to) {
					newPos.Pieces[enemy][pt] = ClearBit(newPos.Pieces[enemy][pt], to)
					break
				}
			}
		}
	}

	// Place piece on destination (or promoted piece)
	if move.IsPromotion() {
		newPos.Pieces[side][move.Promotion()] = SetBit(newPos.Pieces[side][move.Promotion()], to)
	} else {
		newPos.Pieces[side][piece] = SetBit(newPos.Pieces[side][piece], to)
	}

	// Handle castling: move the rook
	if move.IsCastle() {
		switch to {
		case G1: // White kingside
			newPos.Pieces[White][Rook] = ClearBit(newPos.Pieces[White][Rook], H1)
			newPos.Pieces[White][Rook] = SetBit(newPos.Pieces[White][Rook], F1)
		case C1: // White queenside
			newPos.Pieces[White][Rook] = ClearBit(newPos.Pieces[White][Rook], A1)
			newPos.Pieces[White][Rook] = SetBit(newPos.Pieces[White][Rook], D1)
		case G8: // Black kingside
			newPos.Pieces[Black][Rook] = ClearBit(newPos.Pieces[Black][Rook], H8)
			newPos.Pieces[Black][Rook] = SetBit(newPos.Pieces[Black][Rook], F8)
		case C8: // Black queenside
			newPos.Pieces[Black][Rook] = ClearBit(newPos.Pieces[Black][Rook], A8)
			newPos.Pieces[Black][Rook] = SetBit(newPos.Pieces[Black][Rook], D8)
		}
	}

	// Update en passant square
	if move.IsDoublePush() {
		if side == White {
			newPos.EnPassant = from + 8
		} else {
			newPos.EnPassant = from - 8
		}
	} else {
		newPos.EnPassant = NoSquare
	}

	// Update castling rights
	newPos.CastlingRights = updateCastlingRights(newPos.CastlingRights, from, to)

	// Update half-move clock
	if piece == Pawn || move.IsCapture() {
		newPos.HalfMoveClock = 0
	} else {
		newPos.HalfMoveClock++
	}

	// Update full move number
	if side == Black {
		newPos.FullMoveNumber++
	}

	// Switch side to move
	newPos.SideToMove = enemy

	// Recalculate occupancies
	newPos.UpdateOccupancies()

	return newPos
}

// updateCastlingRights adjusts castling rights based on piece movement.
func updateCastlingRights(rights int, from, to int) int {
	// If king moves, lose both castling rights for that side
	if from == E1 || to == E1 {
		rights &^= WhiteKingSide | WhiteQueenSide
	}
	if from == E8 || to == E8 {
		rights &^= BlackKingSide | BlackQueenSide
	}

	// If rook moves from or is captured on its starting square
	if from == A1 || to == A1 {
		rights &^= WhiteQueenSide
	}
	if from == H1 || to == H1 {
		rights &^= WhiteKingSide
	}
	if from == A8 || to == A8 {
		rights &^= BlackQueenSide
	}
	if from == H8 || to == H8 {
		rights &^= BlackKingSide
	}

	return rights
}
