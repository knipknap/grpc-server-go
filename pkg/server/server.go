package server

import (
	"context"
	"time"

	"grpc-server-go.localhost/proto"
	"go.uber.org/zap"
	//"github.com/golang/protobuf/ptypes"
)

// Server is the service that implements the grpc service interface
type server struct {
	logger *zap.Logger
	suggar *zap.SugaredLogger
}

// NewServiceServer returns an implementation of ServiceServer
func NewServiceServer(logger *zap.Logger) proto.ServiceServer {
	return &server{
		logger: logger,
		suggar: logger.Sugar(),
	}
}

func (s *server) GetInfo(ctx context.Context, req *proto.ModelInfoRequest) (*proto.ModelInfo, error) {
	logger := s.suggar.With("method", "GetInfo")
	logger.Infow("call", "req", req)
	return &proto.ModelInfo{
		Name: "demo",
		Version: "v1",
		Hname: "Demo MicroModel",
		InputReferencePeriod: proto.ReferencePeriod_DAY,
	}, nil
}

func (s *server) GetValueFromDateRange(ctx context.Context, req *proto.ModelRequest) (*proto.ModelResult, error) {
	logger := s.suggar.With("method", "GetValueFromDateRange")
	logger.Infow("call", "req", req)

	// Unpack the model-specific options
	//optionsBlob := req.GetOptions()
	//options := proto.Options
	//if err := ptypes.UnmarshalAny(optionsBlob, &options); err != nil {
	//	logger.Fatal("could not deserialize options")
	//}

	// Access other input parameters
	//options := req.GetOptions()
	//input := req.GetInput()
	//fromDate := processTime(req.GetFromDate().AsTime())
	//toDate := processTime(req.GetToDate().AsTime())

	// Implementation goes here

	return &proto.ModelResult{}, nil
}

func (s *server) GetIncomeFromDateRange(ctx context.Context, req *proto.ModelRequest) (*proto.ModelResult, error) {
	logger := s.suggar.With("method", "GetIncomeFromDateRange")
	logger.Infow("call", "req", req)

	// Implementation goes here. Access to parameters is the same as in GetValueFromDateRange()

	return &proto.ModelResult{}, nil
}

func processTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Local().Location())
}
