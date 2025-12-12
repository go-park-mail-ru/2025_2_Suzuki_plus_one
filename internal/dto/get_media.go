package dto

import "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"

type GetMediaInput struct {
	MediaID uint `json:"media_id"`
}

type GetMediaOutput struct {
	entity.Media
	Genres   []GetGenreOutput `json:"genres"`
	Posters  []string         `json:"posters"`  // S3 public URLs
	Trailers []string         `json:"trailers"` // S3 public URLs
	Rating   MediaRating      `json:"user_rating"`
}

type MediaRating struct {
	Likes    uint `json:"likes"`
	Dislikes uint `json:"dislikes"`
}
