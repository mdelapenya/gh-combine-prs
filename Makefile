BINARY_NAME := $(shell basename $(shell pwd))

.PHONY: build
build:
	@echo "Building..."
	@go build -o $(BINARY_NAME) -v

.PHONY: help
help: build
	@gh combine-prs --help

.PHONY: fallback
fallback: build
	@gh combine-prs

.PHONY: query
query: build
	@gh combine-prs --query "author:app/dependabot"

.PHONY: query-interactive
query-interactive: build
	@gh combine-prs --query "author:app/dependabot" --interactive
