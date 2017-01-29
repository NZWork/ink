package main

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	//"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	//"time"
)

var rootPath = "test"

func mdStream() {
	files, _ := filepath.Glob(filepath.Join(rootPath, "*.md"))
	//fmt.Printf("parsing %d files\n", len(files))

	worker := runtime.NumCPU()
	//runtime.GOMAXPROCS(worker)

	worker = 20

	rc := make(chan []byte, worker)
	done := make(chan bool, len(files))

	go mdReadStream(files, rc)

	for i := 0; i < worker; i++ {
		go mdParseStream(rc, done)
	}
	for i := 0; i < len(files); i++ {
		<-done
	}

}

func mdReadStream(files []string, readChan chan []byte) {
	for _, f := range files {
		md, _ := ioutil.ReadFile(f)
		readChan <- md
	}
}

func mdParseStream(readChan chan []byte, done chan bool) {
	for c := range readChan {
		ioutil.WriteFile("test.html", bluemonday.UGCPolicy().SanitizeBytes(blackfriday.MarkdownCommon(c)), 0644)
		done <- true
	}
}
