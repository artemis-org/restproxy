package proxy

import (
	"fmt"
	"strings"
)

const baseUrl = "https://discordapp.com/api/v6"

type Worker struct {
	Receiver chan Request
}

func NewWorker() Worker {
	return Worker{
		Receiver: make(chan Request),
	}
}

func (w *Worker) Start() {
	for {
		req := <- w.Receiver

		url := baseUrl + req.Endpoint

		for _, field := range req.Values {
			url = strings.Replace(url, "{}", field, 1)
		}

		fmt.Println(url)
	}
}
