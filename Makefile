# Makefile

APPNAME=admitCtlr
IMAGE_NAME=autermann/admission-webhook-server

DOCKERFILE=Dockerfile

VERSION_TAG=`git describe 2>/dev/null | cut -f 1 -d '-' 2>/dev/null`
COMMIT_HASH=`git rev-parse --short=8 HEAD 2>/dev/null`
BUILD_TIME=`date +%FT%T%z`
LDFLAGS=-ldflags "-s -w \
    -X main.CommitHash=${COMMIT_HASH} \
    -X main.BuildTime=${BUILD_TIME} \
    -X main.Tag=${VERSION_TAG}"

all: clean test fast

test:
	go test -v ./...

clean:
	go clean
	rm -r ./$(APPNAME) || true

fast:
	go build -o ${APPNAME} ${LDFLAGS}

docker_build:
	docker build . -f $(DOCKERFILE) -t $(IMAGE_NAME) --no-cache

docker: docker_build
