package main

import (
	"github.com/apex/log"
	"github.com/artemis-org/restproxy/amqp"
	"github.com/artemis-org/restproxy/config"
	"github.com/artemis-org/restproxy/proxy"
	"github.com/artemis-org/restproxy/redis"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	provider := config.GetConfigProvider()
	provider.LoadConfig()

	redisClient := redis.NewRedisClient()
	redisClient.Connect(redis.CreateRedisURI(config.Conf.RedisUri))

	worker := proxy.NewWorker(redisClient)
	go worker.Start()

	amqpClient := amqp.NewAmqpClient(&worker)
	defer amqpClient.Close()

	ch, err := amqpClient.CreateChannel(); if err != nil { // TODO: Proper error handling with ticker for retries
		log.Error(err.Error())
		return
	}

	// Outbound as outbound from restproxy to Discord
	consumer, err := amqpClient.CreateQueue(ch, "rpc_outbound"); if err != nil { // TODO: Proper error handling with ticker for retries
		log.Error(err.Error())
		return
	}

	// Inbound as inbound from Discord to restproxy
	publisher, err := amqpClient.CreateQueue(ch, "rpc_inbound"); if err != nil { // TODO: Proper error handling with ticker for retries
		log.Error(err.Error())
		return
	}

	for i := 0; i < config.Conf.ConsumerPoolSize; i++ {
		go func() {
			for {
				amqpClient.Handle(ch, consumer)
			}
		}()
	}

	for i := 0; i < config.Conf.PublisherPoolSize; i++ {
		go func() {
			for {
				amqpClient.StartPublisher(ch, publisher)
			}
		}()
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<- sigs
}