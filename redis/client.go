package redis

import (
	"github.com/artemis/restproxy/config"
	"github.com/go-redis/redis"
)

type RedisClient struct {
	*redis.Client
}

func NewRedisClient() *RedisClient {
	return &RedisClient{}
}

func (c *RedisClient) Connect(uri RedisURI) {
	c.Client = redis.NewClient(&redis.Options{
		Addr: uri.Addr,
		Password: uri.Password,
		DB: 0,
		PoolSize: config.Conf.RedisPoolSize,
	})
}

