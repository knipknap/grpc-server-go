protobuf:
	cd proto; \
	ls -1 *.proto 2>/dev/null | while read FILENAME; do \
	    protoc -I /usr/local/include/ -I. "$$FILENAME" --go_out=plugins=grpc:.; \
	done

build: protobuf
	go mod download
	go build -o cmd/start cmd/start.go

docker-build:
	docker build -t spiff-mm:latest .

docker-run:
	docker run \
		-e GRPC_HOST=0.0.0.0 \
		-e GRPC_PORT=8181 \
		-e DEBUG=1 \
		-p 8181:8181 \
		spiff-mm:latest
