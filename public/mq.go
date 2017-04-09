package public

import (
	"log"

	"github.com/streadway/amqp"
)

var (
	MQConn    *amqp.Connection
	MQChannel *amqp.Channel
	MQQueue   amqp.Queue
)

func MQConnect() {
	var err error
	log.Println("mq connecting")

	MQConn, err = amqp.Dial("amqp://guest:guest@10.1.1.7:30006/")
	failOnError(err, "failed to connect to RabbitMQ")

	MQChannel, err = MQConn.Channel()
	failOnError(err, "failed to open a channel")

	MQQueue, err = MQChannel.QueueDeclare(
		QueueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "failed to declare a queue")
}

func MQClose() {
	MQConn.Close()
	MQChannel.Close()
	log.Println("mq closed")
}
