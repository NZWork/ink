package master

import (
	"ink/public"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

func taskGenrator(repo string) []public.Task {
	var tasks = []public.Task{}
	// read file and pack
	files, _ := ioutil.ReadDir(public.GetRepoOriginPath(repo))
	var f os.FileInfo
	task := public.Task{
		Repo:  repo,
		Files: []string{},
	}
	log.Printf("[%s] %d files in total", repo, len(files))

	for _, f = range files {
		if len(task.Files) > MaxWorkPerTask {
			// pack
			tasks = append(tasks, task)
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

func newTask(repo string) {
	var (
		err  error
		body []byte
		task public.Task
	)
	startTime := time.Now().UnixNano()
	tasks := taskGenrator(repo)

	for _, task = range tasks {
		body, _ = task.JSON()
		err = public.MQChannel.Publish(
			"",                  // exchange
			public.MQQueue.Name, // routing key
			false,               // mandatory
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         body,
			})
		failOnError(err, "failed to publish a task")
		log.Printf("[%s] sent %s", repo, body)
	}
	log.Printf("[%s] cost %f ms", repo, float64(time.Now().UnixNano()-startTime)/1000.0/1000.0)

}
