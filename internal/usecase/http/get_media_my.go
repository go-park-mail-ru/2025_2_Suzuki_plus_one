package http

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

const GET_MEDIA_MY_LIMIT_MIN = 1
const GET_MEDIA_MY_LIMIT_MAX = 20

type GetMediaMyUseCase struct {
	logger          logger.Logger
	mediaRepo       MediaRepository
	getMediaUseCase *GetMediaUseCase
}

func NewGetMediaMyUseCase(
	logger logger.Logger,
	mediaRepo MediaRepository,
	getMediaUseCase *GetMediaUseCase,
) *GetMediaMyUseCase {
	if logger == nil {
		panic("logger is nil")
	}
	if mediaRepo == nil {
		panic("mediaRepo is nil")
	}
	if getMediaUseCase == nil {
		panic("getMediaUseCase is nil")

	}
	return &GetMediaMyUseCase{
		logger:          logger,
		mediaRepo:       mediaRepo,
		getMediaUseCase: getMediaUseCase,
	}
}

func (uc *GetMediaMyUseCase) Execute(
	ctx context.Context,
	input dto.GetMediaMyInput,
) (
	dto.GetMediaMyOutput,
	*dto.Error,
) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(uc.logger, ctx, common.ContextKeyRequestID)
	log.Debug("GetMediaMyUseCase called",
		log.ToString("access_token", input.AccessToken),
		log.ToAny("IsDislike", input.IsDislike),
	)

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/get_media_my",
			entity.ErrGetMediaMyParamsInvalid,
			err.Error(),
		)
		return dto.GetMediaMyOutput{}, &derr
	}

	// Get user ID from access token
	userID, err := common.ValidateToken(input.AccessToken)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_media_my",
			entity.ErrGetAuthSignOutInvalidParams,
			"invalid access token: "+err.Error(),
		)
		log.Error("GetMediaMyUseCase failed to validate access token", log.ToError(err))
		return dto.GetMediaMyOutput{}, &derr
	}

	// Get media IDs from repository
	mediaIDs, err := uc.mediaRepo.GetMediaIDsByLikeStatus(ctx, userID, input.IsDislike, input.Limit, input.Offset)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_media_my",
			entity.ErrGetMediaMyParamsInvalid,
			err.Error(),
		)
		log.Error("GetMediaMyUseCase failed to get media IDs by like status", log.ToError(err))
		return dto.GetMediaMyOutput{}, &derr
	}

	medias := make([]dto.GetMediaOutput, 0, len(mediaIDs))
	for _, mediaID := range mediaIDs {
		getMediaOutput, derr := uc.getMediaUseCase.Execute(ctx, dto.GetMediaInput{
			MediaID: uint(mediaID),
		})
		if derr != nil {
			log.Error("GetMediaMyUseCase failed to get media details", log.ToError(err))
			continue
		}
		medias = append(medias, getMediaOutput)
	}

	return dto.GetMediaMyOutput{
		Medias: medias,
	}, nil
}
