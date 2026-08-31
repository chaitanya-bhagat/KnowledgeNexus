package adapterutils

import (
	"encoding/json"
	"net/http"
)

func WriteJson(w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}
