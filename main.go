package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bndr/gojenkins"
	"github.com/docker/docker/api/types/container"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stevebargelt/Dockhand/docker"
	"github.com/stevebargelt/Dockhand/jenkins"
)

var (
	dockerHostURL    = flag.String("dockerurl", "tcp://abs.harebrained-apps.com:2376", "the full address to the docker Jenkins host: tcp://<address>:<port>")
	dockerTLSFolder  = flag.String("dockertlsfolder", "/users/steve/tlsBuild/", "Path to PEM encoded certificate, Key and CA for secure Docker TLS communication")
	certFile         = flag.String("cert", "/users/steve/tlsBuild/cert.pem", "Path to a PEM encoded certificate file.")
	keyFile          = flag.String("key", "/users/steve/tlsBuild/key.pem", "Path to a PEM encoded private key file.")
	caFile           = flag.String("CA", "/users/steve/tlsBuild/ca.pem", "Path to a PEM encoded CA certificate file.")
	registryURL      = flag.String("registry", "https://abs-registry.harebrained-apps.com", "The URL of the registry of where to find the image we are testing.")
	registryUser     = flag.String("registryuser", "absadmin", "A user with rights to the registry we are pulling the test image from.")
	registryPassword = flag.String("registrypassword", "correcthorsebatteystaple", "The password of the registry user")
	imageName        = flag.String("imagename", "dockerbuild.harebrained-apps.com/jenkins-slavedotnet", "The name of the image we are testing.")
	cloudName        = flag.String("cloudname", "AzureJenkins", "The name of the cloud configuration in Jenkins to use.")
	label            = flag.String("label", "TeamBargelt_DotNetCore23", "The name of the label to use in Jenkins")
	jenkinsURL       = flag.String("jenkinsurl", "http://dockerbuild.harebrained-apps.com", "The URL of the Jenkins Master.")
	jenkinsUser      = flag.String("jenkinsuser", "stevebargelt", "A user with rights to the registry we are pulling the test image from.")
	jenkinsPassword  = flag.String("jenkinspassword", "correcthorsebatteystaple", "The password of the registry user")
	repoURL          = flag.String("repourl", "https://github.com/stevebargelt/simpleDotNet.git", "The repo url.")
	configFile       = flag.String("config", "dockhand.yml", "A config file to use.")

	dockerClient  *docker.Host
	jenkinsClient *gojenkins.Jenkins
)

func main() {

	var err error

	//flag.Parse()

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	fmt.Print("\n********************\n Docker Image and Container Verification Process\n********************\n")
	connectToDockerHost()

	err = dockerClient.BuildDockerImage(*imageName, *repoURL)
	if err != nil {
		panic(err)
	}

	err = dockerClient.PushDockerImage(*imageName, *registryUser, *registryPassword, *registryURL)
	if err != nil {
		panic(err)
	}

	pullDockerImage()
	newContainer, err := createDockerContainer()
	if err != nil {
		panic(err)
	}
	startDockerContainer(newContainer)
	testResult, err := testDockerContainer(newContainer)
	if err != nil {
		panic(err)
	}
	if !testResult {
		fmt.Println("test failed!")
		fmt.Println("Need to handle this (still remove contianer but then quit)!")

	}
	removeDockerContainer(newContainer)

	//TODO: If tests pass we want to pull the image to all hosts in SWARM to save time at build
	//		(advice from Maxfield Stewart of Riot Games )

	fmt.Print("\n\n********************\nAdd Build To Jenkins\n********************\n")

	fmt.Print("Checking that label ", *label, " is unique...")
	labelIsUnique, err := jenkins.CheckLabelIsUnique(*jenkinsURL, *cloudName, *label, *jenkinsUser, *jenkinsPassword)
	if err != nil {
		panic(err)
	}
	if labelIsUnique {
		fmt.Println(" it is unique, continuing.")
	} else {
		fmt.Println(" it is NOT unique! build labels mut be unique.")
		fmt.Println("The label", *label, "is not unique in Jenkins at", *jenkinsURL, "cannot create this build.")
		fmt.Println("exiting...")
		os.Exit(1)
	}

	fmt.Print("Creating docker slave template in ", *cloudName, "... ")
	slaveTemplateCreated, err := jenkins.CreateDockerTemplate(*jenkinsURL, *cloudName, *label, *imageName, *jenkinsUser, *jenkinsPassword)
	if err != nil {
		panic(err)
	}
	if slaveTemplateCreated {
		fmt.Println("success!")
	} else {
		fmt.Println("Failed. Exiting.")
		os.Exit(1)
	}

	fmt.Print("Connecting to Jenkins... ")
	jenkinsClient, err := jenkins.InitClient(*jenkinsURL, *jenkinsUser, *jenkinsPassword)
	if err != nil {
		fmt.Println("failed. Exiting")
		os.Exit(1)
	}
	if jenkinsClient == nil {
		fmt.Println("Failed. Jenkins object is nil. Exiting.")
		os.Exit(1)
	}
	fmt.Println("success. Connected.")
	fmt.Print("Adding jenkins job... ")
	//TODO: Move this to a config file...
	configString := `<?xml version='1.0' encoding='UTF-8'?>
<flow-definition plugin="workflow-job@2.6">
  <actions/>
  <description></description>
  <keepDependencies>false</keepDependencies>
  <properties>
    <com.coravy.hudson.plugins.github.GithubProjectProperty plugin="github@1.21.1">
      <projectUrl>` + *repoURL + `</projectUrl>
      <displayName></displayName>
    </com.coravy.hudson.plugins.github.GithubProjectProperty>
    <org.jenkinsci.plugins.workflow.job.properties.PipelineTriggersJobProperty>
      <triggers>
        <hudson.triggers.TimerTrigger>
          <spec>H */3 * * *</spec>
        </hudson.triggers.TimerTrigger>
        <com.cloudbees.jenkins.GitHubPushTrigger plugin="github@1.21.1">
          <spec></spec>
        </com.cloudbees.jenkins.GitHubPushTrigger>
      </triggers>
    </org.jenkinsci.plugins.workflow.job.properties.PipelineTriggersJobProperty>
  </properties>
  <definition class="org.jenkinsci.plugins.workflow.cps.CpsScmFlowDefinition" plugin="workflow-cps@2.17">
    <scm class="hudson.plugins.git.GitSCM" plugin="git@3.0.0">
      <configVersion>2</configVersion>
      <userRemoteConfigs>
        <hudson.plugins.git.UserRemoteConfig>
          <url>` + *repoURL + `</url>
        </hudson.plugins.git.UserRemoteConfig>
      </userRemoteConfigs>
      <branches>
        <hudson.plugins.git.BranchSpec>
          <name>*/master</name>
        </hudson.plugins.git.BranchSpec>
      </branches>
      <doGenerateSubmoduleConfigurations>false</doGenerateSubmoduleConfigurations>
      <submoduleCfg class="list"/>
      <extensions/>
    </scm>
    <scriptPath>Jenkinsfile</scriptPath>
  </definition>
  <triggers/>
</flow-definition>`

	tempStr := *label + "_JOB"
	newJob, err := jenkinsClient.CreateJob(configString, tempStr)
	if err != nil {
		panic(err)
		//fmt.Println("failed. Exiting")
		//os.Exit(1)
	}
	fmt.Println("success! Created", newJob.GetName(), ".")

	fmt.Print("Kicking first build... ")
	m := make(map[string]string)
	jobResult, err := newJob.InvokeSimple(m)
	if err != nil {
		panic(err)
	}
	if jobResult == true {
		fmt.Println("build Success.")
	} else {
		fmt.Println("build fail.")
	}

}

