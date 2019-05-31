package redis

import (
	"fmt"
	"strconv"
	"time"
)

// We want to keep the limits for a long time
func (c *RedisClient) SetLimit(endpoint string, limit int) {
	c.Set(fmt.Sprintf("ratelimit:limit:%s", endpoint), limit, 10 * time.Minute)
}

func (c *RedisClient) GetLimit(endpoint string) int {
	str, err := c.Get(fmt.Sprintf("ratelimit:limit:%s", endpoint)).Result(); if err != nil {
		return -1
	}

	limit, err := strconv.Atoi(str); if err != nil {
		return -1
	}

	return limit
}

func (c *RedisClient) SetRemaining(endpoint string, remaining int) {
	c.Set(fmt.Sprintf("ratelimit:remaining:%s", endpoint), remaining, 10 * time.Minute)
}

func (c *RedisClient) GetRemaining(endpoint string) int {
	str, err := c.Get(fmt.Sprintf("ratelimit:remaining:%s", endpoint)).Result(); if err != nil {
		return -1
	}

	limit, err := strconv.Atoi(str); if err != nil {
		return -1
	}

	return limit
}

// No point keeping the reset after it has reset
func (c *RedisClient) SetReset(endpoint string, reset int64) {
	expiry := reset - time.Now().Unix()
	c.Set(fmt.Sprintf("ratelimit:reset:%s", endpoint), reset, time.Duration(expiry) * time.Second)
}

func (c *RedisClient) GetReset(endpoint string) int64 {
	str, err := c.Get(fmt.Sprintf("ratelimit:reset:%s", endpoint)).Result(); if err != nil {
		return -1
	}

	limit, err := strconv.ParseInt(str, 10, 64); if err != nil {
		return -1
	}

	return limit
}
