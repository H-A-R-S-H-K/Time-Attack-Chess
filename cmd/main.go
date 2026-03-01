package main

import (
	"chess"
	"fmt"
)

func main() {
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println("  Chess Backend — Bitboard Engine Demo")
	fmt.Println("═══════════════════════════════════════════════════")

	// Initialize a new game
	gs := chess.InitializeGame()

	fmt.Println("\n▸ Initial position:")
	gs.Position.PrintBoard()
	fmt.Printf("  Legal moves available: %d\n\n", len(gs.LegalMoves))

	// Demonstrate playing a few moves (Scholar's Mate attempt)
	type moveStep struct {
		from, to    int
		promotion   int
		timeSpentMs int64
		description string
	}

	// Scholar's Mate: 1. e4 e5 2. Bc4 Nc6 3. Qh5 Nf6?? 4. Qxf7#
	moves := []moveStep{
		{chess.E2, chess.E4, 0, 3000, "1. e4 — White opens with king's pawn"},
		{chess.E7, chess.E5, 0, 2500, "1... e5 — Black responds symmetrically"},
		{chess.F1, chess.C4, 0, 4000, "2. Bc4 — White develops bishop to c4"},
		{chess.B8, chess.C6, 0, 3500, "2... Nc6 — Black develops knight"},
		{chess.D1, chess.H5, 0, 2000, "3. Qh5 — White brings queen out early"},
		{chess.G8, chess.F6, 0, 5000, "3... Nf6?? — Black blunders!"},
		{chess.H5, chess.F7, 0, 1500, "4. Qxf7# — Scholar's Mate!"},
	}

	for _, m := range moves {
		fmt.Println("───────────────────────────────────────────────────")
		fmt.Printf("  %s\n", m.description)

		var err error
		gs, err = chess.ApplyMoveBySquares(gs, m.from, m.to, m.promotion, m.timeSpentMs)
		if err != nil {
			fmt.Printf("  ✗ Error: %v\n", err)
			return
		}

		gs.Position.PrintBoard()

		// Show timer status
		wt, bt := chess.GetRemainingTime(gs.Position)
		fmt.Printf("  ⏱  White: %.1fs  |  Black: %.1fs\n", float64(wt)/1000, float64(bt)/1000)

		if gs.IsCheck {
			fmt.Println("  ⚡ CHECK!")
		}
		if gs.IsCheckmate {
			fmt.Println("  🏆 CHECKMATE!")
			fmt.Printf("  Result: %s (%s)\n", gs.GameResult, gs.GameOverReason)
		}
		if gs.IsStalemate {
			fmt.Println("  🤝 STALEMATE!")
		}
		if gs.IsGameOver {
			fmt.Printf("\n  Game Over: %s\n", gs.GameResult)
			break
		}

		fmt.Printf("  Legal moves for next player: %d\n", len(gs.LegalMoves))
	}

	// --- Timer demonstration ---
	fmt.Println("\n═══════════════════════════════════════════════════")
	fmt.Println("  Timer-Attack Mode Demonstration")
	fmt.Println("═══════════════════════════════════════════════════")

	gs2 := chess.InitializeGameWithTime(15000) // 15 seconds per side
	fmt.Println("\n▸ New game with 15s per king:")

	timerMoves := []moveStep{
		{chess.E2, chess.E4, 0, 5000, "1. e4 (5s thinking)"},
		{chess.E7, chess.E5, 0, 3000, "1... e5 (3s thinking)"},
		{chess.D1, chess.H5, 0, 4000, "2. Qh5 (4s thinking)"},
		{chess.A7, chess.A6, 0, 2000, "2... a6 (2s thinking)"},
		{chess.H5, chess.E5, 0, 1000, "3. Qxe5+ CHECK! (1s thinking)"},
	}

	for _, m := range timerMoves {
		fmt.Printf("\n  %s\n", m.description)
		var err error
		gs2, err = chess.ApplyMoveBySquares(gs2, m.from, m.to, m.promotion, m.timeSpentMs)
		if err != nil {
			fmt.Printf("  ✗ Error: %v\n", err)
			return
		}
		wt, bt := chess.GetRemainingTime(gs2.Position)
		fmt.Printf("  ⏱  White: %.1fs  |  Black: %.1fs", float64(wt)/1000, float64(bt)/1000)
		if gs2.IsCheck {
			fmt.Print("  ⚡ CHECK! (+10s attacker, -5s defender)")
		}
		fmt.Println()
	}

	fmt.Println("\n═══════════════════════════════════════════════════")
	fmt.Println("  Demo Complete!")
	fmt.Println("═══════════════════════════════════════════════════")
}
