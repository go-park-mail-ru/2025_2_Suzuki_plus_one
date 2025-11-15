package controller

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
)

//go:generate mockgen -source=contract.go -destination=./mocks/contract_mock.go -package=mocks
type (
	GetMediaRecommendationsUseCase interface {
		Execute(context.Context, dto.GetMediaRecommendationsInput) (dto.GetMediaRecommendationsOutput, *dto.Error)
	}

	GetObjectUseCase interface {
		Execute(context.Context, dto.GetObjectInput) (dto.GetObjectOutput, *dto.Error)
	}

	PostAuthSignInUseCase interface {
		Execute(context.Context, dto.PostAuthSignInInput) (dto.PostAuthSignInOutput, *dto.Error)
	}

	GetAuthRefreshUseCase interface {
		Execute(context.Context, dto.GetAuthRefreshInput) (dto.GetAuthRefreshOutput, *dto.Error)
	}

	PostAuthSignUpUseCase interface {
		Execute(context.Context, dto.PostAuthSignUpInput) (dto.PostAuthSignUpOutput, *dto.Error)
	}

	GetAuthSignOutUseCase interface {
		Execute(context.Context, dto.GetAuthSignOutInput) (dto.GetAuthSignOutOutput, *dto.Error)
	}

	GetUserMeUseCase interface {
		Execute(context.Context, dto.GetUserMeInput) (dto.GetUserMeOutput, *dto.Error)
	}

	GetActorUseCase interface {
		Execute(context.Context, dto.GetActorInput) (dto.GetActorOutput, *dto.Error)
	}

	GetMediaUseCase interface {
		Execute(context.Context, dto.GetMediaInput) (dto.GetMediaOutput, *dto.Error)
	}

	GetMediaWatchUseCase interface {
		Execute(context.Context, dto.GetMediaWatchInput) (dto.GetMediaWatchOutput, *dto.Error)
	}

	PostUserMeUpdateUseCase interface {
		Execute(context.Context, dto.PostUserMeUpdateInput) (dto.PostUserMeUpdateOutput, *dto.Error)
	}

	PostUserMeUpdateAvatarUseCase interface {
		Execute(context.Context, dto.PostUserMeUpdateAvatarInput) (dto.PostUserMeUpdateAvatarOutput, *dto.Error)
	}

	GetActorMediaUseCase interface {
		Execute(context.Context, dto.GetActorMediaInput) (dto.GetActorMediaOutput, *dto.Error)
	}

	GetMediaActorUseCase interface {
		Execute(context.Context, dto.GetMediaActorInput) (dto.GetMediaActorOutput, *dto.Error)
	}

	PostUserMeUpdatePasswordUseCase interface {
		Execute(context.Context, dto.PostUserMeUpdatePasswordInput) (dto.PostUserMeUpdatePasswordOutput, *dto.Error)
	}

	GetAppealMyUseCase interface {
		Execute(context.Context, dto.GetAppealMyInput) (dto.GetAppealMyOutput, *dto.Error)
	}

	PostAppealNewUseCase interface {
		Execute(context.Context, dto.PostAppealNewInput) (dto.PostAppealNewOutput, *dto.Error)
	}

	GetAppealUseCase interface {
		Execute(context.Context, dto.GetAppealInput) (dto.GetAppealOutput, *dto.Error)
	}

	PutAppealResolveUseCase interface {
		Execute(context.Context, dto.PutAppealResolveInput) (dto.PutAppealResolveOutput, *dto.Error)
	}

	PostAppealMessageUseCase interface {
		Execute(context.Context, dto.PostAppealMessageInput) (dto.PostAppealMessageOutput, *dto.Error)
	}

	GetAppealMessageUseCase interface {
		Execute(context.Context, dto.GetAppealMessageInput) (dto.GetAppealMessageOutput, *dto.Error)
	}
)
