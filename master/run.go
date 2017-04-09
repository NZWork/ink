package master

import (
	"ink/public"
	"log"
)

func Run() {
	log.Println("running as master")
	public.MQConnect()
	newTask("test")
	Close()
}

func Close() {
	public.MQClose()
	log.Println("master closed")
}
