#!/bin/sh

./start &
sleep 2
grpcui -bind 0.0.0.0 -plaintext -port ${GRPCUI_PORT:-8080} ${GRPCUI_SERVER:-}
