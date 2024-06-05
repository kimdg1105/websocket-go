package network

import (
	"log"
	"net/http"
	"time"
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

func (c *Client) Read() {
	defer c.Socket.Close()
	for {
		var msg Message
		err := c.Socket.ReadJSON(&msg)
		if err != nil {
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				break
			} else {
				panic(err)
			}
		} else {
			log.Println("Read: ", msg, "CLIENT: ", c.Name)
			log.Println()
			msg.Time = time.Now().Unix()
			msg.Name = c.Name
			c.Room.Forward <- &msg
		
		}
	}
}

func (c *Client) Write() {
	defer c.Socket.Close()
	for msg := range c.Send {
		err := c.Socket.WriteJSON(msg)
		if err != nil {
			panic(err)
		}
		log.Println("Write: ", msg, "CLIENT: ", c.Name)
		log.Println()
	}
}

func (r *Room) Run() {
	// Room 내 모든 채널값을 받는 역할
	for {
		select {
		case client := <-r.Join:
			r.Clients[client] = true
		case client := <-r.Leave:
			r.Clients[client] = false
			close(client.Send)
			delete(r.Clients, client)
		case msg := <-r.Forward:
			for client := range r.Clients {
				client.Send <- msg
			}
		}
	}
}

func (r *Room) SocketServe(c *gin.Context) {
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

	go client.Write()
	client.Read()


	defer func() { r.Leave <- client }()
}
