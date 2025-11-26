package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

func UnaryRequestIDInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {

		// Extract metadata
		md, ok := metadata.FromIncomingContext(ctx)
		var requestID string

		if ok {
			ids := md.Get(string(common.ContextKeyRequestID))
			if len(ids) > 0 {
				requestID = ids[0]
			}
		}

		// Generate if not present
		if requestID == "" {
			logger.Warn("Generated new id, because request ID header is missing",
				logger.ToString("method", info.FullMethod),
				logger.ToString("key", string(common.ContextKeyRequestID)))
			// requestID = uuid.NewString()
		}

		// Add into context
		ctx = context.WithValue(ctx, common.ContextKeyRequestID, requestID)

		// Logging
		in_banner := "--->"
		logger.Debug(in_banner, logger.ToString("requestID", requestID),
			logger.ToString("method", info.FullMethod))

		// Call handler
		resp, err = handler(ctx, req)

		// Logging after response
		out_banner := "<---"
		logger.Debug(out_banner, logger.ToString("requestID", requestID),
			logger.ToString("method", info.FullMethod),
			logger.ToError(err))

		return resp, err
	}
}
