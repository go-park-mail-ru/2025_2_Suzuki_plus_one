package common

type ContextKey string

// For context that is passed from Controller to Repository layer
const (
	RequestIDContextKey ContextKey = "requestID"
)
