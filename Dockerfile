FROM golang:1.13-alpine as build-env

ENV PROTOC_VER=v3.12.3
ENV PROTOC_ZIP=protoc-3.12.3-linux-x86_64.zip
ENV GO111MODULE=on

RUN apk add --no-cache --update curl && \
    GRPC_HEALTH_PROBE_VERSION=v0.3.1 && \
    curl -OL https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    mv grpc_health_probe* grpc_health_probe && \
    chmod +x grpc_health_probe

RUN apk add --no-cache --update curl unzip build-base autoconf automake libtool ca-certificates \
    && go get google.golang.org/grpc \
    && go get github.com/golang/protobuf/protoc-gen-go \
    && go get -u -v github.com/golang/protobuf/proto \
    && apk add --no-cache --update protoc protobuf \
    && curl -OL https://github.com/protocolbuffers/protobuf/releases/download/$PROTOC_VER/$PROTOC_ZIP \
    && unzip -o $PROTOC_ZIP -d /usr/local 'include/*' \
    && rm -f $PROTOC_ZIP

WORKDIR /app
COPY go.mod Makefile ./
COPY proto proto
COPY server server
COPY service service

RUN make build

################################################
# Production
FROM golang:1.13-alpine
WORKDIR /app
COPY --from=build-env /app/server/cmd/start .
COPY --from=build-env /app/server/cmd/service.so .
CMD ["./start"]
