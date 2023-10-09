#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <stdbool.h>


void *safe_alloc(size_t size) {
    void *ptr = malloc(size);
    if (ptr == NULL) {
        fprintf(stderr, "Memory allocation failed\n");
        exit(EXIT_FAILURE);
    }
    return ptr;
}

void check_balance(char *code, size_t code_len) {
    int balance = 0;
    for (size_t i = 0; i < code_len; ++i) {
        if (code[i] == '[') {
            balance++;
        } else if (code[i] == ']') {
            balance--;
        }
        if (balance < 0) {
            fprintf(stderr, "Unbalanced brackets\n");
            exit(EXIT_FAILURE);
        }
    }
    if (balance != 0) {
        fprintf(stderr, "Unbalanced brackets\n");
        exit(EXIT_FAILURE);
    }
}

char *tape, *ptr;
size_t tape_size;

void resize_tape() {
    size_t new_tape_size = tape_size * 2;
    char *new_tape = realloc(tape, new_tape_size);
    if (new_tape == NULL) {
        fprintf(stderr, "Memory reallocation for tape failed\n");
        exit(EXIT_FAILURE);
    }
    memset(new_tape + tape_size, 0, new_tape_size - tape_size);
    ptr = new_tape + (ptr - tape);
    tape = new_tape;
    tape_size = new_tape_size;
}

char brainfuck(char *code, size_t code_len, bool step_by_step, bool debug_mode) {
    check_balance(code, code_len);
    char output[30000] = {0};
    int output_index = 0;
    int code_index = 0;

    while (code_index < code_len) {
        char c = code[code_index];
        if (debug_mode) {
            printf("Debug: Executing command '%c' at index %d\n", c, code_index);
        }

        if (ptr >= tape + tape_size) {
            resize_tape();
        }

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
                        if (code_index >= code_len) {
                            fprintf(stderr, "Jumped to unbalanced bracket\n");
                            exit(EXIT_FAILURE);
                        }
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
                        if (code_index < 0) {
                            fprintf(stderr, "Jumped to unbalanced bracket\n");
                            exit(EXIT_FAILURE);
                        }
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

        if (step_by_step && output_index > 0) {
            output[output_index] = '\0';  // Null-terminate for printing
            printf("Tape (non-zero values): ");
            for (size_t i = 0; i < tape_size; ++i) {
                if (tape[i] != 0) {
                    printf("[%zu]=%d ", i, tape[i]);
                }
            }
            printf("Output: %s\n", output);
        }

        code_index++;
    }

    output[output_index] = '\0';  // Ensure the output is null-terminated
    printf("%s", output);
    return 0;
}




void read_code_from_stdin(char **code, size_t *code_len) {
    size_t capacity = 4096;
    *code = safe_alloc(capacity * sizeof(char));

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

    *code = safe_alloc(*code_len * sizeof(char));

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
    fprintf(stderr, "Usage: %s [options]\n\n", program_name);
    fprintf(stderr, "Options:\n");
    fprintf(stderr, "  -c /path/to/code.bf        Read Brainfuck code from a file.\n");
    fprintf(stderr, "  -s                         Enable step-by-step execution, showing tape and output.\n");
    fprintf(stderr, "\nExamples:\n");
    fprintf(stderr, "  %s -c /path/to/code.bf     Execute the Brainfuck code in the specified file.\n", program_name);
    fprintf(stderr, "  %s -s                      Execute Brainfuck code from stdin with step-by-step execution.\n", program_name);
    fprintf(stderr, "  %s -c /path/to/code.bf -s  Execute the Brainfuck code in the specified file with step-by-step execution.\n", program_name);
    fprintf(stderr, "\nIf no options are provided, Brainfuck code is read from stdin.\n");
    exit(EXIT_FAILURE);
}


int main(int argc, char *argv[]) {
    char *code;
    size_t code_len;
    bool step_by_step = false;
    bool debug_mode = false;
    char *file_path = NULL;

    int opt;
    while ((opt = getopt(argc, argv, "c:sd")) != -1) {
        switch (opt) {
            case 'c':
                file_path = optarg;
                break;
            case 's':
                step_by_step = true;
                break;
            case 'd':
                debug_mode = true;
                break;
            default:
                show_usage(argv[0]);
                break;
        }
    }

    if (file_path) {
        read_code_from_file(file_path, &code, &code_len);
    } else if (isatty(STDIN_FILENO) && optind == argc) {
        show_usage(argv[0]);
    } else {
        read_code_from_stdin(&code, &code_len);
    }

    tape_size = 30000;
    tape = calloc(tape_size, sizeof(char));
    if (tape == NULL) {
        fprintf(stderr, "Memory allocation for tape failed\n");
        exit(EXIT_FAILURE);
    }
    ptr = tape;

    brainfuck(code, code_len, step_by_step, debug_mode);

    free(tape);
    free(code);

    return 0;
}
