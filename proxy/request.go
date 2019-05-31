package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/apex/log"
	"net/http"
	"strings"
)

var client *HttpClient = nil

type Request struct {
	Endpoint      string   `json:"endpoint"`
	Values        []string `json:"values"`
	RequestType   string   `json:"request_type"`
	Headers       []Header `json:"headers"`
	Content       string   `json:"content"`
	CorrelationId string   `json:"-"` // Used when sending the response back over RabbitMQ
}

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func NewRequest(raw []byte) (Request, error) {
	var req Request
	err := json.Unmarshal(raw, &req)
	return req, err
}

func (r *Request) GetGenericEndpoint() string {
	endpoint := ""

	// TODO: Make this a lot cleaner
	if strings.HasPrefix(r.Endpoint, "/channels/{}") {
		endpoint = fmt.Sprintf("/channels/%s", r.Values[0])
	} else if strings.HasPrefix(r.Endpoint, "/guilds/{}") {
		endpoint = fmt.Sprintf("/guilds/%s", r.Values[0])
	} else if strings.HasPrefix(r.Endpoint, "/invites/{}") {
		endpoint = fmt.Sprintf("/invites/%s", r.Values[0])
	} else if strings.HasPrefix(r.Endpoint, "/users/{}") {
		endpoint = fmt.Sprintf("/users/%s", r.Values[0])
	} else if strings.HasPrefix(r.Endpoint, "/guilds/{}") {
		endpoint = fmt.Sprintf("/guilds/%s", r.Values[0])
	} else if strings.HasPrefix(r.Endpoint, "/webhooks/{}") {
		endpoint = fmt.Sprintf("/webhooks/%s", r.Values[0])
	}

	return endpoint
}

func (r *Request) Handle(w Worker) *http.Response {
	if client == nil {
		c := NewClient()
		client = &c
	}

	url := baseUrl + r.Endpoint

	for _, field := range r.Values {
		url = strings.Replace(url, "{}", field, 1)
	}

	var req *http.Request
	var err error
	if r.Content == "" {
		req, err = http.NewRequest(r.RequestType, url, nil); if err != nil {
			log.Error(err.Error())
			return nil
		}
	} else {
		req, err = http.NewRequest(r.RequestType, url, bytes.NewBuffer([]byte(r.Content))); if err != nil {
			log.Error(err.Error())
			return nil
		}
	}

	for _, header := range r.Headers {
		req.Header.Add(header.Key, header.Value)
	}

	ctx := context.TODO()
	conn, err := client.Get(ctx); if err != nil {
		log.Error(err.Error())
		return nil
	}
	defer client.Put(conn)

	res, err := conn.(HttpConnection).Do(req); if err != nil {
		log.Error(err.Error())
		return nil
	}

	return res
}
