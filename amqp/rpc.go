package amqp

import (
	"context"
	"github.com/apex/log"
	"github.com/artemis/restproxy/proxy"
	"github.com/streadway/amqp"
)

func (c *AmqpClient) CreateChannel() (*amqp.Channel, error) {
	ctx := context.TODO()
	r, err := c.Get(ctx); if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer c.Put(r)

	conn := r.(AmqpConnection)
	ch, err := conn.Channel(); if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return ch, nil
}

func (c *AmqpClient) CreateQueue(ch *amqp.Channel, name string) (amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Error(err.Error())
		return q, err
	}

	if err = ch.Qos(0, 0, false); err != nil {
		log.Error(err.Error())
		return q, err
	}

	return q, nil
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

		req.CorrelationId = payload.CorrelationId

		c.Worker.Receiver <- req
	}
}

func (c *AmqpClient) StartPublisher(channel *amqp.Channel, queue amqp.Queue) {
	for {
		res := <- c.Worker.Responder

		if err := channel.Publish(
			"",
			"rpc_inbound",
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				CorrelationId: res.CorrelationId,
				ReplyTo: queue.Name,
				Body: res.Response,
			},
		); err != nil {
			log.Error(err.Error())
		}
	}
}
