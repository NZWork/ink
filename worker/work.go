package worker

import (
	"encoding/json"
	"ink/public"
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func work() {
	var err error
	err = public.MQChannel.Qos(
		TaskFetch, // prefetch count
		0,         // prefetch size
		false,     // global
	)
	failOnError(err, "failed to set QoS")

	msgs, err := public.MQChannel.Consume(
		public.QueueName, // queue
		"",               // consumer
		false,            // auto-ack
		false,            // exclusive
		false,            // no-local
		false,            // no-wait
		nil,              // args
	)
	failOnError(err, "failed to register a consumer")

	forever := make(chan bool)

	go func() {
		var (
			d         amqp.Delivery
			task      *public.Task
			startTime int64
		)

		for d = range msgs {
			log.Printf("received task: %s", d.Body)
			startTime = public.TimerStart()
			d.Ack(false)
			task = &public.Task{}
			json.Unmarshal(d.Body, task)
			mdStream(task)
			log.Printf("task done: %s, cost %f ms", d.Body, public.TimerStop(startTime))
		}
	}()

	log.Printf("worker has ready to work")
	<-forever
}
