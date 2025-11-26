package http

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type DeleteMediaLikeUseCase struct {
	log      logger.Logger
	likeRepo LikeRepository
}

func NewDeleteMediaLikeUseCase(
	log logger.Logger,
	likeRepo LikeRepository,
) *DeleteMediaLikeUseCase {
	if log == nil {
		panic("log is nil")
	}
	if likeRepo == nil {
		panic("likeRepo is nil")
	}
	return &DeleteMediaLikeUseCase{
		log:      log,
		likeRepo: likeRepo,
	}
}

func (uc *DeleteMediaLikeUseCase) Execute(
	ctx context.Context,
	input dto.DeleteMediaLikeInput,
) (
	dto.DeleteMediaLikeOutput,
	*dto.Error,
) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(uc.log, ctx, common.ContextKeyRequestID)
	log.Debug("DeleteMediaLikeUseCase called",
		log.ToInt("media_id", int(input.MediaID)),
	)

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/get_media_like",
			entity.ErrDeleteMediaLikeInvalidParams,
			err.Error(),
		)
		return dto.DeleteMediaLikeOutput{}, &derr
	}

	// Get access token user ID
	userID, err := common.ValidateToken(input.AccessToken)
	if err != nil {
		derr := dto.NewError(
			"auth/usecase/get_auth_signout",
			entity.ErrGetAuthSignOutInvalidParams,
			"invalid access token: "+err.Error(),
		)
		return dto.DeleteMediaLikeOutput{}, &derr
	}

	// Call Like repository to get like status
	err = uc.likeRepo.DeleteLike(ctx, userID, uint(input.MediaID))
	if err != nil {
		derr := dto.NewError(
			"usecase/delete_media_like",
			entity.ErrDeleteMediaLikeRepositoryFailed,
			"like repository failed to delete like status: "+err.Error(),
		)
		log.Error("DeleteMediaLikeUseCase failed on like repository call", log.ToError(err))
		return dto.DeleteMediaLikeOutput{}, &derr
	}

	// Prepare output
	output := dto.DeleteMediaLikeOutput{}

	log.Debug("DeleteMediaLikeUseCase succeeded")

	return output, nil
}
