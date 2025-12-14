package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type RequestParams struct {
	logger              logger.Logger
	queryParams         []string
	queryParamsStorage  []any
	bodyParams          []string
	bodyParamsStorage   []any
	cookieParams        []string
	cookieParamsStorage []any
	accessToken         string
	accessTokenStorage  *string
	pathParams          []string
	pathParamsStorage   []any
	fileParams          []string
	fileParamsStorage   []*[]byte
	request             *http.Request
	dto                 *dto.DTO
}

func NewRequestParams(l logger.Logger, request *http.Request, dto dto.DTO) *RequestParams {
	return &RequestParams{
		logger:              l,
		queryParams:         make([]string, 0),
		queryParamsStorage:  make([]any, 0),
		bodyParams:          make([]string, 0),
		bodyParamsStorage:   make([]any, 0),
		cookieParams:        make([]string, 0),
		cookieParamsStorage: make([]any, 0),
		pathParams:          make([]string, 0),
		pathParamsStorage:   make([]any, 0),
		fileParams:          make([]string, 0),
		fileParamsStorage:   make([]*[]byte, 0),
		request:             request,
		dto:                 &dto,
	}
}

// Register a query parameter to be parsed with [RequestParams.Parse] into valueStorage (dto field).
// If the parameter is not present in the request, the value in valueStorage becomes zero value.
func (rp *RequestParams) AddQuery(key string, valueStorage any) {
	rp.queryParams = append(rp.queryParams, key)
	rp.queryParamsStorage = append(rp.queryParamsStorage, valueStorage)
}

// This is unused in favor of full dto param.
// However, it is needed to call it for all input body parameters to mark them for parsing.
func (rp *RequestParams) AddBody(key string, valueStorage any) {
	rp.bodyParams = append(rp.bodyParams, key)
	rp.bodyParamsStorage = append(rp.bodyParamsStorage, valueStorage)
}

// Adds a cookie parameter to be parsed with [RequestParams.Parse] into valueStorage (dto field).
func (rp *RequestParams) AddCookie(key string, valueStorage any) {
	rp.cookieParams = append(rp.cookieParams, key)
	rp.cookieParamsStorage = append(rp.cookieParamsStorage, valueStorage)
}

// Parses Authorization header Bearer token into internal string
// Parse methiod have to be called to scan it into given storage.
func (rp *RequestParams) AddAuthHeader(valueStorage *string) {
	rp.accessToken = jwtauth.TokenFromHeader(rp.request)
	rp.accessTokenStorage = valueStorage
}

// Registers a path parameter to be parsed with [RequestParams.Parse] into valueStorage (dto field).
func (rp *RequestParams) AddPath(key string, valueStorage any) {
	rp.pathParams = append(rp.pathParams, key)
	rp.pathParamsStorage = append(rp.pathParamsStorage, valueStorage)
}

func (rp *RequestParams) AddFile(key string, valueStorage *[]byte) {
	rp.fileParams = append(rp.fileParams, key)
	rp.fileParamsStorage = append(rp.fileParamsStorage, valueStorage)
}

