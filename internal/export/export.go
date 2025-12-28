package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"stringinspect/internal/analysis"
)

// Format represents an export format.
type Format int

const (
	FormatText Format = iota
	FormatJSON
	FormatCSV
)

func (f Format) String() string {
	switch f {
	case FormatText:
		return "Text"
	case FormatJSON:
		return "JSON"
	case FormatCSV:
		return "CSV"
	default:
		return "Unknown"
	}
}

// Extension returns the file extension for the format.
func (f Format) Extension() string {
	switch f {
	case FormatText:
		return "txt"
	case FormatJSON:
		return "json"
	case FormatCSV:
		return "csv"
	default:
		return "txt"
	}
}

// Exporter handles exporting character analysis to various formats.
type Exporter struct{}

// NewExporter creates a new Exporter.
func NewExporter() *Exporter {
	return &Exporter{}
}

// Export exports the characters to the specified format and returns the filename.
func (e *Exporter) Export(chars []analysis.Character, format Format) (string, error) {
	if len(chars) == 0 {
		return "", fmt.Errorf("no characters to export")
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("stringinspect-%s.%s", timestamp, format.Extension())

	var err error
	switch format {
	case FormatText:
		err = e.exportText(chars, filename)
	case FormatJSON:
		err = e.exportJSON(chars, filename)
	case FormatCSV:
		err = e.exportCSV(chars, filename)
	default:
		return "", fmt.Errorf("unsupported format: %v", format)
	}

	if err != nil {
		return "", err
	}

	return filename, nil
}

// exportText exports characters to a text file.
func (e *Exporter) exportText(chars []analysis.Character, filename string) error {
	var b strings.Builder

	b.WriteString("StringInspect Export\n")
	b.WriteString("====================\n\n")

	// Original string
	b.WriteString("Original: ")
	for _, c := range chars {
		b.WriteString(c.Char)
	}
	b.WriteString("\n\n")

	// Character table
	b.WriteString(fmt.Sprintf("%-6s %-8s %-6s %-6s %-10s %-10s %-12s\n",
		"Pos", "Char", "Hex", "Dec", "Oct", "Unicode", "UTF-8"))
	b.WriteString(strings.Repeat("-", 70) + "\n")

	for i, c := range chars {
		charDisplay := c.Char
		if len(charDisplay) > 6 {
			charDisplay = charDisplay[:6]
		}
		b.WriteString(fmt.Sprintf("%-6d %-8s %-6s %-6d %-10s %-10s %-12s\n",
			i, charDisplay, c.Hex, c.Dec, c.Oct, c.Unicode, c.UTF8Hex))
	}

	b.WriteString(fmt.Sprintf("\nTotal: %d characters\n", len(chars)))

	return os.WriteFile(filename, []byte(b.String()), 0644)
}

// JSONCharacter is the JSON representation of a character.
type JSONCharacter struct {
	Position   int    `json:"position"`
	Char       string `json:"char"`
	Hex        string `json:"hex"`
	Decimal    int    `json:"decimal"`
	Octal      string `json:"octal"`
	Binary     string `json:"binary"`
	Unicode    string `json:"unicode"`
	UTF8Bytes  string `json:"utf8_bytes"`
	Type       string `json:"type"`
	ByteOffset int    `json:"byte_offset"`
	RuneOffset int    `json:"rune_offset"`
}

// JSONExport is the top-level JSON export structure.
type JSONExport struct {
	Original   string          `json:"original"`
	Count      int             `json:"count"`
	ExportedAt string          `json:"exported_at"`
	Characters []JSONCharacter `json:"characters"`
}

// exportJSON exports characters to a JSON file.
func (e *Exporter) exportJSON(chars []analysis.Character, filename string) error {
	// Build original string
	var original strings.Builder
	for _, c := range chars {
		original.WriteString(c.Char)
	}

	// Convert characters
	jsonChars := make([]JSONCharacter, len(chars))
	for i, c := range chars {
		jsonChars[i] = JSONCharacter{
			Position:   i,
			Char:       c.Char,
			Hex:        c.Hex,
			Decimal:    int(c.Dec),
			Octal:      c.Oct,
			Binary:     c.Bin,
			Unicode:    c.Unicode,
			UTF8Bytes:  c.UTF8Hex,
			Type:       c.Type.String(),
			ByteOffset: c.ByteOffset,
			RuneOffset: c.RuneOffset,
		}
	}

	export := JSONExport{
		Original:   original.String(),
		Count:      len(chars),
		ExportedAt: time.Now().Format(time.RFC3339),
		Characters: jsonChars,
	}

	data, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return os.WriteFile(filename, data, 0644)
}

// exportCSV exports characters to a CSV file.
func (e *Exporter) exportCSV(chars []analysis.Character, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"Position", "Char", "Hex", "Decimal", "Octal", "Binary", "Unicode", "UTF8_Bytes", "Type"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write rows
	for i, c := range chars {
		row := []string{
			fmt.Sprintf("%d", i),
			c.Char,
			c.Hex,
			fmt.Sprintf("%d", c.Dec),
			c.Oct,
			c.Bin,
			c.Unicode,
			c.UTF8Hex,
			c.Type.String(),
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	return nil
}
