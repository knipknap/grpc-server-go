FROM knipknap/grpc-go:latest as build-env
WORKDIR /app
RUN apk add findutils \
    && go get -x github.com/fullstorydev/grpcui \
    && go install -x github.com/fullstorydev/grpcui/cmd/grpcui \
    && go get github.com/fullstorydev/grpcurl \
    && go install github.com/fullstorydev/grpcurl/cmd/grpcurl
COPY go.mod go.sum Makefile ./
COPY config config
COPY proto proto
COPY healthcheck healthcheck
COPY cmd cmd
RUN make build

FROM golang:1.13-alpine
ENV GRPC_PORT=8181
ENV GRPCUI_PORT=8080
WORKDIR /app
COPY entrypoint.sh healthcheck.sh ./
COPY healthchecks healthchecks
COPY --from=build-env /app/start .
COPY --from=build-env /bin/grpc_health_probe /usr/bin
COPY --from=build-env /go/bin/grpcui /usr/bin/grpcui
COPY --from=build-env /go/bin/grpcurl /usr/bin/grpcurl
COPY --from=build-env /usr/bin/find /usr/bin
HEALTHCHECK --interval=30s --timeout=2s --start-period=20s CMD ./healthcheck.sh -addr=:$GRPC_PORT
RUN adduser -S -u 10001 user
USER user
CMD ["./entrypoint.sh"]
