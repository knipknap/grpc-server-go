# grpc-server-go

## Introduction

grpc-server-go provides a production-ready gRPC server for Golang in a **Docker** container.

<img src="./assets/grpc-logo.svg" width="300px" align="right">

This Dockerfile implements a base container for microservices implemented in Go.

The container does the following:

- It includes a main.go for a gRPC server that imports server/main.go for your code.
- **Reflection** is enabled.
- It includes a complete [health check](https://github.com/grpc/grpc/blob/master/doc/health-checking.md) to allow for zero-downtime updates.
- It runs a [grpcui](https://github.com/fullstorydev/grpcui).
- It picks up additional JSON-based health checks and runs them, validating if the RPC call generates the expected result. More below.

## Your Dockerfile

A typical grpc-server-go Dockerfile will contain two stages:

- The first stage is based on [grpc-server-go](https://github.com/knipknap/grpc-go) to provide
  an environment for compiling your Go code; this environment contains the grpc server as cmd/main.go
- The second stage produces a light weight production-ready container with a gRPC server.

Example:

```
FROM knipknap/grpc-server-go:latest as build-env
WORKDIR /app
COPY proto/service.proto proto/
COPY proto/options.proto proto/
COPY service/main.go service/
RUN make build

FROM golang:1.13-alpine

# Make sure to include all these COPY commands, they are required.
COPY --from=build-env /app/start /app/entrypoint.sh /app/healthcheck.sh ./
COPY --from=build-env /app/healthchecks /app/
COPY --from=build-env /bin/grpc_health_probe /bin
COPY --from=build-env /go/bin/grpcui /usr/local/bin/grpcui
COPY --from=build-env /go/bin/grpcurl /usr/local/bin/grpcurl
COPY --from=build-env /usr/bin/find /usr/bin

RUN adduser -S -u 10001 user
USER user
ENV GRPC_PORT=8181
ENV GRPCUI_PORT=8080
HEALTHCHECK --interval=30s --timeout=2s --start-period=20s CMD ./healthcheck.sh -addr=:$GRPC_PORT
CMD ["./entrypoint.sh"]
```

### Supported environment variables

- `GRPC_PORT`: The port of the gRPC server, by default 8181
- `GRPCUI_PORT`: The port of the gRPC user interface, by default 8080
- `DEBUG`: To change the zap logger from Production to Development

## Your Go plugin

Add your code as follows:

- Make sure that your container adds a `/app/service/main.go`. The package could be named "service"
- Place your .proto files in the `/app/proto/` folder. The container includes a Makefile that will
  compile them during the build stage (see "make build" in the example above)
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
