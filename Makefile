# Copyright 2025 Amim Knabben
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

SHELL := /usr/bin/env bash
.DEFAULT_GOAL := help

BINARY_PATH="${PWD}/output/observatio"

## --------------------------------------
## Help
## --------------------------------------
##@ help:

help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

## --------------------------------------
## Linters
## --------------------------------------
##@ lint:

.PHONY: lint-go
lint-go: ## Lint codebase
	docker run --rm -v $(PWD):/app -w /app -it golangci/golangci-lint golangci-lint run -v --fix

## --------------------------------------
## Development targets
## --------------------------------------
##@ run-backend:

.PHONY: run-backend
run-backend:  ## Build the binary using local golang
	pushd webserver; \
	go run . $(what) --dev; \
	popd;

.PHONY: run-tests-backend
run-tests-backend:
	pushd webserver; \
	go test ./... ;\
	popd;

.PHONY: run-frontend
run-frontend:  ## Run the frontend locally
	pushd front; \
	pnpm run dev; \
	popd;

## --------------------------------------
## Building
## --------------------------------------
##@ build:

.PHONY: build
build:
	pushd front; \
	pnpm run build; \
	popd;
	pushd webserver; \
	go build -o ${BINARY_PATH} .; \
	popd;
