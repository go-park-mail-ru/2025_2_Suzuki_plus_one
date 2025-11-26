package auth

import (
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/adapter/service"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/metrics"
	pb "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/proto/auth"
)

type AuthService struct {
	service.Service
	client pb.AuthServiceClient
}

func NewAuthService(l logger.Logger, serverString string) *AuthService {
	s := &AuthService{
		Service: *service.NewService(l, serverString),
	}
	return s
}

func (s *AuthService) Connect() error {
	// Add metrics middleware
	s.Service.Middleware = append(s.Service.Middleware,
		metrics.GRPCClientMetricsInterceptor(metrics.ServiceHTTP, metrics.ServiceSearch),
	)
	err := s.Service.Init()
	if err != nil {
		return err
	}
	s.client = pb.NewAuthServiceClient(s.Connection)
	return nil
}

func (s *AuthService) Close() {
	s.Service.Close()
}