func connectToDockerHost() {

	var err error
	fmt.Print("Connecting to dockerhost... ")
	dockerClient, err = docker.New(*dockerHostURL, *dockerTLSFolder)
	if err != nil {
		panic(err)
	}
	if dockerClient != nil {
		fmt.Println("success!")
	} else {
		fmt.Println("failed. client is nil. Exiting")
		os.Exit(1)
	}

}

func pullDockerImage() {

	var err error

	fmt.Print("Pulling ", *imageName, " from registry ", *registryURL)
	newImage, err := dockerClient.GetDockerImage(*imageName, *registryUser, *registryPassword, *registryURL)
	if err != nil {
		panic(err)
	}
	fmt.Println(" success!\nPulled Image ID:", newImage.ID[7:19])

}

func createDockerContainer() (container.ContainerCreateCreatedBody, error) {

	fmt.Print("Creating continer from ", *imageName, "...")
	//TODO: unique value here for container name? Add GUID? Add LabelName?
	//TODO: create process to kill all containers that start with DockhandTesting??
	newContianer, err := dockerClient.CreateContainer(*imageName, "DockhandTesting"+*label)
	if err != nil {
		return *newContianer, err
	}
	fmt.Println(" success.\nContiner ID: ", newContianer.ID[0:11])
	return *newContianer, nil

}

func startDockerContainer(container container.ContainerCreateCreatedBody) {

	var err error
	fmt.Print("Starting continer ", container.ID[0:11], "...")
	err = dockerClient.StartContainer(container.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println("success.")
}

func testDockerContainer(container container.ContainerCreateCreatedBody) (bool, error) {

	//TODO: Write tests to make sure the container fits company standards
	fmt.Print("Testing continer ", container.ID[0:11], "...")
	containerInfo, err := dockerClient.ContainerInspect(container.ID)
	if err != nil {
		return false, err
	}
	// fmt.Println("CONTAINERINFO:")
	// fmt.Println("Name:", containerInfo.Name)
	// fmt.Println("Status:", containerInfo.State.Status)
	// fmt.Println("Exit Code:", containerInfo.State.ExitCode)
	if containerInfo.State.ExitCode != 0 {
		fmt.Println("failed. Exit code must be 0.")
		return false, nil
	}
	fmt.Println("success.")
	return true, nil
}

func removeDockerContainer(container container.ContainerCreateCreatedBody) error {

	var err error
	fmt.Print("Removing continer ", container.ID[0:11], "...")
	err = dockerClient.ContainerRemove(container.ID)
	if err != nil {
		return err
	}
	fmt.Println("success.")
	return nil
}
