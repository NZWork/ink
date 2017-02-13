package main

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	"io/ioutil"
	"os"
	"sync"
)

var wg sync.WaitGroup

var policy *bluemonday.Policy

func mdStream(path string) {
	files, _ := ioutil.ReadDir(path)
	policy = bluemonday.UGCPolicy()

	var f os.FileInfo
	path += "/"
	for _, f = range files {
		wg.Add(1)
		go func(f string) {
			md, _ := ioutil.ReadFile(path + f)
			mdParseStream(&md)
		}(f.Name())
	}
	wg.Wait()
}

func mdParseStream(c *[]byte) {
	defer wg.Done()
	ioutil.WriteFile("test.html", policy.SanitizeBytes(blackfriday.MarkdownCommon(*c)), 0644)
}
