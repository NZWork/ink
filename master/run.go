package master

import (
	"fmt"
	"ink/public"
	"log"
	"net/http"
	"net/http/httputil"
)

const InvalidAction = `{"stat": 0, "err_msg": "invalid action"}`
const ParsedSuccess = `{"stat": 1, "repo_id": %s, "files": %d, "response_time": %f}`

func Run() {
	log.Println("running as master")
	public.MQConnect()

	http.HandleFunc("/parse", taskHandler)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

func Close() {
	public.MQClose()
	log.Println("master closed")
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	// newTask("test")
	r.ParseForm()
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))
	repo := r.FormValue("repo")

	if r.FormValue("auth") == public.APIKey && repo != "" {
		files, tasks, responseTime := newTask(repo)
		log.Printf("[%s] %d files cost %f ms to parse using %d tasks", repo, files, responseTime, tasks)
		fmt.Fprintf(w, ParsedSuccess, repo, files, responseTime)
	} else {
		fmt.Fprint(w, InvalidAction)
	}
}
