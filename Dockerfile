FROM golang:1.9.1
ENV SERVICE_IP=localhost
WORKDIR /go/src/github.com/morrah77/messagesender
RUN go get -u github.com/golang/dep/cmd/dep
COPY . .
RUN dep ensure
RUN go build -o build/messagesender main.go

#TODO(h.lazar) remove unused files
#TODO(h.lazar) consider to check wheather commservice is up
CMD `sleep 2s && echo 'Run messagesender...' && ./build/messagesender --url=http://$SERVICE_IP:9090/messages`

