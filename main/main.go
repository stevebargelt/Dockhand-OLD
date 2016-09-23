package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bndr/gojenkins"
	"github.com/stevebargelt/harbormaster/docker"
	"github.com/stevebargelt/harbormaster/jenkins"
)

var (
	dockerHostURL    = flag.String("dockerurl", "tcp://dockerbuild.harebrained-apps.com:2376", "the full address to the docker Jenkins host: tcp://<address>:<port>")
	certFile         = flag.String("cert", "/users/steve/tlsBlog/cert.pem", "A PEM eoncoded certificate file.")
	keyFile          = flag.String("key", "/users/steve/tlsBlog/key.pem", "A PEM encoded private key file.")
	caFile           = flag.String("CA", "/users/steve/tlsBlog/ca.pem", "A PEM eoncoded CA's certificate file.")
	registryURL      = flag.String("registry", "https://dockerbuild.harebrained-apps.com", "The URL of the registry of where to find the image we are testing.")
	registryUser     = flag.String("registryuser", "dockerUser", "A user with rights to the registry we are pulling the test image from.")
	registryPassword = flag.String("registrypassword", "notARealPAssword", "The password of the registry user")
	imageName        = flag.String("imagename", "dockerbuild.harebrained-apps.com/jenkins-slavedotnet", "The name of the image we are testing.")
	cloudName        = flag.String("cloudname", "AzureJenkins", "The name of the cloud configuration in Jenkins to use.")
	label            = flag.String("label", "TeamBargelt_DotNetCore", "The name of the label to use in Jenkins")
	jenkinsURL       = flag.String("jenkins", "http://dockerbuild.harebrained-apps.com", "The URL of the Jenkins Master.")
	jenkinsUser      = flag.String("jenkinsuser", "stevebargelt", "A user with rights to the registry we are pulling the test image from.")
	jenkinsPassword  = flag.String("jenkinspassword", "notARealPAssword", "The password of the registry user")

	dockerClient  *docker.Host
	jenkinsClient *gojenkins.Jenkins
)

func main() {

	flag.Parse()
	var err error

	fmt.Print("\n********************\nPulling ", *imageName, " from registry ", *registryURL, "\n********************\n")
	fmt.Print("Connecting to dockerhost... ")

	dockerClient, err = docker.New(*dockerHostURL, *certFile, *keyFile, *caFile)
	if err != nil {
		panic(err)
	}
	if dockerClient != nil {
		fmt.Println("success!")
	} else {
		fmt.Println("failed. client is nil. Exiting")
		os.Exit(1)
	}

	fmt.Print("Pulling image to docker host... ")
	newImage, err := dockerClient.GetDockerImage(*imageName, *registryUser, *registryPassword, *registryURL)
	if err != nil {
		panic(err)
	}
	fmt.Println("success! Pulled Image ID:", newImage.ID)

	//Create a contianer from the image and run container
	//TODO

	//Test the  container
	//TODO

	//Remove container
	//TODO

	fmt.Print("\n\n********************\nAdd Slave Template to Jenkins\n********************\n")
	//Add Docker Slave Template to Jenkins
	fmt.Print("Checking that label ", *label, " is unique...")
	labelIsUnique, err := jenkins.CheckLabelIsUnique(*jenkinsURL, *cloudName, *label, *jenkinsUser, *jenkinsPassword)
	if err != nil {
		panic(err)
	}

	if labelIsUnique {
		fmt.Println(" it is unique, continuing.")
	} else {
		fmt.Println(" it is NOT unique! Labels mut be unique.")
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

	fmt.Print("\n********************\nAdding job to Jenkins\n********************\n")
	//TODO
	fmt.Print("Connecting to Jenkins... ")
	jenkinsClient, err := jenkins.InitClient(*jenkinsURL, *jenkinsURL, *jenkinsPassword)
	if err != nil {
		panic(err)
	}
	if jenkinsClient == nil {
		fmt.Println("Failed. Jenkisn object is nil. Exiting.")
		os.Exit(1)

	}

}
