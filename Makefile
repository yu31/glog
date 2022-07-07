# Copyright (C) 2020 Yunify, Inc.

VERBOSE = no

.PHONY: help
help: ## help for command
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_%-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: format
format: ## Execute go fmt ./...
	@[ ${VERBOSE} = "yes" ] && set -x; go fmt ./...;

.PHONY: vet
vet: ## Execute go vet ./...
	@[ ${VERBOSE} = "yes" ] && set -x; go vet ./...;

.PHONY: lint
lint: ## Execute staticcheck ./...
	@[ ${VERBOSE} = "yes" ] && set -x; staticcheck ./...;

.PHONY: tidy
tidy: ## Execute go mod tidy
	@[ ${VERBOSE} = "yes" ] && set -x; go mod tidy;

.PHONY: check
check: ## Execute tidy format vet lint
check: tidy format vet lint

.PHONY: test
test: ## Run test case
test: check
	@[[ ${VERBOSE} = "yes" ]] && set -x; go test -race -v -test.count=1 -failfast ./...;

.PHONY: bench
bench: ## Run benchmark
	@[[ ${VERBOSE} = "yes" ]] && set -x; go test -test.bench="." -test.run="Benchmark" -benchmem -count=1 ./...;

.DEFAULT_GOAL = help

# Target name % means that it is a rule that matches anything, @: is a recipe;
# the : means do nothing
%:
	@: