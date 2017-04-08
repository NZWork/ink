package main

import (
	"flag"
	"ink/master"
	"ink/worker"
)

type Flag struct {
	isMaster bool
}

func flags() *Flag {
	f := &Flag{}
	flag.BoolVar(&f.isMaster, "m", false, "Run as master")
	flag.Parse()

	return f
}

func main() {
	f := flags()
	if f.isMaster {
		master.Run()
	} else {
		worker.Run()
	}
}
