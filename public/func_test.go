package public

import (
	"log"
	"testing"
)

func TestGetPath(t *testing.T) {
	log.Println(GetRepoOriginPath("test"))
}

func TestTimer(t *testing.T) {
	s := timerStart()
	log.Println(timerStop(s))
}
