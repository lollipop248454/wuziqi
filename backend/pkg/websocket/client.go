package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"strings"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

type Message struct {
	SenderId string `json:"sender_id"`
	Type     int    `json:"type"`
	Body     string `json:"body"`
}

func (c *Client) Read() {
	defer func() {
		// 两个*才能更改c的pool的具体内容
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		msg := string(p)
		tp, _ := strconv.Atoi(strings.Split(msg, " ")[0])
		msg = strings.Join(strings.Split(msg," ")[1:]," ")
		message := Message{SenderId: c.ID, Type: tp, Body: msg}
		c.Pool.Broadcast <- message
		fmt.Printf("Message Received: %+v\n", message)
	}
}
