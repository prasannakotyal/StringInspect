// Package analysis provides character encoding analysis functionality.
package analysis

import (
	"fmt"
	"unicode"
)

// CharType represents the category of a character.
type CharType int

const (
	CharTypePrintable CharType = iota
	CharTypeWhitespace
	CharTypeControl
	CharTypeExtended
)

// String returns the string representation of a CharType.
func (ct CharType) String() string {
	switch ct {
	case CharTypePrintable:
		return "printable"
	case CharTypeWhitespace:
		return "whitespace"
	case CharTypeControl:
		return "control"
	case CharTypeExtended:
		return "extended"
	default:
		return "unknown"
	}
}

// Character holds all encoding representations of a single character.
type Character struct {
	Rune       rune     // The actual rune
	Char       string   // String representation (or placeholder for non-printable)
	Hex        string   // Hexadecimal representation
	Dec        int      // Decimal value
	Bin        string   // Binary representation
	Oct        string   // Octal representation
	Unicode    string   // Unicode codepoint (U+XXXX)
	UTF8Bytes  []byte   // UTF-8 byte sequence
	UTF8Hex    string   // UTF-8 bytes as hex string
	Type       CharType // Character type category
	ByteOffset int      // Position in original byte slice
	RuneOffset int      // Position in rune slice
}

// String returns a display-friendly representation of the character.
func (c Character) String() string {
	return fmt.Sprintf("%s (U+%04X)", c.Char, c.Rune)
}

// IsPrintable returns true if the character is printable.
func (c Character) IsPrintable() bool {
	return c.Type == CharTypePrintable
}

// IsWhitespace returns true if the character is whitespace.
func (c Character) IsWhitespace() bool {
	return c.Type == CharTypeWhitespace
}

// IsControl returns true if the character is a control character.
func (c Character) IsControl() bool {
	return c.Type == CharTypeControl
}

// IsExtended returns true if the character is extended ASCII (>127).
func (c Character) IsExtended() bool {
	return c.Type == CharTypeExtended
}

// classifyRune determines the CharType for a given rune.
func classifyRune(r rune) CharType {
	switch {
	case r == ' ' || r == '\t' || r == '\n' || r == '\r':
		return CharTypeWhitespace
	case unicode.IsControl(r):
		return CharTypeControl
	case r > 127:
		return CharTypeExtended
	default:
		return CharTypePrintable
	}
}

// displayChar returns a display string for a character.
// Non-printable characters get special representations.
func displayChar(r rune) string {
	switch r {
	case ' ':
		return "␣" // Space symbol
	case '\t':
		return "⇥" // Tab symbol
	case '\n':
		return "↵" // Newline symbol
	case '\r':
		return "↩" // Carriage return symbol
	case '\x00':
		return "∅" // Null symbol
	default:
		if unicode.IsControl(r) || !unicode.IsPrint(r) {
			return fmt.Sprintf("<%02X>", r)
		}
		return string(r)
	}
}

// controlCharName returns the name of a control character.
func controlCharName(r rune) string {
	names := map[rune]string{
		0x00: "NUL", 0x01: "SOH", 0x02: "STX", 0x03: "ETX",
		0x04: "EOT", 0x05: "ENQ", 0x06: "ACK", 0x07: "BEL",
		0x08: "BS", 0x09: "TAB", 0x0A: "LF", 0x0B: "VT",
		0x0C: "FF", 0x0D: "CR", 0x0E: "SO", 0x0F: "SI",
		0x10: "DLE", 0x11: "DC1", 0x12: "DC2", 0x13: "DC3",
		0x14: "DC4", 0x15: "NAK", 0x16: "SYN", 0x17: "ETB",
		0x18: "CAN", 0x19: "EM", 0x1A: "SUB", 0x1B: "ESC",
		0x1C: "FS", 0x1D: "GS", 0x1E: "RS", 0x1F: "US",
		0x7F: "DEL",
	}
	if name, ok := names[r]; ok {
		return name
	}
	return ""
}
