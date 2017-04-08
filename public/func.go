package public

func GetRealPathByRepo(repo string) string {
	return WorkDir + repo
}

func GetRepoOriginPath(repo string) string {
	return GetRealPathByRepo(repo) + "/origin/"
}

func GetRepoParsedPath(repo string) string {
	return GetRealPathByRepo(repo) + "/parsed/"
}
