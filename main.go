package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now().UnixNano()
	mdStream("test")
	fmt.Printf("time %vms\n", float64(time.Now().UnixNano()-start)/1000/1000)
}
