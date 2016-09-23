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

// Stuff I used at the start to test the connection to the Docker Host - might come in handy some day.
// TODO: delete after adding to git
// func randomTestStuff() {

// 	fmt.Print("\n\n********************\nList Images\n********************\n")
// 	imageOptions := types.ImageListOptions{All: true}
// 	images, err := cli.ImageList(context.Background(), imageOptions)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, i := range images {
// 		fmt.Println(i.ID, i.Labels)
// 	}

// 	fmt.Print("\n\n********************\nList Containers\n********************\n")
// 	options := types.ContainerListOptions{All: true}
// 	containers, err := cli.ContainerList(context.Background(), options)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, c := range containers {
// 		fmt.Println(c.ID, c.Names)
// 	}

// 	fmt.Print("\n\n********************\nOther Info???\n********************\n")
// 	fmt.Println("cli type=", reflect.TypeOf(cli))

// }
