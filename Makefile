GOLANGCI_LINT := $(GOPATH)/bin/golangci-lint
GOLANGCI_LINT_VERSION := 1.60.3

# Ensure that the shell is bash.
SHELL := /bin/bash

.PHONY: lint
lint:
	@echo "Checking for golangci-lint..."
	@if [ ! -f $(GOLANGCI_LINT) ] || [ $$($(GOLANGCI_LINT) --version | cut -d ' ' -f4) != $(GOLANGCI_LINT_VERSION) ]; then \
		echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION)..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v$(GOLANGCI_LINT_VERSION); \
	else \
		echo "golangci-lint is already installed and up to date."; \
	fi
	@echo "Running golangci-lint..."
	@$(GOLANGCI_LINT) run



