.DEFAULT_GOAL := help

.PHONY: all
all: # Run all recipes.
all: lint test

.PHONY: help
help: # Show help information.
	@grep --extended-regexp "^[a-z-]+: #" "$(MAKEFILE_LIST)" | \
		awk 'BEGIN {FS = ": # "}; {printf "%-10s%s\n", $$1, $$2}'

.PHONY: lint
lint: # Lint the source code.
	@golangci-lint run

.PHONY: test
test: \
	export MISE_ENV=production
test: # Run tests.
	@mise exec -- go test -v ./...
