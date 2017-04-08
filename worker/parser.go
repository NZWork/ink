package worker

import (
	"sync"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

var wg sync.WaitGroup

var policy = bluemonday.UGCPolicy()
var f string
var md []byte

func mdStream(files []string) {
	for _, f = range files {
		wg.Add(1)
		go func(f string) {
			md, _ = readFile(f)
			mdParseStream(f, &md)
		}(f)
	}
	wg.Wait()
}

func mdParseStream(f string, c *[]byte) {
	defer wg.Done()
	writeFile(policy.SanitizeBytes(blackfriday.MarkdownCommon(*c)), f)
}
