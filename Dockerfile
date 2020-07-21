FROM knipknap/grpc-server-go
WORKDIR /app
COPY proto server/proto
COPY go.mod main.go service/
