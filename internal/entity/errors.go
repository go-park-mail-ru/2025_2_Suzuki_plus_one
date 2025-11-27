package entity

import "errors"

var (
	ErrGetMediaRecommendationsParamsInvalid = errors.New("invalid get_movie_recommendations request parameters")
	ErrGetMediaRecommendationsCantReturnAll = errors.New("can't return too many movies, specify limit and offset")

	ErrGetObjectParamsInvalid = errors.New("invalid get_object request parameters")
	ErrGetObjectFailed        = errors.New("failed to get object from repository")

	ErrPostAuthSignInParamsInvalid         = errors.New("invalid email or password")
	ErrPostAuthSignInNewRefreshTokenFailed = errors.New("failed to create new refresh token")
	ErrPostAuthSignInNewAccessTokenFailed  = errors.New("failed to create new access token")
	ErrPostAuthSignInAddSessionFailed      = errors.New("failed to add session to session repository")

	ErrPostAuthSignUpParamsInvalid = errors.New("invalid post_auth_sign_up request parameters")
	ErrPostAuthSignUpAlreadyExists = errors.New("user with given email already exists")

	ErrGetAuthSignOutInvalidParams = errors.New("invalid get_auth_signout request parameters")

	ErrUserNotFound = errors.New("user not found")

	ErrGetAuthRefreshInvalidParams = errors.New("invalid get_auth_refresh request parameters")

	ErrGetUserMeParamsInvalid   = errors.New("invalid get_user_me request parameters")
	ErrGetUserMeSessionNotFound = errors.New("session not found for given access token")

	ErrGetActorParamsInvalid = errors.New("invalid get_actor request parameters")

	ErrGetMediaParamsInvalid = errors.New("invalid get_media request parameters")

	ErrGetMediaPosterFailed = errors.New("failed to get media poster from object storage")

	ErrPostUserMeUpdateAvatarParamsInvalid = errors.New("invalid post_user_me_update_avatar request parameters")

	ErrSessionNotFound = errors.New("session not found for given access token")

	ErrGetMediaActorParamsInvalid = errors.New("invalid get_media_actor request parameters")

	ErrPostUserMeUpdatePasswordCurrentPasswordMismatch = errors.New("current password does not match")

	ErrWrongAccessToken = errors.New("wrong access token provided")

	ErrGetAppealMyFailed = errors.New("failed to get appeals for user")

	ErrGetAppealAllFailed = errors.New("failed to get all appeals")

	ErrPostAppealMessageInvalidParams = errors.New("invalid post_appeal_message request parameters")

	ErrGetAppealMessageInvalidParams = errors.New("invalid get_appeal_message request parameters")
	ErrGetAppealFailed               = errors.New("failed to get details of specific appeal")

	ErrPostAppealNew = errors.New("failed to create a new appeal")

	ErrPutAppealResolve = errors.New("failed to resolve an appeal")

	ErrGetSearchParamsInvalid = errors.New("invalid get_search request parameters")
	ErrGetMediaMyParamsInvalid = errors.New("invalid get_media_my request parameters")

	ErrGetGenreInvalidParams = errors.New("invalid get_genre request parameters")
	ErrGetGenreRepo = errors.New("failed to get genre from repository")
	ErrGetGenreNotFound = errors.New("genre not found")

	ErrGetGenreAllInvalidParams = errors.New("invalid get_genre_all request parameters")
	ErrGetGenreMediaFailed = errors.New("failed to get media related to genre")
	ErrGetAllGenresFailed = errors.New("failed to get all genres")	

	ErrGetMediaEpisodesParamsInvalid = errors.New("invalid get_media_episodes request parameters")

	// Like usecase errors
	ErrGetMediaLikeInvalidParams = errors.New("invalid get_media_like request parameters")
	ErrGetMediaLikeRepositoryFailed = errors.New("like repository failed to get like status")
	ErrPutMediaLikeInvalidParams = errors.New("invalid put_media_like request parameters")
	ErrPutMediaLikeRepositoryFailed = errors.New("like repository failed to toggle like status")
	ErrDeleteMediaLikeInvalidParams = errors.New("invalid delete_media_like request parameters")
	ErrDeleteMediaLikeRepositoryFailed = errors.New("like repository failed to toggle like status")

	// Authentication service errors
	ErrPostAuthSignInAuthServiceFailed = errors.New("authentication service login call failed")
	ErrGetAuthRefreshAuthServiceFailed = errors.New("authentication service refresh call failed")
	ErrGetAuthSignOutAuthServiceFailed = errors.New("authentication service logout call failed")
	ErrPostAuthSignUpAuthServiceFailed = errors.New("authentication service signup call failed")

	// Search service errors
	ErrGetSearchSearchServiceFailed = errors.New("search service call failed")
	ErrorSearchMediaFailed          = errors.New("search media failed")
	ErrorSearchActorFailed          = errors.New("search actor failed")
)
