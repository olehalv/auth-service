package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

var requests []Request

func preventSpam(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		ip := getIp(r)

		var filtered []Request

		for _, req := range requests {
			if req.Time.After(now.Add(-1 * time.Minute)) {
				filtered = append(filtered, req)
			}
		}

		requests = filtered

		requests = append(requests, Request{
			Ip:   ip,
			Time: now,
		})

		var count int

		for _, req := range requests {
			if req.Ip == ip {
				count++
			}
		}

		attempts, err := strconv.Atoi(os.Getenv("MAX_HTTP_REQUESTS_PER_MINUTE"))

		if err != nil {
			returnHttpStatus(w, r, http.StatusInternalServerError, "Internal server error", err)
			return
		}

		if count > attempts {
			returnHttpStatus(w, r, http.StatusTooManyRequests, "Too many requests", errors.New(fmt.Sprintf("%s: Reached maximum HTTP requests per minute", ip)))
			return
		}

		f(w, r)
	}
}
