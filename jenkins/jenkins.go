package jenkins

import (
	"fmt"
	"github.com/bndr/gojenkins"
)

func InitClient() {

	fmt.Print("\n\n********************\nStarting Jenkins\n********************\n")

	jenkins, err := gojenkins.CreateJenkins("http://xdockerbuild.harebrained-apps.com/", "stevebargelt", "fd9faa2e4a5c1e99d99ed5d1c6dea062").Init()
	if err != nil {
		panic("Something Went Wrong")
	}

	jobs, err := jenkins.GetAllJobs()
	if err != nil {
		//return nil, err
		panic(err)
	}
	if len(jobs) == 0 {
		//create error obj
		//return nil, err.
		fmt.Println("Get All Jobs Failed. Jobs Count = ", len(jobs))
	}

	job, err := jenkins.GetJob("testjob")
	if err != nil {
		panic("Job Does Not Exist")
	}

	build, err := job.GetLastSuccessfulBuild()

	fmt.Println("Last run =", build.GetDuration()/1000, "seconds")

}
