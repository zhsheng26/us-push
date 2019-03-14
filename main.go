package main

import (
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
	queue := connectMq.BindQueue("us-push", "#.us.#")
	connectMq.Consume(queue.Name, func(body string) {
		//需要推送的消息
		fmt.Println(body)
	})
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := socket.Upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

		for {
			// Read message from browser
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			// Print the message to the console
			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			// Write message back to browser
			if err = conn.WriteMessage(msgType, msg); err != nil {
				return
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/websockets.gohtml")
	})

	_ = http.ListenAndServe(":8080", nil)

}
