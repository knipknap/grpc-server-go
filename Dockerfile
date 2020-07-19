FROM knipknap/grpc-server-go:latest
WORKDIR /app
COPY pkg pkg
COPY proto/service.proto proto
ENV GRPC_GO_LOG_VERBOSITY_LEVEL=99
ENV GRPC_GO_LOG_SEVERITY_LEVEL=info
