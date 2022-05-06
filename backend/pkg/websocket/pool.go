package websocket

import "fmt"
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
	for {
		select {
		case client := <-pool.Register:
			id := uuid.New().String()
			client.ID = id
			pool.Clients[client] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			fmt.Println("client id: ", client.ID)
			for client, _ := range pool.Clients {
				fmt.Println(client)
				client.Conn.WriteJSON(Message{SenderId: id, Type: 1, Body: "New User Joined..."})
				if len(pool.Clients) == 2 {
					client.Conn.WriteJSON(Message{SenderId: id, Type: 3, Body: "这人先手哦！"})
				}
			}
			// 发消息给第二个进来的告知其为先手

			break
		case client := <-pool.Unregister:
			id := client.ID
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				client.Conn.WriteJSON(Message{SenderId: id, Type: 0, Body: "User Disconnected..."})
			}
			break
		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in Pool")
			for client, _ := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
