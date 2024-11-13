#include <stdio.h>
#include <string.h>
#include <stdbool.h>

void print_help() {
    printf("StringInspect CLI Tool\n");
    printf("Usage:\n");
    printf("  stringinspect [OPTIONS] <string>\n\n");
    printf("Options:\n");
    printf("  -h, --help       Show this help message\n");
    printf("  -v, --version    Show version information\n\n");
    printf("Description:\n");
    printf("  This tool inspects ASCII, Hex, Decimal, and Binary representations of each character in the provided string.\n\n");
    printf("Examples:\n");
    printf("  stringinspect \"Hello\"\n");
    printf("  stringinspect -h\n");
    printf("  stringinspect --version\n");
}

void print_version() {
    printf("StringInspect version 1.0.0\n");
}

void print_character_analysis(const char *input) {
    size_t str_len = strlen(input);

    // ASCII values
    printf("ASCII:");
    for (int i = 0; i < str_len; i++) {
        printf("\t%8c", input[i]);
    }

    // Hex values
    printf("\nHex:");
    for (int i = 0; i < str_len; i++) {
        printf("\t%8X", input[i]);
    }

    // Decimal values
    printf("\nDec:");
    for (int i = 0; i < str_len; i++) {
        printf("\t%8d", input[i]);
    }

    // Binary values
    printf("\nBin:");
    for (int i = 0; i < str_len; i++) {
        printf("\t");
        for (int b = 7; b >= 0; b--) {
            unsigned bit = (input[i] >> b) & 0b00000001;
            printf("%u", bit);
        }
    }
    printf("\n");
}

int main(int argc, char *argv[]) {
    if (argc < 2) {
        printf("Please pass an argument\n");
        return 1;
    }

    // Check for flags
    if (strcmp(argv[1], "-h") == 0 || strcmp(argv[1], "--help") == 0) {
        print_help();
        return 0;
    }

    if (strcmp(argv[1], "-v") == 0 || strcmp(argv[1], "--version") == 0) {
        print_version();
        return 0;
    }

    // If it's not a flag, treat the argument as the input string
    printf("Input string: \"%s\"\n", argv[1]);
    print_character_analysis(argv[1]);

    return 0;
}
