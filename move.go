package chess

import "fmt"

// Move is a 32-bit encoded chess move.
//
// Bit layout:
//   Bits  0-5:   from square (6 bits, 0-63)
//   Bits  6-11:  to square (6 bits, 0-63)
//   Bits 12-15:  piece type (4 bits)
//   Bits 16-19:  promotion piece type (4 bits, 0 if none)
//   Bits 20-23:  flags
//
// Flag bits:
//   Bit 20: capture
//   Bit 21: double pawn push
//   Bit 22: en passant
//   Bit 23: castling

type Move uint32

// Move flag constants
const (
	FlagCapture   = 1 << 20
	FlagDoublePush = 1 << 21
	FlagEnPassant  = 1 << 22
	FlagCastling   = 1 << 23
)

// NewMove creates a new encoded move.
func NewMove(from, to, piece, promotion, flags int) Move {
	return Move(from | (to << 6) | (piece << 12) | (promotion << 16) | flags)
}

// From returns the source square of the move.
func (m Move) From() int {
	return int(m) & 0x3F
}

// To returns the destination square of the move.
func (m Move) To() int {
	return (int(m) >> 6) & 0x3F
}

// Piece returns the moving piece type.
func (m Move) Piece() int {
	return (int(m) >> 12) & 0xF
}

// Promotion returns the promotion piece type (0 if not a promotion).
func (m Move) Promotion() int {
	return (int(m) >> 16) & 0xF
}

// IsCapture returns true if this move is a capture.
func (m Move) IsCapture() bool {
	return int(m)&FlagCapture != 0
}

// IsDoublePush returns true if this move is a double pawn push.
func (m Move) IsDoublePush() bool {
	return int(m)&FlagDoublePush != 0
}

// IsEnPassant returns true if this move is an en passant capture.
func (m Move) IsEnPassant() bool {
	return int(m)&FlagEnPassant != 0
}

// IsCastle returns true if this move is a castling move.
func (m Move) IsCastle() bool {
	return int(m)&FlagCastling != 0
}

// IsPromotion returns true if this move includes a pawn promotion.
func (m Move) IsPromotion() bool {
	return m.Promotion() != 0
}

// String returns a human-readable UCI-style representation of the move.
func (m Move) String() string {
	s := SquareToString[m.From()] + SquareToString[m.To()]
	if m.IsPromotion() {
		promoChars := map[int]string{Knight: "n", Bishop: "b", Rook: "r", Queen: "q"}
		if ch, ok := promoChars[m.Promotion()]; ok {
			s += ch
		}
	}
	return s
}

// MoveString returns a verbose description of the move for debugging.
func (m Move) MoveString() string {
	flags := ""
	if m.IsCapture() {
		flags += " capture"
	}
	if m.IsDoublePush() {
		flags += " double-push"
	}
	if m.IsEnPassant() {
		flags += " en-passant"
	}
	if m.IsCastle() {
		flags += " castle"
	}
	if m.IsPromotion() {
		flags += fmt.Sprintf(" promote=%s", PieceToString[m.Promotion()])
	}
	return fmt.Sprintf("%s%s (%s)%s",
		SquareToString[m.From()],
		SquareToString[m.To()],
		PieceToString[m.Piece()],
		flags,
	)
}
