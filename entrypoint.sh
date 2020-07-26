#!/bin/sh

grpcui -bind 0.0.0.0 -plaintext -port ${GRPCUI_PORT:-8080} ${GRPCUI_SERVER:-} &
./start
