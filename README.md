# StringInspect

Interactive TUI for analyzing character encodings in strings.

![Demo](assets/demo.gif)

## Features

- **Real-time analysis** - Live encoding display as you type
- **Multiple formats** - ASCII, hex, decimal, binary, octal, Unicode
- **Three view modes** - Table, detail, and compact (hex dump)
- **Unicode support** - Full UTF-8 with codepoints and byte sequences
- **Color-coded** - Printable (white), whitespace (cyan), control (pink), extended (yellow)
- **Search** - Find characters by hex (`0x41`), decimal (`65`), or literal (`A`)
- **Export** - Save analysis as text, JSON, or CSV
- **History** - Browse previous inputs with arrow keys
- **Clipboard** - Paste input, copy character info
- **File input** - Analyze files directly

## Installation

```bash
git clone https://github.com/yourusername/stringinspect.git
cd stringinspect
make build
```

## Usage

```bash
./stringinspect              # Interactive mode
./stringinspect -f file.txt  # Analyze file contents
```

## Key Bindings

| Key | Action |
|-----|--------|
| `Tab` | Cycle modes: Input → Table → Detail → Compact |
| `←`/`→`, `h`/`l` | Navigate characters |
| `Home`/`End`, `g`/`G` | Jump to first/last character |
| `PgUp`/`PgDn` | Page navigation |
| `/` | Search by hex, decimal, or character |
| `e` | Export menu (Text/JSON/CSV) |
| `c` | Copy selected character info |
| `Ctrl+V` | Paste from clipboard |
| `↑`/`↓` | History navigation (in input mode) |
| `F1` | Toggle help |
| `Esc`, `Enter` | Return to input mode |
| `q`, `Ctrl+C` | Quit |

## View Modes

**Table** - All characters with encodings in columns  
**Detail** - Single character with full encoding breakdown  
**Compact** - Hex dump view (16 bytes per line)

## Building

```bash
make build          # Build binary
make test           # Run tests
make test-coverage  # Tests with coverage
make fmt            # Format code
make lint           # Lint (requires golangci-lint)
make clean          # Clean artifacts
```

## Requirements

- Go 1.21+

## License

MIT
