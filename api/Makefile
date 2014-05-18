all: build container

build: build/api
	@godep go build -o build/api

container: build Dockerfile
	docker build --no-cache -t ehazlett/dialogue-api .
