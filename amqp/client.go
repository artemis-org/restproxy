package amqp

import(
	"github.com/apex/log"
	"github.com/artemis-org/restproxy/config"
	"github.com/artemis-org/restproxy/proxy"
	"github.com/streadway/amqp"
	"github.com/vitessio/vitess/go/pools"
	"time"
)

type AmqpClient struct {
	Worker *proxy.Worker
	*pools.ResourcePool
}

type AmqpConnection struct {
	*amqp.Connection
}

func (c AmqpConnection) Close() {
	if err := c.Connection.Close(); err != nil {
		log.Error(err.Error())
	}
}

func NewAmqpClient(w *proxy.Worker) AmqpClient {
	return AmqpClient{
		Worker: w,
		ResourcePool: pools.NewResourcePool(func() (pools.Resource, error) {
			c, err := amqp.Dial(config.Conf.AmqpUri)
			return AmqpConnection{c}, err
		}, config.Conf.ConsumerPoolSize + config.Conf.PublisherPoolSize, config.Conf.ConsumerPoolSize + config.Conf.PublisherPoolSize, time.Duration(config.Conf.IdleTimeout) * time.Second),
	}
}
