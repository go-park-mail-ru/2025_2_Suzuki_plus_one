package postgres

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// ASSET

func (db *DataBase) CreateAsset(ctx context.Context, asset entity.Asset) (uint, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("CreateAsset called",
		log.ToString("s3_key", asset.S3Key),
		log.ToString("mime_type", asset.MimeType),
		log.ToAny("file_size_mb", asset.FileSizeMB),
	)

	var assetID uint
	query := `
		INSERT INTO asset (s3_key, mime_type, file_size_mb)
		VALUES ($1, $2, $3)
		RETURNING asset_id
	`
	err := db.conn.QueryRowContext(ctx, query, asset.S3Key, asset.MimeType, asset.FileSizeMB).Scan(&assetID)
	if err != nil {
		log.Error("Failed to create asset: " + err.Error())
		return 0, err
	}
	return assetID, nil
}

func (db *DataBase) GetAssetByID(ctx context.Context, assetID uint) (*entity.Asset, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("GetAssetByID called",
		log.ToInt("asset_id", int(assetID)),
	)

	var asset entity.Asset
	query := `
		SELECT asset_id, s3_key, mime_type, file_size_mb
		FROM asset
		WHERE asset_id = $1
	`
	err := db.conn.QueryRowContext(ctx, query, assetID).Scan(
		&asset.ID,
		&asset.S3Key,
		&asset.MimeType,
		&asset.FileSizeMB,
	)
	if err != nil {
		log.Error("Failed to get asset by ID: " + err.Error())
		return nil, err
	}
	return &asset, nil
}

// ASSET IMAGE

func (db *DataBase) CreateAssetImage(ctx context.Context, assetImage entity.AssetImage) (uint, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("CreateAssetImage called",
		log.ToInt("asset_id", int(assetImage.AssetID)),
	)

	var assetImageID uint
	query := `
		INSERT INTO asset_image (asset_id, resolution_width, resolution_height)
		VALUES ($1, $2, $3)
		RETURNING asset_image_id
	`

	// Scan the resulting asset_image_id into assetImageID
	err := db.conn.QueryRowContext(
		ctx,
		query,
		assetImage.AssetID,
		assetImage.ResolutionWidth,
		assetImage.ResolutionHeight,
	).Scan(&assetImageID)

	if err != nil {
		log.Error("Failed to create asset image: " + err.Error())
		return 0, err
	}
	return assetImageID, nil
}

func (db *DataBase) GetAssetImageByID(ctx context.Context, assetImageID uint) (*entity.AssetImage, error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(db.logger, ctx, common.ContextKeyRequestID)
	log.Debug("GetAssetImageByID called",
		log.ToInt("asset_image_id", int(assetImageID)),
	)

	var assetImage entity.AssetImage
	query := `
		SELECT asset_image_id, asset_id, resolution_width, resolution_height
		FROM asset_image
		WHERE asset_image_id = $1
	`
	err := db.conn.QueryRowContext(ctx, query, assetImageID).Scan(
		&assetImage.ID,
		&assetImage.AssetID,
		&assetImage.ResolutionWidth,
		&assetImage.ResolutionHeight,
	)
	if err != nil {
		log.Error("Failed to get asset image by ID: " + err.Error())
		return nil, err
	}
	return &assetImage, nil
}
