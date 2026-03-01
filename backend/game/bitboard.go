package game

import "fmt"

// Bitboard is a 64-bit unsigned integer representing a chess board.
// Each bit corresponds to a square: bit 0 = A1, bit 7 = H1, bit 56 = A8, bit 63 = H8.
type Bitboard uint64

// Square indices (rank-file mapping: index = rank*8 + file)
const (
	A1 = iota; B1; C1; D1; E1; F1; G1; H1
	A2; B2; C2; D2; E2; F2; G2; H2
	A3; B3; C3; D3; E3; F3; G3; H3
	A4; B4; C4; D4; E4; F4; G4; H4
	A5; B5; C5; D5; E5; F5; G5; H5
	A6; B6; C6; D6; E6; F6; G6; H6
	A7; B7; C7; D7; E7; F7; G7; H7
	A8; B8; C8; D8; E8; F8; G8; H8
	NoSquare = 64
)

// File masks
const (
	FileA Bitboard = 0x0101010101010101
	FileB Bitboard = FileA << 1
	FileC Bitboard = FileA << 2
	FileD Bitboard = FileA << 3
	FileE Bitboard = FileA << 4
	FileF Bitboard = FileA << 5
	FileG Bitboard = FileA << 6
	FileH Bitboard = FileA << 7
)

// Rank masks
const (
	Rank1 Bitboard = 0x00000000000000FF
	Rank2 Bitboard = Rank1 << 8
	Rank3 Bitboard = Rank1 << 16
	Rank4 Bitboard = Rank1 << 24
	Rank5 Bitboard = Rank1 << 32
	Rank6 Bitboard = Rank1 << 40
	Rank7 Bitboard = Rank1 << 48
	Rank8 Bitboard = Rank1 << 56
)

// NotFileA and NotFileH are used to prevent wrapping when shifting.
const (
	NotFileA Bitboard = ^FileA
	NotFileH Bitboard = ^FileH
	NotFileAB Bitboard = ^(FileA | FileB)
	NotFileGH Bitboard = ^(FileG | FileH)
)

// Piece types
const (
	Pawn   = iota
	Knight
	Bishop
	Rook
	Queen
	King
	NoPiece = 6
)

// Sides
const (
	White = 0
	Black = 1
)

// Castling rights flags
const (
	WhiteKingSide  = 1
	WhiteQueenSide = 2
	BlackKingSide  = 4
	BlackQueenSide = 8
)

// SquareToString maps a square index to algebraic notation.
var SquareToString [65]string

// StringToSquare maps algebraic notation to a square index.
var StringToSquare map[string]int

// PieceToString maps piece types to single-character representations.
var PieceToString = [7]string{"P", "N", "B", "R", "Q", "K", "."}

func init() {
	StringToSquare = make(map[string]int)
	files := "abcdefgh"
	ranks := "12345678"
	for r := 0; r < 8; r++ {
		for f := 0; f < 8; f++ {
			sq := r*8 + f
			name := string(files[f]) + string(ranks[r])
			SquareToString[sq] = name
			StringToSquare[name] = sq
		}
	}
	SquareToString[NoSquare] = "-"
}

// SetBit returns a bitboard with the given square set.
func SetBit(bb Bitboard, sq int) Bitboard {
	return bb | (1 << sq)
}

// ClearBit returns a bitboard with the given square cleared.
func ClearBit(bb Bitboard, sq int) Bitboard {
	return bb &^ (1 << sq)
}

// GetBit returns true if the given square is set in the bitboard.
func GetBit(bb Bitboard, sq int) bool {
	return bb&(1<<sq) != 0
}

// PopCount returns the number of set bits (population count) using Brian Kernighan's algorithm.
func PopCount(bb Bitboard) int {
	count := 0
	for bb != 0 {
		bb &= bb - 1
		count++
	}
	return count
}

// LSB returns the index of the least significant set bit.
// Returns 64 if the bitboard is empty.
func LSB(bb Bitboard) int {
	if bb == 0 {
		return 64
	}
	count := 0
	for bb&1 == 0 {
		bb >>= 1
		count++
	}
	return count
}

// PopLSB removes and returns the index of the least significant set bit.
func PopLSB(bb *Bitboard) int {
	sq := LSB(*bb)
	*bb &= *bb - 1
	return sq
}

// PrintBitboard prints a visual representation of a bitboard to stdout.
func PrintBitboard(bb Bitboard) {
	fmt.Println()
	for rank := 7; rank >= 0; rank-- {
		fmt.Printf("  %d  ", rank+1)
		for file := 0; file < 8; file++ {
			sq := rank*8 + file
			if GetBit(bb, sq) {
				fmt.Print("1 ")
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Println()
	}
	fmt.Println("     a b c d e f g h")
	fmt.Printf("\n     Bitboard: 0x%016X\n\n", uint64(bb))
}

// FileOf returns the file (0-7) of a square.
func FileOf(sq int) int {
	return sq & 7
}

// RankOf returns the rank (0-7) of a square.
func RankOf(sq int) int {
	return sq >> 3
}
