FROM knipknap/grpc-server-go
WORKDIR /app
COPY proto proto
COPY go.mod main.go service/
