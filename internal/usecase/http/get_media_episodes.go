package http

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type GetMediaEpisodesUseCase struct {
	logger          logger.Logger
	mediaRepo       MediaRepository
	getMediaUseCase *GetMediaUseCase
}

func NewGetMediaEpisodesUseCase(
	logger logger.Logger,
	mediaRepo MediaRepository,
	getMediaUseCase *GetMediaUseCase,
) *GetMediaEpisodesUseCase {
	if logger == nil {
		panic("logger is nil")
	}
	if mediaRepo == nil {
		panic("mediaRepo is nil")
	}
	if getMediaUseCase == nil {
		panic("getMediaUseCase is nil")

	}
	return &GetMediaEpisodesUseCase{
		logger:          logger,
		mediaRepo:       mediaRepo,
		getMediaUseCase: getMediaUseCase,
	}
}

func (uc *GetMediaEpisodesUseCase) Execute(
	ctx context.Context,
	input dto.GetMediaEpisodesInput,
) (
	dto.GetMediaEpisodesOutput,
	*dto.Error,
) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(uc.logger, ctx, common.ContextKeyRequestID)
	log.Debug("GetMediaEpisodesUseCase called",
		log.ToAny("IsDislike", input.MediaID),
	)

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/get_media_episodes",
			entity.ErrGetMediaEpisodesParamsInvalid,
			err.Error(),
		)
		return dto.GetMediaEpisodesOutput{}, &derr
	}

	// Get media IDs from repository
	episodes, err := uc.mediaRepo.GetEpisodesByMediaID(ctx, input.MediaID)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_media_episodes",
			entity.ErrGetMediaEpisodesParamsInvalid,
			err.Error(),
		)
		log.Error("GetMediaEpisodesUseCase failed to get media IDs (episodes ids) by media ID", log.ToError(err))
		return dto.GetMediaEpisodesOutput{}, &derr
	}

	episodesOutput := make([]dto.GetMediaEpisodeOutput, len(episodes))
	for i, episode := range episodes {

		// Get media for each episode
		getMediaOutput, derr := uc.getMediaUseCase.Execute(ctx, dto.GetMediaInput{
			MediaID: episode.EpisodeID,
		})
		if derr != nil {
			log.Error("GetMediaEpisodesUseCase failed to get media details", log.ToError(err))
			continue
		}

		// Build output
		episodesOutput[i] = dto.GetMediaEpisodeOutput{
			Episode: episode,
			Media:   getMediaOutput,
		}
	}

	return dto.GetMediaEpisodesOutput{
		Episodes: episodesOutput,
	}, nil
}
