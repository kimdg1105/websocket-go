package network

import (
	"net/http"
	. "websocket-go/types"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var Upgrader = &websocket.Upgrader{
	ReadBufferSize:  SocketBufferSize,
	WriteBufferSize: MessageBufferSize,
	CheckOrigin: func(r *http.Request) bool { return true },
}

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
	socket, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	userCookie, err := c.Request.Cookie("auth")
	if err != nil {
		panic(err)
	}

	client := &Client{
		Send:   make(chan *Message, MessageBufferSize),
		Room:   r,
		Name:   userCookie.Value,
		Socket: socket,	
	}

	r.Join <- client

	defer func() { r.Leave <- client }()
}
