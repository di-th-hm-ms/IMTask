package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebsocketHandler struct {}

func New() *WebsocketHandler {
	return &WebsocketHandler{}
}

func (wsh *WebsocketHandler) Handle(rw http.ResponseWriter, req *http.Request) {
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(req *http.Request) bool {
			return true
		},
	}
	if _, err := upgrader.Upgrade(rw, req, nil); err != nil {
		log.Fatal(err)
	}
}
