BINARY_NAME := $(shell basename $(shell pwd))

.PHONY: build
build:
	@echo "Building..."
	@go build -o $(BINARY_NAME) -v
	@gh multi-merge-prs