// Parse all registered parameters from the request into their storages.
func (rp *RequestParams) Parse() error {
	// Check request_id exists
	requestID := middleware.GetReqID(rp.request.Context())
	if requestID == "" {
		rp.logger.Error("Can't get requestID",
			rp.logger.ToString("requestURI", rp.request.URL.String()))
		return errors.New("no request ID found: " + rp.request.URL.String())
	}

	query := rp.request.URL.Query()
	for i := range rp.queryParams {
		param := rp.queryParams[i]
		storage := rp.queryParamsStorage[i]

		// Use all occurrences of the parameter (supports repeated keys)
		vals := query[param]

		// If there are no repeated params but there is a single comma-separated value,
		// the helper will split it into elements.
		if len(vals) == 0 {
			// nothing provided -> leave zero value
			continue
		}

		if err := parseAndSetFromStrings(vals, storage); err != nil {
			rp.logger.Warn("Invalid query parameter",
				rp.logger.ToString("param", param),
				rp.logger.ToString("value", strings.Join(vals, ",")),
				rp.logger.ToError(err))
			// If can't parse just set zero value
		}
	}

	// Read and parse body parameters if any
	// TODO: we force user to list all body parameters even if we parse full dto
	if len(rp.bodyParams) != 0 {
		if err := json.NewDecoder(rp.request.Body).Decode(rp.dto); err != nil {
			rp.logger.Warn("Failed to decode body parameters",
				rp.logger.ToError(err))
			return err
		}
	}

	// Read and parse cookie parameters if any
	for i := range rp.cookieParams {
		param := rp.cookieParams[i]
		storage := rp.cookieParamsStorage[i]
		cookie, err := rp.request.Cookie(param)
		if err != nil {
			rp.logger.Warn("Failed to get cookie parameter",
				rp.logger.ToString("param", param),
				rp.logger.ToError(err))
			continue
		}
		// Scan value into the given storage
		if _, err := fmt.Sscanf(cookie.Value, "%v", storage); err != nil {
			rp.logger.Warn("Invalid cookie parameter",
				rp.logger.ToString("param", param),
				rp.logger.ToString("value", cookie.Value),
				rp.logger.ToError(err))
			// If can't scan just set zero value
		}
	}

	// Parse access token to given storage
	if _, err := fmt.Sscanf(rp.accessToken, "%s", rp.accessTokenStorage); err != nil {
		rp.logger.Warn("Invalid access token",
			rp.logger.ToString("access_token", rp.accessToken),
			rp.logger.ToError(err))
		// If can't scan just set zero value
	}

	// Parse path parameters
	for i := range rp.pathParams {
		param := rp.pathParams[i]
		storage := rp.pathParamsStorage[i]
		val := chi.URLParam(rp.request, param)

		// Scan value into the given storage
		if _, err := fmt.Sscanf(val, "%v", storage); err != nil {
			rp.logger.Warn("Invalid path parameter",
				rp.logger.ToString("param", param),
				rp.logger.ToString("value", val),
				rp.logger.ToError(err))
			// If can't scan just set zero value
		}
	}

	// Parse file parameters
	for i := range rp.fileParams {
		param := rp.fileParams[i]
		storage := rp.fileParamsStorage[i]

		// open file
		file, _, err := rp.request.FormFile(param)
		if err != nil {
			rp.logger.Warn("Failed to get file parameter",
				rp.logger.ToString("param", param),
				rp.logger.ToError(err))
			continue
		}
		defer func() {
			if cerr := file.Close(); cerr != nil {
				rp.logger.Warn("Failed to close uploaded file", rp.logger.ToError(cerr))
			}
		}()

		const maxFileSize = 10 * 1024 * 1024 // 10 MB limit for avatar files
		fileData := make([]byte, 0)
		buf := make([]byte, 1024)
		totalRead := 0
		for {
			n, err := file.Read(buf)
			if n > 0 {
				totalRead += n
				if totalRead > maxFileSize {
					rp.logger.Warn("File size exceeded limit",
						rp.logger.ToString("param", param),
						rp.logger.ToInt("maxFileSize", maxFileSize))
					return fmt.Errorf("file size exceeded limit: %d bytes", maxFileSize)
				}
				fileData = append(fileData, buf[:n]...)
			}
			if err != nil {
				break
			}
		}

		// Scan value into the given storage
		*storage = fileData
	}

	fileParamsLengths := make([]int, 0, len(rp.fileParamsStorage))
	// TODO: make this check for every type of parameter
	for i := range rp.fileParamsStorage {
		fileParamsLengths = append(fileParamsLengths, len(*rp.fileParamsStorage[i]))
		if len(*rp.fileParamsStorage[i]) == 0 {
			rp.logger.Warn("Empty file parameter",
				rp.logger.ToString("param", rp.fileParams[i]))
		}
	}

	rp.logger.Debug(
		"Parsed request parameters successfully",
		rp.logger.ToAny("queryParams", rp.queryParams),
		rp.logger.ToAny("queryParamsStorage", rp.queryParamsStorage),
		rp.logger.ToAny("bodyParams", rp.bodyParams),
		rp.logger.ToAny("bodyParamsStorage", rp.bodyParamsStorage),
		rp.logger.ToAny("cookieParams", rp.cookieParams),
		rp.logger.ToAny("cookieParamsStorage", rp.cookieParamsStorage),
		rp.logger.ToAny("pathParams", rp.pathParams),
		rp.logger.ToAny("pathParamsStorage", rp.pathParamsStorage),
		rp.logger.ToAny("accessToken", rp.accessToken),
		rp.logger.ToAny("fileParams", rp.fileParams),
		rp.logger.ToAny("fileParamsStorage lengths in bytes", fileParamsLengths),
	)
	return nil
}

