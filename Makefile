.DEFAULT_GOAL := help

MARKDOWNLINT := npx --yes markdownlint-cli2@0.23.3

.PHONY: help
help:  ## List the available commands
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z0-9_-]+:.*## / {printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: toolchain
toolchain:  ## Check that the build tools are installed
	@scripts/check-toolchain.sh

.PHONY: lint-md
lint-md:  ## Lint the Markdown files
	@$(MARKDOWNLINT) "**/*.md" "!CLAUDE.md"

.PHONY: pangaea-build
pangaea-build:  ## Compile Pangaea's COBOL programs
	go run ./cmd/pangaea build

.PHONY: ci
ci: toolchain lint-md  ## Full local CI
