FROM knipknap/grpc-go:latest as build-env
WORKDIR /app
RUN GOBIN=/usr/local/bin \
    apk add findutils \
    && go get -x github.com/fullstorydev/grpcui \
    && go install -x github.com/fullstorydev/grpcui/cmd/grpcui \
    && go get github.com/fullstorydev/grpcurl \
    && go install github.com/fullstorydev/grpcurl/cmd/grpcurl
COPY go.mod go.sum entrypoint.sh healthcheck.sh Makefile ./
COPY config config
COPY proto proto
COPY healthcheck healthcheck
COPY healthchecks healthchecks
COPY cmd cmd
COPY service service
RUN make build

ENV GRPC_PORT=8181
ENV GRPCUI_PORT=8080
HEALTHCHECK --interval=30s --timeout=2s --start-period=20s CMD ./healthcheck.sh -addr=:$GRPC_PORT
CMD ["./entrypoint.sh"]
