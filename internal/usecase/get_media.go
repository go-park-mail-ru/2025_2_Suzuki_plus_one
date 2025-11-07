package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type GetMediaUseCase struct {
	logger           logger.Logger
	mediaRepo        MediaRepository
	actorRepo        ActorRepository
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
	// Bind logger with request ID
	log := logger.LoggerWithKey(uc.logger, ctx, common.ContextKeyRequestID)

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
	media, err := uc.mediaRepo.GetMediaByID(ctx, input.MediaID)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_media",
			err,
			"Failed to get media by ID",
		)
		log.Error("Failed to get media by ID", log.ToError(err))
		return dto.GetMediaOutput{}, &derr
	}

	// Get genres
	genres, err := uc.mediaRepo.GetMediaGenres(ctx, media.MediaID)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_media",
			err,
			"Failed to get genres by media ID",
		)
		log.Error("Failed to get genres by media ID", log.ToError(err))
		return dto.GetMediaOutput{}, &derr
	}

	// Convert genres to DTO
	genresDTO := make([]dto.GenreOutput, 0, len(genres))
	for _, genre := range genres {
		genresDTO = append(genresDTO, dto.GenreOutput{Genre: genre})
	}

	// Get posters
	postersKeys, err := uc.mediaRepo.GetMediaPostersKeys(ctx, media.MediaID)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_media",
			err,
			"Failed to get posters by media ID",
		)
		log.Error("Failed to get posters by media ID", log.ToError(err))
		return dto.GetMediaOutput{}, &derr
	}

	// Get poster links from object storage
	postersLinks := make([]string, 0, len(postersKeys))
	for _, s3key := range postersKeys {
		object, err := uc.getObjectUseCase.Execute(ctx, dto.GetObjectInput{
			BucketName: s3key.BucketName,
			Key:        s3key.Key,
		})
		if err != nil {
			log.Error("Failed to get poster object for "+s3key.GetPath(), err)
			continue
		}
		postersLinks = append(postersLinks, object.URL)
	}

	// Get trailers
	trailersKeys, err := uc.mediaRepo.GetMediaTrailersKeys(ctx, media.MediaID)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_media",
			err,
			"Failed to get trailers by media ID",
		)
		log.Error("Failed to get trailers by media ID", log.ToError(err))
		return dto.GetMediaOutput{}, &derr
	}

	// Get trailer links from object storage
	trailersLinks := make([]string, 0, len(trailersKeys))
	for _, s3key := range trailersKeys {
		object, err := uc.getObjectUseCase.Execute(ctx, dto.GetObjectInput{
			BucketName: s3key.BucketName,
			Key:        s3key.Key,
		})
		if err != nil {
			log.Error("Failed to get trailer object for "+s3key.GetPath(), err)
			continue
		}
		trailersLinks = append(trailersLinks, object.URL)
	}

	return dto.GetMediaOutput{Media: *media, Genres: genresDTO, Posters: postersLinks, Trailers: trailersLinks}, nil
}
