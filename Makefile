test:
	go test ./...

# golangci-lint, gometalinter and others doesn't support go modules yet
lint:
	go vet ./...

build: test lint
	go build -o bin/weather-reporter ./cmd/weather-reporter
