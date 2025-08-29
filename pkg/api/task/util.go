package task

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
}
