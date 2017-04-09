package master

import (
	"ink/public"
	"io/ioutil"
	"os"

	"github.com/streadway/amqp"
)

func taskGenrator(repo string) ([]public.Task, int) {
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
			// new one
			task = public.Task{
				Repo:  repo,
				Files: []string{},
			}
		}
		task.Files = append(task.Files, f.Name())
	}
	tasks = append(tasks, task)

	return tasks, len(files)
}

func newTask(repo string) (int, int, float64) {
	var (
		err  error
		body []byte
		task public.Task
	)
	startTime := public.TimerStart()
	tasks, fileCount := taskGenrator(repo)

	if fileCount == 0 {
		return fileCount, len(tasks), public.TimerStop(startTime)
	}

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
	}

	return fileCount, len(tasks), public.TimerStop(startTime)
}
