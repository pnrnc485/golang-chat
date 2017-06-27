package main

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
