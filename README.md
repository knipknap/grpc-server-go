# grpc-server-go

## Introduction

grpc-server-go provides a production-ready gRPC server for Golang in a **Docker** container.

<img src="./assets/grpc-logo.svg" width="300px" align="right">

This Dockerfile implements a base container for microservices implemented in Go.

The container does the following:

- It runs a gRPC server.
- It dynamically loads the Go plugin `/app/service.so`, **which MUST be provided by the container that you build**.
- **Reflection** is enabled.
- It includes a complete [health check](https://github.com/grpc/grpc/blob/master/doc/health-checking.md) to allow for zero-downtime updates.
- It runs a [grpcui](https://github.com/fullstorydev/grpcui).
- It picks up additional JSON-based health checks and runs them, validating if the RPC call generates the expected result. More below.

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

### Supported environment variables

- `GRPC_PORT`: The port of the gRPC server, by default 8181
- `GRPCUI_PORT`: The port of the gRPC user interface, by default 8080
- `DEBUG`: To change the zap logger from Production to Development

## Your Go plugin

Building your code as a Go plugin is easy:

- Make sure that your package is named "main" (this is a Go requirement)
- Compile using `go build -buildmode=plugin -o service.so service.go` (as shown in the Dockerfile above)
- Make sure that your main package includes a RegisterService function with the following signature:

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

## JSON based health checks

In addition to implementing the standard [health check](https://github.com/grpc/grpc/blob/master/doc/health-checking.md)
protocol, this container also provide a simple mechanism for more powerful health checks
based on [grpcurl](https://github.com/fullstorydev/grpcurl).

The container will execute the gRPC calls every time when Docker requests a status from the
healthcheck.

To use this mechanism, simply put JSON files under `/app/healthchecks/`, with the following
structure:

```
healthchecks/
├── check1
│   └── my.host.Service
│       └── MyMethod
│           ├── input.json
│           └── output.json
└── check2
    └── my.host.Service2
        ├── MyMethod1
        │    ├── input.json
        │    └── output.json
        └── MySuperMethod
             ├── input.json
             └── output.json
```

In this example, the healthcheck would call my.host.Service/MyMethod, passing the input
from input.json as parameters. The return value has to match the contents output.json
exactly, otherwise the check is considered failed.

my.host.Service2/MyMethod1 and MySuperMethod are both called and checked in the same
manner.
