package proxy

import (
	"encoding/json"
)

type Request struct {
	Endpoint    string   `json:"endpoint"`
	Values      []string `json:"values"`
	RequestType string   `json:"request_type"`
	Headers     []Header `json:"headers"`
	Content     string   `json:"content"`
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
