package jenkins

import "github.com/bndr/gojenkins"

//Host - a Docker Host Client
// type Host struct {
// 	URL        string
// 	JenkinsCli *gojenkins.Jenkins
// 	//Logf   LogfCallback
// }

//InitClient initializes the Jenkins client - connects to jenkins instance
func InitClient(jenkinsURL string, username string, password string) (*gojenkins.Jenkins, error) {
	jenkins, err := gojenkins.CreateJenkins(jenkinsURL, username, password).Init()
	return jenkins, err
}

// func (j *Host) CreateJob2(labelName string) {

// 	configString := `<?xml version='1.0' encoding='UTF-8'?>
// 						<flow-definition plugin="workflow-job@2.6">
// 						<description></description>
// 						<keepDependencies>false</keepDependencies>
// 						<properties>
// 							<org.jenkinsci.plugins.workflow.job.properties.PipelineTriggersJobProperty>
// 							<triggers/>
// 							</org.jenkinsci.plugins.workflow.job.properties.PipelineTriggersJobProperty>
// 						</properties>
// 						<definition class="org.jenkinsci.plugins.workflow.cps.CpsFlowDefinition" plugin="workflow-cps@2.17">
// 							<script>node (&apos;` + labelName + `&apos;) {

// 						stage (&apos;Stage 1&apos;) {
// 							sh &apos;echo &quot;Hello World!&quot;&apos;
// 						}
// 						}</script>
// 							<sandbox>true</sandbox>
// 						</definition>
// 						<triggers/>
// 						</flow-definition>`

// 	j.JenkinsCli.CreateJob(configString, "someNewJobsName")

// }

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
