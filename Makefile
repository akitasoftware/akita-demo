CLIENT_IMAGE ?= akitasoftware/demo-client
SERVER_IMAGE ?= akitasoftware/demo-server
TAG ?= dev
LATEST ?= false

BUILDER=buildx-multi-arch

INFO_COLOR = \033[0;36m
NO_COLOR   = \033[m

run-demo: build-client build-server ## Run the demo
	DEMO_IMAGE_TAG=$(TAG) docker compose up -d --always-recreate-deps
.PHONY: run-demo

build-client: ## Build the demo client
	docker build --tag=$(CLIENT_IMAGE):$(TAG) -f client/Dockerfile client
.PHONY: build-client

build-server: ## Build the demo server
	docker build --tag=$(SERVER_IMAGE):$(TAG) -f server/Dockerfile server
.PHONY: build-server

prepare-buildx: ## Create buildx builder for multi-arch build, if not exists
	docker buildx inspect $(BUILDER) || docker buildx create --name=$(BUILDER) --driver=docker-container --driver-opt=network=host
.PHONY: prepare-buildx

check-tag: ## Check if the tag already exists for either the client or server images. If so, fail.
	(docker pull $(CLIENT_IMAGE):$(TAG) || docker pull $(SERVER_IMAGE):$(TAG)) && echo "Failure: Tag already exists" && false || true
.PHONY: check-tag

push-client: check-tag prepare-buildx ## Push the demo client image to the registry
ifeq ($(LATEST),true)
	docker buildx build \
		--push \
		--builder=$(BUILDER) \
		--platform=linux/amd64,linux/arm64 \
		--build-arg TAG=$(TAG) \
		--tag=$(CLIENT_IMAGE):$(TAG) \
		--tag=$(CLIENT_IMAGE):latest \
		-f client/Dockerfile client
else
	docker buildx build \
		--push \
		--builder=$(BUILDER) \
		--platform=linux/amd64,linux/arm64 \
		--build-arg TAG=$(TAG) \
		--tag=$(CLIENT_IMAGE):$(TAG) $(LATEST_TAG) \
		-f client/Dockerfile client
endif
.PHONY: push-client

push-server: check-tag prepare-buildx
ifeq ($(LATEST),true)
	docker buildx build \
		--push \
		--builder=$(BUILDER) \
		--platform=linux/amd64,linux/arm64 \
		--build-arg TAG=$(TAG) \
		--tag=$(SERVER_IMAGE):$(TAG) \
		--tag=$(SERVER_IMAGE):latest \
		-f client/Dockerfile client
else
	docker buildx build \
		--push \
		--builder=$(BUILDER) \
		--platform=linux/amd64,linux/arm64 \
		--build-arg TAG=$(TAG) \
		--tag=$(SERVER_IMAGE):$(TAG) $(LATEST_TAG) \
		-f client/Dockerfile client
endif
.PHONY: push-server

push-images: push-client push-server ## Push the demo images to the registry

help: ## Show this help
	@echo Please specify a build target. The choices are:
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "$(INFO_COLOR)%-30s$(NO_COLOR) %s\n", $$1, $$2}'

.PHONY: help
