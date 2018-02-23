package transport

import (
	"net/http"
	"log"
	"time"
	"encoding/json"
	"bytes"
)

type Conf struct{
	SendUrl string
}

type Transport struct {
	client *http.Client
	conf *Conf
	logger *log.Logger
}

func NewTransport(logger *log.Logger, conf *Conf) (*Transport) {
	transport := &http.Transport{
		IdleConnTimeout:30 * time.Second,
	}
	return &Transport{
		client: &http.Client{Transport: transport},
		conf:conf,
		logger:logger,
	}
}

func(tr *Transport) Send(message interface{}) (err error) {
	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(payload)
	req, err := http.NewRequest(http.MethodPost, tr.conf.SendUrl, buf)
	if err != nil {
		return err
	}
	resp, err := tr.client.Do(req)
	buf = bytes.NewBuffer(make([]byte, resp.ContentLength))
	tr.logger.Println(resp.StatusCode, buf.String())
	return err
}