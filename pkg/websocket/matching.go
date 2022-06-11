package websocket

import "log"

func deleteClient(c *Client, pwd int) {
	delete(ClientRoom, c.ID)
	delete(Room[pwd], c)
	delete(IdClient, c.ID)
}

func appendClient(c *Client, pwd int) {
	Room[pwd][c] = 0
	ClientRoom[c.ID] = pwd
	IdClient[c.ID] = c
}

func Match(c *Client, message Message, pwd int) {
	tp := message.Type
	if tp == 100 {
		message.Type = 102
		message.Body = "匹配中"
		if checkFull(pwd) {
			message.Type = 101
			message.Body = "匹配失败，此房间正在对战中，请尝试更换密码"
		} else {
			if Room[pwd] == nil {
				Room[pwd] = make(map[*Client]int)
			}
			appendClient(c, pwd)
			if len(Room[pwd]) == 2 {
				message.Type = 103
				message.Body = "匹配成功"
				for client, _ := range Room[pwd] {
					message.SenderId = client.ID
					client.Conn.WriteJSON(message)
					log.Printf("%s\n", client.ID)
				}
				return
			}
		}
		log.Printf("%+v", message)
		c.Conn.WriteJSON(message)
	}
	if tp == 104 {
		log.Printf("%+v 此房间一名用户离开", pwd)
		deleteClient(c, pwd)
	}
}
