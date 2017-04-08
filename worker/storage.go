package worker

import "io/ioutil"

func readFile(file string) ([]byte, error) {
	return ioutil.ReadFile(repo + "/origin/" + file)
}

func writeFile(content []byte, file string) {
	ioutil.WriteFile(repo+"/parsed/"+file, content, 0644)
}
