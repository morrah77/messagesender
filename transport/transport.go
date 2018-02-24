package transport

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Conf struct {
	SendUrl string
}

type Transport struct {
	client *http.Client
	conf   *Conf
	logger *log.Logger
}

func NewTransport(logger *log.Logger, conf *Conf) *Transport {
	transport := &http.Transport{
		IdleConnTimeout: 30 * time.Second,
	}
	return &Transport{
		client: &http.Client{Transport: transport},
		conf:   conf,
		logger: logger,
	}
}

func (tr *Transport) Send(payloadToSend interface{}, entityToFillFromResponse interface{}) (err error) {
	payload, err := json.Marshal(payloadToSend)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(payload)
	//TODO(h.lazar) TCP connection might be used here, http package is used just for simplicity
	req, err := http.NewRequest(http.MethodPost, tr.conf.SendUrl, buf)
	if err != nil {
		return err
	}
	resp, err := tr.client.Do(req)
	if err != nil {
		tr.logger.Printf(`HTTP processing error: %#v`, err)
		return err
	}
	b := make([]byte, resp.ContentLength)
	_, err = resp.Body.Read(b)
	err = json.Unmarshal(b, entityToFillFromResponse)
	if err != nil {
		return err
	}
	return err
}
