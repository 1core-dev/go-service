package web

import (
	"context"
	"encoding/json"
	"net/http"
)

// Respond convert a Go value to JSON and sends it to client.
func Respond(ctx context.Context, w http.ResponseWriter, r *http.Request, data any, statusCode int) error {
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}
