BINARY_NAME=emu8tor

all: build

build:
	go build ./cmd/${BINARY_NAME} 

run: build
	./${BINARY_NAME}

clean:
	go clean
	rm ./${BINARY_NAME}

test:
	go test ./internal/chip8/
