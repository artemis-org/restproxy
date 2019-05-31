package main

import (
	"github.com/apex/log"
	"github.com/artemis/restproxy/amqp"
	"github.com/artemis/restproxy/config"
	"github.com/artemis/restproxy/proxy"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	provider := config.GetConfigProvider()
	provider.LoadConfig()

	worker := proxy.NewWorker()
	go worker.Start()

	amqpClient := amqp.NewAmqpClient(&worker)
	defer amqpClient.Close()

	for i := 0; i < config.Conf.PoolSize; i++ {
		go func() {
			ch, q, err := amqpClient.CreateQueue(); if err != nil { // TODO: Proper error handling with ticker for retries
				log.Error(err.Error())
				return
			}

			for {
				amqpClient.Handle(ch, q)
			}
		}()
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<- sigs
}