package main

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	"fmt"
	"io/ioutil"
	//	"path/filepath"
	//	"runtime"
	"sync"
	//"time"
)

var wg sync.WaitGroup
var policy *bluemonday.Policy

func mdStream() {
	files, _ := ioutil.ReadDir("test")
	fmt.Printf("parsing %d files\n", len(files))

	policy = bluemonday.UGCPolicy()

	for _, f := range files {
		wg.Add(1)
		go func(f string) {
			md, _ := ioutil.ReadFile("test/" + f)
			mdParseStream(&md)
		}(f.Name())
	}
	wg.Wait()
}

func mdParseStream(c *[]byte) {
	defer wg.Done()
	ioutil.WriteFile("test.html", policy.SanitizeBytes(blackfriday.MarkdownCommon(*c)), 0644)

}
