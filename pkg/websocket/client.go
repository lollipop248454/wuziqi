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

func checkFull(pwd int) bool {
	return len(Room[pwd]) == 2
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
		msg = strings.Join(strings.Split(msg, " ")[1:], " ")
		message := Message{SenderId: c.ID, Type: tp, Body: msg}
		pwd, _ := strconv.Atoi(msg)
		// 100 以上的类型用于表示匹配，无需发送给别人
		if tp >= 100 {
			Match(c, message, pwd)
			continue
		}
		c.Pool.Broadcast <- message
		// 不显示图片的base64
		if tp != 11 {
			fmt.Printf("Message Received: %+v\n", message)
		}
	}
}
