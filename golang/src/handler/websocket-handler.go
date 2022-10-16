package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebsocketHandler struct {
	pageHub *PageHub
}

type PageHub struct {
	Clients			map[*Client]bool
	RegisteredCh	chan *Client
	UnregisteredCh	chan *Client
	// pageId			string
	TaskCh 			chan *WsTask
}
type Client struct {
	ws		*websocket.Conn
	sendCh	chan *WsTask
}

type WsTask struct {
	// Username 	string `json:"username"`
	// PageId		string `json:"pageId"` // 後で
	TaskId		uint32 `json:"taskId"`
	Content		string `json:"content"`
}

func NewHub() *PageHub {
	return &PageHub{
		Clients:		make(map[*Client]bool),
		RegisteredCh:	make(chan *Client),
		UnregisteredCh: make(chan *Client),
		TaskCh:			make(chan *WsTask),
	}
}
func NewClient(ws *websocket.Conn) *Client {
	return &Client{ws: ws, sendCh: make(chan *WsTask)}
}



// just for cutting off from main
func New(pageHub *PageHub) *WebsocketHandler {
	return &WebsocketHandler{ pageHub: pageHub }
}

func (wh *WebsocketHandler) Handle(rw http.ResponseWriter, req *http.Request) {
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(req *http.Request) bool {
			return true
		},
	}
	ws, err := upgrader.Upgrade(rw, req, nil);
	if err != nil {
		log.Fatal(err)
	}

	client := NewClient(ws)
	go client.readTaskLoop(wh.pageHub.TaskCh, wh.pageHub.UnregisteredCh)
	go client.writeTaskLoop(wh.pageHub)
	wh.pageHub.RegisteredCh <- client
}

// Preparation before clients comes.
func (p *PageHub) Loop() {
	for {
		select {
		case client := <- p.RegisteredCh:
			// register client
			log.Printf("new client: %d nin", len(p.Clients))
			p.Clients[client] = true
		case client := <- p.UnregisteredCh:
			// unregister client
			delete(p.Clients, client)
		case task := <- p.TaskCh:
			// let all clients know
			for c := range p.Clients {
				c.sendCh <- task
			}
		}
	}
}

// func (p *PageHub) register(c *Client) {
// 	p.Clients[c] = true
// }

// func ReadTaskRoop(ws *websocket.Conn) {
func (c *Client) readTaskLoop(taskCh chan<- *WsTask, unregisteredCh chan<- *Client) {
	defer func() {
		c.disconnect(unregisteredCh)
	}()
	for {
		var task WsTask
		if err := c.ws.ReadJSON(&task); err != nil {
			log.Printf("ws read err : %v", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}
		// Hubのチャネルで共有
		taskCh <- &task
	}
}

func (c *Client) writeTaskLoop(p *PageHub) {

	for {
		// infinited loop -> stopped until chanel gets new value(pointer)
		task := <- c.sendCh

		if err := c.ws.WriteJSON(task); err != nil {
			log.Printf("ws write json err: %v", err)
			c.disconnect(p.UnregisteredCh)
		}
	}
}

func (c *Client) disconnect(unregistered chan<- *Client) {
	// Hub unregisterチャネルに入れて LoopよりClientsから削除
	unregistered <- c
	c.ws.Close()
}
