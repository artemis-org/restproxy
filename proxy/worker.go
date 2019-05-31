package proxy

import "github.com/artemis/restproxy/redis"

const baseUrl = "https://discordapp.com/api/v6"

type Worker struct {
	Receiver chan Request
	Responder chan Response
	RedisClient *redis.RedisClient
}

func NewWorker(redisClient *redis.RedisClient) Worker {
	return Worker{
		Receiver: make(chan Request),
		Responder: make(chan Response),
		RedisClient: redisClient,
	}
}

func (w *Worker) Start() {
	for {
		req := <- w.Receiver

		go req.Queue(*w)
	}
}
