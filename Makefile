protobuf:
	cd proto; \
	ls -1 *.proto 2>/dev/null | while read FILENAME; do \
	    protoc -I /usr/local/include/ -I. "$$FILENAME" --go_out=plugins=grpc:.; \
	done

build: protobuf
	go mod download
	go build -o start cmd/start.go

docker-build:
	docker build -t knipknap/grpc-server-go .

docker-run:
	docker run \
		-e GRPC_PORT=8181 \
		-e GRPCUI_PORT=8080 \
		-e DEBUG=1 \
		-p 18080:8080 \
		-p 18181:8181 \
		knipknap/grpc-server-go:latest
