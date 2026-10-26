.DEFAULT_GOAL := help

MARKDOWNLINT := npx --yes markdownlint-cli2@0.23.3

.PHONY: help
help:  ## List the available commands
	@awk 'BEGIN {FS = ":.*## "} /^##@ / {printf "\n\033[1m%s\033[0m\n", substr($$0, 5)} /^[a-zA-Z0-9_-]+:.*## / {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

##@ The bank

.PHONY: pangaea-fixture
pangaea-fixture:  ## Run the bank's daily batch on its fixture and show its reports
	@scripts/pangaea-fixture.sh

.PHONY: pangaea-open
pangaea-open: pangaea-build  ## Open the running bank on generated accounts (SEED, ACCOUNTS)
	go run ./pangaea/cmd/pangaea open -seed $(SEED) -accounts $(ACCOUNTS)

.PHONY: pangaea-days
pangaea-days: pangaea-build  ## Run the bank's next N business days
	go run ./pangaea/cmd/pangaea days -n $(N)

# The ports are goroawase, Japanese number puns: 4510 shigoto, "work";
# 2951 fukukoi, "come, fortune"; 4649 yoroshiku, "I leave it in your
# hands"; 2525 niko niko, "all smiles"; 1192 ii kuni, "good country".
PANGAEA_ADDR ?= 127.0.0.1:4510

.PHONY: pangaea-console
pangaea-console: pangaea-build  ## Serve the bank's console in the browser (PANGAEA_ADDR)
	go run ./pangaea/cmd/console -addr $(PANGAEA_ADDR)

ONLINE_ADDR ?= 127.0.0.1:2951

.PHONY: pangaea-online
pangaea-online: pangaea-build  ## Serve the bank's online day alone until interrupted (ONLINE_ADDR)
	go run ./pangaea/cmd/pangaea online -addr $(ONLINE_ADDR)

RATE ?= 5

.PHONY: pangaea-simulate
pangaea-simulate:  ## Send card traffic to the bank console's online channel until interrupted (RATE)
	go run ./pangaea/cmd/pangaea simulate -api http://$(PANGAEA_ADDR)/api -rate $(RATE)

# Compiles Pangaea's programs for the targets that run the bank; not listed.
.PHONY: pangaea-build
pangaea-build:
	go run ./pangaea/cmd/pangaea build

##@ The card processor

.PHONY: laurasia-fixture
laurasia-fixture:  ## Run the card processor's night on its fixture and show its logs and reports
	@scripts/laurasia-fixture.sh

.PHONY: laurasia-open
laurasia-open: laurasia-build  ## Open the running card processor on generated cards (SEED, CARDS)
	go run ./laurasia/cmd/laurasia open -seed $(SEED) -cards $(CARDS)

.PHONY: laurasia-days
laurasia-days: laurasia-build  ## Run the card processor's next N days
	go run ./laurasia/cmd/laurasia days -n $(N)

LAURASIA_ADDR ?= 127.0.0.1:4649

.PHONY: laurasia-console
laurasia-console: laurasia-build  ## Serve the card processor's console in the browser (LAURASIA_ADDR)
	go run ./laurasia/cmd/console -addr $(LAURASIA_ADDR)

# Compiles Laurasia's programs for the targets that run the processor; not listed.
.PHONY: laurasia-build
laurasia-build:
	go run ./laurasia/cmd/laurasia build

##@ The life insurer

.PHONY: gondwana-fixture
gondwana-fixture:  ## Run the life insurer's March on its fixture and show its logs and reports
	@scripts/gondwana-fixture.sh

.PHONY: gondwana-open
gondwana-open: gondwana-build  ## Open the running life insurer on generated policies (SEED, POLICIES)
	go run ./gondwana/cmd/gondwana open -seed $(SEED) -policies $(POLICIES)

.PHONY: gondwana-days
gondwana-days: gondwana-build  ## Run the life insurer's next N business days
	go run ./gondwana/cmd/gondwana days -n $(N)

GONDWANA_ADDR ?= 127.0.0.1:2525

.PHONY: gondwana-console
gondwana-console: gondwana-build  ## Serve the life insurer's console in the browser (GONDWANA_ADDR)
	go run ./gondwana/cmd/console -addr $(GONDWANA_ADDR)

# Compiles Gondwana's programs for the targets that run the insurer; not listed.
.PHONY: gondwana-build
gondwana-build:
	go run ./gondwana/cmd/gondwana build

##@ The survey

