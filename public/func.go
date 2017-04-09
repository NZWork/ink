package public

import (
	"log"
	"time"
)

func GetRealPathByRepo(repo string) string {
	return WorkDir + repo
}

func GetRepoOriginPath(repo string) string {
	return GetRealPathByRepo(repo) + "/origin/"
}

func GetRepoParsedPath(repo string) string {
	return GetRealPathByRepo(repo) + "/parsed/"
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func TimerStart() int64 {
	return time.Now().UnixNano()
}

func TimerStop(start int64) float64 {
	return float64(time.Now().UnixNano()-start) / 1000.0 / 1000.0
}
