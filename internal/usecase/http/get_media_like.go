package http

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
)

type GetMediaLikeUseCase struct {
	log      logger.Logger
	likeRepo LikeRepository
}

func NewGetMediaLikeUseCase(
	log logger.Logger,
	likeRepo LikeRepository,
) *GetMediaLikeUseCase {
	if likeRepo == nil {
		panic("likeRepo is nil")
	}
	if likeRepo == nil {
		panic("likeRepo is nil")
	}
	return &GetMediaLikeUseCase{
		log:      log,
		likeRepo: likeRepo,
	}
}

func (uc *GetMediaLikeUseCase) Execute(
	ctx context.Context,
	input dto.GetMediaLikeInput,
) (
	dto.GetMediaLikeOutput,
	*dto.Error,
) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(uc.log, ctx, common.ContextKeyRequestID)
	log.Debug("GetMediaLikeUseCase called",
		log.ToInt("media_id", int(input.MediaID)),
	)

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/get_media_like",
			entity.ErrGetMediaLikeInvalidParams,
			err.Error(),
		)
		return dto.GetMediaLikeOutput{}, &derr
	}

	// Get access token user ID
	userID, err := common.ValidateToken(input.AccessToken)
	if err != nil {
		derr := dto.NewError(
			"auth/usecase/get_auth_signout",
			entity.ErrGetAuthSignOutInvalidParams,
			"invalid access token: "+err.Error(),
		)
		return dto.GetMediaLikeOutput{}, &derr
	}

	// Call Like repository to get like status
	exists, like, err := uc.likeRepo.GetLike(ctx, userID, input.MediaID)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_media_like",
			entity.ErrGetMediaLikeRepositoryFailed,
			"like repository failed: "+err.Error(),
		)
		log.Error("GetMediaLikeUseCase failed on like repository call", log.ToError(err))
		return dto.GetMediaLikeOutput{}, &derr
	}

	// Prepare output
	output := dto.GetMediaLikeOutput{
		Liked: exists,
		IsDislike: like.IsDislike,
	}

	log.Debug("GetMediaLikeUseCase succeeded",
		log.ToAny("liked", output.Liked),
		log.ToAny("is_dislike", output.IsDislike),
	)

	return output, nil
}
