
DOCKER_IMAGE=obada/registry

docker: docker/build docker/push

docker/push:
	docker push $(DOCKER_IMAGE)

docker/build:
	docker \
		build \
		-t $(DOCKER_IMAGE) \
		-f docker/Dockerfile 
		.

lint:
	cd src &&  golangci-lint --config .golangci.yml run --print-issued-lines --out-format=github-actions ./...

test:
	cd src && go test ./... -v

vendor:
	cd src && go mod tidy && go mod vendor

coverage:
	cd src && go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out
