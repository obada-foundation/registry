
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
	golangci-lint --config .golangci.yml run --print-issued-lines --out-format=github-actions ./...

test:
	go test ./... -v -cover

vendor:
	go mod tidy && go mod vendor

coverage:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

swagger:
	cd src && swag fmt
	cd src && swag init -g main.go
