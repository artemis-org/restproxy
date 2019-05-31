package amqp

import (
	"context"
	"github.com/apex/log"
	"github.com/artemis/restproxy/proxy"
	"github.com/streadway/amqp"
)

func (c *AmqpClient) CreateQueue() (*amqp.Channel, amqp.Queue, error) {
	ctx := context.TODO()
	r, err := c.Get(ctx); if err != nil {
		log.Error(err.Error())
		return nil, amqp.Queue{}, err
	}
	defer c.Put(r)

	conn := r.(AmqpConnection)
	ch, err := conn.Channel(); if err != nil {
		log.Error(err.Error())
		return nil, amqp.Queue{}, err
	}

	q, err := ch.QueueDeclare(
		"rpc_queue",
		true,
		false,
		false,
		false,
		nil,
		)

	if err != nil {
		log.Error(err.Error())
		return ch, q, err
	}

	if err = ch.Qos(0, 0, false); err != nil {
		log.Error(err.Error())
		return ch, q, err
	}

	return ch, q, nil
}

func (c *AmqpClient) Handle(channel *amqp.Channel, queue amqp.Queue) {
	msgs, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	); if err != nil {
		log.Error(err.Error())
		return
	}

	for payload := range msgs {
		data := payload.Body

		req, err := proxy.NewRequest(data); if err != nil {
			log.Warn(err.Error())
			continue
		}

		c.Worker.Receiver <- req
	}
}
