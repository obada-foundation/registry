DOCKER_IMAGE=obada/registry

docker: docker/build docker/push

docker/push:
	docker push $(DOCKER_IMAGE)

docker/build:
	docker \
		build \
		-t $(DOCKER_IMAGE) \
		-f docker/Dockerfile \
		.
