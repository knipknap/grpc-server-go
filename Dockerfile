FROM knipknap/grpc-server-go:latest
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
CMD ["./start"]
