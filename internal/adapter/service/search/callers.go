package search

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/adapter/service"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	pb "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/proto/search"
	"google.golang.org/grpc/metadata"
)

func (s *SearchService) CallSearchMedia(ctx context.Context, query string, limit, offset uint) ([]uint, error) {
	log := logger.LoggerWithKey(s.Logger, ctx, common.ContextKeyRequestID)
	log.Info("SearchService CallSearchMedia method invoked")

	// Add request ID metadata
	md := metadata.Pairs(string(common.ContextKeyRequestID), common.GetRequestIDFromContext(ctx))
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Prepare the gRPC request
	req := &pb.SearchMediaRequest{
		Query:  query,
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, service.ResponseTimeout)
	defer cancel()


	// Make the gRPC call
	r, err := s.client.SearchMedia(ctx, req)
	if err != nil {
		log.Error("could not call SearchMedia", log.ToError(err))
		return nil, err
	}

	// Extract media IDs from the response
	mediaIDs := make([]uint, len(r.GetMedias()))
	for i, x := range r.GetMedias() {
		mediaIDs[i] = uint(x.GetId())
	}

	// Log the successful search
	log.Info("SearchMedia successful", log.ToInt("results_count", len(mediaIDs)))
	return mediaIDs, nil
}


func (s *SearchService) CallSearchActors(ctx context.Context, query string, limit, offset uint) ([]uint, error) {
	log := logger.LoggerWithKey(s.Logger, ctx, common.ContextKeyRequestID)
	log.Info("SearchService CallSearchActors method invoked")

	// Add request ID metadata
	md := metadata.Pairs(string(common.ContextKeyRequestID), common.GetRequestIDFromContext(ctx))
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Prepare the gRPC request
	req := &pb.SearchActorRequest{
		Query:  query,
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, service.ResponseTimeout)
	defer cancel()


	// Make the gRPC call
	r, err := s.client.SearchActor(ctx, req)
	if err != nil {
		log.Error("could not call SearchActor", log.ToError(err))
		return nil, err
	}

	// Extract media IDs from the response
	mediaIDs := make([]uint, len(r.GetActors()))
	for i, x := range r.GetActors() {
		mediaIDs[i] = uint(x.GetId())
	}

	// Log the successful search
	log.Info("SearchActor successful", log.ToInt("results_count", len(mediaIDs)))
	return mediaIDs, nil
}
