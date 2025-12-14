package search

import (
	srv "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/grpc"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	pb "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/proto/search"
)

type SearchServer struct {
	pb.UnimplementedSearchServiceServer
	srv.GRPCServer

	searchMediaUsecase SearchMediaUsecase
	searchActorUsecase SearchActorUsecase
}

func NewSearchServer(log logger.Logger, searchMediaUsecase SearchMediaUsecase, searchActorUsecase SearchActorUsecase) *SearchServer {
	if searchMediaUsecase == nil {
		panic("searchMediaUsecase is nil")
	}
	if searchActorUsecase == nil {
		panic("searchActorUsecase is nil")
	}
	return &SearchServer{
		GRPCServer:         *srv.NewGRPCServer(log),
		searchMediaUsecase: searchMediaUsecase,
		searchActorUsecase: searchActorUsecase,
	}
}

func (s *SearchServer) StartGRPCServer(serveString string) {
	s.InitGRPCServer(serveString)
	lis := s.Lis
	server := s.Server

	pb.RegisterSearchServiceServer(server, s)

	s.Log.Info("server listening at", "lis.addr", lis.Addr())
	if err := server.Serve(lis); err != nil {
		s.Log.Fatal("failed to serve", "error", err)
	}
}
