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
)
