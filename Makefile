BUILD_COMMIT := $(shell git rev-parse HEAD)
BUILD_DIRTY := $(if $(shell git status --porcelain),+CHANGES)
BUILD_COMMIT_FLAG := github.com/rasorp/attila/internal/version.BuildCommit=$(BUILD_COMMIT)$(BUILD_DIRTY)

BUILD_TIME ?= $(shell TZ=UTC0 git show -s --format=%cd --date=format-local:'%Y-%m-%dT%H:%M:%SZ' HEAD)
BUILD_TIME_FLAG := github.com/rasorp/attila/internal/version.BuildTime=$(BUILD_TIME)

# Populate the ldflags using the Git commit information and and build time
# which will be present in the binary version output.
GO_LDFLAGS = -X $(BUILD_COMMIT_FLAG) -X $(BUILD_TIME_FLAG)

bin/%/attila: GO_OUT ?= $@
bin/%/attila: ## Build Attila for GOOS & GOARCH; eg. bin/linux_amd64/attila
	@echo "==> Building $@..."
	@GOOS=$(firstword $(subst _, ,$*)) \
		GOARCH=$(lastword $(subst _, ,$*)) \
		go build \
		-o $(GO_OUT) \
		-trimpath \
		-ldflags "$(GO_LDFLAGS)" \
		internal/cmd/cmd.go
	@echo "==> Done"

.PHONY: build
build: ## Build a development version of Attila
	@echo "==> Building Attila..."
	@go build \
		-o ./bin/attila \
		-trimpath \
		-ldflags "$(GO_LDFLAGS)" \
		internal/cmd/cmd.go
	@echo "==> Done"

HELP_FORMAT="    \033[36m%-25s\033[0m %s\n"
.PHONY: help
help: ## Display this usage information
	@echo "Valid targets:"
	@grep -E '^[^ ]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		sort | \
		awk 'BEGIN {FS = ":.*?## "}; \
			{printf $(HELP_FORMAT), $$1, $$2}'
	@echo ""

.PHONY: lint
lint: ## Lint the Attila code
	@echo "==> Linting Attila source code..."
	@golangci-lint run -c ./build/lint/golangci.yaml ./...
	@echo "==> Done"

	@echo "==> License copywrite check of Attila source code..."
	@copywrite --config build/license/copywrite.hcl headers --plan
	@echo "==> Done"

	@echo "==> Running gopls moderniztion analysis..."
	@go run \
	    golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@v0.20.0 -fix -test ./...
	@if (git status -s | grep -q -e '\.go$$'); then \
	    echo modernize analysis found corrections:; git status -s | grep -e '\.go$$'; exit 1; fi
	@echo "==> Done"

.PHONY: lint-deps
lint-deps: ## Install Attila lint dependencies
	@echo "==> Installing lint dependencies..."
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6
	@go install github.com/hashicorp/copywrite@v0.22.0
	@echo "==> Done"

.PHONY: test
test: ## Test the Attila code
	@echo "==> Testing Attila source code..."
	@gotestsum --format pkgname -- -race -cover ./...
	@echo "==> Done"

.PHONY: test-deps
test-deps: ## Install Attila test dependencies
	@echo "==> Installing test dependencies..."
	@go install gotest.tools/gotestsum@v1.12.0
	@echo "==> Done"
