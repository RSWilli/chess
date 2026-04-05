package chess

import "testing"

func TestPieceSwitchSide(t *testing.T) {
	tests := []struct {
		input    Piece
		expected Piece
	}{
		{White, Black},
		{WhitePawn, BlackPawn},
		{WhiteKnight, BlackKnight},
		{WhiteBishop, BlackBishop},
		{WhiteRook, BlackRook},
		{WhiteQueen, BlackQueen},
		{WhiteKing, BlackKing},
		{Black, White},
		{BlackPawn, WhitePawn},
		{BlackKnight, WhiteKnight},
		{BlackBishop, WhiteBishop},
		{BlackRook, WhiteRook},
		{BlackQueen, WhiteQueen},
		{BlackKing, WhiteKing},
	}

	for _, tt := range tests {
		t.Run(tt.input.String(), func(t *testing.T) {
			got := tt.input.SwitchColor()
			if got != tt.expected {
				t.Errorf("SwitchSide() = %v, want %v", got, tt.expected)
			}
		})
	}
}
