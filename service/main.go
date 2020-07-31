package service

import (
	"github.com/knipknap/grpc-server-go/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func RegisterService(grpcServer *grpc.Server, logger *zap.SugaredLogger) {
	proto.RegisterServiceServer(grpcServer, NewService(logger))
}

// Server is the service that implements the grpc service interface
type Service struct {
	logger *zap.SugaredLogger
}

// NewServiceServer returns an implementation of ServiceServer
func NewService(logger *zap.SugaredLogger) proto.ServiceServer {
	return &Service{
		logger: logger,
	}
}
