package handlers

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

// All possible http responses for PostUserMeUpdateAvatar handler
// Blank fields are not used and are filled in the handler
var (
	ErrPostUserMeUpdateAvatarInvalidParams = ResponseError{
		Code:    http.StatusBadRequest,
		Message: errors.New("invalid parameters for actor"),
	}
	ErrPostUserMeUpdateAvatarNotFound = ResponseError{
		Code: http.StatusBadRequest,
	}
	ResponsePostUserMeUpdateAvatar = Response{
		Code: http.StatusOK,
	}
)

// Input file parameter name
const FileParamPostUserMeUpdateAvatar = "avatar"

// Get all movies from database
func (h *Handlers) PostUserMeUpdateAvatar(w http.ResponseWriter, r *http.Request) {
	// Extract context, bind logger with request ID
	ctx := common.GetContext(r)
	log := logger.LoggerWithKey(h.Logger, ctx, common.ContextKeyRequestID)
	log.Debug("Handler called")

	// Handle input parameters
	input := dto.PostUserMeUpdateAvatarInput{}
	rp := NewRequestParams(log, r, &input)
	rp.AddAuthHeader(&input.AccessToken)
	rp.AddFile(FileParamPostUserMeUpdateAvatar, &input.Bytes)

	// Parse request parameters
	if err := rp.Parse(); err != nil {
		log.Error(
			"Failed to parse query parameters",
			log.ToString("error", err.Error()),
		)
		// Respond with error, if input parameters are invalid
		RespondWithError(log, w, ErrPostUserMeUpdateAvatarInvalidParams, err.Error())
		return
	}

	// MIME type (detected from file bytes)
	input.MimeFormat = http.DetectContentType(input.Bytes)
	// File size in MB
	input.FileSizeMB = float32(len(input.Bytes)) / (1024.0 * 1024.0)

	log.Info("Attempt to upload file",
		log.ToString("mime_type", input.MimeFormat),
		log.ToAny("file_size_mb", input.FileSizeMB),
	)

	// Execute use case
	output, err := h.PostUserMeUpdateAvatarUseCase.Execute(ctx, input)
	if err != nil {
		log.Error("Failed to update user avatar",
			log.ToString("error", err.Message),
		)
		// Respond with error, if use case execution fails
		RespondWithDTOError(log, w, ErrPostUserMeUpdateAvatarInvalidParams, err)
		return
	}

	// Log successful completion and print user avatarURL
	log.Debug(
		"Avatar uploaded successfully",
		log.ToAny("user", output),
	)

	// Respond with output
	Respond(log, w, ResponsePostUserMeUpdateAvatar.Code, output)
}
