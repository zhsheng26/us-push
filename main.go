package main

import (
	"encoding/json"
	"eusunpower.com/us-push/melody"
	"eusunpower.com/us-push/mq"
	"eusunpower.com/us-push/socket"
	"flag"
	"fmt"
	"net/http"
)

var url *string
var host *string
var username string
var password string

func init() {
	url = flag.String("url", "127.0.0.1:5672", "rabbitmq address")
	host = flag.String("host", "us", "rabbitmq vhost")
	flag.StringVar(&username, "u", "zhsheng", "account of rabbitmq admin")
	flag.StringVar(&password, "p", "mingxi", "password of rabbitmq admin")
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
		topic := request.URL.Query().Get("topic")
		_ = m.HandleRequestWithKeys(writer, request, map[string]interface{}{"topic": topic})
	})
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		//客户端发来的消息
		fmt.Println("receive:" + string(msg))
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/test.gohtml")
	})
	_ = http.ListenAndServe(":8080", nil)

}
