
# StringInspect CLI Tool

## 📜 Description

`stringinspect` is a lightweight command-line tool that inspects and displays the ASCII, Hexadecimal, Decimal, and Binary representations of each character in a provided string. It's perfect for developers, hobbyists, and anyone curious about character encoding.

## ⚙️ Features

- Displays ASCII, Hex, Decimal, and Binary values for each character in a string.
- Supports help and version flags for quick reference.
- Easy to install and use on Linux systems.

## 🚀 Usage

```bash
./stringinspect [OPTIONS] <string>
```

## 📌 Options

| Option          | Description              |
|-----------------|--------------------------|
| `-h`, `--help`  | Show help message        |
| `-v`, `--version` | Show version information |

## 🚀 Examples

```bash
# Analyze the string "Hello"
./stringinspect "Hello"

# Display help information
./stringinspect --help

# Display version information
./stringinspect --version
```

## Sample Output

```
Input string: "Hello"
ASCII:       H        e        l        l        o
Hex:        48       65       6C       6C       6F
Dec:        72      101      108      108      111
Bin:  01001000 01100101 01101100 01101100 01101111
```

## 💡 Inspiration

Inspired by [Kay Lack](https://www.youtube.com/@neoeno4242) on YouTube, This is an exercise from the second lecture in the amazing ODE5 series.

## 🛠️ Installation

### For Linux

1. Clone the repository:

    ```bash
    git clone https://github.com/yourusername/stringinspect.git
    cd stringinspect
    ```

2. Build the tool using `make`:

    ```bash
    make
    ```

3. Install it globally (optional):

    ```bash
    sudo make install
    ```

4. Run the tool:

    ```bash
    stringinspect "Your String Here"
    ```

## Uninstallation

To remove the tool from your system:

```bash
sudo make uninstall
```
