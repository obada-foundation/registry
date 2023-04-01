export GO111MODULE=on

DOCKER_IMAGE=obada/registry

BUILD_TAGS=-tags

# Various build flags
V_LDFLAGS_SYMBOL := -s
V_LDFLAGS_COMMON := ${V_LDFLAGS_SYMBOL}
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
	cd src && swag fmt
	cd src && swag init -g main.go
