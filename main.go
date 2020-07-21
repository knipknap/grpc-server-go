package main

import (
	"context"
	"time"

	"github.com/barebaric/spiff-mm/proto"
	"go.uber.org/zap"
)

// Server is the service that implements the grpc service interface
type MicroModel struct {
	logger *zap.SugaredLogger
}

// NewServiceServer returns an implementation of ServiceServer
func New(logger *zap.SugaredLogger) proto.ServiceServer {
	return &MicroModel{
		logger: logger,
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

func processTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Local().Location())
}
