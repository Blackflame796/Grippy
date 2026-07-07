package handlers

import (
	"encoding/json"
	"net/http"

	"Grippy/pkg/logger"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Log.Errorf("writeJSON encode error: %v", err)
	}
}
