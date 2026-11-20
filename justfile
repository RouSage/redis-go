set shell := ["zsh", "-cu"]

default:
    @just --list --unsorted

[group('help')]
_confirm:
    @echo -n 'Are you sure? [y/N] ' && read ans && [ ${ans:-N} = y ]

[group('help')]
_no-dirty:
    test -z "$(shell git status --porcelain)"

# list direct dependencies that have upgrades available
[group('help')]
upgradeable:
    @go list -u -f "$(cat .go-list-upgradeable.tmpl)" -m all

# run quality control checks
[group('lint')]
audit:
    @echo "Checking module dependencies..."
    go mod tidy -diff
    go mod verify
    test -z "$(gofmt -l .)"
    @echo "Vetting code..."
    go vet ./...
    golangci-lint run ./...
    go tool govulncheck ./...

# tidy modfiles and format .go files
[group('dev')]
tidy:
    @echo "Tidying module dependencies..."
    go mod tidy
    @echo "Formatting code..."
    go fmt ./...

# build the application
[group('build')]
[group('dev')]
build: _no-dirty audit
    @echo "Building..."
    go build -ldflags='-s' -o ./bin/redis-go ./cmd/redis-go/main.go
    GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o ./bin/linux_amd64/redis-go ./cmd/redis-go/main.go

# build the application for development
[group('build')]
[group('dev')]
build-dev:
    go build -ldflags='-s' -o ./bin/redis-go ./cmd/redis-go/main.go

# run the application
[group('dev')]
run:
    go run cmd/redis-go/main.go

# clean the binary
[group('build')]
[group('dev')]
clean:
    @echo "Cleaning..."
    rm -rf bin/

# live Reload
[group('dev')]
watch:
    #!/usr/bin/env bash
    if command -v air > /dev/null; then
        air
        @echo "Watching..."
    else
        read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice
        if [[ "$choice" != "n" && "$choice" != "N" ]]; then
            go install github.com/air-verse/air@latest
            air
            @echo "Watching..."
        else
            @echo "You chose not to install air. Exiting..."
            exit 1
        fi
    fi
