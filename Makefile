CLIENT_IMAGE ?= akitasoftware/demo-client
SERVER_IMAGE ?= akitasoftware/demo-server
CONFIG_FILE ?= application.yml
TAG ?= latest
LATEST ?= true

BUILDER=buildx-multi-arch

INFO_COLOR = \033[0;36m
NO_COLOR   = \033[m

build-client: ## Build the demo client
	docker build --tag=$(CLIENT_IMAGE):$(TAG) --secret id=application.yml,src=$(CONFIG_FILE) -f client/Dockerfile client
.PHONY: build-client

build-server: ## Build the demo server
	docker build --tag=$(SERVER_IMAGE):$(TAG) -f server/Dockerfile server
.PHONY: build-server

build-images: build-client build-server ## Build the demo images

prepare-buildx: ## Create buildx builder for multi-arch build, if not exists
	docker buildx inspect $(BUILDER) || docker buildx create --name=$(BUILDER) --driver=docker-container --driver-opt=network=host
.PHONY: prepare-buildx

# Check if the specified tag exists remotely
define check_remote_tag
	@if docker pull $(1):$(TAG) >/dev/null 2>&1; then \
		echo "Error: Image $(1):$(TAG) already exists in the registry."; \
		exit 1; \
	fi
endef

push-client: prepare-buildx ## Push the demo client image to the registry
	$(call check_remote_tag,$(CLIENT_IMAGE))
ifeq ($(LATEST),true)
	docker buildx build \
		--push \
		--builder=$(BUILDER) \
		--platform=linux/amd64,linux/arm64 \
		--secret id=application.yml,src=$(CONFIG_FILE) \
		--build-arg TAG=$(TAG) \
		--tag=$(CLIENT_IMAGE):$(TAG) \
		--tag=$(CLIENT_IMAGE):latest \
		-f client/Dockerfile client
else
	docker buildx build \
		--push \
		--builder=$(BUILDER) \
		--platform=linux/amd64,linux/arm64 \
		--secret id=application.yml,src=$(CONFIG_FILE) \
		--build-arg TAG=$(TAG) \
		--tag=$(CLIENT_IMAGE):$(TAG) \
		-f client/Dockerfile client
endif
.PHONY: push-client

push-server: prepare-buildx
	$(call check_remote_tag,$(SERVER_IMAGE))
ifeq ($(LATEST),true)
	docker buildx build \
		--push \
		--builder=$(BUILDER) \
		--platform=linux/amd64,linux/arm64 \
		--build-arg TAG=$(TAG) \
		--tag=$(SERVER_IMAGE):$(TAG) \
		--tag=$(SERVER_IMAGE):latest \
		-f server/Dockerfile server
else
	docker buildx build \
		--push \
		--builder=$(BUILDER) \
		--platform=linux/amd64,linux/arm64 \
		--build-arg TAG=$(TAG) \
		--tag=$(SERVER_IMAGE):$(TAG) \
		-f server/Dockerfile server
endif
.PHONY: push-server

push-images: push-client push-server ## Push the demo images to the registry
.PHONY: push-images

help: ## Show this help
	@echo Please specify a build target. The choices are:
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "$(INFO_COLOR)%-30s$(NO_COLOR) %s\n", $$1, $$2}'
.PHONY: help
