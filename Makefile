BINARY_NAME := $(shell basename $(shell pwd))

.PHONY: build
build:
	@echo "Building..."
	@go build -o $(BINARY_NAME) -v

.PHONY: help
help: build
	@gh multi-merge-prs --help

.PHONY: fallback
fallback: build
	@gh multi-merge-prs

.PHONY: query
query: build
	@gh multi-merge-prs --query "author:app/dependabot"

.PHONY: query-interactive
query-interactive: build
	@gh multi-merge-prs --query "author:app/dependabot" --interactive
