package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now().Nanosecond()
	readMarkdownFiles()
	fmt.Printf("\ntime %vms\n", float64(time.Now().Nanosecond()-start)/1000/1000)
}
