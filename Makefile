CC=gcc
CFLAGS=-Wall -Werror

all: bli.c
		$(CC) $(CFLAGS) bli.c -o bli 

clean: bli
		rm -f bli bli.js bli.wasm

web: bli.html
		emcc bli.c -o bli.js -s EXPORTED_FUNCTIONS="['_run_brainfuck']" -s EXTRA_EXPORTED_RUNTIME_METHODS="['ccall', 'cwrap', 'UTF8ToString']"
