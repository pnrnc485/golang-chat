package main

type room struct {

	//他のクライアントに転送するためのメッセージを保持するチャネル
	forward chan []byte
}
