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

.PHONY: lint-go
lint-go:  ## Lint the Go code and check its formatting
	golangci-lint run ./...

.PHONY: test
test:  ## Run the tests
	go test ./...

.PHONY: pangaea-build
pangaea-build:  ## Compile Pangaea's COBOL programs
	go run ./cmd/pangaea build

.PHONY: pangaea-daily
pangaea-daily:  ## Run the daily batch on the fixture and show its reports
	@scripts/pangaea-daily.sh

.PHONY: pangaea-golden
pangaea-golden:  ## Rewrite Pangaea's golden outputs from a run on the fixture
	go test ./pangaea -run TestDaily -update

SEED ?= 1
ACCOUNTS ?= 3000
N ?= 1

.PHONY: bank-reset
bank-reset:  ## Reset the running bank to a generated start (SEED, ACCOUNTS)
	go run ./cmd/pangaea reset -seed $(SEED) -accounts $(ACCOUNTS)

.PHONY: bank-days
bank-days: pangaea-build  ## Run the bank's next N business days
	go run ./cmd/pangaea days -n $(N)

ADDR ?= localhost:8080

.PHONY: console
console: pangaea-build  ## Serve the bank console in the browser (ADDR)
	go run ./cmd/pangaea console -addr $(ADDR)

ONLINE_ADDR ?= localhost:8081

.PHONY: online
online: pangaea-build  ## Serve the bank's online day until interrupted (ONLINE_ADDR)
	go run ./cmd/pangaea online -addr $(ONLINE_ADDR)

RATE ?= 5

.PHONY: simulate
simulate:  ## Send card traffic to the console's online channel until interrupted (RATE)
	go run ./cmd/pangaea simulate -api http://$(ADDR)/api -rate $(RATE)

.PHONY: blueprints
blueprints:  ## Regenerate Pangaea's blueprints with Copal
	rm -rf blueprints/pangaea
	go run ./cmd/copal -copybooks pangaea/copybooks -out blueprints/pangaea -title pangaea pangaea/programs/*.cbl pangaea/jobs/*.jcl

.PHONY: ci
ci: toolchain lint-md lint-go test  ## Full local CI
