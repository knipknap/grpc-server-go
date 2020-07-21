FROM knipknap/grpc-server-go
WORKDIR /app
COPY proto proto
COPY main.go service/
