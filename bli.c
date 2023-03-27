#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <emscripten.h>

char tape[30000], *ptr = tape;

char brainfuck(char *code, size_t code_len) {
    char output[30000] = {0};
    int output_index = 0;
    int code_index = 0;

    while (code_index < code_len) {
        char c = code[code_index];

        switch (c) {
            case '>':
                ptr++;
                break;
            case '<':
                ptr--;
                break;
            case '+':
                (*ptr)++;
                break;
            case '-':
                (*ptr)--;
                break;
            case '.':
                output[output_index] = *ptr;
                output_index++;
                break;
            case ',':
                *ptr = getchar();
                break;
            case '[':
                if (*ptr == 0) {
                    int balance = 1;
                    while (balance) {
                        code_index++;
                        char next_c = code[code_index];
                        if (next_c == '[') {
                            balance++;
                        } else if (next_c == ']') {
                            balance--;
                        }
                    }
                }
                break;
            case ']':
                if (*ptr != 0) {
                    int balance = 1;
                    while (balance) {
                        code_index--;
                        char next_c = code[code_index];
                        if (next_c == ']') {
                            balance++;
                        } else if (next_c == '[') {
                            balance--;
                        }
                    }
                }
                break;
            default:
                break;
        }

        code_index++;
    }

    printf("%s", output);
    return 0;
}

int main(int argc, char *argv[]) {
    if (argc < 2) {
        fprintf(stderr, "Usage: %s \"brainfuck-code\"\n", argv[0]);
        exit(EXIT_FAILURE);
    }

    char *code = argv[1];
    size_t code_len = strlen(code);

    brainfuck(code, code_len);

    return 0;
}
