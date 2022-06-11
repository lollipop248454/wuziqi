package websocket

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
)
import "github.com/google/uuid"

// Room 数组存房间，房间k-v => id-(0,1=>没准备，准备)
var Room []map[*Client]int
var ClientRoom map[string]int
var IdClient map[string]*Client

// InitRoom 根据你的密码选择对应的房间
func InitRoom() {
	Room = make([]map[*Client]int, 1000)
	ClientRoom = make(map[string]int)
	IdClient = make(map[string]*Client)
}

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func IsAllReady(pwd int) bool {
	c := 0
	for _, v := range Room[pwd] {
		c += v
	}
	return c == 2
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			id := uuid.New().String()
			client.ID = id
			pool.Clients[client] = true
			log.Println("Size of Connection Pool: ", len(pool.Clients))
			fmt.Println("client id: ", client.ID)
			// 只需要告诉自己
			client.Conn.WriteJSON(Message{SenderId: id, Type: 1, Body: "New User Joined..."})
			//for client, _ := range pool.Clients {
			//	fmt.Println(client)
			//	client.Conn.WriteJSON(Message{SenderId: id, Type: 1, Body: "New User Joined..."})
			//}
			break
		case client := <-pool.Unregister:
			id := client.ID
			delete(pool.Clients, client)
			log.Println("Size of Connection Pool: ", len(pool.Clients))
			// 2种情况 1.匹配断线 2.准备或对战断线
			pwd := ClientRoom[client.ID]
			if len(Room[pwd]) == 1 {
				deleteClient(client, pwd)
			}
			if len(Room[pwd]) == 2 {
				deleteClient(client, pwd)
				for client, _ := range Room[pwd] {
					if client.ID == id {
						continue
					}
					client.Conn.WriteJSON(Message{SenderId: client.ID, Type: 105, Body: "tell this user_id that opponent Disconnected..."})
				}
			}
			//for client, _ := range pool.Clients {
			//	client.Conn.WriteJSON(Message{SenderId: id, Type: 0, Body: "User Disconnected..."})
			//}
			break
		case message := <-pool.Broadcast:
			pwd := ClientRoom[message.SenderId]
			c := IdClient[message.SenderId]
			if message.Type == -1 {
				Room[pwd][c] = 1
			}
			if message.Type == -2 {
				Room[pwd][c] = 0
			}
			for client, _ := range Room[pwd] {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
			if IsAllReady(pwd) {
				message.Type = 3
				rand.Seed(time.Now().UnixNano())
				temp := rand.Intn(2)
				for client, _ := range Room[pwd] {
					message.Body = strconv.Itoa(temp)
					Room[pwd][client] = 0
					temp++
					if err := client.Conn.WriteJSON(message); err != nil {
						fmt.Println(err)
						return
					}
				}
			}
		}
	}
}
