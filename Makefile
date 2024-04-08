BINARY_NAME=emu8tor
OUTDIR=bin

all: output build

output:
	mkdir -p ./bin

build: output
	go build -o bin/ ./cmd/${BINARY_NAME} 

run: build
	./${OUTDIR}/${BINARY_NAME}

clean:
	go clean
	rm -r ./bin

test:
	go test ./internal/chip8/

rayrender: output
	go build -o bin/ ./cmd/rayrender

raychecker: output
	go build -o bin/ ./cmd/raychecker
