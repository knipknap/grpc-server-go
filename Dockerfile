FROM knipknap/grpc-go:latest as build-env
WORKDIR /app
RUN go get -x github.com/fullstorydev/grpcui && \
    go install -x github.com/fullstorydev/grpcui/cmd/grpcui
COPY go.mod Makefile ./
COPY proto proto
COPY healthcheck healthcheck
COPY cmd cmd
RUN make build

FROM golang:1.13-alpine
ENV GRPC_PORT=8181
ENV GRPCUI_PORT=8080
WORKDIR /app
COPY entrypoint.sh .
COPY --from=build-env /app/start .
COPY --from=build-env /bin/grpc_health_probe /usr/bin
COPY --from=build-env /go/bin/grpcui /usr/bin/grpcui
HEALTHCHECK --interval=30s --timeout=2s --start-period=20s CMD grpc_health_probe -addr=:$GRPC_PORT
RUN adduser -S -u 10001 user
USER user
CMD ["./entrypoint.sh"]
