package httpx

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	if status == http.StatusNoContent {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}
