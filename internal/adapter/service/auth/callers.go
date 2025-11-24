package auth

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	logger "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	pb "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/proto/auth"
	"google.golang.org/grpc/metadata"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/adapter/service"

)

func (s *AuthService) CallLogin(ctx context.Context, email string, password string) (accessToken string, refreshToken string, err error) {
	log := logger.LoggerWithKey(s.Logger, ctx, common.ContextKeyRequestID)

	log.Info("AuthService CallLogin method invoked")

	// Add request ID metadata
	md := metadata.Pairs(string(common.ContextKeyRequestID), common.GetRequestIDFromContext(ctx))
	ctx = metadata.NewOutgoingContext(ctx, md)

	ctx, cancel := context.WithTimeout(ctx, service.ResponseTimeout)
	defer cancel()

	r, err := s.client.Login(ctx, &pb.LoginRequest{Email: email, Password: password})
	if err != nil {
		log.Error("could not call Login", log.ToError(err))
		return "", "", err
	}

	if r.GetSuccess() == false {
		log.Warn("Login failed for email", log.ToString("email", email))
		return "", "", fmt.Errorf("login failed")
	}

	log.Info("Login successful for email", log.ToString("email", email))
	return r.GetAccessToken(), r.GetRefreshToken(), nil
}

func (s *AuthService) CallCreateUser(ctx context.Context, email string, username string, password string) (userID uint, accessToken string, refreshToken string, err error) {
	log := logger.LoggerWithKey(s.Logger, ctx, common.ContextKeyRequestID)
	log.Info("AuthService CallCreateUser method invoked")

	// Add request ID metadata
	md := metadata.Pairs(string(common.ContextKeyRequestID), common.GetRequestIDFromContext(ctx))
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), service.ResponseTimeout)
	defer cancel()

	r, err := s.client.CreateUser(ctx, &pb.CreateUserRequest{Email: email, Username: username, Password: password})
	if err != nil {
		log.Error("could not call CreateUser", log.ToError(err))
		return 0, "", "", err
	}

	if r.GetSuccess() == false {
		log.Warn("CreateUser failed for email", log.ToString("email", email))
		return 0, "", "", fmt.Errorf("create user failed")
	}

	log.Info("CreateUser successful for email", log.ToString("email", email))

	userID = uint(r.GetUserId())
	return userID, r.GetAccessToken(), r.GetRefreshToken(), nil
}

func (s *AuthService) CallLogout(ctx context.Context, refreshToken string) error {
	log := logger.LoggerWithKey(s.Logger, ctx, common.ContextKeyRequestID)
	log.Info("AuthService CallLogout method invoked")

	// Add request ID metadata
	md := metadata.Pairs(string(common.ContextKeyRequestID), common.GetRequestIDFromContext(ctx))
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), service.ResponseTimeout)
	defer cancel()

	r, err := s.client.Logout(ctx, &pb.LogoutRequest{RefreshToken: refreshToken})
	if err != nil {
		log.Error("could not call Logout", log.ToError(err))
		return err
	}

	if r.GetSuccess() == false {
		log.Warn("Logout failed for refresh token", log.ToString("refresh_token", refreshToken))
		return fmt.Errorf("logout failed")
	}

	log.Info("Logout successful for refresh token", log.ToString("refresh_token", refreshToken))
	return nil
}

func (s *AuthService) CallRefresh(ctx context.Context, refreshToken string) (newAccessToken string, err error) {
	log := logger.LoggerWithKey(s.Logger, ctx, common.ContextKeyRequestID)
	log.Info("AuthService CallRefreshToken method invoked")

	// Add request ID metadata
	md := metadata.Pairs(string(common.ContextKeyRequestID), common.GetRequestIDFromContext(ctx))
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), service.ResponseTimeout)
	defer cancel()

	r, err := s.client.RefreshToken(ctx, &pb.RefreshTokenRequest{RefreshToken: refreshToken})
	if err != nil {
		log.Error("could not call RefreshToken", log.ToError(err))
		return "", err
	}

	if r.GetSuccess() == false {
		log.Warn("RefreshToken failed for refresh token", log.ToString("refresh_token", refreshToken))
		return "", fmt.Errorf("refresh token failed")
	}

	log.Info("RefreshToken successful for refresh token", log.ToString("refresh_token", refreshToken))
	return r.GetAccessToken(), nil
}
