all: build container

build: *.go
	godep go build -o build/dialogue

container: build Dockerfile
	docker build --no-cache -t ehazlett/dialogue-cli .