.PHONY: blueprints
blueprints:  ## Regenerate the blueprints of the three specimens with Survey
	rm -rf blueprints/pangaea blueprints/laurasia blueprints/gondwana
	go run ./survey/cmd/survey blueprints -copybooks pangaea/copybooks -out blueprints/pangaea -annotations annotations/pangaea -title pangaea pangaea/programs/*.cbl pangaea/jobs/*.jcl
	go run ./survey/cmd/survey blueprints -dialect mf -copybooks laurasia/copybooks -out blueprints/laurasia -annotations annotations/laurasia -title laurasia laurasia/programs/*.cbl laurasia/scripts/*.sh
	go run ./survey/cmd/survey blueprints -dialect fujitsu -copybooks gondwana/copybooks -out blueprints/gondwana -annotations annotations/gondwana -title gondwana gondwana/programs/*.cbl gondwana/scripts/*.BAT

EXPLORE_ADDR ?= 127.0.0.1:1192

.PHONY: explore-pangaea
explore-pangaea:  ## Explore Pangaea's analysis in the browser, with the running bank's runs (EXPLORE_ADDR)
	go run ./survey/cmd/survey explore -copybooks pangaea/copybooks -addr $(EXPLORE_ADDR) -title "the First Bank of Pangaea" -annotations annotations/pangaea -cases oracle/pangaea -runs pangaea/run/bank/RUNS pangaea/programs/*.cbl pangaea/jobs/*.jcl

.PHONY: explore-laurasia
explore-laurasia:  ## Explore Laurasia's analysis in the browser, with the running processor's runs (EXPLORE_ADDR)
	go run ./survey/cmd/survey explore -dialect mf -copybooks laurasia/copybooks -addr $(EXPLORE_ADDR) -title "Laurasia Card Services" -annotations annotations/laurasia -cases oracle/laurasia -runs laurasia/run/processor/RUNS laurasia/programs/*.cbl laurasia/scripts/*.sh

.PHONY: explore-gondwana
explore-gondwana:  ## Explore Gondwana's analysis in the browser, with the running insurer's runs (EXPLORE_ADDR)
	go run ./survey/cmd/survey explore -dialect fujitsu -copybooks gondwana/copybooks -addr $(EXPLORE_ADDR) -title "Gondwana Life Insurance" -annotations annotations/gondwana -cases oracle/gondwana -runs gondwana/run/insurer/RUNS gondwana/programs/*.cbl gondwana/scripts/*.BAT

.PHONY: oracle
oracle:  ## Search the specimens' programs again for cases that reach their rules, and rewrite oracle/
	go run ./survey/cmd/survey oracle -copybooks pangaea/copybooks -out oracle/pangaea -rules ../../blueprints/pangaea/rules -title pangaea -jobs START,ONLINE,DAILY,MONTHLY pangaea/programs/*.cbl pangaea/jobs/*.jcl
	go run ./survey/cmd/survey oracle -dialect mf -copybooks laurasia/copybooks -out oracle/laurasia -rules ../../blueprints/laurasia/rules -title laurasia -jobs open,nightly,cycle laurasia/programs/*.cbl laurasia/scripts/*.sh
	go run ./survey/cmd/survey oracle -dialect fujitsu -copybooks gondwana/copybooks -out oracle/gondwana -rules ../../blueprints/gondwana/rules -title gondwana -jobs OPEN,DAILY,BILLING,RESULTS,CLOSE,ANNUAL gondwana/programs/*.cbl gondwana/scripts/*.BAT

.PHONY: conform-pangaea
conform-pangaea:  ## Run a replacement for one of Pangaea's programs on its recorded cases (PROGRAM, WITH)
	go run ./survey/cmd/survey conform -copybooks pangaea/copybooks -cases oracle/pangaea -program $(PROGRAM) -with $(WITH) pangaea/programs/*.cbl pangaea/jobs/*.jcl

.PHONY: conform-laurasia
conform-laurasia:  ## Run a replacement for one of Laurasia's programs on its recorded cases (PROGRAM, WITH)
	go run ./survey/cmd/survey conform -dialect mf -copybooks laurasia/copybooks -cases oracle/laurasia -program $(PROGRAM) -with $(WITH) laurasia/programs/*.cbl laurasia/scripts/*.sh

.PHONY: conform-gondwana
conform-gondwana:  ## Run a replacement for one of Gondwana's programs on its recorded cases (PROGRAM, WITH)
	go run ./survey/cmd/survey conform -dialect fujitsu -copybooks gondwana/copybooks -cases oracle/gondwana -program $(PROGRAM) -with $(WITH) gondwana/programs/*.cbl gondwana/scripts/*.BAT

.PHONY: conform-demo
conform-demo:  ## Check BANKDATE rewritten in Go against the bank's recorded cases
	go build -buildvcs=false -o build/bankdate ./survey/internal/oracle/testdata/rewrites/bankdate
	$(MAKE) conform-pangaea PROGRAM=BANKDATE WITH=build/bankdate

##@ Golden files

.PHONY: pangaea-golden
pangaea-golden:  ## Rewrite Pangaea's golden outputs from a run on the fixture
	go test ./pangaea -run TestDaily -update

.PHONY: laurasia-golden
laurasia-golden:  ## Rewrite Laurasia's golden outputs from a run of its fixtures
	go test ./laurasia -update

.PHONY: gondwana-golden
gondwana-golden:  ## Rewrite Gondwana's golden outputs from a run of its fixtures
	go test ./gondwana -update

SEED ?= 1
ACCOUNTS ?= 3000
CARDS ?= 1000
POLICIES ?= 2000
N ?= 1

##@ Checks

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

.PHONY: ci
ci: toolchain lint-md lint-go lint-deps test  ## Full local CI

##@ Docs

.PHONY: docs-diagrams
docs-diagrams:  ## Redraw the diagrams of docs/architecture.md
	go run ./survey/cmd/docsdiagrams -out docs/diagrams
