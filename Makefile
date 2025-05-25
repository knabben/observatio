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
BUILD_PATH="../webserver/internal/web/handlers/build/"

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

.PHONY: lint-backend
lint-backend: ## Lint codebase
	docker run --rm -v $(PWD)/webserver:/app -w /app golangci/golangci-lint golangci-lint run -v --fix

.PHONY: lint-frontend
lint-frontend: ## Lint codebase
	cd front; \
	npm run lint


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
run-tests-backend:  ## Run backend tests
	pushd webserver; \
	go test ./... -v -cover ;\
	popd;

.PHONY: run-frontend
run-frontend:  ## Run the frontend locally
	pushd front; \
	pnpm run dev; \
	popd;

.PHONY: run-tests-frontend
run-tests-frontend:  ## Run frontend tests
	pushd front; \
	npm run test; \
	popd;

## --------------------------------------
## Building
## --------------------------------------
##@ build:

.PHONY: build
build:  ## Build the full stack
	pushd front; \
		npm run build; \
		find ${BUILD_PATH}  ! -name 'index.html' ! -name 'build' -type "f,d" -exec rm -fr {} +; \
		mv output/* ${BUILD_PATH} ; \
	popd;
	pushd webserver; \
		CGO_ENABLED=0 go build -o ${BINARY_PATH} .; \
	popd;

