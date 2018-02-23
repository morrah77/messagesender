package transport

import (
	"testing"
	"log"
	"bytes"
	"net/http"
	"bufio"
)

var conf *Conf  = &Conf{
	SendUrl: `localhost:12345/test`,
}

var buf = bytes.NewBuffer(make([]byte, 1024))
var logger = log.New(buf, `test-transport`, log.Flags())

type message struct {
	email string `json:"email"`
	text string `json:"text"`
}

func setupHttpResponder(t *testing.T) {
	http.HandleFunc(`/test`, func(writer http.ResponseWriter, request *http.Request) {
		rd := make([]byte, 1024)
		n, err := request.Body.Read(rd)
		if err != nil {
			t.Error(err.Error())
		}
		if n != 100 {
			t.Error(`Incorrect message received!`)
		}
		writer.Write(rd)
	})
	http.ListenAndServe(`localhost:12345`, nil)
}

func TestNewTransport(t *testing.T) {
	var transport = NewTransport(logger, conf)
	if transport.conf.SendUrl != `localhost:12345/test` {
		t.Error(`Incorrect conf.SendUrl on creation!`)
	}
}

func TestTransport_Send(t *testing.T) {
	var transport = NewTransport(logger, conf)
	setupHttpResponder(t)
	transport.Send(&message{
		email:`foo@bar.baz`,
		text:`Hello!`,
	})
}