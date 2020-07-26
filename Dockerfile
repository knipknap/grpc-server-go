FROM knipknap/grpc-server-go:latest as build-env
WORKDIR /app
COPY go.mod Makefile ./
COPY proto proto
COPY server server
COPY service service
RUN make build

FROM golang:1.13-alpine
WORKDIR /app
COPY --from=build-env /app/server/cmd/start .
COPY --from=build-env /app/server/cmd/service.so .
COPY --from=build-env /bin/grpc_health_probe /bin
#HEALTHCHECK --interval=30s --timeout=2s --start-period=35s CMD grpc_health_probe -addr=:8181
CMD ["./start"]
