package main

import (
	"flag"
	"fmt"

	"github.com/docker/docker/client"
	"github.com/stevebargelt/harbormaster/docker"
	"github.com/stevebargelt/harbormaster/jenkins"
)

var (
	dockerHost       = flag.String("dockerhost", "tcp://dockerbuild.harebrained-apps.com:2376", "the full address to the docker Jenkins host: tcp://<address>:<port>")
	certFile         = flag.String("cert", "/users/steve/tlsBlog/cert.pem", "A PEM eoncoded certificate file.")
	keyFile          = flag.String("key", "/users/steve/tlsBlog/key.pem", "A PEM encoded private key file.")
	caFile           = flag.String("CA", "/users/steve/tlsBlog/ca.pem", "A PEM eoncoded CA's certificate file.")
	registryURL      = flag.String("registry", "https://dockerbuild.harebrained-apps.com", "The URL of the registry of where to find the image we are testing.")
	registryUser     = flag.String("registryuser", "dockerUser", "A user with rights to the registry we are pulling the test image from.")
	registryPassword = flag.String("registrypassword", "steel2000", "The password of the registry user")
	imageName        = flag.String("imagename", "dockerbuild.harebrained-apps.com/jenkins-slavedotnet", "The name of the image we are testing.")
	cloudName        = flag.String("cloudname", "dockerAzure", "The name of the cloud configuration in Jenkins to use.")
	label            = flag.String("label", "testslave", "The name of the label to use in Jenkins")

	dockerClient *client.Client
)

func main() {

	flag.Parse()
	var err error
	dockerClient, err = docker.InitClient(*dockerHost, *certFile, *keyFile, *caFile)
	if err != nil {
		panic(err)
	}
	newImage, err := docker.GetDockerImage(dockerClient, *imageName, *registryUser, *registryPassword, *registryURL)
	if err != nil {
		panic(err)
	}
	fmt.Println("Image ID:", newImage.ID)

	jenkins.InitClient()

}
