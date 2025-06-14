# List available recipes
default:
    @just --list

# Run golangci-lint to check code
lint:
    golangci-lint run

# Build the lexis binary
build: lint
    go build -o lexis

# Run the lexis application
run: lint
    go run .

# Install lexis binary to $GOPATH/bin
install:
    go install .