package network

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	. "websocket-go/types"
)

var Upgrader = &websocket.Upgrader{ReadBufferSize: SocketBufferSize, WriteBufferSize: MessageBufferSize}

type Message struct {
	Name    string
	Message string
	Time    int64
}

type Client struct {
	Send   chan *Message
	Room   *Room
	Name   string
	Socket *websocket.Conn
}

type Room struct {
	Forward chan *Message

	Join  chan *Client
	Leave chan *Client

	Clients map[*Client]bool
}

func NewRoom() *Room {
	return &Room{
		Forward: make(chan *Message),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
		Clients: make(map[*Client]bool),
	}
}

func (r *Room) SocketServce(c *gin.Context) {
	Upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
}
