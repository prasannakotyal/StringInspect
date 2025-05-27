# StringInspect

A command-line utility for analyzing character encoding representations in strings.

## Overview

StringInspect examines each character in a given string and displays its ASCII, hexadecimal, decimal, and binary representations. This tool is useful for debugging character encoding issues, educational purposes, and understanding low-level string representation.

## Usage

```bash
./stringinspect [OPTIONS] <string>
```

### Options

- `-h`, `--help` - Display help information

### Examples

```bash
# Analyze a string
./stringinspect "Hello"

# Display help
./stringinspect --help
```

### Sample Output

```
Input string: "Hello"
ASCII:       H        e        l        l        o
Hex:        48       65       6C       6C       6F
Dec:        72      101      108      108      111
Bin:  01001000 01100101 01101100 01101100 01101111
```

## Building

Compile the program using the provided Makefile:

```bash
make
```

## Installation

To install system-wide:

```bash
sudo make install
```

To uninstall:

```bash
sudo make uninstall
```
