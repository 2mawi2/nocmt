default:
    @just --list

build:
    go build -o nocmt

test:
    go test ./...

lint:
    golangci-lint run ./...

run *args:
    go run main.go {{args}}

clean:
    rm -f nocmt 

bench *args:
    ./benchmark.sh {{args}} 