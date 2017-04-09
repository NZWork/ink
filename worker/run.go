package worker

import (
	"ink/public"
	"log"
)

func Run() {
	log.Println("Running as worker")
	public.MQConnect()
	work()
}

func Close() {
	public.MQClose()
	log.Println("worker closed")
}
