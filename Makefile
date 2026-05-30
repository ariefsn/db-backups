APP_NAME := db-backup
DOCKER_IMAGE := ariefsn/$(APP_NAME)
TAG := latest

.PHONY: all build run docker-build docker-push docker-push-web tag retag clean

all: build

build:
	go build -o $(APP_NAME) cmd/server/main.go

run: build
	./$(APP_NAME)

docker-build:
	docker build -t $(DOCKER_IMAGE):$(TAG) .

# NOTE: Multi-arch publishing is handled by CI (.github/workflows/docker-publish.yml).
# These targets remain for local one-off pushes; building linux/amd64 locally on
# Apple Silicon goes through QEMU and the Go compiler may segfault — prefer CI.
docker-push:
	docker buildx build --platform linux/amd64,linux/arm64 -t $(DOCKER_IMAGE):$(TAG) --push .

docker-push-web:
	docker buildx build --platform linux/amd64,linux/arm64 -t $(DOCKER_IMAGE)-web:$(TAG) --push --no-cache ./web

# Push a git tag to trigger the CI release build (.github/workflows/docker-publish.yml),
# which builds & pushes BOTH images (db-backup and db-backup-web) for that version.
# Usage: make tag VERSION=1.0.0
tag:
ifndef VERSION
	$(error VERSION is required, e.g. make tag VERSION=1.0.0)
endif
	git tag $(VERSION)
	git push origin $(VERSION)

# Re-push an existing tag, overwriting it locally and on the remote (re-triggers CI).
# Usage: make retag VERSION=1.0.0
retag:
ifndef VERSION
	$(error VERSION is required, e.g. make retag VERSION=1.0.0)
endif
	-git tag -d $(VERSION)
	git tag $(VERSION)
	git push --force origin $(VERSION)

clean:
	rm -f $(APP_NAME)
