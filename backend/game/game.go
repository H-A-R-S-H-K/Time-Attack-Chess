package game

import "fmt"

// game.go — High-level game API wrapping all chess backend modules.

// GameResult represents the outcome of a game.
type GameResult string

const (
	ResultOngoing   GameResult = "ongoing"
	ResultWhiteWins GameResult = "white_wins"
	ResultBlackWins GameResult = "black_wins"
	ResultDraw      GameResult = "draw"
)

// GameState represents the full evaluated state of a chess game.
type GameState struct {
	Position    *Position
	LegalMoves  []Move
	IsCheck     bool
	IsCheckmate bool
	IsStalemate bool
	IsGameOver  bool
	GameResult  GameResult
	GameOverReason string
}

// InitializeGame creates a new game with the standard starting position
// and default king timers (60 seconds each).
func InitializeGame() *GameState {
	pos := NewPosition()
	return computeGameState(pos)
}

// InitializeGameWithTime creates a new game with custom king timer values.
// timeMs is the starting time for each king in milliseconds.
func InitializeGameWithTime(timeMs int64) *GameState {
	pos := NewPosition()
	pos.WhiteKingTime = timeMs
	pos.BlackKingTime = timeMs
	return computeGameState(pos)
}

// GetLegalMoves returns all legal moves for the current position.
func GetLegalMoves(pos *Position) []Move {
	return GenerateLegalMoves(pos)
}

// GetGameState evaluates and returns the full game state.
func GetGameState(pos *Position) *GameState {
	return computeGameState(pos)
}

// ValidateAndApplyMove validates and applies a move with timer update.
// Returns the new game state and an error if the move is invalid.
func ValidateAndApplyMove(gs *GameState, move Move, timeSpentMs int64) (*GameState, error) {
	// Check if game is already over
	if gs.IsGameOver {
		return nil, fmt.Errorf("game is already over: %s", gs.GameResult)
	}

	// Find the matching legal move (to fill in correct flags)
	legalMove, found := FindLegalMove(gs.Position, move.From(), move.To(), move.Promotion())
	if !found {
		return nil, fmt.Errorf("illegal move: %s", move.String())
	}

	// Apply move with timer
	newPos := ApplyMoveWithTimer(gs.Position, legalMove, timeSpentMs)

	// Compute new game state
	newGS := computeGameState(newPos)

	return newGS, nil
}

// ApplyMoveBySquares validates and applies a move specified by from/to squares
// and optional promotion piece. This is a convenience wrapper.
func ApplyMoveBySquares(gs *GameState, from, to, promotion int, timeSpentMs int64) (*GameState, error) {
	move := NewMove(from, to, 0, promotion, 0)
	return ValidateAndApplyMove(gs, move, timeSpentMs)
}

// computeGameState evaluates the position and fills in check/checkmate/stalemate/timeout.
func computeGameState(pos *Position) *GameState {
	gs := &GameState{
		Position:   pos,
		GameResult: ResultOngoing,
	}

	// Check for timeout first
	if timeout, losingSide := IsTimeout(pos); timeout {
		gs.IsGameOver = true
		if losingSide == White {
			gs.GameResult = ResultBlackWins
			gs.GameOverReason = "white king timeout"
		} else {
			gs.GameResult = ResultWhiteWins
			gs.GameOverReason = "black king timeout"
		}
		return gs
	}

	// Generate legal moves
	gs.LegalMoves = GenerateLegalMoves(pos)

	// Check detection
	gs.IsCheck = IsKingInCheck(pos.SideToMove, pos)

	if len(gs.LegalMoves) == 0 {
		gs.IsGameOver = true
		if gs.IsCheck {
			// Checkmate
			gs.IsCheckmate = true
			if pos.SideToMove == White {
				gs.GameResult = ResultBlackWins
				gs.GameOverReason = "checkmate"
			} else {
				gs.GameResult = ResultWhiteWins
				gs.GameOverReason = "checkmate"
			}
		} else {
			// Stalemate
			gs.IsStalemate = true
			gs.GameResult = ResultDraw
			gs.GameOverReason = "stalemate"
		}
	}

	// 50-move rule
	if pos.HalfMoveClock >= 100 {
		gs.IsGameOver = true
		gs.GameResult = ResultDraw
		gs.GameOverReason = "50-move rule"
	}

	return gs
}
