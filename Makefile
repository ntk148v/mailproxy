# Use git tag/git branch to tag Docker image.
CURRENT_DIR             ?= $(shell pwd)
DOCKER_IMAGE_TAG        ?= $(subst /,-,$(shell git describe --tags --abbrev=0 || git rev-parse --abbrev-ref HEAD))
DOCKER_REPO             ?= kiennt26
DOCKER_IMAGE_NAME       ?= mailproxy
DOCKER_IMAGE_FULL       ?= $(DOCKER_REPO)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)
DOCKER_IMAGE_LATEST     ?= $(DOCKER_REPO)/$(DOCKER_IMAGE_NAME):latest
DOCKER_CONTAINER_NAME   ?= mailproxy
MAILPROXY_PORT             ?= 9011
MAILPROXY_CONF_DIR        ?= $(CURRENT_DIR)/etc

build:
	docker build -t "$(DOCKER_IMAGE_FULL)" .

build-latest: build
	docker tag "$(DOCKER_IMAGE_FULL)" "$(DOCKER_IMAGE_LATEST)"

push: build
	docker push "$(DOCKER_IMAGE_FULL)"

push-latest: build-latest
	docker push "$(DOCKER_IMAGE_LATEST)"

run:
	docker rm -f "$(DOCKER_CONTAINER_NAME)" || true
	docker run -d --restart always --net host -v "$(MAILPROXY_CONF_DIR)":/etc/mailproxy --name "$(DOCKER_CONTAINER_NAME)" "$(DOCKER_IMAGE_FULL)"
