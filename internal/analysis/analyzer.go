package analysis

import (
	"fmt"
)

// Analyzer handles string analysis operations.
type Analyzer struct{}

// NewAnalyzer creates a new Analyzer instance.
func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

// AnalyzeString examines each character in the input string and returns
// a slice of Character structs containing all encoding representations.
func (a *Analyzer) AnalyzeString(input string) []Character {
	if input == "" {
		return nil
	}

	bytes := []byte(input)
	runes := []rune(input)
	characters := make([]Character, 0, len(runes))

	byteOffset := 0
	for runeOffset, r := range runes {
		// Get UTF-8 bytes for this rune
		runeBytes := []byte(string(r))
		utf8Hex := ""
		for i, b := range runeBytes {
			if i > 0 {
				utf8Hex += " "
			}
			utf8Hex += fmt.Sprintf("%02X", b)
		}

		char := Character{
			Rune:       r,
			Char:       displayChar(r),
			Hex:        fmt.Sprintf("%02X", r),
			Dec:        int(r),
			Bin:        formatBinary(r),
			Oct:        fmt.Sprintf("%o", r),
			Unicode:    fmt.Sprintf("U+%04X", r),
			UTF8Bytes:  runeBytes,
			UTF8Hex:    utf8Hex,
			Type:       classifyRune(r),
			ByteOffset: byteOffset,
			RuneOffset: runeOffset,
		}

		characters = append(characters, char)
		byteOffset += len(runeBytes)
	}

	// Sanity check
	_ = bytes

	return characters
}

// AnalyzeBytes examines each byte and returns Character structs.
// Unlike AnalyzeString, this treats each byte individually.
func (a *Analyzer) AnalyzeBytes(input []byte) []Character {
	characters := make([]Character, 0, len(input))

	for i, b := range input {
		r := rune(b)
		char := Character{
			Rune:       r,
			Char:       displayChar(r),
			Hex:        fmt.Sprintf("%02X", b),
			Dec:        int(b),
			Bin:        formatBinaryByte(b),
			Oct:        fmt.Sprintf("%03o", b),
			Unicode:    fmt.Sprintf("U+%04X", r),
			UTF8Bytes:  []byte{b},
			UTF8Hex:    fmt.Sprintf("%02X", b),
			Type:       classifyRune(r),
			ByteOffset: i,
			RuneOffset: i,
		}
		characters = append(characters, char)
	}

	return characters
}

// formatBinary converts a rune to its binary representation.
// Pads to appropriate width based on value.
func formatBinary(r rune) string {
	if r <= 0xFF {
		return fmt.Sprintf("%08b", r)
	} else if r <= 0xFFFF {
		return fmt.Sprintf("%016b", r)
	}
	return fmt.Sprintf("%021b", r) // Max 21 bits for Unicode
}

// formatBinaryByte converts a byte to 8-bit binary string.
func formatBinaryByte(b byte) string {
	return fmt.Sprintf("%08b", b)
}

// Analyze is a convenience function that creates an analyzer and analyzes a string.
func Analyze(input string) []Character {
	a := NewAnalyzer()
	return a.AnalyzeString(input)
}
