package websocket

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
)
import "github.com/google/uuid"

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

func (pool *Pool) Start() {
	count := 0
	for {
		select {
		case client := <-pool.Register:
			id := uuid.New().String()
			client.ID = id
			pool.Clients[client] = true
			log.Println("Size of Connection Pool: ", len(pool.Clients))
			fmt.Println("client id: ", client.ID)
			for client, _ := range pool.Clients {
				fmt.Println(client)
				client.Conn.WriteJSON(Message{SenderId: id, Type: 1, Body: "New User Joined..."})
			}
			break
		case client := <-pool.Unregister:
			id := client.ID
			delete(pool.Clients, client)
			log.Println("Size of Connection Pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				client.Conn.WriteJSON(Message{SenderId: id, Type: 0, Body: "User Disconnected..."})
			}
			break
		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in Pool")
			if message.Type == -2 {
				count--
			}
			if message.Type == -1 {
				count++
			}
			if message.Type == -3 {
				count = 0
			}
			fmt.Printf("count: %d\n", count)
			for client, _ := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
			if count == 2 {
				count = 0
				message.Type = 3
				rand.Seed(time.Now().UnixNano())
				temp := rand.Intn(2)
				for client, _ := range pool.Clients {
					message.Body = strconv.Itoa(temp)
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
