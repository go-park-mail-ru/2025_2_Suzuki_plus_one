package entity

import "errors"

var (
	ErrGetMovieRecommendationsParamsInvalid = errors.New("invalid get_movie_recommendations request parameters")
	ErrGetObjectParamsInvalid               = errors.New("invalid get_object request parameters")
	ErrGetObjectFailed                      = errors.New("failed to get object from repository")
)
