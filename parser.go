package main

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	//"time"
)

var rootPath = "test"

func mdStream() {
	files, _ := filepath.Glob(filepath.Join(rootPath, "*.md"))
	fmt.Printf("parsing %d files\n", len(files))

	worker := runtime.NumCPU()
	runtime.GOMAXPROCS(worker)

	rc := make(chan []byte, worker)
	wc := make(chan []byte, worker)
	done := make(chan bool, len(files))

	go mdReadStream(files, rc)
	go mdWriteStream(wc, done)

	for i := 0; i < worker*2; i++ {
		go mdParseStream(rc, wc)
	}
	for i := 0; i < len(files); i++ {
		<-done
	}

}

func mdReadStream(files []string, readChan chan []byte) {
	for _, f := range files {
		md, _ := ioutil.ReadFile(f)
		readChan <- md
		//fmt.Printf("rf @ %v\n", time.Now().Nanosecond())
	}
	close(readChan)
}

func mdParseStream(readChan chan []byte, writeChan chan []byte) {
	for c := range readChan {
		//fmt.Printf("cf @ %v\n", time.Now().Nanosecond())
		writeChan <- bluemonday.UGCPolicy().SanitizeBytes(blackfriday.MarkdownCommon(c))
	}
}

func mdWriteStream(writeChan chan []byte, d chan bool) {
	for c := range writeChan {
		//fmt.Printf("wf @ %v\n", time.Now().Nanosecond())
		ioutil.WriteFile("test.html", c, 0644)
		d <- true
	}
}
