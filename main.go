package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"
)

func main() {
	f, _ := os.Create("profile_file")
	pprof.StartCPUProfile(f)     // 开始cpu profile，结果写到文件f中
	defer pprof.StopCPUProfile() // 结束profile

	start := time.Now().UnixNano()
	mdStream()
	fmt.Printf("time %vms\n", float64(time.Now().UnixNano()-start)/1000/1000)
}
