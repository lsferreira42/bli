CC=gcc
CFLAGS=-Wall -Werror

all: bli.c
		$(CC) $(CFLAGS) bli.c -o bli 

clean: bli
		rm -f bli

web: bli.html
		emcc -O3 -s WASM=1 bli.c