package jenkins

import "github.com/bndr/gojenkins"

//InitClient initializes the Jenkins client - connects to jenkins instance
func InitClient(jenkinsURL string, username string, password string) (*gojenkins.Jenkins, error) {
	jenkins, err := gojenkins.CreateJenkins(jenkinsURL, username, password).Init()
	return jenkins, err
}

func stuff() {

	// jobs, err := jenkins.GetAllJobs()
	// if err != nil {
	// 	return nil, err
	// }
	// if len(jobs) == 0 {
	// 	//create error obj
	// 	//return nil, err.
	// 	fmt.Println("Get All Jobs Failed. Jobs Count = ", len(jobs))
	// }

	// job, err := jenkins.GetJob("testjob")
	// if err != nil {
	// 	panic("Job Does Not Exist")
	// }
	// build, err := job.GetLastSuccessfulBuild()

	// fmt.Println("Last run =", build.GetDuration()/1000, "seconds")

}
