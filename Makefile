VERSION=`git rev-parse HEAD`
BUILD=`date +%FT%T%z`
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"
DOCKER_REPO=docker.pkg.github.com/codyseavey/test-app/web
DOCKER_TAG=latest
GO_OPTS="GO111MODULE=on GOFLAGS='-mod=vendor'"

.PHONY: help
help: ## - Displays help message
	@printf "\033[32m\xE2\x9c\x93 usage: make [target]\n\n\033[0m\n\033[0m"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## - Docker build
	@docker build -t $(DOCKER_REPO):$(DOCKER_TAG) .

.PHONY: build-no-cache
build-no-cache: ## - Docker build with no-cache setting
	@docker build --no-cache -t $(DOCKER_REPO):$(DOCKER_TAG) .

.PHONY: ls
ls: ## - List images that were created
	@docker image ls $(DOCKER_REPO)

.PHONY: test
test: ## - Run all the tests
	$(GO_OPTS) go test $(shell go list ./... | grep -v /vendor/ | grep -v /pkg/ | grep -v /hack)

.PHONY: push
push: ## - Pushes the image to docker registry
	@docker push $(DOCKER_REPO):$(DOCKER_TAG)

.PHONY: code-gen
code-gen: ## - Runs the code-gen for the custom kubernetes api as well as pulling the code for vendoring our dependecies
	@$(GO_OPTS) go mod vendor
	@$(GO_OPTS) go mod tidy