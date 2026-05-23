package limiter

import (
	"sync"
	"time"

	"product-catalog-service/internal/models"
)

var (
	UserStore = make(map[string]*models.UserStats)

	Mu sync.Mutex
)

func CheckRateLimit(userID string) (bool, int) {
	Mu.Lock()
	defer Mu.Unlock()

	user, exists := UserStore[userID]

	if !exists {
		UserStore[userID] = &models.UserStats{
			AcceptedCount: 1,
			RejectedCount: 0,
			WindowStart:   time.Now(),
		}

		return true, 1
	}

	// reset after 1 minute
	if time.Since(user.WindowStart) >= time.Minute {
		user.AcceptedCount = 0
		user.RejectedCount = 0
		user.WindowStart = time.Now()
	}

	if user.AcceptedCount >= 5 {
		user.RejectedCount++

		return false, user.AcceptedCount
	}

	user.AcceptedCount++

	return true, user.AcceptedCount
}
