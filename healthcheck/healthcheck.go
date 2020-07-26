package healthcheck

import (
	"context"
	"sync"

	"github.com/knipknap/grpc-server-go/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

// Server is the service that implements the grpc service interface
type HealthCheck struct {
	logger *zap.SugaredLogger

	// statusMap stores the serving status of the services this Server monitors.
	statusMap map[string]proto.HealthCheckResponse_ServingStatus
	statusMapMutex sync.Mutex
}

// NewServiceServer returns an implementation of ServiceServer
func NewHealthCheck(logger *zap.SugaredLogger) *HealthCheck {
	return &HealthCheck{
		logger: logger,
		statusMap: make(map[string]proto.HealthCheckResponse_ServingStatus),
	}
}

func (s *HealthCheck) Check(ctx context.Context, in *proto.HealthCheckRequest) (*proto.HealthCheckResponse, error) {
	s.statusMapMutex.Lock()
	defer s.statusMapMutex.Unlock()
	if in.Service == "" {
		// check the server overall health status.
		return &proto.HealthCheckResponse{
			Status: proto.HealthCheckResponse_SERVING,
		}, nil
	}
	if status, ok := s.statusMap[in.Service]; ok {
		return &proto.HealthCheckResponse{
			Status: status,
		}, nil
	}
	return nil, status.Error(codes.NotFound, "unknown service")
}
