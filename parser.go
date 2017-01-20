package main

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	"fmt"
	"io/ioutil"
	"path/filepath"
)

var rootPath = "test"

func parseMarkdownToHTML(md []byte, f string, fi chan bool) {
	unsafe := blackfriday.MarkdownCommon(md)
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	ioutil.WriteFile(f+".html", html, 0644)
	fi <- true
}

func readMarkdownFiles() {
	files, _ := filepath.Glob(filepath.Join(rootPath, "*.md"))

	c := len(files)
	fi := make(chan bool, c)

	fmt.Printf("%d files in total", c)

	for _, f := range files {
		md, _ := ioutil.ReadFile(f)
		go parseMarkdownToHTML(md, f, fi)
	}
	for i := 0; i < c; i++ {
		<-fi
	}
}

func mdReadStream()  {}
func mdParseStream() {}
func mdWriteStream() {}
