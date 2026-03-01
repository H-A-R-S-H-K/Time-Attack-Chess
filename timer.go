package chess

// timer.go — Time-attack mode with king timers.

// Timer bonus/penalty constants (in milliseconds)
const (
	CheckBonusMs  int64 = 10000 // +10 seconds for attacker on check
	CheckPenaltyMs int64 = 5000  // -5 seconds for defender on check
)

// ApplyMoveWithTimer applies a move and updates king timers.
//
// timeSpentMs is the time the active player spent thinking (in milliseconds).
//
// Timer rules:
//   - Deduct timeSpentMs from active player's king timer
//   - If the resulting position has the opponent in check:
//     - Attacker's king gets +CheckBonusMs
//     - Defender's king gets -CheckPenaltyMs
//   - Timers are clamped to a minimum of 0
func ApplyMoveWithTimer(pos *Position, move Move, timeSpentMs int64) *Position {
	side := pos.SideToMove
	enemy := side ^ 1

	// Apply the move
	newPos := ApplyMove(pos, move)

	// Deduct time spent from the moving player's timer
	if side == White {
		newPos.WhiteKingTime -= timeSpentMs
		if newPos.WhiteKingTime < 0 {
			newPos.WhiteKingTime = 0
		}
	} else {
		newPos.BlackKingTime -= timeSpentMs
		if newPos.BlackKingTime < 0 {
			newPos.BlackKingTime = 0
		}
	}

	// Check if the move gives check to the opponent
	if IsKingInCheck(enemy, newPos) {
		// Attacker bonus
		if side == White {
			newPos.WhiteKingTime += CheckBonusMs
		} else {
			newPos.BlackKingTime += CheckBonusMs
		}
		// Defender penalty
		if enemy == White {
			newPos.WhiteKingTime -= CheckPenaltyMs
			if newPos.WhiteKingTime < 0 {
				newPos.WhiteKingTime = 0
			}
		} else {
			newPos.BlackKingTime -= CheckPenaltyMs
			if newPos.BlackKingTime < 0 {
				newPos.BlackKingTime = 0
			}
		}
	}

	return newPos
}

// IsTimeout returns true if either player's king timer has reached 0.
// Returns (timeout, losingSide).
// If no timeout, losingSide is -1.
func IsTimeout(pos *Position) (bool, int) {
	if pos.WhiteKingTime <= 0 {
		return true, White
	}
	if pos.BlackKingTime <= 0 {
		return true, Black
	}
	return false, -1
}

// GetRemainingTime returns the remaining time for both kings in milliseconds.
// Returns (whiteTimeMs, blackTimeMs).
func GetRemainingTime(pos *Position) (int64, int64) {
	return pos.WhiteKingTime, pos.BlackKingTime
}
