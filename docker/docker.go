package docker

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	dockerClient "github.com/docker/docker/client"
	"golang.org/x/net/context"
)

//Host - a Docker Host Client
type Host struct {
	URL       string
	DockerCli *dockerClient.Client
	//Logf   LogfCallback
}

//New - creates a new Docker Host with given docker host URL and TLS cert file paths
func New(URL string, certFile string, keyFile string, caFile string) (*Host, error) {

	url := strings.TrimSuffix(URL, "/")

	// Load client cert
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	return newClientFromTransport(url, transport)
}

//NewInsecure - creates a new Host with given docker host: this is not secure... please know what you are doing!
func NewInsecure(URL string) (*Host, error) {

	url := strings.TrimSuffix(URL, "/")

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	return newClientFromTransport(url, transport)
}

func newClientFromTransport(url string, transport http.RoundTripper) (*Host, error) {

	httpCli := &http.Client{Transport: transport}

	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := dockerClient.NewClient(url, "v1.24", httpCli, defaultHeaders)
	if err != nil {
		return nil, err
	}

	dockerHost := &Host{
		URL:       url,
		DockerCli: cli,
	}

	return dockerHost, nil

}

//GetDockerImage given imageName and the registry information returns a docker image
func (d *Host) GetDockerImage(imageName, registryUsername, registryPassword, registryURL string) (*types.ImageInspect, error) {

	//TODO: log
	//fmt.Println("Looking for image", imageName, "...")
	image, _, imageErr := d.DockerCli.ImageInspectWithRaw(context.TODO(), imageName)
	newImage, err := d.pullImage(imageName, registryUsername, registryPassword, registryURL)
	if err != nil {
		if imageErr == nil {
			fmt.Println("Cannot pull the latest version of image", imageName, ":", err)
			fmt.Println("Locally found image will be used instead.")
			return &image, nil
		}
		return nil, err
	}
	return newImage, nil
}

func (d *Host) pullImage(imageName, registryUsername, registryPassword, registryURL string) (*types.ImageInspect, error) {

	//TODO: log
	//fmt.Println("Pulling docker image", imageName, "...")

	ref := imageName
	// Add :latest to limit the download results
	if !strings.ContainsAny(ref, ":@") {
		ref += ":latest"
	}

	authConfig := types.AuthConfig{Username: registryUsername, Password: registryPassword, ServerAddress: registryURL}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(authConfig); err != nil {
		return nil, err
	}
	encodedAuth := base64.URLEncoding.EncodeToString(buf.Bytes())

	options := types.ImagePullOptions{RegistryAuth: encodedAuth}
	readCloser, err := d.DockerCli.ImagePull(context.Background(), ref, options)
	if err != nil {
		return nil, err
	}
	defer readCloser.Close()

	if _, err := io.Copy(ioutil.Discard, readCloser); err != nil {
		return nil, fmt.Errorf("Failed to pull image: %s: %s", ref, err)
	}

	image, _, err := d.DockerCli.ImageInspectWithRaw(context.Background(), imageName)
	return &image, err
}

//CreateContainer - creates a container named containerName given an imageName
func (d *Host) CreateContainer(imageName, containerName string) (*types.ContainerCreateResponse, error) {

	container, err := d.DockerCli.ContainerCreate(context.Background(), &container.Config{Image: imageName}, nil, nil, containerName)
	if err != nil {
		return nil, err
	}

	return &container, nil

}

//StartContainer - runs a container named containerName given an imageName
func (d *Host) StartContainer(containerID string) error {

	err := d.DockerCli.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	return nil

}

// ContainerInspect returns the deatils of a container given a containerID
func (d *Host) ContainerInspect(containerID string) (types.ContainerJSON, error) {

	return d.DockerCli.ContainerInspect(context.Background(), containerID)

}

//ContainerRemove - removes a container give a containerID
func (d *Host) ContainerRemove(id string) error {

	return d.DockerCli.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{RemoveVolumes: true, Force: true})

}
