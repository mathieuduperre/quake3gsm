.PHONY: build fast container

EXECUTABLE ?= main
IMAGE ?= bin/$(EXECUTABLE)

all: build

func:
	./func.sh

build:
	CGO_ENABLED=0 go build --ldflags '${EXTLDFLAGS}' -o ${IMAGE} github.com/sigma-dev/sigma/plugin/client

container:
	
	docker build -t quake3 .
