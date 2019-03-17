package test

import (
	"encoding/json"
	"eusunpower.com/us-push/mq"
	"eusunpower.com/us-push/socket"
	"eusunpower.com/us-push/util"
	"github.com/streadway/amqp"
	"testing"
	"time"
)

func TestPublish(t *testing.T) {
	setting := mq.Setting{
		Url:      "amqp://zhsheng:mingxi@127.0.0.1:5672",
		Host:     "us",
		Exchange: "push",
	}
	connectMq := mq.ConnectMq(setting)
	defer connectMq.Close()
	go func() {
		body := &socket.PushMessage{
			Content: "us",
			Topic:   "gis",
		}
		bytes, _ := json.Marshal(body)
		err := connectMq.Ch.Publish(setting.Exchange, "my.us.message", false, false,
			amqp.Publishing{
				ContentType: "json/application",
				Body:        bytes,
			})
		util.FailOnError(err, "Fail to publish message")
	}()
	time.Sleep(2 * time.Second)
}
