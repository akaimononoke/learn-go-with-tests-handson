package concurrency

import (
	"net/http"
	"time"
)

func Racer(a, b string) (winner string) {
	durationA, durationB := measureResponseTime(a), measureResponseTime(b)
	if durationA < durationB {
		return a
	}
	return b
}

func measureResponseTime(url string) time.Duration {
	start := time.Now()
	http.Get(url)
	return time.Since(start)
}
