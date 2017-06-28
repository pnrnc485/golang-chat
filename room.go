package main

import (
	_ "syscall"
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	_ "net/http/cookiejar"
)

type room struct {

	//他のクライアントに転送するためのメッセージを保持するチャネル
	forward chan []byte

	//チャットルームに参加しようとしているクライアントのためのチャネル
	join chan *client

	//チャットルームから退室しようとしているクライアントのためのチャネル
	leave chan *client

	//在室している全てのクライアントが保持される
	clients map[*client]bool
}

//ヘルパー関数
func newRoom() *room {

	return  &room{
		forward: make(chan []byte),
		join: make(chan *client),
		leave: make(chan  *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for  {
		select {
		case client := <-r.join:
			//参加
			r.clients[client] = true
		case client := <-r.leave:
			//退室
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:

			// 全てのクライアントにメッセージを転送
			for client := range r.clients {
				select {
				case client.send <- msg:
					//メッセージを送信
				default:
					//送信に失敗
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}

const (
	socketBufferSize 	= 1024
	messageBufferSize 	= 256
)

var upgrade = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	socket, err := upgrade.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	client := &client{socket: socket, send: make(chan []byte, messageBufferSize), room: r, }

	r.join <- client
	defer  func () { r.leave <- client }()
	go client.write()
	client.read()
}