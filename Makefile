CC=gcc
CFLAGS=-Wall -Werror
EMCC=emcc
EMFLAGS=-s WASM=1 -s NO_EXIT_RUNTIME=1

build: bli.c
		$(CC) $(CFLAGS) bli.c -o bli 

web: bli.c
		$(EMCC) $(EMFLAGS) bli.c -o bli.html

clean:
		rm -f bli bli.html bli.js bli.wasm
