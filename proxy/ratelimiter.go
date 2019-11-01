package proxy

import (
	"github.com/apex/log"
	"io/ioutil"
	"strconv"
	"sync"
	"time"
)

var(
	GlobalRatelimit = false
	GlobalRatelimitMutex = sync.Mutex{}
)

func (r *Request) Queue(w Worker) {
	endpoint := r.GetGenericEndpoint()

	// Check if we're global ratelimited
	GlobalRatelimitMutex.Lock()
	GlobalRatelimitMutex.Unlock()

	epoch := time.Now().Unix() // Seconds
	if epoch > w.RedisClient.GetReset(endpoint) || w.RedisClient.GetRemaining(endpoint) > 0 { // We can run this straight away
		res := r.Handle(w)
		if res == nil { // Something went wrong, let client time out listening
			return
		}
		defer res.Body.Close()

		if res.StatusCode == 429 { // Probably hit global ratelimit, retry
			if res.Header.Get("X-RateLimit-Global") == "true" {
				if millis, err := strconv.Atoi(res.Header.Get("Retry-After")); err == nil {
					// Wait for global ratelimit to expire
					GlobalRatelimitMutex.Lock()
					time.Sleep(time.Duration(millis) * time.Millisecond)
					GlobalRatelimitMutex.Unlock()
					r.Queue(w) // Retry
				}
				return
			} else {
				if millis, err := strconv.Atoi(res.Header.Get("Retry-After")); err == nil {
					// Wait for global ratelimit to expire
					GlobalRatelimitMutex.Lock()
					time.Sleep(time.Duration(millis) * time.Millisecond)
					GlobalRatelimitMutex.Unlock()
					r.Queue(w) // Retry
				}
				return
			}
		}

		// We don't actually have a use for this
		// limit := res.Header.Get("X-RateLimit-Limit")

		if remaining := res.Header.Get("X-RateLimit-Remaining"); remaining != "" {
			if i, err := strconv.Atoi(remaining); err != nil {
				go w.RedisClient.SetRemaining(endpoint, i)
			}
		}

		if reset := res.Header.Get("X-RateLimit-Reset"); reset != "" {
			if i, err := strconv.ParseInt(reset, 10, 64); err != nil {
				go w.RedisClient.SetReset(endpoint, i)
			}
		}

		body, err := ioutil.ReadAll(res.Body); if err != nil {
			log.Error(err.Error())
			return
		}

		resp := NewResponse(*r, body)
		resp.Publish(w)
	} else {
		reset := w.RedisClient.GetReset(endpoint)
		time.Sleep(time.Duration(reset - time.Now().Unix()) * time.Second)
		r.Queue(w)
	}
}
