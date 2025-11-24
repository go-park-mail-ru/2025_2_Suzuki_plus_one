package auth

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	pb "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/proto/auth"
)

var LoginError = fmt.Errorf("login error")

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Register log with request ID
	log := logger.LoggerWithKey(s.Log, ctx, common.ContextKeyRequestID)
	log.Debug("GRPC HANDLER Login called")

	input := dto.PostAuthSignInInput{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	output, err := s.loginUsecase.Execute(ctx, input)
	
	if err != nil {
		log.Error("Can't execute login usecase",
			log.ToString("email", req.GetEmail()),
			log.ToAny("derr", err),
		)
		return &pb.LoginResponse{Success: false}, LoginError
	}

	pbOutput := &pb.LoginResponse{
		Success:      true,
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	}

	log.Debug("LOGIN RESPONSE: %v", pbOutput)

	return pbOutput, nil
}

func (s *AuthServer) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	// Register log with request ID
	log := logger.LoggerWithKey(s.Log, ctx, common.ContextKeyRequestID)
	log.Debug("GRPC HANDLER Login called")

	input := dto.GetAuthRefreshInput{
		RefreshToken: req.GetRefreshToken(),
	}

	output, err := s.refreshUsecase.Execute(ctx, input)

	if err != nil {
		log.Error("Can't execute refresh token usecase",
			log.ToString("refresh_token", req.GetRefreshToken()),
			log.ToAny("derr", err),
		)
		return &pb.RefreshTokenResponse{Success: false}, fmt.Errorf("refresh token error")
	}

	pbOutput := &pb.RefreshTokenResponse{
		Success:     true,
		AccessToken: output.AccessToken,
	}

	log.Debug("REFRESH TOKEN RESPONSE: %v", pbOutput)
	return &pb.RefreshTokenResponse{}, nil
}

func (s *AuthServer) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	// Register log with request ID
	log := logger.LoggerWithKey(s.Log, ctx, common.ContextKeyRequestID)
	log.Debug("GRPC HANDLER Logout called")

	input := dto.GetAuthSignOutInput{
		RefreshToken: req.GetRefreshToken(),
	}

	_, err := s.logoutUsecase.Execute(ctx, input)

	if err != nil {
		log.Error("Can't execute logout usecase",
			log.ToString("refresh_token", req.GetRefreshToken()),
			log.ToAny("derr", err),
		)
		return &pb.LogoutResponse{Success: false}, fmt.Errorf("logout error")
	}

	log.Debug("LOGOUT SUCCESSFUL")
	return &pb.LogoutResponse{Success: true}, nil
}

func (s *AuthServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// Register log with request ID
	log := logger.LoggerWithKey(s.Log, ctx, common.ContextKeyRequestID)
	log.Debug("GRPC HANDLER CreateUser called")

	input := dto.PostAuthSignUpInput{
		Email:    req.GetEmail(),
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}

	output, err := s.createUserUsecase.Execute(ctx, input)

	if err != nil {
		log.Error("Can't execute create user usecase",
			log.ToString("email", req.GetEmail()),
			log.ToAny("derr", err),
		)
		return &pb.CreateUserResponse{Success: false}, fmt.Errorf("create user error")
	}

	pbOutput := &pb.CreateUserResponse{
		UserId:       output.UserID,
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
		Success:      true,
	}

	log.Debug("CREATE USER RESPONSE: %v", pbOutput)

	return pbOutput, nil
}
