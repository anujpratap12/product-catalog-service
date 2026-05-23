package handlers

import (
	"encoding/json"
	"net/http"

	"sourceasia-backend-assignment/internal/limiter"
)

func StatsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := make(map[string]interface{})

	limiter.Mu.Lock()

	for userID, stats := range limiter.UserStore {

		response[userID] = map[string]interface{}{
			"accepted_requests": stats.AcceptedCount,
			"rejected_requests": stats.RejectedCount,
			"window_start":      stats.WindowStart,
		}
	}

	limiter.Mu.Unlock()

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}
