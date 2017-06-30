package main

import (
	"github.com/gorilla/websocket"
	"time"
)

//clientはチャットを行っている一人のユーザを表します
 type client struct {
	 //クライアントのためのWebSocket
	 socket *websocket.Conn

	 // メッセージが送られるチャネル
	 send chan *message

	 // roomはこのクライアントが参加しているチャットルーム
	 room *room

	 //ユーザ情報の保持
	 userData map[string]interface{}
 }

func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
			c.room.forward <- msg
		} else {
			break
		}
	}

	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}

	c.socket.Close()
}