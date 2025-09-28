package models

import (
	"encoding/json"
	"log"
)

// Converts error to JSON { "error": "error message" }
func bytesFromError(err error) []byte {
	log.Println("Error marshaling to JSON:", err)
	return []byte(`{ "error": "` + err.Error() + `" }`)
}

// Converts any model to JSON bytes
func ConvertModelToBytes(model any) []byte {
	bytes, err := json.Marshal(model)
	if err != nil {
		return bytesFromError(err)
	}
	return bytes
}

func ConvertBytesToModel(data []byte, model any) error {
	return json.Unmarshal(data, model)
}
