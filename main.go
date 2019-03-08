package main

import (
	"eusunpower.com/us-push/mq"
	"flag"
	"fmt"
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
	forever := make(chan bool)
	connectMq := mq.ConnectMq(setting)
	defer connectMq.Close()
	queue := connectMq.BindQueue("us-push", "#.us.#")
	connectMq.Consume(queue.Name, func(body string) {
		//需要推送的消息
		fmt.Println(body)
	})
	<-forever
}
