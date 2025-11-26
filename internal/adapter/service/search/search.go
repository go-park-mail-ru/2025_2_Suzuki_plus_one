package search

import (
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/adapter/service"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/metrics"
	pb "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/proto/search"
)

type SearchService struct {
	service.Service
	client pb.SearchServiceClient
}

func NewSearchService(l logger.Logger, serverString string) *SearchService {
	s := &SearchService{
		Service: *service.NewService(l, serverString),
	}
	return s
}

func (s *SearchService) Connect() error {
	// Add metrics middleware
	s.Service.Middleware = append(s.Service.Middleware,
		metrics.GRPCClientMetricsInterceptor(metrics.ServiceHTTP, metrics.ServiceSearch),
	)
	err := s.Service.Init()
	if err != nil {
		return err
	}
	s.client = pb.NewSearchServiceClient(s.Connection)
	return nil
}

func (s *SearchService) Close() {
	s.Service.Close()
}
