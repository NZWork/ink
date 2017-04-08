package worker

import (
	"ink/public"
	"io/ioutil"
)

func readFile(repo, file string) ([]byte, error) {
	return ioutil.ReadFile(public.GetRepoOriginPath(repo) + file)
}

func writeFile(repo, file string, content []byte) {
	ioutil.WriteFile(public.GetRepoParsedPath(repo)+file, content, 0644)
}
