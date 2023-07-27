FROM golang:1.20-alpine3.18 AS build


RUN apk add --no-cache git

WORKDIR /usr/src/app

COPY ./ ./

RUN set -ex \
    && ls -lA \
    && git rev-parse --short=8 HEAD \
    && COMMIT_HASH="$(git rev-parse --short=8 HEAD 2>/dev/null)" \
    && BUILD_TIME="$(date +%FT%T%z)" \
    && TAG="$(git describe 2>/dev/null | cut -f 1 -d '-' 2>/dev/null)" \
    && export CGO_ENABLED=0 \
    && export GOOS=linux \
    && go build -a -installsuffix cgo -o admitCtlr \
        -ldflags "-s -w -X main.CommitHash=${COMMIT_HASH} -X main.BuildTime=${BUILD_TIME} -X main.Tag=${TAG}" .

FROM alpine:3.18

RUN apk add --no-cache ca-certificates

COPY --from=build /usr/src/app/admitCtlr /admitCtlr

ENTRYPOINT ["/admitCtlr"]
