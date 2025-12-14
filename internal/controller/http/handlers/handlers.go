package handlers

import (
	ctrlhttp "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/controller/http"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type Handlers struct {
	Logger logger.Logger

	GetMediaRecommendationsUseCase  ctrlhttp.GetMediaRecommendationsUseCase
	GetObjectMediaUseCase           ctrlhttp.GetObjectUseCase
	PostAuthSignInUseCase           ctrlhttp.PostAuthSignInUseCase
	GetAuthRefreshUseCase           ctrlhttp.GetAuthRefreshUseCase
	PostAuthSignUpUseCase           ctrlhttp.PostAuthSignUpUseCase
	GetAuthSignOutUseCase           ctrlhttp.GetAuthSignOutUseCase
	GetUserMeUseCase                ctrlhttp.GetUserMeUseCase
	GetActorUseCase                 ctrlhttp.GetActorUseCase
	GetMediaUseCase                 ctrlhttp.GetMediaUseCase
	GetMediaWatchUseCase            ctrlhttp.GetMediaWatchUseCase
	PostUserMeUpdateUseCase         ctrlhttp.PostUserMeUpdateUseCase
	PostUserMeUpdateAvatarUseCase   ctrlhttp.PostUserMeUpdateAvatarUseCase
	GetActorMediaUseCase            ctrlhttp.GetActorMediaUseCase
	GetMediaActorUseCase            ctrlhttp.GetMediaActorUseCase
	PostUserMeUpdatePasswordUseCase ctrlhttp.PostUserMeUpdatePasswordUseCase
	GetAppealMyUseCase              ctrlhttp.GetAppealMyUseCase
	PostAppealNewUseCase            ctrlhttp.PostAppealNewUseCase
	GetAppealUseCase                ctrlhttp.GetAppealUseCase
	PutAppealResolveUseCase         ctrlhttp.PutAppealResolveUseCase
	PostAppealMessageUseCase        ctrlhttp.PostAppealMessageUseCase
	GetAppealMessageUseCase         ctrlhttp.GetAppealMessageUseCase
	GetAppealAllUseCase             ctrlhttp.GetAppealAllUseCase
	GetSearchUseCase                ctrlhttp.GetSearchUseCase
	GetMediaLikeUseCase             ctrlhttp.GetMediaLikeUseCase
	PutMediaLikeUseCase             ctrlhttp.PutMediaLikeUseCase
	DeleteMediaLikeUseCase          ctrlhttp.DeleteMediaLikeUseCase
	GetMediaMyUseCase               ctrlhttp.GetMediaMyUseCase
	GetGenreUseCase                 ctrlhttp.GetGenreUseCase
	GetGenreAllUseCase              ctrlhttp.GetGenreAllUseCase
	GetGenreMediaUseCase            ctrlhttp.GetGenreMediaUseCase
	GetMediaEpisodesUseCase         ctrlhttp.GetMediaEpisodesUseCase
	PostPaymentCompletedUseCase     ctrlhttp.PostPaymentCompletedUseCase
	PostPaymentNewUseCase           ctrlhttp.PostPaymentNewUseCase
}

func NewHandlers(
	logger logger.Logger,
	GetMediaRecommendationsUseCase ctrlhttp.GetMediaRecommendationsUseCase,
	getObjectMediaUseCase ctrlhttp.GetObjectUseCase,
	postAuthSignInUseCase ctrlhttp.PostAuthSignInUseCase,
	getAuthRefreshUseCase ctrlhttp.GetAuthRefreshUseCase,
	postAuthSignupUseCase ctrlhttp.PostAuthSignUpUseCase,
	getAuthSignOutUseCase ctrlhttp.GetAuthSignOutUseCase,
	getUserMeUseCase ctrlhttp.GetUserMeUseCase,
	GetActorUseCase ctrlhttp.GetActorUseCase,
	getMediaUseCase ctrlhttp.GetMediaUseCase,
	getMediaWatchUseCase ctrlhttp.GetMediaWatchUseCase,
	PostUserMeUpdateUseCase ctrlhttp.PostUserMeUpdateUseCase,
	PostUserMeUpdateAvatarUseCase ctrlhttp.PostUserMeUpdateAvatarUseCase,
	GetActorMediaUseCase ctrlhttp.GetActorMediaUseCase,
	GetMediaActorUseCase ctrlhttp.GetMediaActorUseCase,
	PostUserMeUpdatePasswordUseCase ctrlhttp.PostUserMeUpdatePasswordUseCase,
	GetAppealMyUseCase ctrlhttp.GetAppealMyUseCase,
	PostAppealNewUseCase ctrlhttp.PostAppealNewUseCase,
	GetAppealUseCase ctrlhttp.GetAppealUseCase,
	PutAppealResolveUseCase ctrlhttp.PutAppealResolveUseCase,
	PostAppealMessageUseCase ctrlhttp.PostAppealMessageUseCase,
	GetAppealMessageUseCase ctrlhttp.GetAppealMessageUseCase,
	GetAppealAllUseCase ctrlhttp.GetAppealAllUseCase,
	GetSearchUseCase ctrlhttp.GetSearchUseCase,
	GetMediaLikeUseCase ctrlhttp.GetMediaLikeUseCase,
	PutMediaLikeUseCase ctrlhttp.PutMediaLikeUseCase,
	DeleteMediaLikeUseCase ctrlhttp.DeleteMediaLikeUseCase,
	GetMediaMyUseCase ctrlhttp.GetMediaMyUseCase,
	GetGenreUseCase ctrlhttp.GetGenreUseCase,
	GetGenreAllUseCase ctrlhttp.GetGenreAllUseCase,
	GetGenreMediaUseCase ctrlhttp.GetGenreMediaUseCase,
	GetMediaEpisodesUseCase ctrlhttp.GetMediaEpisodesUseCase,
	PostPaymentCompletedUseCase ctrlhttp.PostPaymentCompletedUseCase,
	PostPaymentNewUseCase ctrlhttp.PostPaymentNewUseCase,
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
		GetMediaLikeUseCase:             GetMediaLikeUseCase,
		PutMediaLikeUseCase:             PutMediaLikeUseCase,
		DeleteMediaLikeUseCase:          DeleteMediaLikeUseCase,
		GetMediaMyUseCase:               GetMediaMyUseCase,
		GetGenreUseCase:                 GetGenreUseCase,
		GetGenreAllUseCase:              GetGenreAllUseCase,
		GetGenreMediaUseCase:            GetGenreMediaUseCase,
		GetMediaEpisodesUseCase:         GetMediaEpisodesUseCase,
		PostPaymentCompletedUseCase:     PostPaymentCompletedUseCase,
		PostPaymentNewUseCase:           PostPaymentNewUseCase,
	}
}
