all: build container

build: build/dialogue
	@godep go build -o build/dialogue

container: build Dockerfile
	docker build --no-cache -t ehazlett/dialogue-cli .
