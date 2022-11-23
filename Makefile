MAKEFLAGS += --warn-undefined-variables
SHELL := bash
.SHELLFLAGS := -euo pipefail -c
.DEFAULT_GOAL := all

RUN_GO_GROUPS := go run oss.indeed.com/go/go-groups@v1.1.3
RUN_GO_OPINE := go run oss.indeed.com/go/go-opine@v1.3.0
RUN_GOLANGCI_LINT := go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.49.0

.PHONY: deps
deps: ## download go modules
	go mod download

.PHONY: fmt
fmt: ## ensure consistent code style
	$(RUN_GO_GROUPS) -w .
	$(RUN_GOLANGCI_LINT) run --fix > /dev/null 2>&1 || true
	go mod tidy

.PHONY: lint
lint: ## run golangci-lint
	$(RUN_GOLANGCI_LINT) run
	@if [ -n "$$($(RUN_GO_GROUPS) -l .)" ]; then \
		echo -e "\033[0;33mdetected fmt problems: run \`\033[0;32mmake fmt\033[0m\033[0;33m\`\033[0m"; \
		exit 1; \
	fi

.PHONY: test
test: lint ## run go tests
	$(RUN_GO_OPINE) test -coverprofile=cover.out -junit=report.xml

.PHONY: all
all: test

.PHONY: help
help: ## displays this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-12s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
