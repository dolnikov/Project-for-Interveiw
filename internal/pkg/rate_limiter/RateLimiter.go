package rate_limiter

import (
	"sync"
	"time"
)

type RateLimiter struct {
	requests map[string]map[string][]time.Time // Method Name; Key for defining restrictions (like user_id or IP)
	limits   map[string]int                    // Method name and limits (Request per minute)
	mu       sync.Mutex
}

func NewRateLimiter(limits map[string]int) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string]map[string][]time.Time),
		limits:   limits,
	}
}

func (rl *RateLimiter) Allow(methodName string, key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limit, ok := rl.limits[methodName]
	if !ok {
		return true
	}

	if _, ok = rl.requests[methodName]; !ok {
		rl.requests[methodName] = make(map[string][]time.Time)
	}

	currentRequests := rl.requests[methodName][key]
	currentTime := time.Now()

	var newRequests []time.Time
	for _, reqTime := range currentRequests {
		if reqTime.Add(time.Minute).After(currentTime) {
			newRequests = append(newRequests, reqTime)
		}
	}

	if len(newRequests) >= limit {
		return false
	}

	newRequests = append(newRequests, currentTime)
	rl.requests[methodName][key] = newRequests

	return true
}
