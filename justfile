# List available recipes
default:
    @just --list

# Run the lexis application
run:
    go run .

# Build the lexis binary
build:
    go build -o lexis

# Install lexis binary to $GOPATH/bin
install:
    go install .