package entity

import "errors"

var (
	ErrGetMovieRecommendationsParamsInvalid = errors.New("invalid get_movie_recommendations request parameters")
	ErrGetObjectParamsInvalid               = errors.New("invalid get_object request parameters")
	ErrGetObjectFailed                      = errors.New("failed to get object from repository")
	ErrPostAuthSignInParamsInvalid          = errors.New("invalid post_auth_sign_in request parameters")
	ErrPostAuthSignInNewRefreshTokenFailed  = errors.New("failed to create new refresh token")
	ErrPostAuthSignInNewAccessTokenFailed   = errors.New("failed to create new access token")
	ErrPostAuthSignInAddSessionFailed       = errors.New("failed to add session to session repository")

	ErrUserNotFound = errors.New("user not found")

	ErrGetAuthRefreshInvalidParams = errors.New("invalid get_auth_refresh request parameters")
)
