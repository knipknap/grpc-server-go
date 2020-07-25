package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"plugin"
	"path/filepath"

	"github.com/barebaric/spiff-mm/proto"
	"github.com/oklog/oklog/pkg/group"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 8181, "The server port")
)

var log grpclog.LoggerV2

func init() {
	log = grpclog.NewLoggerV2(os.Stdout, os.Stderr, os.Stderr)
	grpclog.SetLoggerV2(log)
}

func main() {
	// Set up logging.
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	flag.Parse()

	// Find directory of this file
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	// Load the service as a plugin
	plugin_filename := filepath.Join(exPath, "service.so")
	service, err := plugin.Open(plugin_filename)
	if err != nil {
		sugar.Panicw("Failed to load service plugin", plugin_filename, err)
	}
	NewServiceSymbol, err := service.Lookup("NewMicroModel")
	if err != nil {
		sugar.Panicw("Service plugin is missing 'New' method", err)
	}
	NewService := NewServiceSymbol.(func(logger *zap.SugaredLogger) proto.ServiceServer)

	// clearly demarcates the scope in which each listener/socket may be used.
	var g group.Group
	{
		// The gRPC listener mounts the Go kit gRPC server we created.
		grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
		if err != nil {
			sugar.Errorw("failed to listen grpc", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			sugar.Infow("grpc address", "addr", fmt.Sprintf(":%d", *port))

			var opts []grpc.ServerOption
			grpcServer := grpc.NewServer(opts...)
			proto.RegisterServiceServer(grpcServer, NewService(sugar))
			reflection.Register(grpcServer)
			sugar.Infow("starting server")
			return grpcServer.Serve(grpcListener)
		}, func(error) {
			grpcListener.Close()
		})
	}
	{
		// This function just sits and waits for ctrl-C.
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}

	sugar.Infow("exit", "reason", g.Run())
}
