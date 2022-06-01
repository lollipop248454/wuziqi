package main

import (
	"backend/pkg/util"
	"fmt"
	"net/http"
	"strings"

	"backend/pkg/websocket"
)

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	ip := r.RemoteAddr
	if strings.Contains(ip, ":") {
		ip = strings.Split(ip, ":")[0]
	}
	util.AskIpAddr(ip)
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
	http.HandleFunc("/http", func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if strings.Contains(ip, ":") {
			ip = strings.Split(ip, ":")[0]
		}
		util.AskIpAddr(ip)
		w.Write([]byte("hello world\n"))
	})
}

func main() {
	fmt.Println("Distributed Chat App v0.01")
	setupRoutes()
	_ = http.ListenAndServe(":8000", nil)
}
