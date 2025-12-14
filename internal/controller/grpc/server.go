package grpc

import (
	"net"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	Log        logger.Logger
	Server     *grpc.Server
	Lis        net.Listener
	Middleware []grpc.UnaryServerInterceptor
}

func NewGRPCServer(log logger.Logger) *GRPCServer {
	if log == nil {
		panic("logger is nil")
	}
	return &GRPCServer{
		Log:        log,
		Middleware: make([]grpc.UnaryServerInterceptor, 0),
	}
}

func (s *GRPCServer) InitGRPCServer(serveString string) {
	var err error
	s.Lis, err = net.Listen("tcp", serveString)

	if err != nil {
		s.Log.Fatal("failed to listen", "error", err)
	}

	// Create gRPC server with middleware
	s.Server = grpc.NewServer(
		grpc.UnaryInterceptor(
			UnaryRequestIDInterceptor(s.Log),
		),
		grpc.ChainUnaryInterceptor(s.Middleware...),
	)
}
