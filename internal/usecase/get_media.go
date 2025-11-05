package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type GetMediaUseCase struct {
	logger           logger.Logger
	mediaRepo        MediaRepository
	getObjectUseCase *GetObjectUseCase
}

func NewGetMediaUseCase(
	logger logger.Logger,
	mediaRepo MediaRepository,
	getObjectUseCase *GetObjectUseCase,
) *GetMediaUseCase {
	return &GetMediaUseCase{
		logger:           logger,
		mediaRepo:        mediaRepo,
		getObjectUseCase: getObjectUseCase,
	}
}

func (uc *GetMediaUseCase) Execute(ctx context.Context, input dto.GetMediaInput) (dto.GetMediaOutput, *dto.Error) {
	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/get_media",
			entity.ErrGetMediaParamsInvalid,
			err.Error(),
		)
		return dto.GetMediaOutput{}, &derr
	}

	// Get media from repository
	media, err := uc.mediaRepo.GetMedia(ctx, input.MediaID)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_media",
			err,
			"Failed to get media by ID",
		)
		uc.logger.Error("Failed to get media by ID", uc.logger.ToError(err))
		return dto.GetMediaOutput{}, &derr
	}

	// Convert entity.Media to dto.MediaOutput

	// Get genres
	genres, err := uc.mediaRepo.GetMediaGenres(ctx, media.MediaID)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_media",
			err,
			"Failed to get genres by media ID",
		)
		uc.logger.Error("Failed to get genres by media ID", uc.logger.ToError(err))
		return dto.GetMediaOutput{}, &derr
	}

	// Get posters
	postersKeys, err := uc.mediaRepo.GetMediaPostersKeys(ctx, media.MediaID)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_media",
			err,
			"Failed to get posters by media ID",
		)
		uc.logger.Error("Failed to get posters by media ID", uc.logger.ToError(err))
		return dto.GetMediaOutput{}, &derr
	}

	// Get poster links from object storage
	postersLinks := make([]string, 0, len(postersKeys))
	for _, s3key := range postersKeys {
		bucket, key := splitToBucketAndKey(s3key)
		object, err := uc.getObjectUseCase.Execute(ctx, dto.GetObjectInput{
			BucketName: bucket,
			Key:        key,
		})
		if err != nil {
			uc.logger.Error("Failed to get poster object for "+key, err)
			continue
		}
		postersLinks = append(postersLinks, object.URL)
	}

	// Get actors
	actors, err := uc.mediaRepo.GetActorsByMediaID(ctx, media.MediaID)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_media",
			err,
			"Failed to get actors by media ID",
		)
		uc.logger.Error("Failed to get actors by media ID", uc.logger.ToError(err))
		return dto.GetMediaOutput{}, &derr
	}

	// Convert genres to DTO
	genresDTO := make([]dto.GenreOutput, 0, len(genres))
	for _, genre := range genres {
		genresDTO = append(genresDTO, dto.GenreOutput{Genre: genre})
	}

	// Skip movies conversion for actors
	actorsDTO := make([]dto.GetActorOutput, 0, len(actors))
	for _, actor := range actors {
		actorsDTO = append(actorsDTO, dto.GetActorOutput{
			Actor: actor,
		})
	}
	return dto.GetMediaOutput{Media: *media, Genres: genresDTO, Posters: postersLinks, Actors: actorsDTO}, nil
}
