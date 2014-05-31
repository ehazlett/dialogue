all: build container

build:
	@godep restore
	@godep go build -o build/dialogue

container: build Dockerfile
	docker build --no-cache -t ehazlett/dialogue-cli .


.PHONY: all build container
