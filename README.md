# gRPC server for Go

## Introduction

This Dockerfile implements a base container for microservices implemented in Go.
The container does the following:

- It runs a gRPC server on port 8181. Reflection is enabled.
- It includes a a ready-to-use health check to allow for zero-downtime updates.
- It dynamically loads the Go plugin /app/service.so, which must be provided by the container that you build.
- It runs a [grpcui](https://github.com/fullstorydev/grpcui) on port 8080.

```
FROM knipknap/grpc-go:latest as build-env
WORKDIR /app
COPY go.mod Makefile ./
COPY proto proto
COPY service service
RUN go build -buildmode=plugin -o service.so service.go

FROM knipknap/grpc-server-go:latest
COPY --from=build-env /app/service.so .
```

Important: The Go plugin MUST offer a function like this:

```
func RegisterService(grpcServer *grpc.Server, logger *zap.SugaredLogger) {
	// You should register your gRPC service here like this:
	proto.RegisterServiceServer(grpcServer, NewService(logger))
}
```
