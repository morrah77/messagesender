package transport

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"testing"
)

var conf *Conf = &Conf{
	SendUrl: `http://localhost:12345/test`,
}

var buf = bytes.NewBuffer(make([]byte, 1024))
var logger = log.New(buf, `test-transport`, log.Flags())

type message struct {
	Email string `json:"email"`
	Text  string `json:"text"`
	Paid  bool   `json:"paid, omitempty"`
}

var sentRequest = &message{
	Email: `foo@bar.baz`,
	Text:  `Hello!`,
}

var expectedResponse = &message{
	Email: `foo@bar.baz`,
	Text:  `Hello!`,
	Paid:  true,
}

func setupHttpResponder(t *testing.T) {
	http.HandleFunc(`/test`, func(writer http.ResponseWriter, request *http.Request) {
		resp, _ := json.Marshal(expectedResponse)
		rd := make([]byte, 1024)
		_, err := request.Body.Read(rd)
		if err != nil {
			t.Error(err.Error())
		}
		req := &message{}
		err = json.Unmarshal(rd, req)
		if err != nil {
			t.Error(`Could not form response!`)
		}
		writer.Write(resp)
	})
	http.ListenAndServe(`localhost:12345`, nil)
}

func TestNewTransport(t *testing.T) {
	var transport = NewTransport(logger, conf)
	if transport.conf.SendUrl != `http://localhost:12345/test` {
		t.Error(`Incorrect conf.SendUrl on creation!`)
	}
}

//func TestTransport_Send(t *testing.T) {
//	var transport = NewTransport(logger, conf)
//	go setupHttpResponder(t)
//	req := sentRequest
//	resp := &message{}
//	err := transport.Send(
//		req,
//		resp,
//	)
//	if err != nil {
//		t.Error(`Could not send message!`)
//	}
//	if resp.Email != expectedResponse.Email || resp.Text != expectedResponse.Text || resp.Paid != expectedResponse.Paid {
//		t.Error(`Resonse parse error!`)
//	}
//}
