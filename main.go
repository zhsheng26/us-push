package main

import (
	"encoding/json"
	"eusunpower.com/us-push/melody"
	"eusunpower.com/us-push/mq"
	"eusunpower.com/us-push/socket"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var url *string
var host *string
var username string
var password string

func init() {
	url = flag.String("url", "192.168.1.2:5672", "rabbitmq address")
	host = flag.String("host", "us-push-test", "rabbitmq vhost")
	flag.StringVar(&username, "u", "us", "account of rabbitmq admin")
	flag.StringVar(&password, "p", "1234rewq!", "password of rabbitmq admin")
}

func main() {
	flag.Parse()
	address := "amqp://" + username + ":" + password + "@" + *url
	setting := mq.Setting{
		Url:      address,
		Host:     *host,
		Exchange: "push",
	}
	connectMq := mq.ConnectMq(setting)
	defer connectMq.Close()
	//socket.Start()
	m := melody.New()
	var mutex sync.Mutex
	pairs := make(map[*melody.Session]*melody.Session)
	group := make(map[string]*melody.Session)
	queue := connectMq.BindQueue("us-push", "#.us.#")
	connectMq.Consume(queue.Name, func(body []byte) {
		//需要推送给客户端的消息
		_ = m.BroadcastFilter(body, func(session *melody.Session) bool {
			content := socket.PushMessage{}
			_ = json.Unmarshal(body, &content)
			return session.Keys["topic"] == content.Topic
		})
	})
	http.HandleFunc("/us-push", func(writer http.ResponseWriter, request *http.Request) {
		userId := request.URL.Query().Get("userId")
		log.Printf("userId = %s", userId)
		_ = m.HandleRequestWithKeys(writer, request, map[string]interface{}{"userId": userId})
	})
	m.HandleConnect(func(s *melody.Session) {
		key, exists := s.Get("userId")
		if !exists {
			return
		}
		userId := key.(string)
		group[userId] = s
		var list []User
		for k := range group {
			list = append(list, User{UserId: k})
		}
		msg := Msg{
			Topic:   "group",
			Content: list,
		}
		data, _ := json.Marshal(msg)
		_ = m.Broadcast(data)
	})
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		mutex.Lock()
		pairs[s] = s
		get, exists := s.Get("userId")
		if exists {
			fmt.Println(get)
		}
		fmt.Println("receive:" + string(msg))
		mutex.Unlock()
	})

	m.HandleDisconnect(func(s *melody.Session) {
		fmt.Println(s.Get("userId"))
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/test.gohtml")
	})
	_ = http.ListenAndServe(":8080", nil)

}

type User struct {
	UserId string `json:"userId"`
}

type Msg struct {
	Topic   string      `json:"topic"`
	Content interface{} `json:"content"`
}

type ReqMsg struct {
	Topic  string `json:"topic"`
	Msg    string `json:"msg"`
	ToUser string `json:"toUser"`
}
