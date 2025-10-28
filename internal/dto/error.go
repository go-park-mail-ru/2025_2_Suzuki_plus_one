package dto

type Error struct {
	Type    string `json:"type"`
	Message string `json:"message,omitempty"`
	Details string `json:"details,omitempty"` // Internal error details, optional
}

func NewError(errType string, err error, details string) Error {
	return Error{
		Type:    errType,
		Message: err.Error(),
		Details: details,
	}
}

