#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#ifdef __EMSCRIPTEN__
#include <emscripten.h>
#endif


char *tape, *ptr;

char *brainfuck(char *code) {
    char output[30000] = {0};
    int output_index = 0;
    int code_index = 0;
    size_t code_len = strlen(code);

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
        while (code_index < code_len) {
        char c = code[code_index];

        switch (c) {
            // ...
        }
        code_index++;
    }

        code_index++;
    }

    char *result = malloc(output_index + 1);
    memcpy(result, output, output_index);
    output[output_index] = '\0';
    printf("%s", output);
    return result;
}

void read_code_from_stdin(char **code, size_t *code_len) {
    size_t capacity = 4096;
    *code = malloc(capacity * sizeof(char));
    if (*code == NULL) {
        fprintf(stderr, "Memory allocation for code buffer failed\n");
        exit(EXIT_FAILURE);
    }

    char c;
    *code_len = 0;
    while ((c = getchar()) != EOF) {
        if (*code_len >= capacity) {
            capacity *= 2;
            *code = realloc(*code, capacity * sizeof(char));
            if (*code == NULL) {
                fprintf(stderr, "Memory reallocation for code buffer failed\n");
                exit(EXIT_FAILURE);
            }
        }
        (*code)[(*code_len)++] = c;
    }
}

void read_code_from_file(const char *file_path, char **code, size_t *code_len) {
    FILE *file = fopen(file_path, "r");
    if (file == NULL) {
        perror("Error opening file");
        exit(EXIT_FAILURE);
    }

    fseek(file, 0, SEEK_END);
    *code_len = ftell(file);
    fseek(file, 0, SEEK_SET);

    *code = malloc(*code_len * sizeof(char));
    if (*code == NULL) {
        fprintf(stderr, "Memory allocation for code buffer failed\n");
        exit(EXIT_FAILURE);
    }

    fread(*code, sizeof(char), *code_len, file);
    fclose(file);
}

int is_stdin_empty() {
    int c = getchar();
    if (c == EOF) {
        return 1;
    }
    ungetc(c, stdin);
    return 0;
}

void show_usage(const char *program_name) {
    fprintf(stderr, "Usage: %s [-c /path/to/code.bf]\n", program_name);
    exit(EXIT_FAILURE);
}

#ifdef __EMSCRIPTEN__
EMSCRIPTEN_KEEPALIVE
char *run_brainfuck(const char *code_str) {
    size_t code_len = strlen(code_str);
    char *code = (char *)code_str;

    tape = calloc(30000, sizeof(char));
    if (tape == NULL) {
        fprintf(stderr, "Memory allocation for tape failed\n");
        exit(EXIT_FAILURE);
    }
    ptr = tape;

    char *output = brainfuck(code);

    free(tape);

    return output;
}

#else
int main(int argc, char *argv[]) {
    char *code;
    size_t code_len;

    if (argc == 1) {
        if (is_stdin_empty()) {
            show_usage(argv[0]);
        }
        read_code_from_stdin(&code, &code_len);
    } else if (argc == 3 && strcmp(argv[1], "-c") == 0) {
        read_code_from_file(argv[2], &code, &code_len);
    } else {
        fprintf(stderr, "Usage: %s [-c /path/to/code.bf]\n", argv[0]);
        exit(EXIT_FAILURE);
    }

    tape = calloc(30000, sizeof(char));
    if (tape == NULL) {
        fprintf(stderr, "Memory allocation for tape failed\n");
        exit(EXIT_FAILURE);
    }
    ptr = tape;

    char *output = brainfuck(code);
    printf("%s\n", output);

    free(tape);
    free(code);

    return 0;
}
#endif