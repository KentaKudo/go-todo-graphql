SERVICE := todo-graphql

GIT_HASH := $(shell git rev-parse HEAD)
LINKFLAGS := -X main.gitHash=$(GIT_HASH)

PROTO_DIR := $(GOPATH)/src/github.com/KentaKudo/go-todo-service/proto
GENERATED_DIR := ./internal/pb
GENERATED_SERVICE_DIR := $(GENERATED_DIR)/service

.PHONY: protos
protos:
	mkdir -pv $(GENERATED_DIR) $(GENERATED_SERVICE_DIR)
	protoc \
		-I $(PROTO_DIR) \
		--gogoslick_out=plugins=grpc:$(GENERATED_SERVICE_DIR) \
		service.proto


.PHONY: gql
gql:
	go get github.com/99designs/gqlgen
	gqlgen --verbose


.PHONY: install
install:
	go get -v ./...

LINTER_EXE := golangci-lint
LINTER := $(GOBIN)/$(LINTER_EXE)

$(LINTER):
	go get github.com/golangci/golangci-lint/cmd/golangci-lint

LINT_FLAGS :=--enable golint,unconvert,unparam,gofmt

.PHONY: lint
lint: $(LINTER)
	$(LINTER) run $(LINT_FLAGS)

TEST_FLAGS := -v -cover -timeout 30s

.PHONY: test
test:
	go test $(TEST_FLAGS) ./...

$(SERVICE):
	go build -ldflags '$(LINKFLAGS)' ./cmd/$(SERVICE)

.PHONY: build
build: $(SERVICE)

.PHONY: clean
clean:
	@rm -f $(SERVICE)

.PHONY: all
all: install lint test clean build


DOCKER_ID=kentakudo
DOCKER_REGISTRY=docker.io
DOCKER_REPOSITORY_NAMESPACE=kentakudo
DOCKER_REPOSITORY_IMAGE=$(SERVICE)
DOCKER_REPOSITORY=$(DOCKER_REGISTRY)/$(DOCKER_REPOSITORY_NAMESPACE)/$(DOCKER_REPOSITORY_IMAGE)
DOCKER_IMAGE_TAG=$(GIT_HASH)

.PHONY: docker-image
docker-image:
	docker build -t $(DOCKER_REPOSITORY):$(DOCKER_IMAGE_TAG) . \
	  --build-arg SERVICE=$(SERVICE)

.PHONY: docker-auth
docker-auth:
	@docker login -u $(DOCKER_ID) -p $(DOCKER_PASSWORD) $(DOCKER_REGISTRY)

.PHONY: docker-build
docker-build: docker-image docker-auth
	docker tag $(DOCKER_REPOSITORY):$(DOCKER_IMAGE_TAG) $(DOCKER_REPOSITORY):latest
	docker push $(DOCKER_REPOSITORY)
