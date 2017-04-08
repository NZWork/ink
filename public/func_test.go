package public

import (
	"log"
	"testing"
)

func TestGetPath(t *testing.T) {
	log.Println(GetRepoOriginPath("test"))
}
