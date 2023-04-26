export GO111MODULE=on

DOCKER_IMAGE=obada/registry

BUILD_TAGS=-tags

# Various build flags
V_LDFLAGS_SYMBOL := -s
V_LDFLAGS_BUILD := -X "google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn"
V_LDFLAGS_COMMON := ${V_LDFLAGS_SYMBOL} ${V_LDFLAGS_BUILD}
V_LDFLAGS_STATIC := ${V_LDFLAGS_COMMON}

GO ?= go

docker: docker/build docker/push
.PHONY: docker


docker/push:
	docker push $(DOCKER_IMAGE)

docker/build:
	docker \
		build \
		-t $(DOCKER_IMAGE) \
		-f docker/Dockerfile \
		.

bin/registry:
	CGO_ENABLED=0 $(GO) build $(BUILD_TAGS) -a -ldflags '$(V_LDFLAGS_STATIC)' ./
.PHONY: bin/registry

lint:
	golangci-lint --config .golangci.yml run --print-issued-lines --out-format=github-actions ./...

test:
	go test ./... -v -cover

vendor:
	go mod tidy && go mod vendor
.PHONY: vendor

coverage:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

swagger:
	docker run \
		-p 80:8080 \
		-e SWAGGER_JSON=/openapi/api.swagger.json \
		-v $$(pwd)/openapi/:/openapi \
		swaggerapi/swagger-ui

mockgen:
	./scripts/mockgen.sh
.PHONY: vendor
