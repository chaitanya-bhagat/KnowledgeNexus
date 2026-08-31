package health

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	adapterutils "github.com/chaitanya-bhagat/knowledge-nexus/adapters/utils"
)

type DatabaseCheck interface {
	Ping(ctx context.Context) error
}

type HealthHandler struct {
	db DatabaseCheck
}

func NewHealthHandler(db DatabaseCheck) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

func (hh *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)

	adapterutils.WriteJson(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (hh *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := hh.db.Ping(ctx); err != nil {
		adapterutils.WriteJson(w, http.StatusServiceUnavailable, map[string]string{"status": "not_ready"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}

// func writeJson(w http.ResponseWriter, statusCode int, body any) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(statusCode)
// 	json.NewEncoder(w).Encode(body)
// }
