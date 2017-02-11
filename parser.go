package main

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"sync"
	//"time"
)

var rootPath = "test"
var wg sync.WaitGroup
var policy *bluemonday.Policy

func mdStream() {
	files, _ := filepath.Glob(filepath.Join(rootPath, "*.md"))
	fmt.Printf("parsing %d files\n", len(files))

	policy = bluemonday.UGCPolicy()

	worker := runtime.NumCPU()
	runtime.GOMAXPROCS(worker)

	rc := make(chan *[]byte, worker)

	for _, f := range files {
		wg.Add(1)
		go func(f *string) {
			md, _ := ioutil.ReadFile(*f)
			rc <- &md
		}(&f)
	}

	for i := 0; i < worker*2; i++ {
		go mdParseStream(rc)
	}
	wg.Wait()
}

func mdParseStream(readChan chan *[]byte) {
	var c *[]uint8
	for {
		select {
		case c = <-readChan:
			ioutil.WriteFile("test.html", policy.SanitizeBytes(blackfriday.MarkdownCommon(*c)), 0644)
			wg.Done()
		}
	}
}
