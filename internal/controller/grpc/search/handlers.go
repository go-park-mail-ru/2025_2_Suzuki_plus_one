package search

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	pb "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/proto/search"
)

var _ pb.SearchServiceServer = (*SearchServer)(nil)

var searchMediaError = fmt.Errorf("search media error")

func (s *SearchServer) SearchMedia(ctx context.Context, req *pb.SearchMediaRequest) (*pb.SearchMediaResponse, error) {
	// Register log with request ID
	log := logger.LoggerWithKey(s.Log, ctx, common.ContextKeyRequestID)
	log.Debug("GRPC HANDLER SearchMedia called")

	input := dto.GetSearchInput{
		Query:  req.GetQuery(),
		Limit:  uint(req.GetLimit()),
		Offset: uint(req.GetOffset()),
	}

	output, errDto := s.searchMediaUsecase.Execute(ctx, input)
	if errDto != nil {
		log.Error("SearchMedia usecase error: ", errDto)
		return nil, searchMediaError
	}

	pbOutputs := &pb.SearchMediaResponse{}
	// Note: we ignore actors here, because they are empty
	// TODO: create a separate dto for grpc
	medias := output.Medias
	for _, media := range medias {
		pbOutput := &pb.Media{
			Id:          int64(media.MediaID),
		}
		pbOutputs.Medias = append(pbOutputs.Medias, pbOutput)
	}
	return pbOutputs, nil
}

func (s *SearchServer) SearchActor(ctx context.Context, req *pb.SearchActorRequest) (*pb.SearchActorResponse, error) {
	// Register log with request ID
	log := logger.LoggerWithKey(s.Log, ctx, common.ContextKeyRequestID)
	log.Debug("GRPC HANDLER SearchActors called")

	input := dto.GetSearchInput{
		Query:  req.GetQuery(),
		Limit:  uint(req.GetLimit()),
		Offset: uint(req.GetOffset()),
	}

	outputs, errDto := s.searchActorUsecase.Execute(ctx, input)
	if errDto != nil {
		log.Error("SearchActors usecase error: ", errDto)
		return nil, fmt.Errorf("search actors error")
	}

	pbOutputs := &pb.SearchActorResponse{}
	// Note: we ignore actors here, because they are empty
	// TODO: create a separate dto for grpc
	actors := outputs.Actors
	for _, actor := range actors {
		pbOutput := &pb.Actor{
			Id:          int64(actor.ID),
		}
		pbOutputs.Actors = append(pbOutputs.Actors, pbOutput)
	}
	return pbOutputs, nil
}
