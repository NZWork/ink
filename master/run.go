package master

import "log"

func Run() {
	log.Println("Running as master")
	newTask()
	Close()
}

func Close() {
	log.Println("master closed")
}
