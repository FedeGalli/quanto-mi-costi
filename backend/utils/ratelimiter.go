package utils

import (
	"sync"
	"time"
)

type safeMap struct {
	mu       sync.Mutex
	visitors map[string]uint
}

var users = safeMap{visitors: make(map[string]uint)}

const RateLimit uint = 100

func (visitors *safeMap) visitorCall(ip string) bool {
	visitors.mu.Lock()
	defer visitors.mu.Unlock()
	calls, ok := visitors.visitors[ip]
	if ok {
		if calls > RateLimit {
			return false
		}
		visitors.visitors[ip] += 1
	} else {
		visitors.visitors[ip] = 1
	}
	return true
}

func InizializeRateLimiter() {
	ticker := time.Tick(60 * time.Second)

	go func() {
		for {
			select {
			case <-ticker:
				users.mu.Lock()
				users.visitors = make(map[string]uint)
				users.mu.Unlock()
			}
		}
	}()
}

func IsAllowed(ip string) bool {
	return users.visitorCall(ip)
}
