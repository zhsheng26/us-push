package socket

import "github.com/gorilla/websocket"

type HandleReceivedMsg func(body string)

type PushMessage struct {
	Content string
	Code    int
}

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
