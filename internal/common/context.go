package common

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

type ContextKey string

// For context that is passed from Controller to Repository layer
const (
	ContextKeyRequestID ContextKey = "requestID"
)

// GetContext extracts and returns a context with request ID from the HTTP request
func GetContext(r *http.Request) context.Context {
	// request_id is injected by chi, but we don't want to depend on chi internal naming
	// so we create our own context with common.RequestIDContextKey
	ctx := context.WithValue(
		r.Context(),
		ContextKeyRequestID,
		middleware.GetReqID(r.Context()),
	)
	// We can add more fields to context here if needed

	return ctx
}

func GetRequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return "Context is nil"
	}
	if reqID, ok := ctx.Value(ContextKeyRequestID).(string); ok {
		return reqID
	}
	return "Context key not found"
}
