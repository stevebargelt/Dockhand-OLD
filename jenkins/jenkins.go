package jenkins

import "github.com/bndr/gojenkins"

//InitClient initializes the Jenkins client - connects to jenkins instance
func InitClient(jenkinsURL string, username string, password string) (*gojenkins.Jenkins, error) {
	jenkins, err := gojenkins.CreateJenkins(jenkinsURL, username, password).Init()
	return jenkins, err
}
