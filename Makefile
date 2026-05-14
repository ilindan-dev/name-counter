APP_NAME = name-counter
BUILD_DIR = bin
MAIN_PKG = ./cmd/name-counter

.PHONY: build
build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PKG)

.PHONY: test
test:
	@echo "Running tests..."
	go test -v -race ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)
	rm -f name-counter-e2e

.PHONY: tools
tools:
	@echo "Installing golangci-lint..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $$(go env GOPATH)/bin v2.12.0