package worker

import (
	"encoding/json"
	"ink/public"
	"log"
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
	failOnError(err, "Failed to set QoS")

	msgs, err := public.MQChannel.Consume(
		public.MQQueue.Name, // queue
		"",                  // consumer
		false,               // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			d.Ack(false)
			task := &public.Task{}
			json.Unmarshal(d.Body, task)
			mdStream(task)
			log.Println("Done", task.Debug())
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
