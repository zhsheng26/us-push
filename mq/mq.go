package mq

import (
	"eusunpower.com/us-push/util"
	"github.com/streadway/amqp"
)

type Setting struct {
	Url      string `json:"amqp://zhsheng:mingxi@127.0.0.1:5672/"`
	Host     string `json:"us"`
	Exchange string `json:"amq.us-push"`
}
type MQ struct {
	Con  *amqp.Connection
	Ch   *amqp.Channel
	Info Setting
}
type HandleReceivedMsg func(body string)

var mq = &MQ{}

func ConnectMq(setting Setting) *MQ {
	conn, err := amqp.Dial(setting.Url + "/" + setting.Host)
	util.FailOnError(err, "Fail to dial mq")
	channel, err := conn.Channel()
	util.FailOnError(err, "Fail to build channel")
	//声明交换机
	err = channel.ExchangeDeclare(setting.Exchange, "topic",
		true, false, false, false, nil)
	util.FailOnError(err, "Fail to declare exchange")
	mq.Info = setting
	mq.Con = conn
	mq.Ch = channel
	return mq
}

//声明消息队列
//绑定队列到Exchange上,指定routingKey。发送到这个exchange上的消息，根据routingKey路由到这个队列
func (mq *MQ) BindQueue(queueName string, routingKey string) amqp.Queue {
	queue, err := mq.Ch.QueueDeclare(queueName, false, false, false, false, nil)
	util.FailOnError(err, "Fail to declare queue")
	err = mq.Ch.QueueBind(queue.Name, routingKey, mq.Info.Exchange, false, nil)
	util.FailOnError(err, "Fail to bind queue")
	return queue
}

func (mq *MQ) Consume(queueName string, handle HandleReceivedMsg) {
	deliveries, err := mq.Ch.Consume(queueName, "", true, false, false, false, nil)
	util.FailOnError(err, "Fail to start consume")
	go func() {
		for data := range deliveries {
			handle(string(data.Body))
		}
	}()
}

func (mq *MQ) Close() {
	_ = mq.Con.Close()
	_ = mq.Ch.Close()
}
