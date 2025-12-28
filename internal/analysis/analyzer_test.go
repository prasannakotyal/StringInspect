package analysis

import (
	"testing"
)

func TestAnalyzeString(t *testing.T) {
	a := NewAnalyzer()

	tests := []struct {
		name    string
		input   string
		wantLen int
		wantHex string // First character's hex
		wantDec int    // First character's decimal
	}{
		{
			name:    "simple ASCII",
			input:   "A",
			wantLen: 1,
			wantHex: "41",
			wantDec: 65,
		},
		{
			name:    "hello",
			input:   "Hello",
			wantLen: 5,
			wantHex: "48",
			wantDec: 72,
		},
		{
			name:    "empty string",
			input:   "",
			wantLen: 0,
		},
		{
			name:    "space",
			input:   " ",
			wantLen: 1,
			wantHex: "20",
			wantDec: 32,
		},
		{
			name:    "unicode emoji",
			input:   "😀",
			wantLen: 1,
			wantHex: "1F600",
			wantDec: 128512,
		},
		{
			name:    "mixed ASCII and Unicode",
			input:   "Hi🌍",
			wantLen: 3,
			wantHex: "48",
			wantDec: 72,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chars := a.AnalyzeString(tt.input)

			if len(chars) != tt.wantLen {
				t.Errorf("AnalyzeString(%q) len = %d, want %d", tt.input, len(chars), tt.wantLen)
			}

			if tt.wantLen > 0 {
				if chars[0].Hex != tt.wantHex {
					t.Errorf("AnalyzeString(%q)[0].Hex = %s, want %s", tt.input, chars[0].Hex, tt.wantHex)
				}
				if chars[0].Dec != tt.wantDec {
					t.Errorf("AnalyzeString(%q)[0].Dec = %d, want %d", tt.input, chars[0].Dec, tt.wantDec)
				}
			}
		})
	}
}

func TestCharacterType(t *testing.T) {
	a := NewAnalyzer()

	tests := []struct {
		input    string
		wantType CharType
	}{
		{"A", CharTypePrintable},
		{" ", CharTypeWhitespace},
		{"\t", CharTypeWhitespace},
		{"\n", CharTypeWhitespace},
		{"\x00", CharTypeControl},
		{"\x1B", CharTypeControl}, // ESC
		{"é", CharTypeExtended},
		{"日", CharTypeExtended},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			chars := a.AnalyzeString(tt.input)
			if len(chars) != 1 {
				t.Fatalf("expected 1 character, got %d", len(chars))
			}
			if chars[0].Type != tt.wantType {
				t.Errorf("type = %v, want %v", chars[0].Type, tt.wantType)
			}
		})
	}
}

func TestBinaryFormat(t *testing.T) {
	a := NewAnalyzer()

	chars := a.AnalyzeString("A")
	if len(chars) != 1 {
		t.Fatal("expected 1 character")
	}

	expected := "01000001"
	if chars[0].Bin != expected {
		t.Errorf("binary = %s, want %s", chars[0].Bin, expected)
	}
}

func TestUTF8Bytes(t *testing.T) {
	a := NewAnalyzer()

	tests := []struct {
		input   string
		wantLen int // Number of UTF-8 bytes
	}{
		{"A", 1}, // ASCII
		{"é", 2}, // 2-byte UTF-8
		{"中", 3}, // 3-byte UTF-8
		{"😀", 4}, // 4-byte UTF-8
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			chars := a.AnalyzeString(tt.input)
			if len(chars) != 1 {
				t.Fatalf("expected 1 character, got %d", len(chars))
			}
			if len(chars[0].UTF8Bytes) != tt.wantLen {
				t.Errorf("UTF8Bytes len = %d, want %d", len(chars[0].UTF8Bytes), tt.wantLen)
			}
		})
	}
}

func TestAnalyze(t *testing.T) {
	// Test the convenience function
	chars := Analyze("test")
	if len(chars) != 4 {
		t.Errorf("Analyze('test') len = %d, want 4", len(chars))
	}
}
