package handlers

import (
	. "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type Handlers struct {
	Logger logger.Logger

	GetMediaRecommendationsUseCase  GetMediaRecommendationsUseCase
	GetObjectMediaUseCase           GetObjectUseCase
	PostAuthSignInUseCase           PostAuthSignInUseCase
	GetAuthRefreshUseCase           GetAuthRefreshUseCase
	PostAuthSignUpUseCase           PostAuthSignUpUseCase
	GetAuthSignOutUseCase           GetAuthSignOutUseCase
	GetUserMeUseCase                GetUserMeUseCase
	GetActorUseCase                 GetActorUseCase
	GetMediaUseCase                 GetMediaUseCase
	GetMediaWatchUseCase            GetMediaWatchUseCase
	PostUserMeUpdateUseCase         PostUserMeUpdateUseCase
	PostUserMeUpdateAvatarUseCase   PostUserMeUpdateAvatarUseCase
	GetActorMediaUseCase            GetActorMediaUseCase
	GetMediaActorUseCase            GetMediaActorUseCase
	PostUserMeUpdatePasswordUseCase PostUserMeUpdatePasswordUseCase
	GetAppealMyUseCase              GetAppealMyUseCase
	PostAppealNewUseCase            PostAppealNewUseCase
	GetAppealUseCase                GetAppealUseCase
	PutAppealResolveUseCase         PutAppealResolveUseCase
	PostAppealMessageUseCase        PostAppealMessageUseCase
	GetAppealMessageUseCase         GetAppealMessageUseCase
	GetAppealAllUseCase             GetAppealAllUseCase
	GetSearchUseCase                GetSearchUseCase
}

func NewHandlers(
	logger logger.Logger,
	GetMediaRecommendationsUseCase GetMediaRecommendationsUseCase,
	getObjectMediaUseCase GetObjectUseCase,
	postAuthSignInUseCase PostAuthSignInUseCase,
	getAuthRefreshUseCase GetAuthRefreshUseCase,
	postAuthSignupUseCase PostAuthSignUpUseCase,
	getAuthSignOutUseCase GetAuthSignOutUseCase,
	getUserMeUseCase GetUserMeUseCase,
	GetActorUseCase GetActorUseCase,
	getMediaUseCase GetMediaUseCase,
	getMediaWatchUseCase GetMediaWatchUseCase,
	PostUserMeUpdateUseCase PostUserMeUpdateUseCase,
	PostUserMeUpdateAvatarUseCase PostUserMeUpdateAvatarUseCase,
	GetActorMediaUseCase GetActorMediaUseCase,
	GetMediaActorUseCase GetMediaActorUseCase,
	PostUserMeUpdatePasswordUseCase PostUserMeUpdatePasswordUseCase,
	GetAppealMyUseCase GetAppealMyUseCase,
	PostAppealNewUseCase PostAppealNewUseCase,
	GetAppealUseCase GetAppealUseCase,
	PutAppealResolveUseCase PutAppealResolveUseCase,
	PostAppealMessageUseCase PostAppealMessageUseCase,
	GetAppealMessageUseCase GetAppealMessageUseCase,
	GetAppealAllUseCase GetAppealAllUseCase,
	GetSearchUseCase GetSearchUseCase,
) *Handlers {
	return &Handlers{
		Logger:                          logger,
		GetMediaRecommendationsUseCase:  GetMediaRecommendationsUseCase,
		GetObjectMediaUseCase:           getObjectMediaUseCase,
		PostAuthSignInUseCase:           postAuthSignInUseCase,
		GetAuthRefreshUseCase:           getAuthRefreshUseCase,
		PostAuthSignUpUseCase:           postAuthSignupUseCase,
		GetAuthSignOutUseCase:           getAuthSignOutUseCase,
		GetUserMeUseCase:                getUserMeUseCase,
		GetActorUseCase:                 GetActorUseCase,
		GetMediaUseCase:                 getMediaUseCase,
		GetMediaWatchUseCase:            getMediaWatchUseCase,
		PostUserMeUpdateUseCase:         PostUserMeUpdateUseCase,
		PostUserMeUpdateAvatarUseCase:   PostUserMeUpdateAvatarUseCase,
		GetActorMediaUseCase:            GetActorMediaUseCase,
		GetMediaActorUseCase:            GetMediaActorUseCase,
		PostUserMeUpdatePasswordUseCase: PostUserMeUpdatePasswordUseCase,
		GetAppealMyUseCase:              GetAppealMyUseCase,
		PostAppealNewUseCase:            PostAppealNewUseCase,
		GetAppealUseCase:                GetAppealUseCase,
		PutAppealResolveUseCase:         PutAppealResolveUseCase,
		PostAppealMessageUseCase:        PostAppealMessageUseCase,
		GetAppealMessageUseCase:         GetAppealMessageUseCase,
		GetAppealAllUseCase:             GetAppealAllUseCase,
		GetSearchUseCase:                GetSearchUseCase,
	}
}
