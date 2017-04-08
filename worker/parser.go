package worker

import (
	"ink/public"
	"sync"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

var policy = bluemonday.UGCPolicy()

func mdStream(task *public.Task) {
	var wg sync.WaitGroup
	var f string
	var md []byte

	for _, f = range task.Files {
		wg.Add(1)
		go func(f string) {
			md, _ = readFile(task.Repo, f)
			mdParseStream(task.Repo, f, &md, &wg)
		}(f)
	}
	wg.Wait()
}

func mdParseStream(repo, f string, c *[]byte, wg *sync.WaitGroup) {
	defer wg.Done()
	writeFile(repo, f, policy.SanitizeBytes(blackfriday.MarkdownCommon(*c)))
}
