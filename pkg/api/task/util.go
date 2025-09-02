package task

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("[WriteJSONStatus] json.Encode error: %v", err)
	}
}

func WriteJSONError(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	err := json.NewEncoder(w).Encode(map[string]string{"error": msg})
	if err != nil {
		log.Printf("[WriteJSONError] marshal/send error: %v (original msg: %s)", err, msg)

	}
}