// parseAndSetFromStrings parses one or more string values into the provided storage pointer.
// Storage must be a pointer to a supported kind: basic scalar (string,int,uint,bool,float)
// or a slice of those types (e.g. *[]uint, *[]int, *[]string). It supports repeated query
// parameters (e.g. ?id=1&id=2) and comma-separated lists (e.g. ?id=1,2,3).
func parseAndSetFromStrings(vals []string, storage any) error {
	if storage == nil {
		return fmt.Errorf("nil storage")
	}

	rv := reflect.ValueOf(storage)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("storage must be pointer")
	}
	ev := rv.Elem()

	// Helper to build flat list of items from repeated params and comma-separated values
	items := make([]string, 0, len(vals))
	for _, v := range vals {
		if v == "" {
			continue
		}
		// split comma-separated values
		if strings.Contains(v, ",") {
			for _, part := range strings.Split(v, ",") {
				if s := strings.TrimSpace(part); s != "" {
					items = append(items, s)
				}
			}
		} else {
			items = append(items, strings.TrimSpace(v))
		}
	}

	// If target is a slice, fill it with parsed elements
	if ev.Kind() == reflect.Slice {
		elemType := ev.Type().Elem()
		newSlice := reflect.MakeSlice(ev.Type(), 0, len(items))
		for _, it := range items {
			parsed, err := parseSingleStringToType(it, elemType)
			if err != nil {
				return err
			}
			newSlice = reflect.Append(newSlice, parsed)
		}
		ev.Set(newSlice)
		return nil
	}

	// Not a slice: take first item only
	var s string
	if len(items) > 0 {
		s = items[0]
	} else {
		s = ""
	}
	parsed, err := parseSingleStringToType(s, ev.Type())
	if err != nil {
		return err
	}
	ev.Set(parsed)
	return nil
}

// parseSingleStringToType parses a single string into a reflect.Value of the provided type.
func parseSingleStringToType(s string, t reflect.Type) (reflect.Value, error) {
	if t.Kind() == reflect.String {
		return reflect.ValueOf(s).Convert(t), nil
	}
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if s == "" {
			return reflect.Zero(t), nil
		}
		v, err := strconv.ParseInt(s, 10, t.Bits())
		if err != nil {
			return reflect.Value{}, err
		}
		rv := reflect.New(t).Elem()
		rv.SetInt(v)
		return rv, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if s == "" {
			return reflect.Zero(t), nil
		}
		v, err := strconv.ParseUint(s, 10, t.Bits())
		if err != nil {
			return reflect.Value{}, err
		}
		rv := reflect.New(t).Elem()
		rv.SetUint(v)
		return rv, nil
	case reflect.Bool:
		if s == "" {
			return reflect.Zero(t), nil
		}
		v, err := strconv.ParseBool(s)
		if err != nil {
			return reflect.Value{}, err
		}
		rv := reflect.New(t).Elem()
		rv.SetBool(v)
		return rv, nil
	case reflect.Float32, reflect.Float64:
		if s == "" {
			return reflect.Zero(t), nil
		}
		v, err := strconv.ParseFloat(s, t.Bits())
		if err != nil {
			return reflect.Value{}, err
		}
		rv := reflect.New(t).Elem()
		rv.SetFloat(v)
		return rv, nil
	default:
		return reflect.Value{}, fmt.Errorf("unsupported type: %s", t.Kind())
	}
}
