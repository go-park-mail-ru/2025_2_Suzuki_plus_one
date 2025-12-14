package service

import (
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const ResponseTimeout = 1 * time.Second

type Service struct {
	Logger       logger.Logger
	Connection   *grpc.ClientConn
	serverString string
	Middleware   []grpc.UnaryClientInterceptor
}

func NewService(l logger.Logger, serverString string) *Service {
	return &Service{
		Logger:       l,
		serverString: serverString,
		Middleware:   make([]grpc.UnaryClientInterceptor, 0),
	}
}

func (s *Service) Init() error {
	conn, err := grpc.NewClient(s.serverString,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(s.Middleware...),
	)
	if err != nil {
		s.Logger.Error("Failed to connect to gRPC server", s.Logger.ToError(err))
		return err
	}
	s.Connection = conn
	return nil
}

func (s *Service) Close() {
	if s.Connection != nil {
		if err := s.Connection.Close(); err != nil {
			s.Logger.Error("Failed to close gRPC connection", s.Logger.ToError(err))
		}
	}
}
