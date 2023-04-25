CLIENT_IMAGE ?= akitasoftware/demo-client
SERVER_IMAGE ?= akitasoftware/demo-server
TAG ?= dev
LATEST ?= false

BUILDER=buildx-multi-arch

run-demo: build-client build-server ## Run the demo
	docker compose up -d --always-recreate-deps
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

push-client: check-tag prepare-buildx
	LATEST_TAG=''
ifeq ($(LATEST),true)
	LATEST_TAG='--tag=$(CLIENT_IMAGE):latest'
endif
	docker buildx build \
		--push \
		--builder=$(BUILDER) \
		--platform=linux/amd64,linux/arm64 \
		--build-arg TAG=$(TAG)
		--tag=$(CLIENT_IMAGE):$(TAG) \
		$(LATEST_TAG) \
		-f client/Dockerfile client
.PHONY: push-client

push-server: check-tag prepare-buildx
	LATEST_TAG=''
ifeq ($(LATEST),true)
	LATEST_TAG='--tag=$(SERVER_IMAGE):latest'
endif
	docker buildx build \
		--push \
		--builder=$(BUILDER) \
		--platform=linux/amd64,linux/arm64 \
		--build-arg TAG=$(TAG)
		--tag=$(SERVER_IMAGE):$(TAG) \
		$(LATEST_TAG) \
		-f server/Dockerfile server
.PHONY: push-server

push-images: push-client push-server ## Push the demo images to the registry
