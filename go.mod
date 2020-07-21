module github.com/barebaric/spiff-mm

go 1.14

require (
	github.com/barebaric/grpc-server-go latest
	github.com/golang/protobuf v1.4.2
	github.com/oklog/oklog v0.3.2
	github.com/oklog/run v1.1.0 // indirect
	go.uber.org/zap v1.15.0
	golang.org/x/sys v0.0.0-20190624142023-c5567b49c5d0 // indirect
	golang.org/x/text v0.3.2 // indirect
	golang.org/x/tools v0.0.0-20191112195655-aa38f8e97acc // indirect
	google.golang.org/grpc v1.30.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/yaml.v2 v2.2.4 // indirect
)

replace github.com/barebaric/grpc-server-go => ../server
