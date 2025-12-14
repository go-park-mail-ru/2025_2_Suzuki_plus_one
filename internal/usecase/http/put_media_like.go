package http

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type PutMediaLikeUseCase struct {
	log      logger.Logger
	likeRepo LikeRepository
}

func NewPutMediaLikeUseCase(
	log logger.Logger,
	likeRepo LikeRepository,
) *PutMediaLikeUseCase {
	if log == nil {
		panic("log is nil")
	}
	if likeRepo == nil {
		panic("likeRepo is nil")
	}
	return &PutMediaLikeUseCase{
		log:      log,
		likeRepo: likeRepo,
	}
}

func (uc *PutMediaLikeUseCase) Execute(
	ctx context.Context,
	input dto.PutMediaLikeInput,
) (
	dto.PutMediaLikeOutput,
	*dto.Error,
) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(uc.log, ctx, common.ContextKeyRequestID)
	log.Debug("PutMediaLikeUseCase called",
		log.ToInt("media_id", int(input.MediaID)),
	)

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/get_media_like",
			entity.ErrPutMediaLikeInvalidParams,
			err.Error(),
		)
		return dto.PutMediaLikeOutput{}, &derr
	}

	// Get access token user ID
	userID, err := common.ValidateToken(input.AccessToken)
	if err != nil {
		derr := dto.NewError(
			"auth/usecase/get_auth_signout",
			entity.ErrGetAuthSignOutInvalidParams,
			"invalid access token: "+err.Error(),
		)
		return dto.PutMediaLikeOutput{}, &derr
	}

	// Call Like repository to set like status
	isDislike, err := uc.likeRepo.ToggleLike(ctx, userID, input.MediaID)
	if err != nil {
		derr := dto.NewError(
			"usecase/put_media_like",
			entity.ErrPutMediaLikeRepositoryFailed,
			"like repository failed to toggle like status: "+err.Error(),
		)
		log.Error("PutMediaLikeUseCase failed on like repository call", log.ToError(err))
		return dto.PutMediaLikeOutput{}, &derr
	}

	// toggle again to match the requested state
	if isDislike != input.IsDislike {
		isDislike, err = uc.likeRepo.ToggleLike(ctx, userID, input.MediaID)
		if err != nil {
			derr := dto.NewError(
				"usecase/put_media_like",
				entity.ErrPutMediaLikeRepositoryFailed,
				"like repository failed to toggle like status: "+err.Error(),
			)
			log.Error("PutMediaLikeUseCase failed on like repository call", log.ToError(err))
			return dto.PutMediaLikeOutput{}, &derr
		}
	}

	// Prepare output
	output := dto.PutMediaLikeOutput{
		Liked:     true,
		IsDislike: isDislike,
	}

	log.Debug("PutMediaLikeUseCase succeeded",
		log.ToAny("liked", output.Liked),
		log.ToAny("is_dislike", output.IsDislike),
	)

	return output, nil
}
