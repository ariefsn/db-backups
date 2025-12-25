APP_NAME := db-backup
DOCKER_IMAGE := ariefsn/$(APP_NAME)
TAG := latest

.PHONY: all build run docker-build docker-push clean

all: build

build:
	go build -o $(APP_NAME) cmd/server/main.go

run: build
	./$(APP_NAME)

docker-build:
	docker build -t $(DOCKER_IMAGE):$(TAG) .

docker-push:
	docker buildx build --platform linux/amd64,linux/arm64 -t $(DOCKER_IMAGE):$(TAG) --push .

clean:
	rm -f $(APP_NAME)
