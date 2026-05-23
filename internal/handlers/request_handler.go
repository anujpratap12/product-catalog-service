package handlers

import (
	"encoding/json"
	"net/http"

	"product-catalog-service/internal/limiter"
	"product-catalog-service/internal/models"
	"product-catalog-service/internal/utils"
)

func RequestHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var request models.RequestBody

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if request.UserID == "" {
		utils.WriteError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	if request.Payload == nil {
		utils.WriteError(w, http.StatusBadRequest, "payload is required")
		return
	}

	allowed, count := limiter.CheckRateLimit(request.UserID)

	if !allowed {

		utils.WriteJSON(w, http.StatusTooManyRequests, map[string]interface{}{
			"error":             "Rate limit exceeded",
			"accepted_requests": count,
		})

		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message":           "Request accepted",
		"accepted_requests": count,
	})
}
