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

BINARY_PATH       := $(PWD)/output/observatio
FRONTEND_DIR      := front
FRONTEND_OUT      := $(FRONTEND_DIR)/output
EMBED_TARGET      := webserver/internal/web/handlers/build
PREREQS_SCRIPT    := scripts/check-prereqs.sh

## --------------------------------------
## Help
## --------------------------------------
##@ help:

help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

## --------------------------------------
## Prerequisites
## --------------------------------------
##@ prereqs:

.PHONY: check-prereqs
check-prereqs: ## Validate all required tools (go ≥1.24, node ≥22, pnpm)
	@bash $(PREREQS_SCRIPT)

## --------------------------------------
## Linters
## --------------------------------------
##@ lint:

.PHONY: lint-backend
lint-backend: ## Lint backend codebase
	docker run --rm -v $(PWD)/webserver:/app -w /app golangci/golangci-lint golangci-lint run -v --fix

.PHONY: lint-frontend
lint-frontend: ## Lint frontend codebase
	cd $(FRONTEND_DIR) && pnpm run lint

## --------------------------------------
## Development
## --------------------------------------
##@ development:

.PHONY: run-backend
run-backend: ## Run backend server in development mode (no static hosting)
	pushd webserver && go run . $(what) --dev && popd

.PHONY: run-frontend
run-frontend: ## Run frontend dev server (http://localhost:3000)
	pushd $(FRONTEND_DIR) && pnpm run dev && popd

## --------------------------------------
## Tests
## --------------------------------------
##@ test:

.PHONY: run-tests-backend
run-tests-backend: ## Run backend test suite
	pushd webserver && go test ./... -v -cover && popd

.PHONY: run-tests-frontend
run-tests-frontend: ## Run frontend test suite
	pushd $(FRONTEND_DIR) && pnpm run test && popd

.PHONY: test
test: run-tests-backend run-tests-frontend ## Run all tests (backend + frontend)

## --------------------------------------
## Build — independent stages
## --------------------------------------
##@ build:

.PHONY: build-frontend
build-frontend: ## Build frontend bundle independently (requires node ≥22 + pnpm)
	@bash $(PREREQS_SCRIPT) --node
	pushd $(FRONTEND_DIR) && pnpm install --frozen-lockfile && pnpm run build && popd

.PHONY: build-backend
build-backend: ## Build backend binary independently (requires go ≥1.24; run build-frontend first for full embed)
	@bash $(PREREQS_SCRIPT) --go
	mkdir -p output
	pushd webserver && CGO_ENABLED=0 go build -o $(BINARY_PATH) . && popd

.PHONY: build
build: check-prereqs build-frontend ## Build full stack — validates prereqs, builds frontend, copies assets, compiles binary
	@echo "[build] Copying frontend assets to embed target..."
	rm -rf $(EMBED_TARGET)
	mkdir -p $(EMBED_TARGET)
	cp -r $(FRONTEND_OUT)/. $(EMBED_TARGET)/
	@echo "[build] Compiling backend binary..."
	mkdir -p output
	pushd webserver && CGO_ENABLED=0 go build -o $(BINARY_PATH) . && popd
	@echo "[build] Done — binary at $(BINARY_PATH)"
=======
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
lint-frontend: ## Lint frontend codebase
	cd front && pnpm run lint


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
run-tests-frontend: ## Run frontend tests
	cd front && pnpm run test

## --------------------------------------
## Building
## --------------------------------------
##@ build:

.PHONY: build
build: ## Build the full stack (frontend bundle embedded in Go binary)
	cd front && pnpm install --frozen-lockfile && pnpm run build
	rm -rf ${BUILD_PATH} && mkdir -p ${BUILD_PATH}
	cp -r front/output/. ${BUILD_PATH}/
	pushd webserver && CGO_ENABLED=0 go build -o ${BINARY_PATH} . && popd

