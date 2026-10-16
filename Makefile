.DEFAULT_GOAL := help

.PHONY: help
help:  ## List the available commands
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z0-9_-]+:.*## / {printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: toolchain
toolchain:  ## Check that the build tools are installed
	@scripts/check-toolchain.sh

.PHONY: ci
ci: toolchain  ## Full local CI
