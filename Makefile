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

.PHONY: lint-deps
lint-deps:  ## Check that the dependencies between the parts run one way
	@scripts/check-deps.sh

.PHONY: test
test:  ## Run the tests
	go test ./...

.PHONY: pangaea-build
pangaea-build:  ## Compile Pangaea's COBOL programs
	go run ./pangaea/cmd/pangaea build

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
	go run ./pangaea/cmd/pangaea reset -seed $(SEED) -accounts $(ACCOUNTS)

.PHONY: bank-days
bank-days: pangaea-build  ## Run the bank's next N business days
	go run ./pangaea/cmd/pangaea days -n $(N)

# The ports are goroawase, Japanese number puns: 4510 shigoto, "work";
# 2951 fukukoi, "come, fortune"; 1192 ii kuni, "good country".
ADDR ?= 127.0.0.1:4510

.PHONY: console
console: pangaea-build  ## Serve the bank console in the browser (ADDR)
	go run ./console/cmd/console -addr $(ADDR)

ONLINE_ADDR ?= 127.0.0.1:2951

.PHONY: online
online: pangaea-build  ## Serve the bank's online day until interrupted (ONLINE_ADDR)
	go run ./pangaea/cmd/pangaea online -addr $(ONLINE_ADDR)

RATE ?= 5

.PHONY: simulate
simulate:  ## Send card traffic to the console's online channel until interrupted (RATE)
	go run ./pangaea/cmd/pangaea simulate -api http://$(ADDR)/api -rate $(RATE)

.PHONY: oracle
oracle:  ## Search Pangaea's programs again for cases that reach their rules, and rewrite oracle/pangaea
	go run ./survey/cmd/survey oracle -copybooks pangaea/copybooks -out oracle/pangaea -rules ../../blueprints/pangaea/rules -title pangaea pangaea/programs/*.cbl

.PHONY: blueprints
blueprints:  ## Regenerate Pangaea's blueprints with Survey
	rm -rf blueprints/pangaea
	go run ./survey/cmd/survey blueprints -copybooks pangaea/copybooks -out blueprints/pangaea -title pangaea pangaea/programs/*.cbl pangaea/jobs/*.jcl

EXPLORE_ADDR ?= 127.0.0.1:1192

.PHONY: explore
explore:  ## Explore Pangaea's analysis in the browser (EXPLORE_ADDR)
	go run ./survey/cmd/survey explore -copybooks pangaea/copybooks -addr $(EXPLORE_ADDR) -title "the First Bank of Pangaea" pangaea/programs/*.cbl pangaea/jobs/*.jcl

.PHONY: docs-diagrams
docs-diagrams:  ## Redraw the diagrams of docs/architecture.md
	go run ./survey/cmd/docsdiagrams -out docs/diagrams

.PHONY: ci
ci: toolchain lint-md lint-go lint-deps test  ## Full local CI
