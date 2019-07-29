package socket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type PushMessage struct {
	Content string
	Topic   string
}

type clientHub struct {
	clients              map[string]*websocket.Conn
	registerClientChan   chan *websocket.Conn
	unRegisterClientChan chan *websocket.Conn
	BroadcastChan        chan PushMessage
}

func (hub *clientHub) register(conn *websocket.Conn) {
	hub.clients[conn.RemoteAddr().String()] = conn
}

func (hub *clientHub) unRegister(conn *websocket.Conn) {
	delete(hub.clients, conn.RemoteAddr().String())
}

func (hub *clientHub) broadcast(message PushMessage) {
	bytes, _ := json.Marshal(message)
	for _, conn := range hub.clients {
		if err := conn.WriteJSON(string(bytes)); err != nil {
			return
		}
	}
}

var Hub *clientHub

var upgrader websocket.Upgrader

func init() {
	Hub = createHub()
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
}

func Start() {
	http.HandleFunc("/us-push", func(writer http.ResponseWriter, request *http.Request) {
		conn, _ := upgrader.Upgrade(writer, request, nil)
		handler(conn, Hub)
	})
	go func() {
		for {
			select {
			case conn := <-Hub.registerClientChan:
				Hub.register(conn)
			case conn := <-Hub.unRegisterClientChan:
				Hub.unRegister(conn)
			case message := <-Hub.BroadcastChan:
				Hub.broadcast(message)
			}

		}
	}()
}

func handler(conn *websocket.Conn, hub *clientHub) {
	hub.registerClientChan <- conn
	//怎么知道conn连接断开了
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("%v", err)
			hub.unRegister(conn)
		}
	}
}

func createHub() *clientHub {
	return &clientHub{
		clients:              make(map[string]*websocket.Conn),
		registerClientChan:   make(chan *websocket.Conn),
		unRegisterClientChan: make(chan *websocket.Conn),
		BroadcastChan:        make(chan PushMessage),
	}
}
