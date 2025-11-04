package entity

import "errors"

var (
	// ErrNoTokenProvided  = models.ErrorResponse{Type: "auth", Message: "no token provided"}
	// ErrInvalidOrExpired = models.ErrorResponse{Type: "auth", Message: "invalid or expired token"}
	// ErrSignInWrongData  = models.ErrorResponse{Type: "auth", Message: "wrong email or password"}
	// ErrSignInInternal   = models.ErrorResponse{Type: "auth", Message: "internal server error"}
	// ErrSignUpWrongData  = models.ErrorResponse{Type: "auth", Message: "wrong email or password"}
	// ErrSignUpUserExists = models.ErrorResponse{Type: "auth", Message: "user does not exist"}
	// ErrSignUpInternal   = models.ErrorResponse{Type: "auth", Message: "internal server error"}

	// ErrMoviesInvalidParams = models.ErrorResponse{Type: "content", Message: "invalid request parameters"}

	ErrGetMovieRecommendationsParamsInvalid = errors.New("invalid get_movie_recommendations request parameters")
)
