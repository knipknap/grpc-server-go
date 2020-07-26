# grpc-server-go

## Introduction

grpc-server-go provides a production-ready gRPC server in a **Docker** container.

<img src="./assets/grpc-logo.svg" width="300px" align="right">

This Dockerfile implements a base container for microservices implemented in Go.

The container does the following:

- It runs a gRPC server on **port 8181**.
- **Reflection** is enabled.
- It includes a complete [health check](https://github.com/grpc/grpc/blob/master/doc/health-checking.md) to allow for zero-downtime updates.
- It dynamically loads the Go plugin `/app/service.so`, **which MUST be provided by the container that you build**.
- It runs a [grpcui](https://github.com/fullstorydev/grpcui) on port 8080.

## Your Dockerfile

A typical grpc-server-go Dockerfile will contain two stages:

- The first stage is based on [grpc-go](https://github.com/knipknap/grpc-go) to provide an environment for compiling your Go code
- The second stage produces a light weight production-ready container with a gRPC server.

Example:

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

# Your Go plugin

Building your code as a Go plugin is easy:

- Make sure that your package is named "main" (this is a Go requirement)
- Compile using `go build -buildmode=plugin -o service.so service.go` (as shown in the Dockerfile above)
- Make sure that your main package includes a function with the following signature:

```go
package main

import (
	"yourmodule.com/path/to/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func RegisterService(grpcServer *grpc.Server, logger *zap.SugaredLogger) {
	// You should register your gRPC service here like this:
	proto.RegisterServiceServer(grpcServer, NewYourService(logger))
}

type YourService struct {
	logger *zap.SugaredLogger
}

func NewYourService(logger *zap.SugaredLogger) proto.ServiceServer {
	return &YourService{
		logger: logger,
	}
}
```
