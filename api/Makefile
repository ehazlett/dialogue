all: build container

build: *.go
	godep go build -o build/api

container: build Dockerfile
	docker build --no-cache -t ehazlett/dialogue-api .
