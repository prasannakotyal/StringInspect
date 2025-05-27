#include <stdio.h>
#include <string.h>

void print_help() {
    printf("StringInspect - Character encoding analyzer\n\n");
    printf("Usage: stringinspect [OPTIONS] <string>\n\n");
    printf("Options:\n");
    printf("  -h, --help    Show this help message\n\n");
    printf("Displays ASCII, hexadecimal, decimal, and binary representations\n");
    printf("of each character in the provided string.\n\n");
    printf("Example:\n");
    printf("  stringinspect \"Hello\"\n");
}

void analyze_string(const char *input) {
    size_t len = strlen(input);
    
    printf("Input string: \"%s\"\n", input);
    
    // ASCII row
    printf("ASCII:");
    for (size_t i = 0; i < len; i++) {
        printf("%9c", input[i]);
    }
    printf("\n");
    
    // Hexadecimal row
    printf("Hex:  ");
    for (size_t i = 0; i < len; i++) {
        printf("%9X", (unsigned char)input[i]);
    }
    printf("\n");
    
    // Decimal row
    printf("Dec:  ");
    for (size_t i = 0; i < len; i++) {
        printf("%9d", (unsigned char)input[i]);
    }
    printf("\n");
    
    // Binary row
    printf("Bin:  ");
    for (size_t i = 0; i < len; i++) {
        for (int bit = 7; bit >= 0; bit--) {
            printf("%d", ((unsigned char)input[i] >> bit) & 1);
        }
        printf(" ");
    }
    printf("\n");
}

int main(int argc, char *argv[]) {
    if (argc != 2) {
        printf("Error: Expected exactly one argument\n");
        printf("Use -h or --help for usage information\n");
        return 1;
    }
    
    if (strcmp(argv[1], "-h") == 0 || strcmp(argv[1], "--help") == 0) {
        print_help();
        return 0;
    }
    
    analyze_string(argv[1]);
    return 0;
}