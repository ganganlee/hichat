package common

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

//定义消息队列名称
const queueName = "messageQueue"

var (
	MqQueue map[string]*amqp.Channel
)

//连接rabbitMq服务器
func InitRabbitMq(cfg *Config) {

	MqQueue = make(map[string]*amqp.Channel)
	mqConfig := cfg.RabbitMq

	//循环数组，链接mq
	for _, host := range mqConfig {
		conn, err := amqp.Dial("amqp://guest:guest@" + host)
		if err != nil {
			panic(err)
		}

		ch, err := conn.Channel()
		if err != nil {
			panic(err)
		}

		//申明队列
		_, err = ch.QueueDeclare(
			queueName, //Queue name
			true,      //durable
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			panic(err)
		}

		//将管道保存进入mq队列中
		MqQueue[host] = ch
	}
}

//向rabbitMq服务器发送消息
func Publish(ch *amqp.Channel, body []byte) error {
	return ch.Publish(
		"",        //exchange
		queueName, //routing key(queue name)
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, //Msg set as persistent
			ContentType:  "text/plain",
			Body:         body,
		})
}

func Consume(ch *amqp.Channel) (mg string, err error) {
	msgs, err := ch.Consume(
		queueName,
		"MsgWorkConsumer",
		false, //Auto Ack
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	for msg := range msgs {
		fmt.Println(string(msg.Body))
		msg.Ack(false) //Ack
	}

	return "", nil
}
