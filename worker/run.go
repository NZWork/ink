package worker

import "log"

func Run() {
	log.Println("Running as worker")
	listenMQ()
}

func Close() {
	log.Println("worker closed")
}
