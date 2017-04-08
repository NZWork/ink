package public

import (
	"encoding/json"
	"fmt"
)

type Task struct {
	Repo  string   `json:"repo"`
	Files []string `json:"works"`
}

func (t *Task) Debug() string {
	return fmt.Sprintf("[%s] files: %v", t.Repo, t.Files)
}

func (t *Task) JSON() ([]byte, error) {
	return json.Marshal(t)
}
