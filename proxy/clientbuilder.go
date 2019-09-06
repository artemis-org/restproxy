package proxy

import (
	"github.com/artemis-org/restproxy/config"
	"github.com/vitessio/vitess/go/pools"
	"net/http"
	"time"
)

type HttpClient struct {
	*pools.ResourcePool
}

type HttpConnection struct {
	*http.Client
}

func (c HttpConnection) Close() {

}

func NewClient() HttpClient {
	return HttpClient{
		pools.NewResourcePool(func() (pools.Resource, error) {
			c := http.Client{
				Timeout: time.Duration(config.Conf.HttpTimeout) * time.Second,
			}
			return HttpConnection{&c}, nil
		}, config.Conf.HttpPoolSize, config.Conf.HttpPoolSize, time.Duration(config.Conf.HttpTimeout) * time.Second),
	}
}
