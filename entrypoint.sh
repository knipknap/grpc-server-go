#!/bin/sh

./start -port ${GRPC_PORT:-8181} &
sleep 2
grpcui -bind 0.0.0.0 -plaintext -port ${GRPCUI_PORT:-8080} localhost:${GRPC_PORT:-8181}
