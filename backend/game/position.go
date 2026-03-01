package game

import "fmt"

// Position represents the full state of a chess game.
type Position struct {
	// Pieces[side][pieceType] — bitboard for each piece type per side.
	// side: 0=White, 1=Black
	// pieceType: Pawn=0, Knight=1, Bishop=2, Rook=3, Queen=4, King=5
	Pieces [2][6]Bitboard

	// Occupancy[0]=white, Occupancy[1]=black, Occupancy[2]=all
	Occupancy [3]Bitboard

	// Side to move: White=0, Black=1
	SideToMove int

	// Castling rights (4-bit bitmask)
	CastlingRights int

	// En passant target square (NoSquare if none)
	EnPassant int

	// Half-move clock (for 50-move rule)
	HalfMoveClock int

	// Full move number (starts at 1, incremented after Black's move)
	FullMoveNumber int

	// King timers (time-attack mode) in milliseconds
	WhiteKingTime int64
	BlackKingTime int64

	// Timestamp of when the last move was made (Unix millis, 0 if unused)
	LastMoveTimestamp int64
}

// DefaultKingTime is the default starting time for each king (60 seconds).
const DefaultKingTime int64 = 60000

// NewPosition creates a standard starting position.
func NewPosition() *Position {
	pos := &Position{
		SideToMove:     White,
		CastlingRights: WhiteKingSide | WhiteQueenSide | BlackKingSide | BlackQueenSide,
		EnPassant:      NoSquare,
		HalfMoveClock:  0,
		FullMoveNumber: 1,
		WhiteKingTime:  DefaultKingTime,
		BlackKingTime:  DefaultKingTime,
	}

	// White pieces
	pos.Pieces[White][Pawn] = Rank2
	pos.Pieces[White][Knight] = SetBit(0, B1) | SetBit(0, G1)
	pos.Pieces[White][Bishop] = SetBit(0, C1) | SetBit(0, F1)
	pos.Pieces[White][Rook] = SetBit(0, A1) | SetBit(0, H1)
	pos.Pieces[White][Queen] = SetBit(0, D1)
	pos.Pieces[White][King] = SetBit(0, E1)

	// Black pieces
	pos.Pieces[Black][Pawn] = Rank7
	pos.Pieces[Black][Knight] = SetBit(0, B8) | SetBit(0, G8)
	pos.Pieces[Black][Bishop] = SetBit(0, C8) | SetBit(0, F8)
	pos.Pieces[Black][Rook] = SetBit(0, A8) | SetBit(0, H8)
	pos.Pieces[Black][Queen] = SetBit(0, D8)
	pos.Pieces[Black][King] = SetBit(0, E8)

	pos.UpdateOccupancies()
	return pos
}

// Copy creates a deep copy of the position.
func (pos *Position) Copy() *Position {
	newPos := *pos
	return &newPos
}

// UpdateOccupancies recalculates the occupancy bitboards from the piece bitboards.
func (pos *Position) UpdateOccupancies() {
	pos.Occupancy[White] = 0
	pos.Occupancy[Black] = 0
	for pt := Pawn; pt <= King; pt++ {
		pos.Occupancy[White] |= pos.Pieces[White][pt]
		pos.Occupancy[Black] |= pos.Pieces[Black][pt]
	}
	pos.Occupancy[2] = pos.Occupancy[White] | pos.Occupancy[Black]
}

// PieceOnSquare returns the side and piece type on the given square.
// Returns (-1, NoPiece) if the square is empty.
func (pos *Position) PieceOnSquare(sq int) (side int, piece int) {
	sqBit := Bitboard(1) << sq
	for s := White; s <= Black; s++ {
		for pt := Pawn; pt <= King; pt++ {
			if pos.Pieces[s][pt]&sqBit != 0 {
				return s, pt
			}
		}
	}
	return -1, NoPiece
}

// PrintBoard prints a human-readable board representation to stdout.
func (pos *Position) PrintBoard() {
	pieceChars := [2][6]string{
		{"P", "N", "B", "R", "Q", "K"}, // White (uppercase)
		{"p", "n", "b", "r", "q", "k"}, // Black (lowercase)
	}

	fmt.Println()
	fmt.Println("   ┌───┬───┬───┬───┬───┬───┬───┬───┐")
	for rank := 7; rank >= 0; rank-- {
		fmt.Printf(" %d │", rank+1)
		for file := 0; file < 8; file++ {
			sq := rank*8 + file
			side, pt := pos.PieceOnSquare(sq)
			if pt != NoPiece {
				fmt.Printf(" %s │", pieceChars[side][pt])
			} else {
				fmt.Print("   │")
			}
		}
		fmt.Println()
		if rank > 0 {
			fmt.Println("   ├───┼───┼───┼───┼───┼───┼───┼───┤")
		}
	}
	fmt.Println("   └───┴───┴───┴───┴───┴───┴───┴───┘")
	fmt.Println("     a   b   c   d   e   f   g   h")
	fmt.Println()

	sideStr := "White"
	if pos.SideToMove == Black {
		sideStr = "Black"
	}
	fmt.Printf("  Side to move: %s\n", sideStr)

	castling := ""
	if pos.CastlingRights&WhiteKingSide != 0 {
		castling += "K"
	}
	if pos.CastlingRights&WhiteQueenSide != 0 {
		castling += "Q"
	}
	if pos.CastlingRights&BlackKingSide != 0 {
		castling += "k"
	}
	if pos.CastlingRights&BlackQueenSide != 0 {
		castling += "q"
	}
	if castling == "" {
		castling = "-"
	}
	fmt.Printf("  Castling: %s\n", castling)
	fmt.Printf("  En passant: %s\n", SquareToString[pos.EnPassant])
	fmt.Printf("  Move: %d\n", pos.FullMoveNumber)
	fmt.Printf("  White King Time: %.1fs\n", float64(pos.WhiteKingTime)/1000)
	fmt.Printf("  Black King Time: %.1fs\n\n", float64(pos.BlackKingTime)/1000)
}
