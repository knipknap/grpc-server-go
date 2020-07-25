package main

import (
	"context"
	"time"
	"sync"

	"github.com/barebaric/spiff-mm/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

// Server is the service that implements the grpc service interface
type MicroModel struct {
	logger *zap.SugaredLogger

	// statusMap stores the serving status of the services this Server monitors.
	statusMap map[string]proto.HealthCheckResponse_ServingStatus
	statusMapMutex sync.Mutex
}

// NewServiceServer returns an implementation of ServiceServer
func NewMicroModel(logger *zap.SugaredLogger) proto.ServiceServer {
	return &MicroModel{
		logger: logger,
		statusMap: make(map[string]proto.HealthCheckResponse_ServingStatus),
	}
}

func (s *MicroModel) GetInfo(ctx context.Context, req *proto.ModelInfoRequest) (*proto.ModelInfo, error) {
	logger := s.logger.With("method", "GetInfo")
	logger.Infow("call", "req", req)
	return &proto.ModelInfo{
		Name: "demo",
		Version: "v1",
		Hname: "Demo MicroModel",
		InputReferencePeriod: proto.ReferencePeriod_DAY,
	}, nil
}

func (s *MicroModel) GetValueFromDateRange(ctx context.Context, req *proto.ModelRequest) (*proto.ModelResult, error) {
	logger := s.logger.With("method", "GetValueFromDateRange")
	logger.Infow("call", "req", req)

	// Access other input parameters
	//options := req.GetOptions()
	input := req.GetInput()
	//fromDate := processTime(req.GetFromDate().AsTime())
	//toDate := processTime(req.GetToDate().AsTime())

	// Implementation goes here
	// The default implementation just echos the input.

	return &proto.ModelResult{Streams: input.GetStreams()}, nil
}

func (s *MicroModel) GetIncomeFromDateRange(ctx context.Context, req *proto.ModelRequest) (*proto.ModelResult, error) {
	logger := s.logger.With("method", "GetIncomeFromDateRange")
	logger.Infow("call", "req", req)

	// Implementation goes here. Access to parameters is the same as in GetValueFromDateRange()
	// The default implementation just echos the input.
	input := req.GetInput()

	return &proto.ModelResult{Streams: input.GetStreams()}, nil
}

func (s *MicroModel) Check(ctx context.Context, in *proto.HealthCheckRequest) (*proto.HealthCheckResponse, error) {
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

func processTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Local().Location())
}
