package public

import "log"

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
