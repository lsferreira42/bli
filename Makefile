GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOPHERJSCMD=gopherjs build
BINARY_NAME=bli
SOURCE_FILE=bli.go
WEB_BINARY_NAME=brainfuck.js

build:
	$(GOBUILD) -o $(BINARY_NAME) $(SOURCE_FILE)

web:
	$(GOPHERJSCMD) -o $(WEB_BINARY_NAME) $(SOURCE_FILE)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME) $(WEB_BINARY_NAME)

test:
	$(GOTEST) -v ./...
