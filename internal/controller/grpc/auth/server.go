package auth

import (
	srv "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/grpc"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	pb "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/proto/auth"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	srv.GRPCServer
	loginUsecase      LoginUsecase
	refreshUsecase    RefreshUsecase
	logoutUsecase     LogoutUsecase
	createUserUsecase CreateUserUsecase
}

func NewAuthServer(log logger.Logger, loginUsecase LoginUsecase, refreshUsecase RefreshUsecase, logoutUsecase LogoutUsecase, createUserUsecase CreateUserUsecase) *AuthServer {
	if loginUsecase == nil {
		panic("loginUsecase is nil")
	}
	if refreshUsecase == nil {
		panic("refreshUsecase is nil")
	}
	if logoutUsecase == nil {
		panic("logoutUsecase is nil")
	}
	if createUserUsecase == nil {
		panic("createUserUsecase is nil")
	}
	return &AuthServer{
		GRPCServer:        *srv.NewGRPCServer(log),
		loginUsecase:      loginUsecase,
		refreshUsecase:    refreshUsecase,
		logoutUsecase:     logoutUsecase,
		createUserUsecase: createUserUsecase,
	}
}

func (s *AuthServer) StartGRPCServer(serveString string) {
	s.InitGRPCServer(serveString)
	lis := s.Lis
	server := s.Server

	pb.RegisterAuthServiceServer(server, s)

	s.Log.Info("server listening at", "lis.addr", lis.Addr())
	if err := server.Serve(lis); err != nil {
		s.Log.Fatal("failed to serve", "error", err)
	}
}
