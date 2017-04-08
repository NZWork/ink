package master

import (
	"ink/public"
	"io/ioutil"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func taskGenrator(repo string) []public.Task {
	var tasks = []public.Task{}
	// read file and pack
	files, _ := ioutil.ReadDir(public.GetRepoOriginPath(repo))
	var f os.FileInfo
	task := public.Task{
		Repo:  repo,
		Files: []string{},
	}

	for _, f = range files {
		if len(task.Files) > MaxWorkPerTask {
			// pack
			tasks = append(tasks, task)
			log.Println(task.Debug())
			// new one
			task = public.Task{
				Repo:  repo,
				Files: []string{},
			}
		}
		task.Files = append(task.Files, f.Name())
	}
	tasks = append(tasks, task)

	return tasks
}

func newTask() {
	conn, err := amqp.Dial("amqp://guest:guest@10.1.1.7:30006/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		public.QueueName, // name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	failOnError(err, "Failed to declare a queue")

	tasks := taskGenrator("test")

	for _, task := range tasks {
		body, _ := task.JSON()
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         body,
			})
		failOnError(err, "Failed to publish a message")
		log.Printf(" [x] Sent %s", body)
	}

}
