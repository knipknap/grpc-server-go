# gRPC server for Go

## Introduction

This Dockerfile implements a base container for microservices implemented in Go.
The container does the following:

- It runs a gRPC server on **port 8181**.
- **Reflection** is enabled.
- It includes a a ready-to-use [health check](https://github.com/grpc/grpc/blob/master/doc/health-checking.md) to allow for zero-downtime updates.
- It dynamically loads the Go plugin `/app/service.so`, **which MUST be provided by the container that you build**.
- It runs a [grpcui](https://github.com/fullstorydev/grpcui) on port 8080.

## Your Dockerfile

A typical Dockerfile could look like this: Two stages to allow for compilation of your Go code,
while still producing a light weight production-ready container.

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

- Make sure that your package is named "main"
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
